package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTaskStoreCreateListReadDelete(t *testing.T) {
	store, err := NewTaskStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	service, err := store.Create(CreateTaskInput{
		Name:             "axolotl-worker",
		Kind:             "service",
		Description:      "Worker de ejemplo",
		ExecStart:        "/usr/bin/php /var/www/app/artisan queue:work",
		WorkingDirectory: "/var/www/app",
		User:             "www-data",
		Restart:          "on-failure",
		After:            "network.target",
	})
	if err != nil {
		t.Fatal(err)
	}
	if service.Filename != "axolotl-worker.service" {
		t.Fatalf("filename inesperado: %s", service.Filename)
	}
	if !strings.Contains(service.Content, "ExecStart=/usr/bin/php /var/www/app/artisan queue:work") {
		t.Fatalf("contenido service incompleto: %s", service.Content)
	}

	timer, err := store.Create(CreateTaskInput{
		Name:        "axolotl-daily",
		Kind:        "timer",
		Description: "Tarea diaria",
		OnCalendar:  "*-*-* 03:00:00",
		TargetUnit:  "axolotl-worker.service",
		Persistent:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if timer.Filename != "axolotl-daily.timer" {
		t.Fatalf("filename timer inesperado: %s", timer.Filename)
	}
	if !strings.Contains(timer.Content, "Persistent=true") {
		t.Fatalf("timer debería ser persistent")
	}

	items, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("se esperaban 2 plantillas, recibido %d", len(items))
	}

	if err := store.Delete("axolotl-worker.service"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Read("axolotl-worker.service"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("se esperaba archivo eliminado, recibido %v", err)
	}
}

func TestTaskStoreRejectsUnsafeInput(t *testing.T) {
	root := t.TempDir()
	store, err := NewTaskStore(root)
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{
		"../escape",
		"bad/name",
		"bad\\name",
		".hidden",
		"two..dots",
		"bad name",
	} {
		if _, err := store.Create(CreateTaskInput{
			Name:      name,
			Kind:      "service",
			ExecStart: "/usr/bin/true",
		}); !errors.Is(err, ErrInvalidTask) {
			t.Fatalf("nombre inseguro aceptado %q: %v", name, err)
		}
	}

	for _, input := range []CreateTaskInput{
		{Name: "newline", Kind: "service", ExecStart: "/usr/bin/true\nEnvironment=BAD=1"},
		{Name: "relative", Kind: "service", ExecStart: "php artisan"},
		{Name: "bad-after", Kind: "service", ExecStart: "/usr/bin/true", After: "network.target\nExecStart=/bin/sh"},
		{Name: "timer", Kind: "timer", OnCalendar: "daily\nUnit=bad.service", TargetUnit: "safe.service"},
		{Name: "target", Kind: "timer", OnCalendar: "daily", TargetUnit: "../bad.service"},
	} {
		if _, err := store.Create(input); !errors.Is(err, ErrInvalidTask) {
			t.Fatalf("input inseguro aceptado %+v: %v", input, err)
		}
	}

	outside := filepath.Join(filepath.Dir(root), "outside.service")
	if _, err := store.Read(outside); !errors.Is(err, ErrInvalidTask) {
		t.Fatalf("lectura externa aceptada: %v", err)
	}
	if err := store.Delete("../outside.service"); !errors.Is(err, ErrInvalidTask) {
		t.Fatalf("delete externo aceptado: %v", err)
	}
}

func TestTaskStoreDoesNotFollowSymlinkFiles(t *testing.T) {
	root := t.TempDir()
	store, err := NewTaskStore(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Ensure(); err != nil {
		t.Fatal(err)
	}

	outside := filepath.Join(t.TempDir(), "outside.service")
	if err := os.WriteFile(outside, []byte("[Service]\nExecStart=/usr/bin/true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "linked.service")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink no disponible en este entorno: %v", err)
	}

	items, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("el symlink no debe listarse")
	}
	if _, err := store.Read("linked.service"); !errors.Is(err, ErrInvalidTask) {
		t.Fatalf("el symlink no debe abrirse: %v", err)
	}
}
