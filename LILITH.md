# Reglas del proyecto Axolotl Dashboard Template

## Stack y arquitectura
- El backend y la aplicación deben ser 100% Go.
- Framework principal: Goravel, siguiendo su estructura Laravel-like.
- HTTP: driver oficial de Goravel (Gin salvo decisión documentada posterior).
- Vistas: motor nativo `html/template` de Goravel, organizado con layouts y partials como equivalente funcional a Blade.
- ORM: Goravel ORM/GORM.
- Base de datos por defecto: SQLite embebido.
- La conexión de base de datos se selecciona con variables `.env` estilo Laravel, especialmente `DB_CONNECTION` y `DB_DATABASE`.
- Mantener capacidad futura de cambiar de driver mediante configuración, sin acoplar dominio a SQLite.
- Autenticación base: guard de sesión de Goravel + tabla `users`.
- Contraseñas: Argon2id; nunca guardar ni registrar contraseñas en texto plano.
- El primer rol `admin` solo puede crearse mediante `/setup` con el PIN temporal del arranque.
- El registro público solo crea rol `user`; nunca confiar en un rol enviado por el cliente.
- Mantener CSRF en mutaciones web y rate limiting en setup/login/registro.
- Todas las áreas privadas deben usar `RequireAuth` o un middleware de autorización más específico.
- El explorador de archivos solo puede operar dentro de `FILES_ROOT`; nunca exponer acceso arbitrario al filesystem del host.
- Mantener protección contra path traversal y symlinks en cualquier nueva operación de archivos.
- Create/save/delete del explorador requieren sesión y CSRF; no relajar estas protecciones para simplificar la UI.

## Investigación y documentación
- NO usar Context7/ctx7 mientras existan herramientas de búsqueda web, documentación oficial, repositorios oficiales, Knowledge local u otras fuentes disponibles.
- ctx7 es únicamente último recurso cuando no haya otra fuente viable.
- Priorizar documentación oficial de Goravel, GORM y paquetes oficiales.
- Verificar versiones y APIs actuales antes de decisiones de framework o dependencias.

## Metodología
- Aplicar siempre la skill `ponytail-development`.
- Priorizar simplicidad, seguridad, pocas dependencias y arquitectura mantenible.
- No crear capas vacías, interfaces de una sola implementación o abstracciones especulativas.
- Toda lógica no trivial debe poder validarse con pruebas o smoke tests.

## Persistencia de contexto
- Mantener `contexto/` actualizado después de decisiones o cambios importantes.
- Mantener `tareas/` con exactamente una tarea `en-proceso-` a la vez.
- Guardar también contexto durable en el MCP de PL mediante continuity/memory/checkpoint.
- Antes de trabajo sustancial en otra sesión, recuperar continuidad del workspace.

## Git
- Rama principal: `main`.
- Usar la identidad Git global existente del usuario YahirHub; no sobrescribir `user.name` ni `user.email`.
- Commits en español.
- No reinicializar Git después del inicio.
- No reescribir historia publicada salvo petición explícita.
- Antes de commit: revisar status y diff staged.
- Después de commit: validar status y pruebas.

## Frontend
- Preservar el design system definido en `design.md`.
- El dashboard debe seguir siendo multipropósito/general.
- No dejar controles aparentemente clicables sin función.
- Mantener responsive, accesibilidad básica, skeletons, modales, buscador, submenús y tokens de tema.
- La referencia original se conserva en `docs/reference/index-original.html` como baseline visual y de regresión; la implementación activa ya está modularizada.
