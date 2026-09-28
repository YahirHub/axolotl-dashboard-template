# Tarea 11 — Corregir excepción Chart.js con gradientes

## Estado
En proceso.

## Error
`Uncaught TypeError: this._fn is not a function` desde `index.umd.ts`.

## Causa
Chart.js intenta animar propiedades de color cuyo valor es un CanvasGradient generado por los helpers Prism. Ese tipo no tiene interpolador válido dentro del animator.

## Solución
Desactivar únicamente la animación de colores en los charts que usan gradientes scriptables, conservando las animaciones numéricas y las animaciones RAF/CSS del resto de la interfaz.

## Validación
- node --check.
- go test ./...
- go vet ./...
- smoke HTTP.
