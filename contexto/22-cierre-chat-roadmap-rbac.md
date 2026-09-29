# Cierre de contexto — Server, usuarios y RBAC

## Fecha
2026-09-29

## Objetivo
Dejar suficiente contexto para continuar Axolotl en un chat nuevo sin perder las decisiones de arquitectura ni el trabajo pendiente sobre usuarios, roles y permisos.

## Estado actual que queda versionado
- El sidebar sigue la regla **route-only**: no debe contener links con `#` ni acciones JS que no representen una página.
- El dropdown **Server** agrupa páginas reales:
  - `/servidor`
  - `/servidor/tareas`
  - `/servidor/entorno`
- `/servidor/tareas` administra plantillas locales `.service` y `.timer`.
- Las plantillas se guardan en `SERVER_TASKS_ROOT` (default `storage/server-tasks`).
- Axolotl NO ejecuta `systemctl`, no instala unidades y no escribe en `/etc/systemd`.
- `/servidor/entorno` se simplificó después de revisión visual:
  - sin iconos decorativos grandes;
  - resumen técnico compacto;
  - layout de dos paneles;
  - pares label/value;
  - responsive.
- CSP permite conexiones a los CDN ya autorizados para evitar errores de source maps de Chart.js/Lenis en consola.
- El módulo Server tiene tests unitarios y feature.
- El usuario pidió en este turno **commit de todo lo pendiente**, por lo que queda anulada la decisión anterior de “no commit hasta aprobación”.

## Pendiente principal: CRUD de usuarios + multirol profesional

### Situación actual
El modelo `User` todavía tiene un único campo string `Role` con valores base `admin` y `user`. Esto sirve para el setup inicial, pero NO es el diseño final para un panel multipropósito multirol.

### Conclusión de investigación sobre Goravel
La documentación oficial actual de Goravel incluye autorización nativa mediante **Gates** y **Policies**. Es el mecanismo recomendado para decidir si un usuario puede ejecutar una acción o modificar un recurso.

Goravel no ofrece en su core un CRUD completo de roles/permisos estilo “Spatie Permission”. Existe un proyecto comunitario Goravel-Admin, pero su propio checklist aún muestra partes del CRUD incompletas, por lo que no conviene usarlo como dependencia central de autorización de Axolotl.

Casbin es una alternativa madura en Go para RBAC/ABAC y soporta roles jerárquicos y políticas, pero para este starter añadiría una capa extra y duplicaría parte de la autorización que Goravel ya resuelve con Gates/Policies. La propuesta preferida para Axolotl es: **RBAC propio en SQLite + Gates/Policies de Goravel para enforcement**.

Fuentes verificadas el 2026-09-29:
- https://www.goravel.dev/security/authorization.html — Goravel documenta Gates y Policies como mecanismos nativos de autorización.
- https://www.goravel.dev/the-basics/validation.html — los Form Requests incluyen Authorize() y Rules().
- https://github.com/goravel-community/goravel-admin — proyecto comunitario; su checklist público todavía marca CRUD list/edit/create como incompleto.
- https://www.casbin.org/docs/rbac — Casbin soporta RBAC y jerarquías de roles, pero no se añadirá mientras el alcance actual pueda resolverse con RBAC propio + Gates/Policies.

### Modelo multirol propuesto
Un usuario puede tener **cero, uno o varios roles**. Los permisos efectivos son la unión de los permisos de todos sus roles.

Tablas propuestas:

```text
users
roles
permissions
user_roles
role_permissions
```

Campos sugeridos:

```text
roles
- id
- name
- slug (unique)
- description
- is_system
- active
- timestamps

permissions
- id
- name
- slug (unique)
- module
- description
- timestamps

user_roles
- user_id
- role_id
- unique(user_id, role_id)

role_permissions
- role_id
- permission_id
- unique(role_id, permission_id)
```

No empezar con permisos directos por usuario. Mantener inicialmente el modelo limpio: **usuario -> roles -> permisos**.

### Roles iniciales sugeridos
- `super-admin`: rol de sistema protegido; puede realizar cualquier acción.
- `admin`: administración normal del panel según permisos explícitos.
- `operator`: operación cotidiana sin configuración crítica.
- `viewer`: solo lectura.

Los nombres son seeds de ejemplo y deben poder adaptarse a cada producto.

### Permisos
Usar slugs explícitos y estables, por ejemplo:

```text
users.view
users.create
users.update
users.delete
users.assign_roles

roles.view
roles.create
roles.update
roles.delete
roles.assign_permissions

server.view
server.tasks.view
server.tasks.create
server.tasks.delete

files.view
files.create
files.update
files.delete

settings.view
settings.update
audit.view
```

Evitar usar nombres de rutas como permisos. Las rutas pueden cambiar; la capacidad de negocio debe conservar un slug estable.

