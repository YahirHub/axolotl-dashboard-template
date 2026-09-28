# Tarea 09 — Convertir inicio en dashboard multipropósito completo

## Estado
En proceso.

## Objetivo
Eliminar la repetición de marca del bloque inferior del sidebar y reorganizar la página principal en módulos reutilizables, variados y alimentados desde el backend Goravel.

## Criterios de aceptación
- El SVG solo aparece como marca principal en la parte superior del sidebar.
- El contenido principal deja de estar centrado en cripto/mercados.
- KPIs, proyectos, tareas, actividad, equipo, agenda e integraciones muestran datos de demostración distintos.
- Los datos repetibles principales se entregan desde DashboardController a html/template.
- Se conservan analítica, métricas de servidor y laboratorio UI como ejemplos reutilizables.
- Sidebar y búsqueda apuntan a secciones genéricas.
- Responsive y acciones existentes siguen funcionando.
- go test ./..., go vet ./... y smoke HTTP pasan.
- contexto/, design.md y memoria PL quedan actualizados.
