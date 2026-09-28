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
	"unicode"
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

func (w *Workspace) Normalize(rel string) (string, error) {
	_, clean, err := w.resolve(rel)
	return clean, err
}

func (w *Workspace) EnsureDemo() error {
	if err := os.MkdirAll(w.root, 0o755); err != nil {
		return err
	}
	if err := w.verifyRoot(); err != nil {
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
		if itemInfo.Mode()&os.ModeSymlink != 0 || isReparsePoint(itemInfo) {
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
	if clean == "" || clean == ".gitkeep" {
		return fmt.Errorf("%w: se requiere archivo editable", ErrInvalidPath)
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
	if clean == "" || clean == ".gitkeep" {
		return fmt.Errorf("%w: no se puede eliminar elemento protegido", ErrInvalidPath)
	}
	if err := w.rejectSymlinks(full); err != nil {
		return err
	}
	return os.RemoveAll(full)
}

func (w *Workspace) resolve(rel string) (full string, clean string, err error) {
	if err := w.verifyRoot(); err != nil {
		return "", "", err
	}

	raw := rel
	if strings.ContainsRune(raw, '\x00') || strings.Contains(raw, "\\") || strings.HasPrefix(raw, "/") {
		return "", "", ErrInvalidPath
	}
	if len(raw) > 1024 || strings.Contains(raw, "%") {
		return "", "", ErrInvalidPath
	}
	if raw == "" {
		return w.root, "", nil
	}

	segments := strings.Split(raw, "/")
	if len(segments) > 32 {
		return "", "", ErrInvalidPath
	}
	for _, segment := range segments {
		if segment == "" {
			return "", "", ErrInvalidPath
		}
		if err := validateName(segment); err != nil {
			return "", "", err
		}
	}

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
		return w.verifyRoot()
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
		if info.Mode()&os.ModeSymlink != 0 || isReparsePoint(info) {
			return fmt.Errorf("%w: enlace/reparse point no permitido", ErrInvalidPath)
		}
	}
	return nil
}

func (w *Workspace) verifyRoot() error {
	info, err := os.Lstat(w.root)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || isReparsePoint(info) {
		return fmt.Errorf("%w: root no puede ser enlace", ErrInvalidPath)
	}

	evaluated, err := filepath.EvalSymlinks(w.root)
	if err != nil {
		return err
	}
	evaluated, err = filepath.Abs(evaluated)
	if err != nil {
		return err
	}
	if !samePath(filepath.Clean(evaluated), filepath.Clean(w.root)) {
		return fmt.Errorf("%w: root resuelto fuera de su ruta física", ErrInvalidPath)
	}
	return nil
}

func samePath(a, b string) bool {
	if os.PathSeparator == '\\' {
		return strings.EqualFold(a, b)
	}
	return a == b
}

func validateName(name string) error {
	if name == "" || name != strings.TrimSpace(name) || name == "." || name == ".." {
		return fmt.Errorf("%w: nombre inválido", ErrInvalidPath)
	}
	if len(name) > 120 || !utf8.ValidString(name) {
		return fmt.Errorf("%w: nombre demasiado largo o inválido", ErrInvalidPath)
	}
	if strings.ContainsAny(name, `/\\:*?"<>|%`) {
		return fmt.Errorf("%w: caracteres no permitidos", ErrInvalidPath)
	}
	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return fmt.Errorf("%w: terminación no permitida", ErrInvalidPath)
	}

	base := strings.ToUpper(strings.SplitN(name, ".", 2)[0])
	switch base {
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return fmt.Errorf("%w: nombre reservado", ErrInvalidPath)
	}

	for _, r := range name {
		if unicode.IsControl(r) ||
			(r >= 0x202A && r <= 0x202E) ||
			(r >= 0x2066 && r <= 0x2069) {
			return fmt.Errorf("%w: caracteres invisibles no permitidos", ErrInvalidPath)
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
