# Axolotl — Panel multipropósito en Go

**Axolotl** es un panel administrativo multipropósito construido como aplicación modular **100% Go** sobre Goravel. Está pensado como starter para sistemas internos, paneles administrativos, herramientas operativas, backoffices, portales de gestión y proyectos que necesiten autenticación, SQLite, vistas server-side y componentes reutilizables sin depender de un frontend SPA.

La aplicación incluye dashboard, cuentas, métricas de servidor, catálogo de componentes, explorador de archivos seguro, autenticación con setup inicial, recuperación opcional por correo y una arquitectura multipágina preparada para crecer por rutas independientes.

La interfaz original se conserva como referencia visual, pero la aplicación activa utiliza controllers, rutas, middleware, vistas server-side, configuración por `.env`, ORM y SQLite.

## Requisitos

Para ejecutarlo localmente solo necesitas:

- **Go 1.25 o superior**.
- Git únicamente si vas a clonar/versionar el proyecto.
- Un navegador moderno.
- No necesitas instalar SQLite por separado: el driver es embebido.
- Node/npm no son necesarios para arrancar la aplicación; los assets frontend de la demo se sirven directamente y algunas librerías visuales se cargan por CDN.

Comprobar Go:

```bash
go version
```

## Stack

- **Go 1.25+** — probado con Go 1.26.4.
- **Goravel 1.18**.
- **Gin** como driver HTTP oficial de Goravel.
- **html/template** como motor de vistas nativo, organizado en páginas y partials.
- **Goravel ORM / GORM**.
- **SQLite** embebido mediante el driver oficial `github.com/goravel/sqlite`.
- Chart.js, Lenis, AOS, Highlight.js y Font Awesome por CDN para la demo frontend.

## Inicio rápido

> Axolotl **no ejecuta migraciones automáticamente** al arrancar. En una instalación nueva debes ejecutar `go run . artisan migrate` antes de `go run .`.

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

Los tres pasos hacen lo siguiente:

1. crean tu archivo local `.env`;
2. generan `APP_KEY`;
3. crean/actualizan las tablas de SQLite mediante las migraciones del proyecto.

Si ya existe un `.env` válido, no necesitas copiarlo de nuevo.

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

### Reiniciar desde una base limpia

Detén el servidor antes de borrar SQLite.

PowerShell:

```powershell
Remove-Item .\storage\axolotl.sqlite -Force -ErrorAction SilentlyContinue
go run . artisan migrate
go run .
```

Linux/macOS:

```bash
rm -f storage/axolotl.sqlite
go run . artisan migrate
go run .
```

Al crear una base nueva se generará nuevamente el flujo de primer administrador y un PIN temporal distinto.

## Configuración

La configuración sigue el enfoque de Laravel mediante variables de entorno.

```dotenv
APP_NAME=Axolotl
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost:3000
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

MAIL_HOST=
MAIL_PORT=1025
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_ENCRYPTION=auto
MAIL_FROM_ADDRESS=
MAIL_FROM_NAME=Axolotl
PASSWORD_RESET_TTL=30
```

Variables principales:

| Variable | Uso |
| --- | --- |
| `APP_NAME` | Nombre de la aplicación. |
| `APP_ENV` | Entorno, por ejemplo `local` o `production`. |
| `APP_DEBUG` | Activa/desactiva información de depuración. En producción debe ser `false`. |
| `APP_URL` | URL pública base. Se usa, entre otras cosas, para enlaces de recuperación. |
| `APP_HOST` | Dirección donde escucha el servidor HTTP. |
| `APP_PORT` | Puerto HTTP. |
| `DB_CONNECTION` | Driver de base de datos. El starter usa `sqlite`. |
| `DB_DATABASE` | Archivo SQLite, relativo al proyecto o absoluto. |
| `SESSION_DRIVER` | Backend de sesiones. Por defecto `file`. |
| `SESSION_LIFETIME` | Duración de sesión en minutos. |
| `SESSION_SECURE` | Debe ser `true` cuando la aplicación funciona bajo HTTPS. |
| `AUTH_ALLOW_REGISTRATION` | Habilita o deshabilita el registro público de usuarios estándar. |
| `FILES_ROOT` | Carpeta raíz permitida para el explorador de archivos. |
| `FILES_MAX_EDIT_BYTES` | Tamaño máximo de un archivo editable desde el explorador. |
| `MAIL_HOST` | Servidor SMTP. Si queda vacío, recuperación por correo permanece deshabilitada. |
| `MAIL_PORT` | Puerto SMTP. |
| `MAIL_USERNAME` / `MAIL_PASSWORD` | Credenciales SMTP, cuando sean necesarias. |
| `MAIL_ENCRYPTION` | `auto`, `tls`, `starttls` o `none` solo para loopback/desarrollo. |
| `MAIL_FROM_ADDRESS` | Dirección remitente. Debe configurarse para habilitar recovery. |
| `MAIL_FROM_NAME` | Nombre visible del remitente. |
| `PASSWORD_RESET_TTL` | Minutos de vigencia del enlace de recuperación. |

La recuperación de contraseña aparece únicamente cuando SMTP está configurado. En producción, `APP_URL` debe usar HTTPS. `MAIL_ENCRYPTION=auto` usa TLS directo en 465 o STARTTLS cuando el servidor lo ofrece; un SMTP remoto sin cifrado se rechaza. `MAIL_ENCRYPTION=none` solo se admite para servidores loopback de desarrollo.

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
│   ├── password_reset_token.go
│   └── user.go
└── security/
    ├── auth_limiter.go
    ├── password_reset.go
    ├── password_reset_mail.go
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
GET  /register          registro de usuario estándar
POST /register          crea usuario con rol user
GET  /forgot-password   solicita recuperación, solo si SMTP está configurado
POST /forgot-password   genera y envía enlace temporal
GET  /reset-password    formulario de nueva contraseña
POST /reset-password    consume token de un solo uso
POST /logout            cierra sesión

