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
- La referencia original se conserva en `docs/reference/index-original.html` mientras se completa la modularización.
