# Tarea 19 — Corregir setup tras eliminar SQLite

## Estado
Completada.

## Resultado
Se identificó que el setup fallaba porque la base nueva no tenía la tabla `users`.

No se implementó auto-migrate. El proyecto exige ejecutar manualmente:

```bash
go run . artisan migrate
go run .
```

El arranque ahora detecta un esquema inexistente y muestra una instrucción clara en log en lugar de dejar que el usuario llegue a un setup roto.
