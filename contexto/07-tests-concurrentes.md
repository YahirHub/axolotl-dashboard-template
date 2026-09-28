# Fecha
2026-09-28

# Objetivo
Evitar que dos ejecuciones concurrentes de la suite compartan el mismo archivo SQLite temporal.

# Cambio
`tests/test_case.go` genera ahora la base temporal usando el PID del proceso:

```go
fmt.Sprintf("axolotl-dashboard-template-test-%d.sqlite", os.Getpid())
```

# Motivo
El archivo fijo anterior era suficiente para una ejecución aislada, pero podía colisionar si dos procesos de tests corrían en paralelo en el mismo host.

# Validación
- gofmt exitoso.
- go test ./... exitoso.
- go vet ./... exitoso.
- validación completa PL exitosa.

# Estado
Cambio pequeño y autocontenido. No modifica producción ni la base SQLite de desarrollo.
