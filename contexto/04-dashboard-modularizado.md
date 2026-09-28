# Fecha
2026-09-28

# Objetivo
Migrar la referencia Axolotl desde un único HTML monolítico a vistas y assets modulares servidos por Goravel, sin perder el diseño ni la interactividad existente.

# Decisiones tomadas
- El motor de vistas es el `html/template` nativo integrado por Goravel/Gin.
- Se usa una vista de página `dashboard.tmpl` y parciales nombrados, en lugar de introducir un motor de plantillas externo.
- CSS y JavaScript salieron del HTML y viven en `public/css/app.css` y `public/js/app.js`.
- Las librerías frontend continúan por CDN desde el parcial de scripts/head.
- La referencia original permanece intacta en `docs/reference/index-original.html`.
- Los tokens RGB de la paleta Axolotl fueron corregidos para coincidir con sus colores hex.
- Se eliminaron la vista Welcome y el UserController de demo del scaffold.
- La ruta raíz ahora usa `DashboardController`.
- Los assets se sirven en `/assets/*`.

# Arquitectura actual
```text
resources/views/
├── dashboard.tmpl
└── partials/
    ├── head.tmpl
    ├── scripts.tmpl
    ├── sidebar.tmpl
    ├── dashboard-main.tmpl
    ├── topbar.tmpl
    ├── stats.tmpl
    ├── dashboard-grid.tmpl
    ├── server-metrics.tmpl
    ├── code-lab.tmpl
    └── overlays.tmpl

public/
├── css/app.css
└── js/app.js
```

# Librerías usadas
- Go html/template a través de Goravel/Gin.
- Chart.js CDN.
- Lenis CDN.
- AOS CDN.
- Highlight.js CDN.
- Font Awesome CDN.

# Archivos importantes modificados
- app/http/controllers/dashboard_controller.go
- routes/web.go
- resources/views/dashboard.tmpl
- resources/views/partials/*
- public/css/app.css
- public/js/app.js

# Problemas encontrados
- La referencia monolítica mezclaba CSS, markup y JS en un solo archivo.
- El primer proceso de desarrollo con `go run .` dejó un binario hijo usando el puerto 3000; se localizó y cerró por PID.
- La paleta tenía valores `--*-rgb` de una versión turquesa anterior aunque los hex ya eran rosados.

# Soluciones implementadas
- Separación automática inicial y posterior división por secciones.
- Render server-side funcionando desde Goravel.
- Smoke test HTTP: `/`, `/assets/css/app.css` y `/assets/js/app.js` responden 200.
- El HTML renderizado contiene Sidebar, Métricas del servidor, Laboratorio UI y scripts esperados.
- `go test ./...` y `go vet ./...` continúan pasando.

# Pendientes
- Implementar middleware global y rutas de salud.
- Añadir pruebas HTTP más útiles.
- Actualizar README final.
- Validación final de build y documentación.

# Próximos pasos
1. Crear middleware global de headers seguros.
2. Añadir `/health`.
3. Probar rutas y headers.
4. Cerrar README y primera versión.
