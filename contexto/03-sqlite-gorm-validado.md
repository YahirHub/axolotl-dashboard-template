# Fecha
2026-09-28

# Objetivo
Validar que SQLite y Goravel ORM/GORM funcionan realmente en el proyecto, no solo a nivel de compilación.

# Decisiones tomadas
- Se creó el modelo general `AppSetting` como entidad mínima de validación.
- Se creó la migración `app_settings` con clave única y valor de texto.
- Las rutas relativas de `DB_DATABASE` se resuelven desde la raíz del proyecto usando `path.Base`; esto evita fallos cuando los tests cambian el directorio de trabajo.
- El archivo `.env` local se crea desde `.env.example` y permanece ignorado por Git.
- Goravel v1.18 no acepta `migrate --force`; el comando correcto usado localmente es `go run . artisan migrate`.

# Arquitectura actual
SQLite es el driver de persistencia por defecto y el ORM está listo para modelos reales del futuro dashboard.

# Librerías usadas
- github.com/goravel/sqlite v1.18.0
- gorm.io/gorm v1.31.2 (indirecta)
- github.com/ncruces/go-sqlite3 y gormlite a través del driver oficial

# Archivos importantes modificados
- app/models/app_setting.go
- database/migrations/20260928000100_create_app_settings_table.go
- bootstrap/migrations.go
- config/database.go
- tests/feature/database_test.go

# Problemas encontrados
1. SQLite inicialmente no podía abrir `storage/axolotl.sqlite` durante `go test` porque el working directory de la prueba no era la raíz.
2. Artisan requiere `APP_KEY` inicializada.
3. La opción `--force` no existe en el comando migrate de esta versión.

# Soluciones implementadas
- Conversión de DB_DATABASE relativa a ruta absoluta del proyecto.
- Generación local de APP_KEY con `go run . artisan key:generate`.
- Migraciones ejecutadas exitosamente sin `--force`.
- Prueba round-trip ORM: create, query y save.
- `go test ./...` y `go vet ./...` exitosos.

# Pendientes
- Modularizar la referencia HTML.
- Crear DashboardController y rutas limpias.
- Añadir middleware base.
- Cerrar README y smoke test HTTP.

# Próximos pasos
1. Extraer CSS y JS a `public/`.
2. Convertir el markup a vistas `.tmpl`.
3. Mantener design.md como fuente de verdad.
