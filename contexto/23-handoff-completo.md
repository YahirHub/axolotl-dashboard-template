# Handoff completo — Axolotl

## Fecha
2026-09-29

## Proyecto
Ruta local:
`C:\Users\Admin\Desktop\GitHub\axolotl-dashboard-template`

Rama:
`main`

Stack:
- Go 1.25+
- Goravel 1.18
- Gin
- SQLite / GORM
- html/template
- frontend ligero con JS/CSS propios + CDN para Chart.js, Highlight.js, AOS, Lenis y Font Awesome

## Regla de arquitectura principal
Axolotl es un **panel multipropósito en Go**.

Las áreas independientes deben tener **routes reales**. El sidebar no debe volver a mezclar:
- enlaces con hash;
- scroll-spy;
- acciones JS sin página.

Un dropdown sí puede existir como control visual, pero sus hijos deben apuntar a routes reales.

## Routes principales privadas
- /dashboard
- /cuenta
- /servidor
- /servidor/tareas
- /servidor/entorno
- /componentes
- /archivos
- /archivos/*path

## Sidebar actual
- Dashboard
- Server
  - Servidor
  - Administrador de tareas
  - Entorno
- Archivos
- Cuenta
- Componentes

El estado active depende de ActiveRoute/pathname.

## Base de datos
No existe auto-migrate.

Para una DB nueva:
```bash
go run . artisan migrate
go run .
```

Si no existe la tabla users, InitializeAdminPIN() detiene el arranque y el log indica ejecutar migrate.

## Autenticación y seguridad
Implementado:
- setup inicial con PIN temporal de 6 dígitos;
- admin inicial;
- login/logout;
- registro opcional;
- Argon2id;
- CSRF;
- sesiones;
- auth_version para invalidar sesiones antiguas;
- rate limit progresivo por IP/identidad;
- logs de seguridad;
- recuperación de contraseña opcional por SMTP;
- token de recovery de 256 bits;
- SQLite guarda solo hash del token;
- expiración y un solo uso;
- token en fragmento #token para evitar access logs;
- SMTP remoto requiere TLS/STARTTLS;
- recovery oculto si correo no está configurado.

## Explorador
Ruta:
`/archivos/*path`

Workspace:
`FILES_ROOT`, default `storage/demo-files`

Protecciones:
- traversal;
- rutas absolutas;
- backslashes;
- percent-encoding ambiguo;
- symlinks;
- junctions/reparse points;
- nombres reservados;
- controles invisibles/bidi;
- límites de longitud/profundidad;
- solo texto UTF-8 editable;
- tamaño máximo configurable;
- create/save/delete con sesión + CSRF.

Editor:
- estilo VS Code ligero;
- Highlight.js;
- números de línea;
- Ctrl/Cmd+S;
- dirty state;
- modal al cerrar sin guardar;
- fallback visible si Highlight.js falla.

Pendiente para más adelante en Explorador:
- visor de imágenes;
- reproductor de video;
- copiar;
- cortar;
- pegar;
- selección múltiple con Ctrl;
- selección mediante checkboxes tipo Windows.

No implementar estas funciones de golpe: diseñar primero modelo de selección/clipboard y permisos backend.

## Módulo Server
### /servidor
Resumen técnico actual.

### /servidor/tareas
Administrador de **plantillas** systemd:
- listar;
- crear .service;
- crear .timer;
- ver;
- eliminar.

Storage:
`SERVER_TASKS_ROOT=storage/server-tasks`

Regla de seguridad:
- NO ejecutar systemctl;
- NO escribir en /etc/systemd;
- NO instalar/habilitar/iniciar/detener servicios;
- solo generar plantillas validadas dentro del storage configurado.

### /servidor/entorno
Después de revisión visual se simplificó:
- menos decoración;
- sin tarjetas/iconos grandes por cada dato;
- resumen compacto;
- paneles técnicos;
- datos label/value;
- solo información no sensible.

## CSP
CSP permite connect-src a:
- self
- cdn.jsdelivr.net
- cdnjs.cloudflare.com

Esto evita warnings por source maps de librerías autorizadas sin abrir conexiones arbitrarias.

## Pendiente principal: CRUD de usuarios + multirol/RBAC

Archivo de tarea:
`tareas/23-pendiente-crud-usuarios-rbac.md`

Contexto detallado:
`contexto/22-cierre-chat-roadmap-rbac.md`

### Investigación verificada
Goravel tiene autorización nativa mediante:
- Gates
- Policies

Los Form Requests incluyen:
- Authorize()
- Rules()

Fuentes oficiales:
- https://www.goravel.dev/security/authorization.html
- https://www.goravel.dev/the-basics/validation.html

Goravel-Admin existe como proyecto comunitario, pero su checklist público todavía muestra partes del CRUD incompletas:
- https://github.com/goravel-community/goravel-admin

Casbin soporta RBAC jerárquico:
- https://www.casbin.org/docs/rbac

Decisión actual:
**No instalar Casbin ni Goravel-Admin todavía.**
Implementar RBAC propio en SQLite y usar Gates/Policies de Goravel para enforcement.

### Modelo propuesto
```text
users
roles
permissions
user_roles
role_permissions
```

Un usuario puede tener 0..N roles.
Los permisos efectivos son la unión de los permisos de sus roles.

No empezar con permisos directos por usuario.

### Roles iniciales sugeridos
- super-admin
- admin
- operator
- viewer

super-admin debe ser rol de sistema protegido.

### Reglas críticas RBAC
- backend siempre autoriza;
- ocultar botón no es seguridad;
- no permitir eliminar/desactivar el último super-admin;
- no permitir autoescalado;
- cambios de roles/permisos dentro de transacción;
- slugs únicos y estables;
- incrementar auth_version cuando cambien roles/permisos del usuario;
- registrar cambios sensibles en auditoría;
- CSRF en mutaciones;
- filtrar sidebar por permiso, pero volver a autorizar en backend.

### CRUD de usuarios
Ruta recomendada:
`/usuarios`

Debe permitir:
- listar/buscar/paginar;
- crear;
- editar;
- activar/desactivar;
- recovery/reset;
- asignar múltiples roles;
- badges/chips;
- ver resumen de permisos efectivos;
- auditoría de acciones sensibles.

Selector multirol profesional:
- dropdown searchable;
- checkboxes;
- chips de roles seleccionados;
- nombre + descripción por rol;
- indicador de rol de sistema;
- resumen de permisos efectivos antes de guardar;
- UI móvil tipo panel/lista.

### CRUD de roles
Ruta recomendada:
`/roles`

Debe permitir:
- listar;
- crear/editar;
- activar/desactivar cuando aplique;
- eliminar solo roles no protegidos;
- matriz de permisos agrupada por módulo;
- seleccionar módulo completo;
- búsqueda de permisos;
- contador;
- usuarios asociados.

## Otras páginas recomendadas después de RBAC
1. Auditoría
2. Logs
3. Jobs/colas
4. Backups
5. Configuración
6. API Keys/Tokens
7. Webhooks
8. Notificaciones/Correo
9. Integraciones
10. Sesiones activas

Prioridad recomendada:
**Usuarios + Roles/Permisos antes de seguir agregando muchas páginas administrativas**, porque la autorización condicionará el resto.

## Metodología
Usar Ponytail:
- contexto/ para decisiones duraderas;
- tareas/ para estado;
- una sola tarea en-proceso durante trabajo sustantivo;
- validar;
- commit en español;
- actualizar continuidad.

## Git
El usuario pidió al cerrar este chat:
- guardar el contexto en `contexto/`;
- comitear todo lo que estuviera pendiente;
- dejar RBAC/CRUD de usuarios como pendiente para el próximo chat.

La instrucción anterior de no comitear el módulo Server quedó reemplazada por esta solicitud posterior.
