# Fecha
2026-09-28

# Objetivo
Cerrar la primera versión utilizable de Axolotl Dashboard Template como starter modular 100% Go basado en Goravel.

# Decisiones tomadas
- El proyecto queda en rama `main` con historial Git incremental y autoría YahirHub.
- Goravel v1.18.0 es la base del starter.
- Go mínimo práctico del módulo: 1.25; validado localmente con Go 1.26.4.
- Gin es el driver HTTP.
- `html/template` integrado por Goravel/Gin es el motor de vistas, organizado en templates nombrados y partials.
- SQLite es la base embebida por defecto, seleccionada por `.env`.
- Las pruebas usan un SQLite temporal separado del archivo de desarrollo.
- `/health` verifica aplicación y conexión SQLite.
- El middleware global aplica headers de seguridad compatibles con la demo basada en CDN.
- La referencia original se mantiene como baseline en `docs/reference/index-original.html`.
- `design.md` es la fuente de verdad visual.
- Docker no copia `.env` dentro de la imagen; usa variables de runtime y un volumen para `storage/`.
- ctx7 permanece prohibido salvo último recurso; toda investigación de framework se hizo mediante documentación oficial y herramientas disponibles.

# Arquitectura actual
```text
Route
  -> Middleware global
    -> Controller
      -> Goravel ORM / GORM cuando aplica
        -> SQLite
      -> html/template o JSON
```

Frontend:
```text
resources/views/*.tmpl
public/css/app.css
public/js/app.js
```

Persistencia:
```text
config/database.go
app/models/
database/migrations/
storage/axolotl.sqlite
```

# Librerías usadas
Directas:
- github.com/goravel/framework v1.18.0
- github.com/goravel/gin v1.18.0
- github.com/goravel/sqlite v1.18.0
- github.com/stretchr/testify v1.11.1

Frontend por CDN:
- Chart.js
- Lenis
- AOS
- Highlight.js
- Font Awesome

# Archivos importantes modificados
- README.md
- design.md
- LILITH.md
- .gitattributes
- .dockerignore
- Dockerfile
- docker-compose.yml
- tests/test_case.go

# Problemas encontrados
- El scaffold de Goravel incluía PostgreSQL/OpenAI que no eran necesarios para esta plantilla.
- SQLite necesitaba resolver rutas relativas desde la raíz para funcionar también bajo tests.
- Los tests inicialmente compartían el archivo SQLite de desarrollo.
- La referencia frontend tenía valores RGB heredados de una paleta turquesa.
- Los procesos hijos de `go run .` en Windows pueden mantener temporalmente el puerto 3000 después de cancelar el wrapper.
- Docker no está instalado en el entorno PL actual, por lo que Dockerfile/Compose no pudieron ejecutarse localmente.
- PL no expone navegador gráfico; la validación frontend realizada fue HTTP/HTML/asset smoke testing, no una auditoría visual automatizada.

# Soluciones implementadas
- Dependencias directas reducidas a framework, Gin, SQLite y testing.
- SQLite aislado para tests mediante archivo temporal del sistema.
- Dashboard renderizado por Goravel y assets servidos bajo `/assets/*`.
- Pruebas HTTP para dashboard, health y headers.
- Prueba ORM round-trip.
- Dockerfile multi-stage sin secretos, usuario no root y storage persistente.
- README y design system actualizados.

# Validaciones ejecutadas
Exitosas:
- `gofmt`
- `go mod tidy`
- `go test ./... -count=1`
- `go vet ./...`
- `go build ./...`
- `go run . artisan migrate:fresh`
- `go run . artisan migrate:status`
- `go run . artisan route:list`
- smoke test real de `/`, `/health`, CSS y JavaScript
- validación completa del MCP PL (`code_validate`)

No ejecutada por falta de herramienta:
- `docker compose build/up`
- auditoría visual en navegador real desde PL

# Pendientes
La primera versión queda funcional. Mejoras futuras opcionales:
- mover dependencias CDN a assets locales y agregar CSP estricta;
- añadir autenticación/roles solo cuando un producto concreto lo requiera;
- conectar métricas del servidor a datos reales;
- crear CI en GitHub cuando exista remoto;
- añadir más páginas/módulos según el producto que use la plantilla.

# Próximos pasos
Usar esta versión como baseline. Toda feature futura debe abrir una nueva tarea en `tareas/`, actualizar `contexto/` y conservar las reglas de `LILITH.md` y `design.md`.
