# Tarea 10 — Separar Axolotl en páginas y rutas reales

## Estado
En proceso.

## Objetivo
Corregir la arquitectura para que Axolotl sea una plantilla multipágina real sobre Goravel, no un único documento largo.

## Rutas base
- /dashboard — overview general.
- /cuenta — perfil, preferencias y seguridad.
- /servidor — métricas de CPU, RAM, red y almacenamiento.
- /componentes — laboratorio y catálogo UI.
- / — mismo dashboard de entrada.

## Criterios de aceptación
- El dashboard deja de incluir servidor y laboratorio UI.
- Sidebar navega mediante rutas reales.
- Cada página conserva layout, marca, tema y navegación.
- Topbar usa título/subtítulo según la ruta.
- /cuenta tiene contenido funcional de cuenta/configuración.
- /servidor contiene el monitor técnico.
- /componentes contiene el laboratorio UI.
- Search y menú de usuario navegan correctamente entre páginas.
- route:list, go test, go vet, go build y node --check pasan.
- Smoke HTTP de cada ruta devuelve 200.
- contexto/, design.md y memoria PL actualizados.
