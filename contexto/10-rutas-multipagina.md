# Fecha
2026-09-28

# Objetivo
Corregir Axolotl para que sea una aplicación Goravel multipágina real, separando el overview, la cuenta, las métricas de servidor y el laboratorio de componentes en rutas independientes.

# Decisiones tomadas
- `/dashboard` contiene únicamente el overview multipropósito.
- `/cuenta` contiene perfil, preferencias, apariencia y seguridad.
- `/servidor` contiene observabilidad técnica y el monitor CPU/RAM/red/almacenamiento.
- `/componentes` contiene el Laboratorio UI y catálogo de componentes.
- `/` se conserva como alias del dashboard de entrada.
- Las páginas comparten sidebar, topbar, overlays, scripts, tema y branding.
- El sidebar usa enlaces reales para cambiar de página y anchors únicamente para secciones internas de `/dashboard`.
- El topbar recibe `PageTitle` y `PageSubtitle` desde el backend.
- Se creó `PageMeta` para compartir metadata mínima sin introducir una capa adicional.
- El JavaScript compartido detecta si un módulo existe antes de inicializar charts o listeners específicos.
- La búsqueda global navega entre rutas en lugar de asumir que todos los módulos viven en el mismo documento.
- La cuenta demo persiste perfil y preferencias visuales en localStorage para sobrevivir navegación entre páginas.

# Arquitectura actual
```text
routes/web.go
├── /              -> DashboardController.Index
├── /dashboard     -> DashboardController.Index
├── /cuenta        -> PageController.Account
├── /servidor      -> PageController.Server
├── /componentes   -> PageController.Components
├── /health        -> HealthController.Show
└── /assets/*      -> public/

resources/views/
├── dashboard.tmpl
├── account.tmpl
├── server.tmpl
├── components.tmpl
└── partials/
    ├── dashboard-main.tmpl
    ├── account-main.tmpl
    ├── server-main.tmpl
    ├── components-main.tmpl
    ├── sidebar.tmpl
    ├── topbar.tmpl
    ├── server-metrics.tmpl
    ├── code-lab.tmpl
    └── ...
```

# Librerías usadas
No se añadieron dependencias.

Se conserva:
- Goravel 1.18
- Gin
- html/template
- SQLite/GORM
- Chart.js
- Lenis
- AOS
- Highlight.js
- Font Awesome

# Archivos importantes modificados
- routes/web.go
- app/http/controllers/dashboard_controller.go
- app/http/controllers/page_controller.go
- resources/views/dashboard.tmpl
- resources/views/account.tmpl
- resources/views/server.tmpl
- resources/views/components.tmpl
- resources/views/partials/dashboard-main.tmpl
- resources/views/partials/account-main.tmpl
- resources/views/partials/server-main.tmpl
- resources/views/partials/components-main.tmpl
- resources/views/partials/sidebar.tmpl
- resources/views/partials/topbar.tmpl
- resources/views/partials/head.tmpl
- public/css/app.css
- public/js/app.js
- tests/feature/http_test.go
- README.md
- design.md

# Problemas encontrados
- El dashboard anterior seguía concentrando overview, servidor y laboratorio UI en la misma página.
- El JavaScript compartido asumía que canvases y botones del dashboard siempre existían, lo que habría provocado errores al separar rutas.
- El menú y el buscador usaban scroll local incluso para módulos que ahora pertenecen a otras páginas.
- El título HTML del head seguía usando la clave anterior `.title` y no la propiedad exportada `.Title`.

# Soluciones implementadas
- Separación física de templates y rutas.
- Guard en `initCharts()` cuando no existe el dashboard.
- Listeners específicos convertidos a opcionales.
- Scroll spy limitado a `/dashboard`.
- Búsqueda global con navegación cross-page.
- Sidebar con rutas reales.
- Menú de usuario enlazado a `/cuenta` y `/componentes`.
- Página de cuenta interactiva y persistencia local demo.
- Tests HTTP independientes para cada página.

# Validación
- `gofmt`: exitoso.
- `go test ./...`: exitoso.
- `go vet ./...`: exitoso.
- `go build ./...`: exitoso.
- `node --check public/js/app.js`: exitoso.
- `artisan route:list`: 7 rutas registradas.
- HTTP:
  - /dashboard -> 200
  - /cuenta -> 200
  - /servidor -> 200
  - /componentes -> 200
  - /health -> 200
- /dashboard ya no contiene Métricas del servidor ni Laboratorio UI.
- /servidor contiene CPU, RAM, red y almacenamiento.
- /componentes contiene Laboratorio UI.
- /cuenta contiene perfil, preferencias y seguridad.

# Pendientes
- La cuenta es una demo frontend persistida en localStorage; autenticación real y persistencia de usuario deben implementarse cuando un proyecto consumidor lo requiera.
- Validación visual manual en navegador real.

# Próximos pasos
Usar esta separación como regla de arquitectura: cuando un módulo represente un área funcional independiente, crear una ruta/página en lugar de seguir alargando `/dashboard`.
