# Fecha
2026-09-28

# Objetivo
Convertir la página inicial de Axolotl en una plantilla dashboard realmente multipropósito, con módulos no repetitivos, datos demostrativos entregados desde Goravel y una sola presencia visual del logotipo en el sidebar.

# Decisiones tomadas
- El SVG oficial permanece visible únicamente en el botón superior del sidebar.
- Se eliminó el bloque inferior que repetía la marca; el footer del sidebar ahora muestra solo estado textual del backend.
- La página principal dejó de estar orientada a cripto/mercados.
- Se conservaron IDs maduros de secciones y canvases cuando permiten reutilizar interacciones/gráficas existentes sin reescribir código estable.
- Los datos repetibles de KPIs, proyectos, tareas, actividad, equipo, agenda e integraciones se entregan desde DashboardController.
- Se añadió un partial independiente workspace-grid.tmpl para equipo, agenda, integraciones y acciones rápidas.
- Métricas del servidor y Laboratorio UI se mantienen como ejemplos técnicos reutilizables.
- La búsqueda global y el sidebar ahora navegan por conceptos genéricos de gestión, analítica y workspace.

# Arquitectura actual
```text
DashboardController
  -> DashboardPageData
     -> Metrics
     -> Projects
     -> Tasks
     -> Activities
     -> Team
     -> Schedule
     -> Integrations
  -> dashboard.tmpl
     -> stats.tmpl
     -> dashboard-grid.tmpl
     -> workspace-grid.tmpl
     -> server-metrics.tmpl
     -> code-lab.tmpl
```

# Librerías usadas
No se añadieron dependencias nuevas. Se conserva:
- Goravel/html/template
- Chart.js
- Lenis
- AOS
- Highlight.js
- Font Awesome

# Archivos importantes modificados
- app/http/controllers/dashboard_controller.go
- resources/views/partials/sidebar.tmpl
- resources/views/partials/topbar.tmpl
- resources/views/partials/stats.tmpl
- resources/views/partials/dashboard-grid.tmpl
- resources/views/partials/workspace-grid.tmpl
- resources/views/partials/dashboard-main.tmpl
- public/css/app.css
- public/js/app.js
- tests/feature/http_test.go
- README.md
- design.md

# Problemas encontrados
- La página original seguía heredando muchos conceptos de una demo financiera/cripto.
- El logo SVG se repetía visualmente en la parte superior e inferior del sidebar.
- Parte del contenido estaba hardcodeado exclusivamente en frontend, reduciendo el valor de la plantilla como ejemplo full-stack.
- El servidor del puerto 3000 podía mantener un proceso anterior de go run en Windows.

# Soluciones implementadas
- Una sola marca visible en el sidebar.
- Sidebar reorganizado en Gestión, Analítica, Workspace, Notificaciones y Ajustes.
- Nuevos módulos: proyectos, tareas, actividad, equipo, agenda, integraciones y acciones rápidas.
- KPIs y colecciones principales generadas por backend.
- Modales, búsqueda y textos visibles generalizados.
- Reinicio explícito del proceso que ocupa el puerto 3000 antes del smoke test.
- Smoke HTTP confirmó que no aparecen términos cripto y que no existe el antiguo bloque ai-brand-image/ai-box en el HTML renderizado.

# Validación
- gofmt: exitoso.
- go test ./...: exitoso.
- go vet ./...: exitoso.
- go build ./...: exitoso.
- node --check public/js/app.js: exitoso.
- GET /: 200 y contiene Centro de operaciones, Proyectos, Tareas, Equipo, Agenda e Integraciones.
- GET /health: 200 con database=ok.
- El HTML contiene exactamente un elemento visual panel-brand-image.
- No existe marca SVG inferior duplicada.
- Sin palabras visibles de la antigua demo cripto en el HTML servido.

# Pendientes
- Validación visual manual en navegador real; PL no expone browser UI.
- En un proyecto consumidor, reemplazar los datos demo del controller por ORM/APIs reales.

# Próximos pasos
Usar cada bloque como módulo opcional: eliminar, mover o reemplazar partials según el producto concreto sin modificar el sistema base de layout, tema, sidebar, búsqueda y componentes.
