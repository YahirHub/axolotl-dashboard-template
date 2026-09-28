# Fecha
2026-09-28

# Objetivo
Establecer la base modular 100% Go de Axolotl usando Goravel y dejar SQLite como driver de base de datos por defecto.

# Decisiones tomadas
- Se adopta Goravel v1.18.0.
- El módulo Go es `github.com/YahirHub/axolotl-dashboard-template`.
- HTTP usa el driver oficial `github.com/goravel/gin` v1.18.0.
- Base de datos usa el driver oficial `github.com/goravel/sqlite` v1.18.0.
- Se eliminó del scaffold la configuración OpenAI y el driver Postgres porque no forman parte del objetivo actual.
- SQLite se selecciona con `DB_CONNECTION=sqlite`.
- Archivo de DB por defecto: `storage/axolotl.sqlite`.
- El pool SQLite se mantiene en una sola conexión abierta/idle para evitar complejidad innecesaria en el starter.
- Se preserva la estructura completa de Goravel para disponer de rutas, middleware, auth, colas, eventos, migraciones, mail, testing y vistas sin reinventar framework.
- El motor de vistas será el `html/template` oficial de Goravel; la modularización del dashboard se hará con layouts/partials.

# Arquitectura actual
- `app/`: controllers y facades.
- `bootstrap/`: arranque y providers.
- `config/`: configuración centralizada.
- `database/migrations/`: migraciones.
- `public/`: assets públicos.
- `resources/views/`: vistas Goravel.
- `routes/`: rutas.
- `tests/`: pruebas.
- `docs/reference/index-original.html`: referencia visual completa.

# Librerías usadas
- github.com/goravel/framework v1.18.0
- github.com/goravel/gin v1.18.0
- github.com/goravel/sqlite v1.18.0
- github.com/stretchr/testify v1.11.1
- GORM v1.31.2 indirectamente a través de Goravel/SQLite.

El driver SQLite oficial actual trae dependencias puras en Go basadas en `github.com/ncruces/go-sqlite3` / gormlite, evitando depender del binario externo de SQLite.

# Archivos importantes modificados
- go.mod
- go.sum
- bootstrap/providers.go
- config/database.go
- .env.example
- estructura base Goravel

# Problemas encontrados
El scaffold completo de Goravel incluye OpenAI y Postgres por defecto. Mantenerlos agregaba dependencias innecesarias.

# Soluciones implementadas
- Eliminación de `config/ai.go`.
- Eliminación de providers OpenAI/Postgres.
- Registro del provider SQLite.
- Actualización de imports al módulo YahirHub.
- `go mod tidy`.
- `gofmt`.
- `go test ./...` exitoso.

# Pendientes
- Crear migración/modelo de validación SQLite.
- Probar migrate + CRUD ORM.
- Modularizar frontend.
- Crear controller dashboard y middleware de aplicación.

# Próximos pasos
1. Validar SQLite con migración/modelo real.
2. Extraer CSS/JS del HTML de referencia.
3. Convertir markup en vistas/partials.
4. Reemplazar ruta welcome por DashboardController.