GET  /dashboard     overview general [auth]
GET  /cuenta        perfil, preferencias y seguridad [auth]
GET  /servidor      métricas técnicas [auth]
GET  /componentes   laboratorio UI y catálogo de componentes [auth]
GET  /archivos          raíz del explorador [auth]
GET  /archivos/*path    navegación segura por carpetas [auth]

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

## Comandos útiles

```bash
# arrancar el sitio
go run .

# crear/actualizar tablas
go run . artisan migrate

# ver migraciones
go run . artisan migrate:status

# ver todas las rutas
go run . artisan route:list

# generar APP_KEY
go run . artisan key:generate

# pruebas
go test ./...

# análisis estático
go vet ./...

# compilar
go build ./...
```

## Base de datos y migraciones

Las migraciones son **manuales**. La aplicación no ejecuta `migrate` al iniciar.

Ejecutar migraciones:

```bash
go run . artisan migrate
```

Ver estado:

```bash
go run . artisan migrate:status
```

La migración `users` almacena las cuentas con correo único, contraseña hasheada, rol, estado activo y `auth_version`. El modelo `User` define los roles base `admin` y `user`. La versión de autenticación permite invalidar sesiones antiguas cuando se restablece la contraseña.

`password_reset_tokens` almacena únicamente el hash SHA-256 del token, su expiración y estado de uso. El token en claro solo existe en el enlace enviado al usuario y se entrega en el fragmento `#token=...`, de modo que no se envía en la URL HTTP ni aparece en el access log del servidor.

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
- bloqueos progresivos tras intentos fallidos;
- recuperación por correo con token temporal de un solo uso;
- invalidación de sesiones previas al cambiar contraseña;
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

Como las migraciones son manuales, en una instalación Docker nueva ejecuta primero:

```bash
docker compose build
docker compose run --rm axolotl artisan migrate
docker compose up
```

En ejecuciones posteriores, si no hay migraciones pendientes, basta con:

```bash
docker compose up
```

Docker fuerza `APP_HOST=0.0.0.0` dentro del contenedor y monta `./storage` para persistir SQLite.

La imagen **no incluye** el archivo `.env`; las variables se cargan en runtime.

## Seguridad base

La autenticación incluye:

- guard de sesión de Goravel;
- regeneración de sesión al iniciar y cerrar sesión;
- contraseñas con **Argon2id**;
- CSRF global para formularios POST;
- throttling para setup, registro y recuperación;
- rate limit adaptativo de login por IP real, identidad e IP+identidad, con bloqueos temporales progresivos;
- la IP de seguridad se toma de la conexión TCP y no de cabeceras reenviadas no confiables;
- PIN inicial generado con `crypto/rand` y comparación constant-time;
- PIN guardado solo en memoria;
- rol `admin` asignable únicamente durante el setup inicial;
- registro público configurable con `AUTH_ALLOW_REGISTRATION`;
- rutas privadas protegidas por `RequireAuth`;
- explorador confinado a `FILES_ROOT`, sin acceso arbitrario al filesystem del host;
- rechazo de traversal, rutas ambiguas, nombres reservados, symlinks, junctions y reparse points;
- límites de longitud/profundidad y rechazo de controles invisibles/bidi en nombres;
- edición limitada a texto UTF-8 y al tamaño `FILES_MAX_EDIT_BYTES`;
- mutaciones de archivos protegidas con sesión y CSRF;
- recuperación de contraseña opcional por SMTP, con respuesta anti-enumeración y cooldown por identidad;
- tokens de recuperación aleatorios de 256 bits, almacenados solo como hash, con expiración y uso único;
- el token de recovery viaja en fragmento URL para no aparecer en access logs HTTP;
- reset de contraseña incrementa `auth_version` e invalida sesiones anteriores;
- vistas de autenticación usan `Cache-Control: no-store` y `Referrer-Policy: no-referrer`;
- canal separado `storage/logs/security.log` para eventos de seguridad.

En producción configura HTTPS y cambia:

```dotenv
SESSION_SECURE=true
```

El middleware global agrega:

- `X-Content-Type-Options: nosniff`;
- `X-Frame-Options: SAMEORIGIN`;
- `Referrer-Policy: strict-origin-when-cross-origin`;
- `Permissions-Policy` deshabilitando cámara, micrófono, geolocalización, pagos y USB por defecto;
- `Content-Security-Policy` limitada al propio origen y a los CDN explícitos de la demo;
- `X-Permitted-Cross-Domain-Policies: none`;
- `Origin-Agent-Cluster: ?1`;
- HSTS cuando la petición llega por TLS.

La CSP todavía permite `style-src 'unsafe-inline'` porque la plantilla contiene estilos inline heredados. El siguiente endurecimiento opcional es localizarlos y eliminar esa excepción.

## Desarrollo del proyecto

- Rama principal: `main`.
- Commits en español.
- `contexto/`: decisiones técnicas persistentes.
- `tareas/`: estado de trabajo según metodología Ponytail.
- `LILITH.md`: reglas obligatorias para agentes/herramientas que continúen el proyecto.
- No usar ctx7 salvo último recurso; priorizar documentación oficial y herramientas de búsqueda disponibles.

## Licencia

La base Goravel conserva su licencia MIT en `LICENSE`. Revisa la licencia de los recursos externos utilizados antes de redistribuir una versión comercial cerrada.
