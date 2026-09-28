# Axolotl Dashboard Template

Plantilla de dashboard multipropósito con identidad visual Axolotl, en proceso de migración desde una referencia HTML monolítica hacia una aplicación modular 100% Go basada en Goravel.

## Objetivo técnico

- Goravel como framework principal.
- Estructura inspirada en Laravel: rutas, middlewares, controllers, configuración y vistas.
- Vistas server-side con `html/template`.
- Goravel ORM/GORM.
- SQLite como base de datos embebida por defecto.
- Configuración mediante `.env` estilo Laravel.
- Frontend modular conservando el design system de `design.md`.

## Estado

Fase de bootstrap del repositorio.

La referencia original se encuentra en:

```text
docs/reference/index-original.html
```

Las decisiones persistentes del proyecto viven en:

```text
contexto/
tareas/
LILITH.md
```

## Requisitos previstos

- Go >= 1.23.
- Git.

El entorno actual de desarrollo dispone de Go 1.26.4.

## Metodología

El proyecto sigue la metodología Ponytail: arquitectura profesional, simple y segura, evitando dependencias y abstracciones innecesarias.
