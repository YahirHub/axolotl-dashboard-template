# Fecha
2026-09-29

# Objetivo
Simplificar el sidebar para que solo navegue a páginas con route real y agrupar las páginas técnicas bajo un dropdown Server.

# Decisiones tomadas
- El sidebar ya no contiene enlaces con hash a secciones internas del dashboard.
- El sidebar ya no contiene acciones sin route como Notificaciones.
- Los dropdowns pueden ser controles de navegación, pero sus hijos deben ser páginas con route.
- Se crea el grupo Server con:
  - /servidor
  - /servidor/tareas
  - /servidor/entorno
- /servidor/tareas administra únicamente plantillas systemd .service y .timer.
- Las plantillas no se instalan, habilitan ni ejecutan.
- Se guardan en SERVER_TASKS_ROOT, por defecto storage/server-tasks.
- No se llama a systemctl ni se escribe en /etc/systemd.
- El usuario pidió explícitamente NO crear commit todavía para revisar la propuesta.

# Arquitectura actual
- ServerController concentra páginas Server y API de plantillas.
- app/server/tasks.go valida nombres/campos y genera archivos de unidad.
- Mutaciones usan sesión + CSRF mediante las rutas web existentes.
- El storage de plantillas está aislado en tests con SERVER_TASKS_ROOT temporal.

# Librerías usadas
Sin dependencias nuevas. Solo librería estándar de Go y componentes existentes de Axolotl.

# Archivos importantes modificados
- resources/views/partials/sidebar.tmpl
- public/js/app.js
- routes/web.go
- app/http/controllers/server_controller.go
- app/server/tasks.go
- config/server.go
- resources/views/server_tasks.tmpl
- resources/views/server_environment.tmpl
- public/css/server.css
- public/js/server_tasks.js
- tests/feature/server_tasks_test.go
- app/server/tasks_test.go

# Problemas encontrados
El sidebar anterior mezclaba routes con anclas # y acciones JS, por lo que el estado active dependía del scroll y no exclusivamente de la página actual.

# Soluciones implementadas
- Navegación active basada solo en pathname/data-route.
- Dropdown Server abre automáticamente en cualquiera de sus rutas hijas.
- Eliminado scroll-spy que cambiaba el active del sidebar.
- Administrador de tareas con listar/crear/ver/eliminar plantillas.
- Entorno muestra únicamente runtime/configuración no sensible.

# Estado posterior
El usuario revisó visualmente el módulo, pidió ajustes en la página Entorno y posteriormente indicó explícitamente que se versionaran todos los cambios pendientes antes de cerrar el chat.

La instrucción anterior de “no hacer commit” quedó anulada por esa orden posterior.

# Próximos pasos
Continuar desde el roadmap RBAC documentado en `contexto/22-cierre-chat-roadmap-rbac.md`.
