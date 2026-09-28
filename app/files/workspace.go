package files

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidPath = errors.New("ruta inválida")
	ErrBinaryFile  = errors.New("archivo no editable")
	ErrTooLarge    = errors.New("archivo demasiado grande")
)

type Entry struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Type     string    `json:"type"`
	Ext      string    `json:"ext,omitempty"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
	Editable bool      `json:"editable"`
}

type FileContent struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Ext      string `json:"ext,omitempty"`
	Language string `json:"language"`
	Size     int64  `json:"size"`
	Content  string `json:"content"`
}

type Workspace struct {
	root     string
	maxBytes int64
}

func NewWorkspace(root string, maxBytes int64) (*Workspace, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("%w: root vacío", ErrInvalidPath)
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	absolute = filepath.Clean(absolute)
	if maxBytes <= 0 {
		maxBytes = 2 * 1024 * 1024
	}
	return &Workspace{root: absolute, maxBytes: maxBytes}, nil
}

func (w *Workspace) Root() string { return w.root }

func (w *Workspace) EnsureDemo() error {
	if err := os.MkdirAll(w.root, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(w.root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() != ".gitkeep" {
			return nil
		}
	}

	seed := map[string]string{
		"README.md":            "# Carpeta demo Axolotl\n\nEsta carpeta está controlada por el backend Go.\n\nPuedes crear, editar y eliminar archivos sin salir de este espacio.\n",
		"src/main.go":          "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hola desde Axolotl\")\n}\n",
		"web/app.js":           "const app = {\n  name: \"Axolotl\",\n  mode: \"demo-files\"\n};\n\nconsole.log(app);\n",
		"config/settings.json": "{\n  \"theme\": \"axolotl\",\n  \"environment\": \"demo\",\n  \"features\": [\"explorer\", \"editor\", \"syntax-highlight\"]\n}\n",
		"notes/todo.txt":       "Pendientes de ejemplo\n- Conectar un servicio real\n- Añadir permisos por carpeta\n- Crear más archivos desde el explorador\n",
	}

	for rel, content := range seed {
		full, _, err := w.resolve(rel)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func (w *Workspace) List(rel string) ([]Entry, error) {
	full, clean, err := w.resolve(rel)
	if err != nil {
		return nil, err
	}
	if err := w.rejectSymlinks(full); err != nil {
		return nil, err
	}
	info, err := os.Stat(full)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%w: no es carpeta", ErrInvalidPath)
	}
	dirEntries, err := os.ReadDir(full)
	if err != nil {
		return nil, err
	}
	out := make([]Entry, 0, len(dirEntries))
	for _, item := range dirEntries {
		if clean == "" && item.Name() == ".gitkeep" {
			continue
		}
		itemPath := filepath.Join(full, item.Name())
		itemInfo, err := os.Lstat(itemPath)
		if err != nil {
			return nil, err
		}
		if itemInfo.Mode()&os.ModeSymlink != 0 {
			continue
		}
		relPath := item.Name()
		if clean != "" {
			relPath = filepath.ToSlash(filepath.Join(clean, item.Name()))
		}
		entryType := "file"
		editable := false
		if itemInfo.IsDir() {
			entryType = "folder"
		} else if itemInfo.Mode().IsRegular() {
			editable = itemInfo.Size() <= w.maxBytes
		}
		out = append(out, Entry{Name: item.Name(), Path: filepath.ToSlash(relPath), Type: entryType, Ext: strings.TrimPrefix(strings.ToLower(filepath.Ext(item.Name())), "."), Size: itemInfo.Size(), Modified: itemInfo.ModTime(), Editable: editable})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Type != out[j].Type {
			return out[i].Type == "folder"
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func (w *Workspace) Read(rel string) (FileContent, error) {
	full, clean, err := w.resolve(rel)
	if err != nil {
		return FileContent{}, err
	}
	if clean == "" {
		return FileContent{}, fmt.Errorf("%w: se requiere archivo", ErrInvalidPath)
	}
	if err := w.rejectSymlinks(full); err != nil {
		return FileContent{}, err
	}
	info, err := os.Stat(full)
	if err != nil {
		return FileContent{}, err
	}
	if !info.Mode().IsRegular() {
		return FileContent{}, fmt.Errorf("%w: no es archivo", ErrInvalidPath)
	}
	if info.Size() > w.maxBytes {
		return FileContent{}, ErrTooLarge
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return FileContent{}, err
	}
	if bytes.IndexByte(data, 0) >= 0 || !utf8.Valid(data) {
		return FileContent{}, ErrBinaryFile
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(info.Name())), ".")
	return FileContent{Name: info.Name(), Path: filepath.ToSlash(clean), Ext: ext, Language: LanguageForPath(info.Name()), Size: info.Size(), Content: string(data)}, nil
}

func (w *Workspace) Save(rel, content string) error {
	if int64(len([]byte(content))) > w.maxBytes {
		return ErrTooLarge
	}
	full, clean, err := w.resolve(rel)
	if err != nil {
		return err
	}
	if clean == "" {
		return fmt.Errorf("%w: se requiere archivo", ErrInvalidPath)
	}
	if err := w.rejectSymlinks(full); err != nil {
		return err
	}
	info, err := os.Stat(full)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%w: no es archivo", ErrInvalidPath)
	}
	return os.WriteFile(full, []byte(content), info.Mode().Perm())
}

func (w *Workspace) Create(parent, name, kind string) (Entry, error) {
	if err := validateName(name); err != nil {
		return Entry{}, err
	}
	if kind != "file" && kind != "folder" {
		return Entry{}, fmt.Errorf("%w: tipo inválido", ErrInvalidPath)
	}
	parentFull, parentClean, err := w.resolve(parent)
	if err != nil {
		return Entry{}, err
	}
	if err := w.rejectSymlinks(parentFull); err != nil {
		return Entry{}, err
	}
	parentInfo, err := os.Stat(parentFull)
	if err != nil {
		return Entry{}, err
	}
	if !parentInfo.IsDir() {
		return Entry{}, fmt.Errorf("%w: parent no es carpeta", ErrInvalidPath)
	}
	targetRel := name
	if parentClean != "" {
		targetRel = filepath.ToSlash(filepath.Join(parentClean, name))
	}
	targetFull, targetClean, err := w.resolve(targetRel)
	if err != nil {
		return Entry{}, err
	}
	if _, err := os.Lstat(targetFull); err == nil {
		return Entry{}, fs.ErrExist
	} else if !errors.Is(err, fs.ErrNotExist) {
		return Entry{}, err
	}
	if kind == "folder" {
		if err := os.Mkdir(targetFull, 0o755); err != nil {
			return Entry{}, err
		}
	} else {
		file, err := os.OpenFile(targetFull, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err != nil {
			return Entry{}, err
		}
		if err := file.Close(); err != nil {
			return Entry{}, err
		}
	}
	info, err := os.Stat(targetFull)
	if err != nil {
		return Entry{}, err
	}
	return Entry{Name: name, Path: filepath.ToSlash(targetClean), Type: kind, Ext: strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), "."), Size: info.Size(), Modified: info.ModTime(), Editable: kind == "file"}, nil
}

func (w *Workspace) Delete(rel string) error {
	full, clean, err := w.resolve(rel)
	if err != nil {
		return err
	}
	if clean == "" {
		return fmt.Errorf("%w: no se puede eliminar root", ErrInvalidPath)
	}
	if err := w.rejectSymlinks(full); err != nil {
		return err
	}
	return os.RemoveAll(full)
}

func (w *Workspace) resolve(rel string) (full string, clean string, err error) {
	raw := strings.TrimSpace(strings.ReplaceAll(rel, "\\", "/"))
	if strings.ContainsRune(raw, '\x00') {
		return "", "", ErrInvalidPath
	}
	raw = strings.TrimPrefix(raw, "/")
	clean = filepath.ToSlash(filepath.Clean(filepath.FromSlash(raw)))
	if clean == "." {
		clean = ""
	}
	osRelative := filepath.FromSlash(clean)
	if filepath.IsAbs(osRelative) || filepath.VolumeName(osRelative) != "" {
		return "", "", ErrInvalidPath
	}
	full = filepath.Clean(filepath.Join(w.root, osRelative))
	relCheck, err := filepath.Rel(w.root, full)
	if err != nil {
		return "", "", err
	}
	if relCheck == ".." || strings.HasPrefix(relCheck, ".."+string(os.PathSeparator)) {
		return "", "", ErrInvalidPath
	}
	if relCheck == "." {
		clean = ""
	} else {
		clean = filepath.ToSlash(relCheck)
	}
	return full, clean, nil
}

func (w *Workspace) rejectSymlinks(full string) error {
	rel, err := filepath.Rel(w.root, full)
	if err != nil {
		return err
	}
	if rel == "." {
		return nil
	}
	current := w.root
	for _, part := range strings.Split(rel, string(os.PathSeparator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%w: symlink no permitido", ErrInvalidPath)
		}
	}
	return nil
}

func validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("%w: nombre vacío", ErrInvalidPath)
	}
	if len(name) > 120 {
		return fmt.Errorf("%w: nombre demasiado largo", ErrInvalidPath)
	}
	if strings.ContainsAny(name, `/\\:*?"<>|`) {
		return fmt.Errorf("%w: caracteres no permitidos", ErrInvalidPath)
	}
	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return fmt.Errorf("%w: terminación no permitida", ErrInvalidPath)
	}
	for _, r := range name {
		if r < 32 {
			return fmt.Errorf("%w: caracteres de control", ErrInvalidPath)
		}
	}
	return nil
}

func LanguageForPath(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".go":
		return "go"
	case ".js", ".mjs", ".cjs":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".json":
		return "json"
	case ".html", ".htm", ".tmpl":
		return "html"
	case ".css":
		return "css"
	case ".md":
		return "markdown"
	case ".sh", ".bash":
		return "bash"
	case ".py":
		return "python"
	case ".sql":
		return "sql"
	case ".xml":
		return "xml"
	case ".yaml", ".yml":
		return "yaml"
	case ".php":
		return "php"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cpp", ".cc", ".cxx", ".hpp":
		return "cpp"
	case ".cs":
		return "csharp"
	case ".rs":
		return "rust"
	case ".rb":
		return "ruby"
	case ".toml":
		return "ini"
	default:
		return "plaintext"
	}
}
