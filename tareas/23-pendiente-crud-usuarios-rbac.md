# Tarea 23 — CRUD de usuarios + RBAC multirol

## Estado
Pendiente para el próximo chat.

## Prioridad
Alta. Debe implementarse antes de seguir agregando muchas páginas administrativas.

## Objetivo
Convertir el esquema actual de rol único en un sistema multirol profesional y crear CRUD 100% funcional de usuarios y roles/permisos.

## Arquitectura acordada
- users
- roles
- permissions
- user_roles
- role_permissions
- RBAC propio en SQLite.
- Goravel Gates + Policies para autorización real.
- super-admin protegido.
- roles múltiples por usuario.
- permisos efectivos = unión de permisos de sus roles.
- incrementar auth_version al cambiar roles/permisos de un usuario.
- sidebar filtrado por permiso, pero backend siempre vuelve a autorizar.

## UI usuarios
- tabla buscable/paginada;
- crear/editar/activar/desactivar;
- selector multirol searchable con checkboxes y chips;
- badges de roles;
- resumen de permisos efectivos;
- recovery/reset;
- acciones sensibles auditadas.

## UI roles
- CRUD de roles;
- matriz de permisos agrupada por módulo;
- seleccionar todo por módulo;
- búsqueda;
- contador de permisos;
- indicar roles de sistema;
- impedir eliminar roles críticos.

## Seguridad a probar
- un usuario no puede autoasignarse permisos;
- no se puede eliminar/desactivar el último super-admin;
- autorización se valida en backend;
- CSRF;
- transacciones en sincronización de roles;
- cambios invalidan sesiones anteriores con auth_version;
- roles/permisos eliminados no dejan relaciones huérfanas.

## Investigación
Ver contexto/22-cierre-chat-roadmap-rbac.md.
