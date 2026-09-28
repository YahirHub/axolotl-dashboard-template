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

### 3. Primer inicio

Si todavía no existe un usuario activo con rol `admin`, el arranque genera un **PIN temporal de 6 dígitos** usando `crypto/rand`.

El PIN se muestra únicamente en el log/terminal del servidor y también queda en:

```text
storage/logs/setup.log
```

Ejemplo:

```text
AXOLOTL SETUP PIN: 123456
No existe un administrador. Abre /setup para crear la cuenta inicial.
```

Al abrir `/`:

1. si no existe administrador, redirige a `/setup`;
2. `/setup` solicita nombre, correo, alias, contraseña y los 6 dígitos del PIN;
3. el rol `admin` se asigna exclusivamente desde el backend;
4. al terminar, Axolotl crea la cuenta e inicia la sesión automáticamente;
5. el PIN queda invalidado y `/setup` deja de permitir otra creación administrativa.

El PIN no se guarda en SQLite: vive únicamente en memoria y cambia al reiniciar mientras el setup siga pendiente.

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

SESSION_DRIVER=file
SESSION_LIFETIME=120
SESSION_SECURE=false

AUTH_ALLOW_REGISTRATION=true

FILES_ROOT=storage/demo-files
FILES_MAX_EDIT_BYTES=2097152
```

`DB_DATABASE` acepta una ruta relativa al proyecto o una ruta absoluta.

La plantilla incluye únicamente SQLite para mantener el starter pequeño. Para usar PostgreSQL, MySQL o SQL Server se debe instalar el driver oficial de Goravel correspondiente y agregar su conexión a `config/database.go`.

## Estructura

```text
app/
├── facades/
├── files/
│   ├── configured.go
│   └── workspace.go
├── http/
│   ├── controllers/
│   │   ├── auth_controller.go
│   │   ├── dashboard_controller.go
│   │   ├── file_controller.go
│   │   ├── page_controller.go
│   │   └── health_controller.go
│   └── middleware/
│       ├── require_auth.go
│       └── security_headers.go
├── models/
│   ├── app_setting.go
│   └── user.go
└── security/
    └── setup.go

bootstrap/
├── app.go
├── migrations.go
└── providers.go

config/
database/
├── migrations/
public/
├── brand/
│   ├── axolotl.svg
│   └── axolotl.png
├── css/
│   ├── app.css
│   ├── auth.css
│   └── files.css
└── js/
    ├── app.js
    ├── auth.js
    └── files.js
resources/
└── views/
    ├── auth.tmpl
    ├── dashboard.tmpl
    ├── account.tmpl
    ├── server.tmpl
    ├── components.tmpl
    ├── files.tmpl
    └── partials/
        ├── head.tmpl
        ├── scripts.tmpl
        ├── sidebar.tmpl
        ├── dashboard-main.tmpl
        ├── topbar.tmpl
        ├── stats.tmpl
        ├── dashboard-grid.tmpl
        ├── workspace-grid.tmpl
        ├── files-main.tmpl
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

La aplicación es multipágina. `/dashboard` concentra el overview de negocio; `/cuenta`, `/servidor`, `/componentes` y `/archivos` separan configuración, observabilidad, catálogo UI y workspace de archivos. El `DashboardController` entrega datos demostrativos para KPIs, proyectos, tareas, actividad, equipo, agenda e integraciones; esos datos pueden reemplazarse después por ORM, APIs o servicios sin reconstruir el layout.

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
GET  /              decide setup, login o dashboard
GET  /setup         setup del primer administrador
POST /setup         crea el primer administrador con PIN temporal
GET  /login         inicio de sesión
POST /login         autentica usuario
GET  /register      registro de usuario estándar
POST /register      crea usuario con rol user
POST /logout        cierra sesión

GET  /dashboard     overview general [auth]
GET  /cuenta        perfil, preferencias y seguridad [auth]
GET  /servidor      métricas técnicas [auth]
GET  /componentes   laboratorio UI y catálogo de componentes [auth]
GET  /archivos      explorador del workspace demo [auth]

GET  /api/files         lista una carpeta del workspace [auth]
GET  /api/files/read    abre un archivo de texto/código [auth]
POST /api/files/create  crea archivo o carpeta [auth + CSRF]
POST /api/files/save    guarda contenido de archivo [auth + CSRF]
POST /api/files/delete  elimina archivo o carpeta [auth + CSRF]

GET  /health        healthcheck con prueba de SQLite
/assets/*           archivos estáticos
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

La migración `users` almacena las cuentas con correo único, contraseña hasheada, rol y estado activo. El modelo `User` define los roles base `admin` y `user`.

La migración de ejemplo `app_settings` y el modelo `AppSetting` siguen disponibles como ejemplo de persistencia general.

Los archivos SQLite locales se guardan en `storage/` y están ignorados por Git.

El explorador usa por defecto `storage/demo-files`. Al iniciar por primera vez, el backend agrega archivos de ejemplo si la carpeta está vacía. El contenido mutable de ese workspace también está ignorado por Git.

## Pruebas y validación

```bash
gofmt -w app bootstrap config database routes tests main.go
go test ./...
go vet ./...
go build ./...
```

Las pruebas actuales verifican:

- round-trip ORM sobre SQLite;
- redirección del primer inicio hacia `/setup`;
- rechazo de PIN de setup incorrecto;
- creación del primer administrador con PIN correcto;
- inicio automático de sesión después del setup;
- login y logout con sesión;
- registro de usuarios estándar;
- que un cliente no puede autoasignarse rol `admin` desde `/register`;
- protección de las páginas privadas;
- render de dashboard, cuenta, servidor, componentes y archivos;
- CRUD del workspace demo mediante backend Go;
- rechazo de path traversal, symlinks y archivos no editables;
- límites de tamaño para edición;
- healthcheck y headers globales de seguridad.

## Docker

Crear primero el `.env` y generar `APP_KEY`.

```bash
docker compose up --build
```

Docker fuerza `APP_HOST=0.0.0.0` dentro del contenedor y monta `./storage` para persistir SQLite.

La imagen **no incluye** el archivo `.env`; las variables se cargan en runtime.

## Seguridad base

La autenticación incluye:

- guard de sesión de Goravel;
- regeneración de sesión al iniciar y cerrar sesión;
- contraseñas con **Argon2id**;
- CSRF global para formularios POST;
- rate limit por IP para setup, login y registro;
- PIN inicial generado con `crypto/rand` y comparación constant-time;
- PIN guardado solo en memoria;
- rol `admin` asignable únicamente durante el setup inicial;
- registro público configurable con `AUTH_ALLOW_REGISTRATION`;
- rutas privadas protegidas por `RequireAuth`;
- explorador confinado a `FILES_ROOT`, sin acceso arbitrario al filesystem del host;
- rechazo de traversal (`../`) y symlinks;
- edición limitada a texto UTF-8 y al tamaño `FILES_MAX_EDIT_BYTES`;
- mutaciones de archivos protegidas con sesión y CSRF.

En producción configura HTTPS y cambia:

```dotenv
SESSION_SECURE=true
```

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
