# Fecha
2026-09-28

# Objetivo
Añadir un explorador de archivos multipropósito, con backend Go real y editor de texto/código integrado, manteniendo todas las operaciones confinadas a una carpeta demo segura.

# Decisiones tomadas
- Nueva página privada: `/archivos`.
- Workspace por defecto: `storage/demo-files`.
- Configuración por `.env`: `FILES_ROOT` y `FILES_MAX_EDIT_BYTES`.
- El frontend nunca recibe ni controla una ruta absoluta del host.
- Las operaciones se implementan con librería estándar de Go; no se añadió dependencia de filesystem.
- El backend rechaza path traversal y symlinks.
- El editor acepta únicamente contenido UTF-8 y aplica límite de tamaño.
- Highlight.js existente se reutiliza para coloreado de código; no se añadió Monaco/CodeMirror.
- El editor usa textarea transparente sobre un bloque resaltado, con números de línea y scroll sincronizado.
- Los cambios sin guardar reutilizan el modal global de Axolotl.
- El contenido mutable de `storage/demo-files` se ignora en Git.
- Al iniciar el módulo, si el workspace está vacío salvo por `.gitkeep`, se crean archivos demo.
- El active del sidebar se corrigió antes de esta tarea y `Archivos` usa el mismo mecanismo `data-route="files"`.

# Arquitectura actual
```text
/archivos [RequireAuth]
  -> FileController.Index
     -> files.tmpl
        -> files-main.tmpl
        -> files.css
        -> files.js

/api/files [RequireAuth]
/api/files/read [RequireAuth]
/api/files/create [RequireAuth + CSRF]
/api/files/save [RequireAuth + CSRF]
/api/files/delete [RequireAuth + CSRF]

app/files/
├── configured.go
└── workspace.go
```

# Librerías usadas
No se añadieron dependencias.

Se reutiliza:
- Go standard library para filesystem;
- Goravel para routing, auth, sesiones y CSRF;
- Highlight.js 11.11.1 para syntax highlighting;
- Font Awesome para iconos;
- design system y modales existentes de Axolotl.

# Archivos importantes modificados
- app/files/workspace.go
- app/files/configured.go
- app/files/workspace_test.go
- app/http/controllers/file_controller.go
- config/files.go
- routes/web.go
- resources/views/files.tmpl
- resources/views/partials/files-main.tmpl
- resources/views/partials/sidebar.tmpl
- resources/views/partials/head.tmpl
- resources/views/partials/scripts.tmpl
- public/css/files.css
- public/js/files.js
- public/js/app.js
- tests/test_case.go
- tests/feature/files_test.go
- tests/feature/http_test.go
- .env.example
- .gitignore
- README.md
- design.md
- LILITH.md

# Problemas encontrados
- `.gitkeep` hacía que el workspace pareciera no vacío e impedía sembrar los archivos demo.
- Los tests debían usar un root temporal para no modificar la carpeta demo real.
- La lógica active anterior dependía de `data-target`, pero los subitems de dashboard ya eran links con hash.
- Un editor enriquecido completo como Monaco habría agregado una dependencia grande innecesaria para el starter.

# Soluciones implementadas
- `EnsureDemo()` ignora `.gitkeep`.
- `List()` oculta `.gitkeep` al usuario.
- Los tests configuran `FILES_ROOT` temporal por PID.
- Backend normaliza y verifica cada path con `filepath.Rel`.
- Se recorren componentes del path con `Lstat` para bloquear symlinks.
- Editor ligero tipo VS Code con Highlight.js, números de línea, guardado y dirty-state.
- Modal de tres acciones: Cancelar / Cerrar sin guardar / Guardar y cerrar.
- `beforeunload` protege frente a cerrar la página con cambios pendientes.
- Nuevo archivo se abre automáticamente en el editor.

# Validación
- `go test ./...`: exitoso.
- `go vet ./...`: exitoso.
- `go build ./...`: exitoso.
- `node --check public/js/app.js`: exitoso.
- `node --check public/js/auth.js`: exitoso.
- `node --check public/js/files.js`: exitoso.
- `go run . artisan route:list`: 20 rutas registradas.
- Tests del workspace: CRUD, traversal, symlink, binarios y límite de tamaño.
- Test feature API: create/save/read/delete y traversal -> 400.
- Test de páginas privadas: `/archivos` devuelve 200 y queda marcado active.
- No existe browser MCP en PL; la comprobación visual interactiva final queda para navegador manual.

# Pendientes
- Opcional: rename/move/upload/download si un consumidor real los necesita.
- Opcional: permisos por carpeta o usuario.
- Opcional: auditoría visual manual en navegador con la instalación autenticada.

# Próximos pasos
Mantener el explorador como workspace confinado. No convertirlo en un file manager arbitrario del servidor. Nuevas operaciones deben pasar por el mismo resolver seguro y contar con tests de traversal.
