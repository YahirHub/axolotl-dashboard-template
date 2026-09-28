# Hardening de rutas, autenticación y recuperación

## Estado
Completado el 2026-09-28.

## Explorador
- Las carpetas se navegan mediante rutas reales `/archivos/*path`.
- La route se valida contra el workspace antes de renderizar.
- El sidebar inicia colapsado en Archivos.
- Se bloquearon traversal, rutas absolutas, backslashes, percent-encoding ambiguo, profundidad excesiva, nombres reservados de Windows, controles/bidi, symlinks, junctions y reparse points.
- CSP global y pruebas hostiles quedaron activas.

## Login
- `security.AuthLimiter` aplica bloqueos progresivos por IP, identidad e IP+identidad.
- La dirección para seguridad sale de la conexión TCP, no de cabeceras reenviadas arbitrarias.
- Setup, registro y recovery conservan throttling.
- Eventos de seguridad se escriben en `storage/logs/security.log`.

## Recuperación de contraseña
- Solo se muestra si SMTP está configurado.
- En producción `APP_URL` debe usar HTTPS.
- `POST /forgot-password` siempre responde de forma genérica para evitar enumeración.
- Token: 32 bytes de `crypto/rand` (256 bits), SQLite almacena solo SHA-256, TTL configurable, un solo uso.
- El enlace usa `/reset-password#token=...` para evitar que el secreto llegue al access log HTTP.
- `auth.js` mueve el token del fragmento al formulario y limpia la URL.
- Al resetear se incrementa `users.auth_version`; `RequireAuth` invalida sesiones anteriores.
- Vistas auth usan `Cache-Control: no-store` y `Referrer-Policy: no-referrer`.
- SMTP remoto exige TLS/STARTTLS; modo sin cifrar solo se permite en loopback.

## Migración
- `20260928000300_add_auth_version_and_password_reset_tokens` aplicada a la SQLite local.
- Agrega `users.auth_version` y tabla `password_reset_tokens`.

## Validación
- `go test ./...`
- `go vet ./...`
- `go build ./...`
- `node --check public/js/auth.js`
- fake SMTP feature test
- recovery one-time + expiración + invalidación de sesiones
