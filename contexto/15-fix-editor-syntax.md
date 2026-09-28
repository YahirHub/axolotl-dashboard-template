# Fecha
2026-09-28

# Objetivo
Corregir el editor de archivos cuando el área de código quedaba negra/invisible, observado al editar un archivo PHP.

# Decisiones tomadas
- El overlay interno del editor no debe pasar por `enhanceCodeBlocks()`.
- El textarea debe mostrar texto visible por defecto.
- La transparencia del textarea solo se activa cuando Highlight.js produjo realmente un resultado.
- Si Highlight.js no está disponible, el lenguaje no existe o el coloreado falla, el usuario debe seguir viendo y editando texto plano.
- Se mantienen Highlight.js y el editor ligero actual; no se añade otra dependencia.
- Se amplían colores específicos para tokens usados por PHP.
- Se incrementan versiones de assets para evitar que el navegador use JS/CSS anterior en caché.

# Arquitectura actual
El editor conserva dos capas:
1. `textarea.editor-input` para interacción/cursor.
2. `pre.editor-highlight` para syntax highlighting.

`.editor-code-wrap.highlight-ready` se añade únicamente si el resaltado se generó correctamente. Sin esa clase, el overlay se oculta y el textarea queda visible.

# Librerías usadas
Sin dependencias nuevas. Se mantiene Highlight.js 11.11.1 por CDN.

# Archivos importantes modificados
- public/js/app.js
- public/js/files.js
- public/css/files.css
- resources/views/partials/files-main.tmpl
- resources/views/partials/head.tmpl
- resources/views/partials/scripts.tmpl
- app/files/workspace_test.go

# Problemas encontrados
El enhancer global de Axolotl recorría todos los `pre code`, incluyendo el `pre` interno usado como overlay del editor. Además, el textarea siempre tenía texto transparente, por lo que cualquier fallo o ausencia de resaltado dejaba el código invisible sobre el fondo oscuro.

# Soluciones implementadas
- `enhanceCodeBlocks()` ignora `.editor-highlight` y `[data-no-code-enhance]`.
- El editor usa `highlight-ready` para decidir cuándo ocultar el texto del textarea.
- Se agregó fallback visible automático.
- Se agregó `highlightAuto` como respaldo cuando el lenguaje solicitado no está registrado.
- Se añadieron colores para variables, meta/tags, nombres, atributos y parámetros.
- Cache bust actualizado:
  - app.js -> 20260928-3
  - files.js -> 20260928-2
  - files.css -> 20260928-2
- Se añadió prueba para `.php -> php`.

# Validación
- `go test ./app/files ./tests/feature`: OK.
- `go vet ./...`: OK.
- `go build ./...`: OK.
- `node --check public/js/app.js`: OK.
- `node --check public/js/files.js`: OK.
- El archivo de prueba existente `storage/demo-files/prueba/a.php` se conserva sin modificar.

# Pendientes
Validación visual manual en navegador después de recargar `/archivos`.

# Próximos pasos
Si un lenguaje no está incluido en el bundle CDN, el editor debe seguir mostrando texto visible gracias al fallback; se puede añadir su grammar específica después solo si es necesaria.
