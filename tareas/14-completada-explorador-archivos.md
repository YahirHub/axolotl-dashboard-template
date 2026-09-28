# Tarea 14 — Explorador de archivos backend

## Estado
En proceso.

## Objetivo
Añadir una página `/archivos` con explorador tipo Windows conectado al backend Go, limitado a una carpeta demo segura.

## Alcance
- Navegar carpetas.
- Crear archivos.
- Crear carpetas.
- Eliminar archivos/carpetas con confirmación.
- Abrir archivos de texto/código.
- Editor tipo VS Code con coloreado Highlight.js.
- Guardar cambios.
- Advertencia modal al cerrar un archivo modificado sin guardar.
- Breadcrumbs y navegación hacia atrás.
- Ruta protegida por sesión.
- Backend confinado a `storage/demo-files`.

## Seguridad
- No permitir path traversal.
- No permitir symlinks.
- Límite de tamaño editable.
- Solo UTF-8/texto para edición.
- CSRF en mutaciones.
- Root configurable por .env.
- Registro de errores y respuestas JSON controladas.

## Validación
- tests unitarios/feature para navegación y CRUD;
- go test/vet/build;
- node --check;
- smoke HTTP;
- active del sidebar en /archivos.
