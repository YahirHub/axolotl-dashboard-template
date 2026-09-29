package server

import (
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

var ErrInvalidTask = errors.New("plantilla de tarea inválida")

type TaskTemplate struct {
	Name       string    `json:"name"`
	Filename   string    `json:"filename"`
	Kind       string    `json:"kind"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modified_at"`
	Content    string    `json:"content,omitempty"`
}

type CreateTaskInput struct {
	Name             string
	Kind             string
	Description      string
	ExecStart        string
	WorkingDirectory string
	User             string
	Restart          string
	After            string
	OnCalendar       string
	TargetUnit       string
	Persistent       bool
}

type TaskStore struct {
	root string
}

func NewTaskStore(root string) (*TaskStore, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("%w: root vacío", ErrInvalidTask)
	}

	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &TaskStore{root: filepath.Clean(absolute)}, nil
}

func (s *TaskStore) Ensure() error {
	return os.MkdirAll(s.root, 0o755)
}

func (s *TaskStore) List() ([]TaskTemplate, error) {
	if err := s.Ensure(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.root)
	if err != nil {
		return nil, err
	}

	out := make([]TaskTemplate, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		kind := taskKindFromFilename(filename)
		if kind == "" {
			continue
		}

		full := filepath.Join(s.root, filename)
		info, err := os.Lstat(full)
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			continue
		}

		out = append(out, TaskTemplate{
			Name:       strings.TrimSuffix(filename, filepath.Ext(filename)),
			Filename:   filename,
			Kind:       kind,
			Size:       info.Size(),
			ModifiedAt: info.ModTime(),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Filename) < strings.ToLower(out[j].Filename)
	})
	return out, nil
}

func (s *TaskStore) Read(filename string) (TaskTemplate, error) {
	full, clean, kind, err := s.resolveFilename(filename)
	if err != nil {
		return TaskTemplate{}, err
	}

	info, err := os.Lstat(full)
	if err != nil {
		return TaskTemplate{}, err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return TaskTemplate{}, ErrInvalidTask
	}
	if info.Size() > 256*1024 {
		return TaskTemplate{}, fmt.Errorf("%w: archivo demasiado grande", ErrInvalidTask)
	}

	data, err := os.ReadFile(full)
	if err != nil {
		return TaskTemplate{}, err
	}
	if !utf8.Valid(data) {
		return TaskTemplate{}, fmt.Errorf("%w: contenido no UTF-8", ErrInvalidTask)
	}

	return TaskTemplate{
		Name:       strings.TrimSuffix(clean, filepath.Ext(clean)),
		Filename:   clean,
		Kind:       kind,
		Size:       info.Size(),
		ModifiedAt: info.ModTime(),
		Content:    string(data),
	}, nil
}

func (s *TaskStore) Create(input CreateTaskInput) (TaskTemplate, error) {
	if err := s.Ensure(); err != nil {
		return TaskTemplate{}, err
	}

	kind := strings.ToLower(strings.TrimSpace(input.Kind))
	if kind != "service" && kind != "timer" {
		return TaskTemplate{}, fmt.Errorf("%w: tipo no permitido", ErrInvalidTask)
	}

	name, err := normalizeUnitName(input.Name, kind)
	if err != nil {
		return TaskTemplate{}, err
	}

	description, err := safeLine(input.Description, 160, false)
	if err != nil {
		return TaskTemplate{}, err
	}
	if description == "" {
		description = "Axolotl managed template"
	}

	var content string
	switch kind {
	case "service":
		content, err = buildServiceTemplate(input, description)
	case "timer":
		content, err = buildTimerTemplate(input, description)
	}
	if err != nil {
		return TaskTemplate{}, err
	}

	filename := name + "." + kind
	full := filepath.Join(s.root, filename)
	file, err := os.OpenFile(full, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return TaskTemplate{}, err
	}

	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		_ = os.Remove(full)
		return TaskTemplate{}, err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(full)
		return TaskTemplate{}, err
	}

	return s.Read(filename)
}

func (s *TaskStore) Delete(filename string) error {
	full, _, _, err := s.resolveFilename(filename)
	if err != nil {
		return err
	}

	info, err := os.Lstat(full)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return ErrInvalidTask
	}
	return os.Remove(full)
}

func (s *TaskStore) resolveFilename(filename string) (full, clean, kind string, err error) {
	clean = strings.TrimSpace(filename)
	if clean == "" || clean != filepath.Base(clean) || strings.ContainsAny(clean, `/\\`) {
		return "", "", "", ErrInvalidTask
	}

	kind = taskKindFromFilename(clean)
	if kind == "" {
		return "", "", "", ErrInvalidTask
	}

	name := strings.TrimSuffix(clean, "."+kind)
	if _, err := normalizeUnitName(name, kind); err != nil {
		return "", "", "", err
	}

	full = filepath.Join(s.root, clean)
	return full, clean, kind, nil
}

