# Tarea 13 — Active del sidebar por ruta

## Estado
Completada.

## Objetivo
Corregir el estado activo del sidebar para que refleje la ruta actual y las subsecciones del dashboard.

## Cambios
- Links principales usan `data-route` y `aria-current="page"`.
- El JavaScript sincroniza active con `window.location.pathname`.
- En `/dashboard`, el hash activa la subsección correspondiente y abre su grupo.
- El scroll spy ya no elimina el active del dashboard sin reemplazo.
- Se añadieron pruebas de render para Dashboard, Cuenta, Servidor y Componentes.

## Validación
- `go test ./tests/feature -run TestAuthenticatedPages -v` OK.
- `node --check public/js/app.js` OK.
