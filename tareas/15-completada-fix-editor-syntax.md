# Tarea 15 — Corregir coloreado del editor

## Estado
Completada.

## Implementado
- Se excluye el overlay del editor del enhancer global.
- El textarea queda visible por defecto.
- Solo se vuelve transparente cuando Highlight.js genera el resaltado.
- Fallback de texto plano visible y highlightAuto.
- Colores adicionales para PHP.
- Cache bust actualizado.
- Prueba de detección .php -> php.

## Validación
- go test ./app/files ./tests/feature OK.
- go vet ./... OK.
- go build ./... OK.
- node --check public/js/app.js OK.
- node --check public/js/files.js OK.