func buildServiceTemplate(input CreateTaskInput, description string) (string, error) {
	execStart, err := safeLine(input.ExecStart, 500, true)
	if err != nil {
		return "", err
	}
	if execStart == "" || !strings.HasPrefix(execStart, "/") {
		return "", fmt.Errorf("%w: ExecStart debe comenzar con /", ErrInvalidTask)
	}

	workingDir, err := safeLine(input.WorkingDirectory, 240, false)
	if err != nil {
		return "", err
	}
	if workingDir != "" && !strings.HasPrefix(workingDir, "/") {
		return "", fmt.Errorf("%w: WorkingDirectory debe ser absoluta", ErrInvalidTask)
	}

	user, err := safeIdentifier(input.User, 64, true)
	if err != nil {
		return "", err
	}

	restart := strings.TrimSpace(input.Restart)
	switch restart {
	case "", "no":
		restart = "no"
	case "on-failure", "always":
	default:
		return "", fmt.Errorf("%w: política Restart no permitida", ErrInvalidTask)
	}

	after := strings.TrimSpace(input.After)
	if after == "" {
		after = "network.target"
	}
	if err := validateUnitList(after); err != nil {
		return "", err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "[Unit]\nDescription=%s\nAfter=%s\n\n", description, after)
	b.WriteString("[Service]\nType=simple\n")
	if user != "" {
		fmt.Fprintf(&b, "User=%s\n", user)
	}
	if workingDir != "" {
		fmt.Fprintf(&b, "WorkingDirectory=%s\n", workingDir)
	}
	fmt.Fprintf(&b, "ExecStart=%s\nRestart=%s\n", execStart, restart)
	if restart != "no" {
		b.WriteString("RestartSec=5\n")
	}
	b.WriteString("\n[Install]\nWantedBy=multi-user.target\n")
	return b.String(), nil
}

func buildTimerTemplate(input CreateTaskInput, description string) (string, error) {
	onCalendar, err := safeLine(input.OnCalendar, 160, true)
	if err != nil {
		return "", err
	}

	target := strings.TrimSpace(input.TargetUnit)
	if target == "" {
		return "", fmt.Errorf("%w: falta Unit objetivo", ErrInvalidTask)
	}
	if !strings.HasSuffix(target, ".service") {
		target += ".service"
	}
	if taskKindFromFilename(target) != "service" {
		return "", fmt.Errorf("%w: Unit objetivo inválida", ErrInvalidTask)
	}
	base := strings.TrimSuffix(target, ".service")
	if _, err := normalizeUnitName(base, "service"); err != nil {
		return "", fmt.Errorf("%w: Unit objetivo inválida", ErrInvalidTask)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "[Unit]\nDescription=%s\n\n", description)
	fmt.Fprintf(&b, "[Timer]\nOnCalendar=%s\nUnit=%s\n", onCalendar, target)
	if input.Persistent {
		b.WriteString("Persistent=true\n")
	}
	b.WriteString("\n[Install]\nWantedBy=timers.target\n")
	return b.String(), nil
}

func normalizeUnitName(value, kind string) (string, error) {
	name := strings.TrimSpace(value)
	name = strings.TrimSuffix(name, ".service")
	name = strings.TrimSuffix(name, ".timer")
	if name == "" || len(name) > 64 || strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") || strings.Contains(name, "..") {
		return "", fmt.Errorf("%w: nombre inválido", ErrInvalidTask)
	}

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' || r == '@' {
			continue
		}
		return "", fmt.Errorf("%w: caracteres no permitidos en nombre", ErrInvalidTask)
	}

	if kind != "service" && kind != "timer" {
		return "", ErrInvalidTask
	}
	return name, nil
}

func safeLine(value string, max int, required bool) (string, error) {
	value = strings.TrimSpace(value)
	if required && value == "" {
		return "", fmt.Errorf("%w: campo requerido", ErrInvalidTask)
	}
	if len(value) > max || strings.ContainsAny(value, "\r\n\x00") || !utf8.ValidString(value) {
		return "", fmt.Errorf("%w: valor inválido", ErrInvalidTask)
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return "", fmt.Errorf("%w: caracteres de control no permitidos", ErrInvalidTask)
		}
	}
	return value, nil
}

func safeIdentifier(value string, max int, optional bool) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" && optional {
		return "", nil
	}
	if value == "" || len(value) > max {
		return "", ErrInvalidTask
	}
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' {
			continue
		}
		return "", ErrInvalidTask
	}
	return value, nil
}

func validateUnitList(value string) error {
	for _, token := range strings.Fields(value) {
		if token == "" || len(token) > 100 {
			return ErrInvalidTask
		}
		for _, r := range token {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' || r == '@' {
				continue
			}
			return ErrInvalidTask
		}
	}
	return nil
}

func taskKindFromFilename(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".service"):
		return "service"
	case strings.HasSuffix(filename, ".timer"):
		return "timer"
	default:
		return ""
	}
}

func IsConflict(err error) bool {
	return errors.Is(err, fs.ErrExist)
}
