# Tarea 12 — Setup inicial, login y registro

## Estado
Completada.

## Objetivo
Implementar autenticación real basada en sesión para Axolotl, con bootstrap seguro del primer administrador mediante PIN temporal de 6 dígitos generado al iniciar la aplicación.

## Implementado
- Detección de administrador activo al arrancar.
- PIN temporal criptográfico de 6 dígitos, solo en memoria.
- Canal de log `setup` con impresión en terminal y archivo `storage/logs/setup.log`.
- `/setup` con OTP de seis cuadros, validación backend y creación exclusiva del primer `admin`.
- Login con sesiones Goravel y Argon2id.
- Registro público configurable que solo crea rol `user`.
- Logout real por POST.
- CSRF global.
- Rate limiting para setup/login/registro.
- Middleware `RequireAuth` para páginas privadas.
- Usuario real conectado al topbar y página Cuenta.
- Sesiones runtime excluidas de Git.

## Validación
- `go test ./...` OK.
- `go vet ./...` OK.
- `go build ./...` OK.
- `node --check public/js/app.js` OK.
- `node --check public/js/auth.js` OK.
- Prueba de PIN incorrecto: no crea usuario.
- Prueba de PIN correcto: crea admin e inicia sesión.
- Prueba de registro: ignora `role=admin` enviado por cliente y crea `user`.
- Prueba de login/logout con sesión real.
- Pruebas de protección y render de rutas privadas.

## Resultado
La autenticación base queda lista como módulo reutilizable del template. Recuperación de contraseña, verificación de correo, 2FA y administración de usuarios quedan como extensiones posteriores, no como requisito de esta tarea.
