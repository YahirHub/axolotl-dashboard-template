# Fecha
2026-09-28

# Objetivo
Añadir a Axolotl un sistema de autenticación reutilizable para aplicaciones multipropósito, incluyendo bootstrap seguro del primer administrador, login, registro de usuarios estándar, sesiones y protección de rutas.

# Flujo de primer inicio
1. La aplicación arranca y consulta si existe un usuario activo con rol `admin`.
2. Si no existe, genera un PIN de 6 dígitos mediante `crypto/rand`.
3. El PIN se guarda únicamente en memoria.
4. El canal de log `setup` imprime el PIN en consola y lo guarda en `storage/logs/setup.log`.
5. `GET /` redirige a `/setup`.
6. `/setup` solicita nombre, correo, alias, contraseña, confirmación y 6 inputs OTP separados.
7. El backend concatena el OTP y lo compara usando `subtle.ConstantTimeCompare`.
8. Solo el backend asigna `role=admin`.
9. El usuario se crea con contraseña Argon2id y sesión iniciada automáticamente.
10. El PIN queda invalidado.
11. Una vez que existe admin, `/setup` redirige a login/dashboard y no permite crear otro administrador.

Si el administrador desaparece posteriormente y el proceso continúa vivo, visitar `/setup` asegura que vuelva a existir un PIN pendiente sin regenerarlo en cada recarga.

# Login
- Ruta: GET/POST `/login`.
- Email normalizado a minúsculas.
- Solo usuarios activos.
- Respuesta genérica ante credenciales incorrectas.
- Password máximo 128 bytes en login.
- Cuando el email no existe se ejecuta una operación Argon2id dummy para reducir diferencias de timing que faciliten enumeración.
- Si el hash requiere actualización, se rehashea después de un login válido.
- El guard de sesión de Goravel regenera ID al iniciar sesión.

# Registro
- Ruta: GET/POST `/register`.
- Disponible después de existir el primer administrador.
- Controlado por `AUTH_ALLOW_REGISTRATION`.
- El cliente no puede elegir rol.
- Incluso si envía `role=admin`, el backend fuerza `role=user`.
- Email único.
- Password Argon2id.
- Sesión automática después de registro.

# Logout
- Ruta POST `/logout`.
- Protegida por CSRF.
- Usa el guard de sesión de Goravel.

# Seguridad
- Guard: session.
- Hash: Argon2id.
- CSRF global: `VerifyCsrfToken`.
- Session middleware global.
- Rate limits:
  - setup: 5/min/IP
  - register: 5/min/IP
  - login: 10/min/IP
- Rutas privadas con `RequireAuth`.
- Rol admin nunca deriva de input del navegador.
- PIN no se persiste en DB.
- PIN temporal no se muestra en HTML.
- En producción: HTTPS + `SESSION_SECURE=true`.

# Rutas
```text
GET  /              -> decide setup/login/dashboard
GET  /setup         -> primer admin
POST /setup         -> crea admin con PIN
GET  /login
POST /login
GET  /register
POST /register
POST /logout

GET  /dashboard     [auth]
GET  /cuenta        [auth]
GET  /servidor      [auth]
GET  /componentes   [auth]

GET  /health
/assets/*
```

# Persistencia
Tabla `users`:
- id
- name
- email unique
- alias
- password
- role
- active
- timestamps

Roles base:
- `admin`
- `user`

# UI
Archivos:
- `resources/views/auth.tmpl`
- `public/css/auth.css`
- `public/js/auth.js`

El OTP usa seis cajas, autoavance, retroceso, flechas y pegado de 6 dígitos.
Login/setup/registro cargan una UI ligera separada del bundle visual del dashboard.

# Integración con la UI privada
`PageMeta` entrega:
- CSRF token
- CurrentUser
- PageTitle/PageSubtitle
- ActiveRoute

Topbar muestra identidad real de sesión.
El logout del topbar es un formulario POST real con CSRF.
La página Cuenta muestra identidad proveniente de SQLite/sesión y conserva preferencias visuales locales como demo.

# Validación realizada
- Migración `users` aplicada correctamente.
- Test: primer inicio redirige a setup.
- Test: PIN incorrecto no crea usuario.
- Test: PIN correcto crea admin e inicia sesión.
- Test: password almacenado no coincide con texto plano y valida mediante Argon2id.
- Test: login y logout reales con sesión.
- Test: registro crea exclusivamente rol `user`.
- Test: intento de enviar `role=admin` desde registro se ignora.
- Test: páginas privadas redirigen a login.
- Test: páginas privadas renderizan al tener sesión.
- go test ./... exitoso.
- go vet ./... exitoso.
- go build ./... exitoso.
- node --check en app.js y auth.js exitoso.

# Archivos principales
- app/models/user.go
- app/security/setup.go
- app/http/controllers/auth_controller.go
- app/http/middleware/require_auth.go
- database/migrations/20260928000200_create_users_table.go
- config/auth.go
- config/hashing.go
- config/session.go
- config/logging.go
- bootstrap/app.go
- routes/web.go
- resources/views/auth.tmpl
- public/css/auth.css
- public/js/auth.js
- tests/feature/auth_flow_test.go
- tests/feature/http_test.go

# Ajustes finales de integración
- Goravel/Gin convierte structs de vista a un mapa superficial; los campos de structs embebidos no se aplanan automáticamente.
- Por ello `DashboardPageData` expone explícitamente `Title`, `PageTitle`, `PageSubtitle`, `ActiveRoute`, `CSRFToken` y `CurrentUser`.
- Esto corrigió el caso donde el dashboard autenticado renderizaba título/usuario vacíos aunque la sesión fuera válida.
- `storage/framework/*` está ignorado por Git; solo se conserva `storage/framework/.gitkeep`.
- El servidor de desarrollo fue reiniciado con la implementación final en `http://127.0.0.1:3000`.
- En la base local actual ya existe al menos un administrador, por lo que `/` redirige a `/login`; en una base nueva el mismo código redirige a `/setup`.

# Próximos pasos opcionales
- recuperación/restablecimiento de contraseña;
- verificación de correo;
- administración de usuarios/roles;
- 2FA real;
- permisos por capacidad/policy;
- persistencia backend de preferencias de perfil.
