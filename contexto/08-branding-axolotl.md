# Fecha
2026-09-28

# Objetivo
Integrar los archivos oficiales de marca Axolotl al starter y reemplazar iconos genéricos por la identidad visual real.

# Decisiones tomadas
- Los archivos `axolotl.svg` y `axolotl.png` se movieron desde el Desktop al repositorio.
- Ruta pública: `public/brand/`.
- El SVG es la marca principal para UI porque escala sin pérdida.
- El PNG se usa como favicon fallback/apple touch.
- El botón de collapse/expand del sidebar conserva su comportamiento actual; solo cambia la marca mostrada.
- El bloque Axolotl AI también reutiliza el SVG para reforzar identidad.

# Arquitectura actual
```text
public/brand/
├── axolotl.svg
└── axolotl.png
```

Se sirven mediante la ruta estática existente:
```text
/assets/brand/axolotl.svg
/assets/brand/axolotl.png
```

# Librerías usadas
Ninguna nueva.

# Archivos importantes modificados
- public/brand/axolotl.svg
- public/brand/axolotl.png
- resources/views/partials/head.tmpl
- resources/views/partials/sidebar.tmpl
- public/css/app.css
- tests/feature/http_test.go
- README.md
- design.md

# Problemas encontrados
Ninguno bloqueante.

# Soluciones implementadas
- Favicon SVG + PNG.
- Apple touch icon.
- Logo real dentro del botón interactivo del sidebar.
- Logo real en el componente Axolotl AI.
- Estilos responsive sin alterar collapse/hover.

# Pendientes
- Validar visualmente en navegador real cuando el entorno PL disponga de browser UI.

# Próximos pasos
Mantener estos dos archivos como fuente oficial de marca y no volver a introducir iconos genéricos como logotipo principal.
