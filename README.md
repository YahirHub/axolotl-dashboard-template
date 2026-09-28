# Axolotl Dashboard Template

Plantilla de dashboard multipropósito con identidad visual Axolotl, construida como aplicación modular **100% Go** sobre [Goravel](https://www.goravel.dev/).

La interfaz original se conserva como referencia visual, pero la aplicación activa ya utiliza controllers, rutas, middleware, vistas server-side, configuración por `.env`, ORM y SQLite.

## Stack

- **Go 1.25+** — probado con Go 1.26.4.
- **Goravel 1.18**.
- **Gin** como driver HTTP oficial de Goravel.
- **html/template** como motor de vistas nativo, organizado en páginas y partials.
- **Goravel ORM / GORM**.
- **SQLite** embebido mediante el driver oficial `github.com/goravel/sqlite`.
- Chart.js, Lenis, AOS, Highlight.js y Font Awesome por CDN para la demo frontend.

## Inicio rápido

### 1. Preparar el entorno

En PowerShell:

```powershell
Copy-Item .env.example .env
go run . artisan key:generate
go run . artisan migrate
```

En Linux/macOS:

```bash
cp .env.example .env
go run . artisan key:generate
go run . artisan migrate
```

### 2. Ejecutar

```bash
go run .
```

Abrir:

```text
http://127.0.0.1:3000
```

Healthcheck:

```text
http://127.0.0.1:3000/health
```

## Configuración

La configuración sigue el enfoque de Laravel mediante variables de entorno.

```dotenv
APP_NAME=Axolotl
APP_ENV=local
APP_DEBUG=true
APP_HOST=127.0.0.1
APP_PORT=3000

DB_CONNECTION=sqlite
DB_DATABASE=storage/axolotl.sqlite
```

`DB_DATABASE` acepta una ruta relativa al proyecto o una ruta absoluta.

La plantilla incluye únicamente SQLite para mantener el starter pequeño. Para usar PostgreSQL, MySQL o SQL Server se debe instalar el driver oficial de Goravel correspondiente y agregar su conexión a `config/database.go`.

## Estructura

```text
app/
├── facades/
├── http/
│   ├── controllers/
│   │   ├── dashboard_controller.go
│   │   └── health_controller.go
│   └── middleware/
│       └── security_headers.go
└── models/
    └── app_setting.go

bootstrap/
├── app.go
├── migrations.go
└── providers.go

config/
database/
├── migrations/
public/
├── css/app.css
└── js/app.js
resources/
└── views/
    ├── dashboard.tmpl
    └── partials/
        ├── head.tmpl
        ├── scripts.tmpl
        ├── sidebar.tmpl
        ├── dashboard-main.tmpl
        ├── topbar.tmpl
        ├── stats.tmpl
        ├── dashboard-grid.tmpl
        ├── server-metrics.tmpl
        ├── code-lab.tmpl
        └── overlays.tmpl
routes/
├── web.go
└── grpc.go
tests/
└── feature/
contexto/
tareas/
design.md
LILITH.md
```

## Vistas: equivalente funcional a Blade

Goravel integra `html/template`. Axolotl usa templates nombrados y partials:

```gotemplate
{{ template "partials/sidebar" . }}
{{ template "partials/dashboard-main" . }}
```

No se agregó otro motor de plantillas: el objetivo es conservar una solución nativa, rápida y con pocas dependencias.

La fuente de verdad visual y las reglas para crear nuevos componentes están en:

```text
design.md
```

La referencia HTML previa a la modularización permanece en:

```text
docs/reference/index-original.html
```

## Rutas

```text
GET /           dashboard
GET /health     healthcheck con prueba de SQLite
/assets/*       archivos estáticos
```

Consultar las rutas registradas:

```bash
go run . artisan route:list
```

## Base de datos y migraciones

Ejecutar migraciones:

```bash
go run . artisan migrate
```

Ver estado:

```bash
go run . artisan migrate:status
```

La migración de ejemplo `app_settings` y el modelo `AppSetting` validan la integración de Goravel ORM con SQLite.

Los archivos SQLite locales se guardan en `storage/` y están ignorados por Git.

## Pruebas y validación

```bash
gofmt -w app bootstrap config database routes tests main.go
go test ./...
go vet ./...
go build ./...
```

Las pruebas actuales verifican:

- round-trip ORM sobre SQLite;
- render del dashboard;
- referencias a CSS/JS;
- healthcheck;
- headers globales de seguridad.

## Docker

Crear primero el `.env` y generar `APP_KEY`.

```bash
docker compose up --build
```

Docker fuerza `APP_HOST=0.0.0.0` dentro del contenedor y monta `./storage` para persistir SQLite.

La imagen **no incluye** el archivo `.env`; las variables se cargan en runtime.

## Seguridad base

El middleware global agrega:

- `X-Content-Type-Options: nosniff`;
- `X-Frame-Options: SAMEORIGIN`;
- `Referrer-Policy: strict-origin-when-cross-origin`;
- `Permissions-Policy` deshabilitando cámara, micrófono y geolocalización por defecto.

No se aplica una CSP estricta todavía porque la demo usa recursos CDN. Si los assets externos se migran a local, `Content-Security-Policy` debe añadirse como siguiente endurecimiento.

## Desarrollo del proyecto

- Rama principal: `main`.
- Commits en español.
- `contexto/`: decisiones técnicas persistentes.
- `tareas/`: estado de trabajo según metodología Ponytail.
- `LILITH.md`: reglas obligatorias para agentes/herramientas que continúen el proyecto.
- No usar ctx7 salvo último recurso; priorizar documentación oficial y herramientas de búsqueda disponibles.

## Licencia

La base Goravel conserva su licencia MIT en `LICENSE`. Revisa la licencia de los recursos externos utilizados antes de redistribuir una versión comercial cerrada.
