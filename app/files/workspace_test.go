package files

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkspaceCRUD(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitkeep"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	workspace, err := NewWorkspace(root, 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
	if err := workspace.EnsureDemo(); err != nil {
		t.Fatal(err)
	}

	entries, err := workspace.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("se esperaban archivos demo")
	}

	if _, err := workspace.Create("", "pruebas", "folder"); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Create("pruebas", "hola.go", "file"); err != nil {
		t.Fatal(err)
	}

	source := "package pruebas\n\nconst Nombre = \"Axolotl\"\n"
	if err := workspace.Save("pruebas/hola.go", source); err != nil {
		t.Fatal(err)
	}

	file, err := workspace.Read("pruebas/hola.go")
	if err != nil {
		t.Fatal(err)
	}
	if file.Content != source {
		t.Fatalf("contenido inesperado: %q", file.Content)
	}
	if file.Language != "go" {
		t.Fatalf("lenguaje inesperado: %q", file.Language)
	}

	if err := workspace.Delete("pruebas/hola.go"); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Read("pruebas/hola.go"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("se esperaba archivo inexistente, recibido %v", err)
	}

	if err := workspace.Delete("pruebas"); err != nil {
		t.Fatal(err)
	}
}

func TestWorkspaceRejectsTraversalAndInvalidNames(t *testing.T) {
	root := t.TempDir()
	workspace, err := NewWorkspace(root, 1024)
	if err != nil {
		t.Fatal(err)
	}
	if err := workspace.EnsureDemo(); err != nil {
		t.Fatal(err)
	}

	if _, err := workspace.List("../"); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("se esperaba ErrInvalidPath para traversal, recibido %v", err)
	}
	if _, err := workspace.Read("../fuera.txt"); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("se esperaba ErrInvalidPath para lectura fuera del root, recibido %v", err)
	}
	if _, err := workspace.Create("", "../fuera.txt", "file"); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("se esperaba ErrInvalidPath para nombre inválido, recibido %v", err)
	}

	outside := t.TempDir()
	link := filepath.Join(root, "link")
	if err := os.Symlink(outside, link); err == nil {
		if _, err := workspace.List("link"); !errors.Is(err, ErrInvalidPath) {
			t.Fatalf("se esperaba rechazo de symlink, recibido %v", err)
		}
	}
}

func TestWorkspaceRejectsBinaryAndOversizedFiles(t *testing.T) {
	root := t.TempDir()
	workspace, err := NewWorkspace(root, 16)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(root, "bin.dat"), []byte{0x00, 0x01, 0x02}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Read("bin.dat"); !errors.Is(err, ErrBinaryFile) {
		t.Fatalf("se esperaba ErrBinaryFile, recibido %v", err)
	}

	if err := os.WriteFile(filepath.Join(root, "large.txt"), []byte(strings.Repeat("x", 20)), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Read("large.txt"); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("se esperaba ErrTooLarge, recibido %v", err)
	}
	if err := workspace.Save("bin.dat", strings.Repeat("y", 20)); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("se esperaba ErrTooLarge al guardar, recibido %v", err)
	}
}
