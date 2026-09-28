# Fecha
2026-09-28

# Objetivo
Corregir la excepción del navegador `Uncaught TypeError: this._fn is not a function` originada dentro del animator de Chart.js.

# Decisiones tomadas
- Mantener Chart.js y los gradientes Prism existentes.
- No desactivar todas las animaciones del dashboard.
- Desactivar únicamente `animations.colors` en los charts que usan `CanvasGradient` scriptable.
- Mantener animaciones numéricas de Chart.js y las animaciones RAF/CSS del servidor.
- Añadir cache busting al asset `app.js` para evitar que el navegador siga usando la versión anterior.

# Arquitectura actual
Sin cambios estructurales. El fix afecta únicamente configuración de Chart.js y carga del asset JS.

# Librerías usadas
- Chart.js 4.4.7.

# Archivos importantes modificados
- public/js/app.js
- resources/views/partials/scripts.tmpl

# Problemas encontrados
Las funciones `prismLineGradient()` y `prismAreaGradient()` devuelven `CanvasGradient`. Chart.js intentaba interpolar esos valores dentro de su animación de colores; el animator no encontraba un interpolador válido y terminaba con `this._fn is not a function`.

# Soluciones implementadas
- `animations: { colors: false }` en sparklines.
- Mismo guard en gráfica de rendimiento.
- Mismo guard en gráfica de tendencia.
- Mismo guard en gráfica de red del servidor.
- Versionado del script compartido a `/assets/js/app.js?v=20260928-2`.

# Validación
- node --check public/js/app.js: exitoso.
- go test ./...: exitoso.
- go vet ./...: exitoso.
- go build ./...: exitoso.
- Chrome headless en /dashboard: cargó app.js?v=20260928-2 sin la excepción Chart.js.
- Chrome headless en /servidor: monitor iniciado y permaneció activo varias actualizaciones sin `this._fn is not a function`.
- Los únicos mensajes ERROR observados en Chrome headless fueron internos de GCM/registro de Chrome y no pertenecen a la aplicación.

# Pendientes
Confirmación visual del usuario después de recargar la nueva versión.

# Próximos pasos
Si aparece otra excepción del animator, revisar primero cualquier propiedad scriptable que devuelva objetos no interpolables antes de desactivar animaciones globalmente.