### Enforcement
- Backend siempre manda: esconder un botón en frontend NO es seguridad.
- Crear middleware/servicio reutilizable para permisos de ruta cuando aplique.
- Usar `facades.Gate()` para capacidades globales.
- Usar **Policies** para reglas sobre un recurso concreto.
- Usar `Gate.Before` o equivalente para el bypass de `super-admin`.
- Form Requests pueden usar `Authorize()` para combinar validación + autorización.
- El sidebar se filtra por permisos, pero todas las rutas/controladores siguen validando autorización en backend.
- Cuando un administrador modifica los roles/permisos de un usuario, incrementar `auth_version` de ese usuario para invalidar sesiones antiguas y aplicar permisos nuevos inmediatamente.

### Reglas críticas
- No permitir eliminar/desactivar el último `super-admin` activo.
- No permitir que un usuario se quite a sí mismo el último rol administrativo si deja el sistema sin administrador.
- Roles de sistema pueden editarse con límites o marcarse como no eliminables.
- Cambios de roles/permisos deben ejecutarse dentro de transacciones.
- Slugs de roles/permisos son únicos e inmutables salvo migración controlada.
- Registrar en auditoría creación, edición, activación/desactivación y cambios de roles.
- CSRF en todas las mutaciones.
- Nunca confiar en un role/permission enviado por el navegador sin verificarlo en backend.

### CRUD de usuarios propuesto
Ruta principal recomendada: `/usuarios`.

Funciones:
- listado paginado/buscable;
- crear usuario;
- editar nombre, correo, alias y estado;
- activar/desactivar;
- restablecer contraseña o enviar recovery;
- asignar múltiples roles;
- ver roles como badges/chips;
- impedir autoescalado de permisos;
- cambios sensibles auditados.

Selector de roles profesional:
- no usar el `select multiple` nativo como UI principal;
- dropdown searchable con checkboxes;
- roles elegidos se muestran como chips/badges;
- cada opción incluye nombre + breve descripción;
- roles de sistema llevan indicador “Sistema”;
- mostrar resumen de permisos efectivos antes de guardar;
- en móvil debe convertirse en panel/lista cómoda de tocar.

### CRUD de roles propuesto
Ruta recomendada: `/roles`.

Funciones:
- listar roles;
- crear/editar;
- activar/desactivar cuando corresponda;
- eliminar solo roles no protegidos y no críticos;
- asignar permisos mediante **matriz agrupada por módulo**;
- checkbox “Seleccionar módulo”;
- contador de permisos seleccionados;
- búsqueda de permisos;
- resumen de usuarios que usan el rol.

La matriz de permisos es preferible a una lista plana cuando Axolotl crezca.

## Otras páginas recomendadas para el roadmap
Después de Usuarios/Roles, las páginas con más valor para un panel multipropósito son:

1. **Auditoría** — acciones administrativas, cambios de roles, accesos y operaciones sensibles.
2. **Logs** — visor seguro de logs de aplicación/seguridad, probablemente dentro de Server.
3. **Jobs / Colas** — estado de trabajos en background cuando se incorporen queues.
4. **Backups** — estado y ejecución controlada de backups, sin exponer rutas arbitrarias.
5. **Configuración** — settings persistentes del producto, con secciones y permisos.
6. **API Keys / Tokens** — emisión, revocación, scopes y último uso.
7. **Webhooks** — endpoints, eventos, secretos rotables e historial de entregas.
8. **Notificaciones / Correo** — plantillas, cola y estado de entrega.
9. **Integraciones** — conectores externos por producto.
10. **Auditoría de sesiones** — sesiones activas, cerrar otras sesiones y eventos de login.

No implementar todas de golpe. Usuarios + Roles/Permisos debe ser la siguiente feature estructural porque condiciona la autorización de las demás páginas.

## Orden recomendado para el próximo chat
1. Recuperar este contexto y revisar `git status`.
2. Crear una tarea Ponytail nueva para RBAC.
3. Diseñar migraciones y seeds de roles/permisos.
4. Migrar gradualmente `users.role` al esquema multirol sin romper setup/login.
5. Implementar Gate/Policies y middleware de permisos.
6. Crear CRUD de usuarios.
7. Crear CRUD de roles con matriz de permisos.
8. Filtrar sidebar por permiso.
9. Añadir auditoría de cambios de seguridad.
10. Tests: escalación, último admin, multirol, autorización backend, sesiones invalidadas y CSRF.

## No hacer todavía
- No instalar Casbin salvo que aparezca una necesidad real de políticas complejas, multi-tenant o ABAC.
- No depender de Goravel-Admin para el core de permisos.
- No implementar permisos solo en frontend.
- No eliminar `auth_version`; será útil al cambiar roles.
