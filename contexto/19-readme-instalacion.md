# README de instalación y configuración

## Fecha
2026-09-29

## Objetivo
Convertir README.md en la guía principal para instalar, configurar y arrancar Axolotl.

## Resultado
- Título y descripción definen Axolotl como panel multipropósito 100% Go.
- Se documentan requisitos.
- Se explica instalación inicial en Windows y Linux/macOS.
- Se deja explícito que NO existe auto-migrate.
- Se documenta el orden correcto: .env -> key:generate -> migrate -> go run.
- Se documenta cómo reiniciar con una SQLite limpia.
- Se describen las variables de entorno principales.
- Se agregan comandos frecuentes.
- Docker también usa migraciones manuales antes del primer arranque.
- Se mantienen setup del primer admin, SMTP/recovery, rutas, seguridad y pruebas.

## Regla importante
No introducir auto-migrate en la documentación ni en el código salvo que el usuario cambie explícitamente esa decisión.
