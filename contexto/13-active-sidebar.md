# Fecha
2026-09-28

# Objetivo
Corregir el marcado active de navegación después de separar Axolotl en rutas multipágina.

# Decisiones tomadas
- Mantener active renderizado desde backend para cada ruta.
- Añadir sincronización frontend como fallback y para cambios de hash.
- En dashboard, los subitems se resuelven por su href `/dashboard#<id>`.

# Arquitectura actual
El backend entrega `ActiveRoute` y el template marca la página activa. El frontend llama `syncActiveRoute()` al cargar y en `hashchange`.

# Librerías usadas
Sin dependencias nuevas.

# Archivos importantes modificados
- resources/views/partials/sidebar.tmpl
- public/js/app.js
- tests/feature/http_test.go

# Problemas encontrados
El scroll spy llamaba `activateSidebarTarget()`, que eliminaba todos los active pero ya no encontraba elementos porque los subitems habían pasado de `data-target` a enlaces href.

# Soluciones implementadas
- Lookup por `href="/dashboard#..."`.
- Fallback al link Dashboard.
- `aria-current` para accesibilidad.
- Tests de ruta activa.

# Pendientes
Ninguno.

# Próximos pasos
Añadir la nueva ruta de Explorador y reutilizar el mismo mecanismo `data-route="files"`.
