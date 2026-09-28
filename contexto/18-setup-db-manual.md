# Setup de base de datos manual

## Fecha
2026-09-28

## Problema
Después de eliminar SQLite, el setup mostraba "No fue posible completar el setup" porque la tabla `users` no existía.

## Causa
La aplicación arrancaba sin comprobar de forma explícita que las migraciones hubieran sido ejecutadas.

## Decisión
No habrá auto-migrate al iniciar la aplicación.

El flujo correcto para una base nueva es manual:

```bash
go run . artisan migrate
go run .
```

Si el esquema no existe, el arranque falla de forma explícita y el log indica que deben ejecutarse las migraciones.

## Cambios
- `InitializeAdminPIN()` verifica que exista la tabla `users`.
- `main.go` muestra en log el comando `go run . artisan migrate` si la base no está preparada.
- `StoreSetup` registra el error real en el canal de seguridad.
- No existe `DB_AUTO_MIGRATE` ni llamada automática a `artisan migrate`.

## Validación
- `go test ./...` OK.
- `go vet ./...` OK.
- `go build ./...` OK.
