# Fecha
2026-09-28

# Objetivo
Convertir la plantilla visual Axolotl entregada como un único HTML en una plantilla modular y reutilizable construida sobre Goravel, 100% Go, con estructura similar a Laravel, vistas server-side, rutas, middlewares, ORM GORM y SQLite embebido/configurable por `.env`.

# Decisiones tomadas
- Repositorio raíz: `C:\Users\Admin\Desktop\GitHub\axolotl-dashboard-template`.
- Rama principal: `main`.
- Identidad Git: usar configuración global YahirHub.
- Framework elegido: Goravel.
- Se verificó documentación oficial actual: Goravel v1.18.0 es la versión estable más reciente encontrada.
- Goravel requiere Go >= 1.23; el entorno dispone de Go 1.26.4.
- Vistas: `html/template` de Goravel en `resources/views`, con layouts y partials.
- Base de datos: SQLite como default, con configuración `DB_CONNECTION=sqlite` y `DB_DATABASE` en `.env`.
- ORM: Goravel ORM, respaldado por GORM.
- No usar ctx7 salvo último recurso; usar primero búsqueda/documentación oficial.
- La plantilla original se conserva en `docs/reference/index-original.html`.
- `design.md` se conserva como fuente de verdad visual.

# Arquitectura actual
Proyecto aún en fase de bootstrap. Se importó la referencia original y se prepara el scaffold Goravel.

# Librerías usadas
Todavía no instaladas en el repositorio. Objetivo inicial:
- github.com/goravel/framework
- driver HTTP oficial github.com/goravel/gin
- driver SQLite oficial github.com/goravel/sqlite
- GORM indirectamente a través de Goravel ORM/driver

# Archivos importantes modificados
- LILITH.md
- design.md
- docs/reference/index-original.html
- contexto/01-contexto-inicial.md
- tareas/*

# Problemas encontrados
- El ZIP original contenía una carpeta anidada `axolotl-template/`; se normalizó al root del nuevo repositorio.
- La regla global del host favorece ctx7 para librerías, pero la instrucción explícita del proyecto exige usarlo solo como último recurso. La regla del proyecto prevalece para este repositorio.

# Soluciones implementadas
- Creación del workspace PL.
- Inicialización de Git en rama main.
- Creación de reglas persistentes, contexto y tareas Ponytail.
- Preservación de la referencia frontend original.

# Pendientes
- Crear scaffold Goravel.
- Instalar/configurar SQLite.
- Modularizar frontend en assets + templates.
- Crear controller/rutas/middleware base.
- Añadir migración/modelo demo para validar GORM/SQLite.
- Ejecutar gofmt, go test, go vet y smoke test HTTP.
- Crear commits incrementales.

# Próximos pasos
1. Commit inicial con referencia y metodología.
2. Integrar scaffold Goravel v1.18.x.
3. Configurar SQLite y entorno.
4. Modularizar dashboard.
5. Validar y documentar.
