# Axolotl Design System

> Documento de referencia para mantener consistencia visual, funcional y responsive en cualquier módulo nuevo de la plantilla Axolotl.

## 1. Objetivo

Axolotl es una plantilla multipropósito de dashboard oscuro. El diseño debe sentirse técnico, limpio y compacto, con una identidad Neo/Prism moderada: superficies oscuras, contraste alto, acentos turquesa/rosa y gradientes luminosos reservados principalmente para gráficas, estados activos y elementos interactivos.

La regla principal es **no crear estilos aislados por componente si existe un token semántico o utilidad reutilizable**.

## 2. Dependencias

La plantilla actual usa:

- Font Awesome por CDN para iconografía.
- Chart.js por CDN para gráficas.
- Lenis por CDN para smooth scrolling.
- AOS por CDN para animaciones de entrada.
- Highlight.js por CDN para detección automática de lenguaje y syntax highlighting.
- HTML, CSS y JavaScript vanilla.
- No depende de Bootstrap; sin embargo, expone una API de utilidades semánticas inspirada en Bootstrap.

## 3. Tokens globales

Los colores base viven en `:root`.

```css
:root {
  --primary: #FC779B;
  --secondary: #FCBBD5;
  --success: #8EDBC7;
  --danger: #FF5F7F;
  --warning: #F5B86F;
  --info: #D49AF6;

  --body-bg: #09070B;
  --sidebar-bg: #0D0910;
  --surface: #151018;
  --surface-2: #1B131E;
  --border: #352330;
  --border-soft: #281B25;
  --body-color: #FFF2F7;
  --muted: #B58C9B;
}
```

También existen equivalentes RGB (`--primary-rgb`, `--secondary-rgb`, etc.) para transparencias con `rgba()`.

### Compatibilidad

La plantilla mantiene aliases internos para componentes antiguos:

```css
--teal: var(--primary);
--pink: var(--secondary);
--green: var(--success);
--red: var(--danger);
--orange: var(--warning);
```

Para componentes nuevos se deben usar directamente los nombres semánticos: `--primary`, `--success`, etc.

## 4. Utilidades estilo Bootstrap

### Fondos

```html
<div class="bg-primary">...</div>
<div class="bg-primary-soft">...</div>
<div class="bg-secondary">...</div>
<div class="bg-success">...</div>
<div class="bg-danger">...</div>
<div class="bg-warning">...</div>
<div class="bg-info">...</div>
<div class="bg-surface">...</div>
```

### Texto

```html
<span class="text-primary">Primary</span>
<span class="text-secondary">Secondary</span>
<span class="text-success">Success</span>
<span class="text-danger">Danger</span>
<span class="text-warning">Warning</span>
<span class="text-info">Info</span>
<span class="text-muted">Texto secundario</span>
```

### Bordes

```html
<div class="border-primary">...</div>
<div class="border-danger">...</div>
```

### Botones

Todos los botones semánticos deben incluir `btn-ui`.

```html
<button class="btn-ui btn-primary">Guardar</button>
<button class="btn-ui btn-secondary">Continuar</button>
<button class="btn-ui btn-success">Aceptar</button>
<button class="btn-ui btn-danger">Eliminar</button>
<button class="btn-ui btn-warning">Revisar</button>
<button class="btn-ui btn-outline-primary">Más información</button>
<button class="btn-ui btn-ghost">Cancelar</button>
```

### Espaciado y layout

La plantilla ofrece utilidades básicas:

- `p-1`, `p-2`, `p-3`, `p-4`
- `m-0`
- `gap-1`, `gap-2`, `gap-3`
- `d-flex`
- `align-center`
- `justify-between`
- `w-100`
- `rounded`, `rounded-sm`, `rounded-lg`
- `shadow-ui`

No se pretende reemplazar un framework completo. Si una pantalla requiere una utilidad nueva, añadirla aquí y documentarla.

## 5. Superficies y tarjetas

El componente base es `.card`.

```html
<article class="card section">
  <div class="card-head">
    <div class="card-title">Título</div>
  </div>
  ...
</article>
```

Reglas:

1. Fondo oscuro.
2. Borde sutil.
3. Radio global `--radius`.
4. Sombras suaves.
5. No añadir glows fuertes a tarjetas normales.
6. El efecto Neo/Prism se reserva para datos visuales e interacción.

## 6. Sidebar

El sidebar tiene dos estados:

- Expandido: icono + texto.
- Colapsado: iconos centrados.

En móvil funciona como drawer off-canvas.

El icono superior es a la vez identidad visual y control de expandir/colapsar. No colocar texto de marca junto a él salvo que el diseño general vuelva a requerirlo.

Todo icono del sidebar usa una caja de alineación uniforme para evitar desplazamientos visuales.

## 7. Header

Elementos:

1. Título/contexto.
2. Selector de moneda de demostración.
3. Buscador global.
4. Notificaciones.
5. User chip rectangular.

El user chip muestra:

- Alias/avatar.
- Nombre de usuario.
- Rol.

El dropdown de usuario **no debe repetir la misma tarjeta de identidad**. Debe contener solamente acciones: perfil, preferencias, tema, componentes, seguridad y cierre de sesión.

## 8. Buscador global

El buscador funciona como command palette multipropósito.

Atajos:

- `Ctrl + K` / `Cmd + K`: abrir/cerrar.
- `/`: abrir cuando no se está escribiendo.
- `↑` / `↓`: navegar.
- `Enter`: ejecutar.
- `Esc`: cerrar.

Categorías:

- Secciones.
- Datos.
- Acciones.
- Novedades.

Los resultados deben ejecutar una acción real o una simulación explícita. No añadir controles sin función.

## 9. Gráficas

### Líneas y sparklines

Usan lenguaje Neo/Prism:

- Gradiente `primary → secondary → info`.
- Glow ligero.
- Área degradada suave.
- Tooltip interactivo.
- Punto de hover.
- Línea vertical de seguimiento cuando corresponde.

No sobrecargar el glow.

### Donut / pie

El donut de distribución **permanece plano por fuera**: no debe tener `drop-shadow`, glow externo ni un halo sobre los segmentos.

El único efecto Neo permitido en este componente va **dentro del hueco central** mediante un `radial-gradient` y `box-shadow: inset`. Así el brillo queda contenido dentro de la dona y no invade el borde exterior.

La geometría del doughnut es estrictamente **1:1**. Si necesita más margen, se reduce su tamaño completo; nunca se modifica ancho y alto de forma independiente.

Usa colores sólidos semánticos:

- Warning
- Azul de serie
- Primary
- Secondary
- Muted strong

Esto crea contraste visual frente a las líneas Neo/Prism.

Regla de implementación:

```css
.donut-wrap {
  --donut-size: 160px;
  width: min(var(--donut-size), 100%);
  aspect-ratio: 1 / 1;
}

.donut {
  position: absolute !important;
  inset: 0 !important;
  width: 100% !important;
  height: 100% !important;
  filter: none !important;
}

.donut-center::before {
  /* brillo únicamente dentro del hueco */
  aspect-ratio: 1 / 1;
  box-shadow:
    inset 0 0 10px rgba(var(--primary-rgb), .18),
    inset 0 0 21px rgba(var(--primary-rgb), .09);
}
```

En Chart.js, el doughnut debe usar:

```js
maintainAspectRatio: true,
aspectRatio: 1
```

Si hay que ganar margen, cambiar `--donut-size`; nunca estirar el canvas.

El plugin de glow de Chart.js solo debe actuar en datasets que definan explícitamente `prismGlowColor`. Un doughnut no debe definir esa propiedad.

### Regla de color

Para gráficas nuevas, consumir tokens mediante JavaScript:

```js
cssVar('--primary')
cssVar('--secondary')
cssVar('--success')
```

No fijar colores hex si el componente debe responder al cambiador global de tema.

## 10. Sistema de modales

La plantilla expone dos APIs.

### Modal compatible existente

```js
openModal(
  'Título',
  'Subtítulo',
  '<div>Contenido</div>',
  false
);
```

### Modal de componente recomendado

```js
openUiModal({
  title: 'Nuevo elemento',
  subtitle: 'Descripción breve',
  icon: 'fa-plus',
  size: 'sm',
  body: '<div class="modal-callout">Contenido</div>',
  footer: `
    <button class="btn-ui btn-ghost">Cancelar</button>
    <button class="btn-ui btn-primary">Guardar</button>
  `
});
```

Tamaños válidos:

- `sm`
- `md` implícito
- `lg`
- `xl`

Reglas del modal:

- Backdrop oscuro con blur.
- Barra superior Prism discreta.
- Superficie `#07171d` compatible con el dashboard.
- Botones semánticos.
- Footer separado por borde.
- Nunca introducir un modal blanco.

## 11. Bloques de código

Todo `pre > code` se mejora automáticamente.

Ejemplo:

```html
<pre><code>&lt;button class="btn-ui btn-primary"&gt;Guardar&lt;/button&gt;</code></pre>
```

La función `enhanceCodeBlocks()` añade:

- Contenedor oscuro.
- Botón Copiar.
- Fallback para navegadores sin Clipboard API.
- Cambio temporal de icono a ✓ y texto `Copiado`.

Cuando un modal inserta código, `openModal()` ejecuta automáticamente el enhancer.

## 12. Temas globales

El menú de usuario incluye **Personalizar tema**.

Se pueden cambiar en ejecución:

- Primary
- Secondary
- Success
- Danger
- Warning
- Info

El tema se guarda en `localStorage` bajo:

```text
axolotl-theme
```

Los componentes que usan tokens se actualizan automáticamente. Las gráficas Neo/Prism también vuelven a renderizarse.

### Regla para nuevos componentes

Correcto:

```css
.my-component {
  color: var(--primary);
  border-color: rgba(var(--primary-rgb), .2);
}
```

Incorrecto:

```css
.my-component {
  color: #FC779B;
}
```

El segundo caso impide que el cambiador de tema funcione.

## 13. Tipografía

Fuente base:

```css
Inter, Segoe UI, Arial, sans-serif
```

Jerarquía aproximada:

- Page title: 18–22 px.
- Card title: 13 px.
- Métrica principal: 17–21 px.
- Texto normal de dashboard: 10–12 px.
- Metadata: 8–10 px.

La interfaz es deliberadamente compacta.

## 14. Radios

Tokens:

- `--radius-sm: 9px`
- `--radius: 14px`
- `--radius-lg: 18px`

No mezclar radios arbitrarios salvo componentes especiales como pills (`999px`) o avatares.

## 15. Estados

Usar semántica consistente:

- Primary: selección, foco y acción principal.
- Secondary: identidad visual complementaria.
- Success: estados positivos y confirmaciones.
- Danger: errores, negativos y acciones destructivas.
- Warning: advertencias y series específicas.
- Info: información o series complementarias.

## 16. Tooltips, toasts y notificaciones

### Tooltip

Para ayuda contextual pequeña.

### Toast

Para confirmar acciones breves:

```js
toast('Perfil actualizado', 'Los cambios se aplicaron correctamente.');
```

No usar toast para información que requiera una decisión.

### Drawer de notificaciones

Se usa para eventos persistentes, alertas y actividad reciente.

## 17. Responsive

Breakpoints principales:

- `> 1400px`: escritorio amplio.
- `1180px`: sidebar compacto y reorganización de contenido.
- `780px`: sidebar off-canvas, layout móvil.
- `520px`: tarjetas y grids a una columna.

Todo componente nuevo debe funcionar sin scroll horizontal de página.

## 18. Interactividad obligatoria

La plantilla tiene una regla estricta: **si un elemento parece clicable, debe hacer algo**.

No crear:

- “Ver detalles” sin destino.
- Botones decorativos.
- Iconos de dropdown que no abren nada.
- Filtros sin efecto.

Si una función todavía no existe, ocultar o deshabilitar el control con una explicación.

## 19. Accesibilidad

Requisitos:

- `button` para acciones, no `div` clicables.
- `aria-label` en botones con solo icono.
- Estados `aria-expanded` para dropdowns.
- Foco visible.
- Escape cierra modales/drawers.
- Contraste suficiente.
- Componentes usables con teclado.

## 20. Cómo añadir una nueva sección

Ejemplo recomendado:

```html
<article class="card section" id="my-section">
  <div class="card-head">
    <div class="card-title">Nueva sección</div>
    <button class="btn-ui btn-outline-primary" type="button">
      Acción
    </button>
  </div>

  <div class="bg-primary-soft border-primary rounded p-3">
    Contenido
  </div>
</article>
```

Para mantener persistencia visual:

1. Usar `card section`.
2. Consumir tokens.
3. Reutilizar botones semánticos.
4. Mantener tipografía compacta.
5. Probar en los cuatro breakpoints.
6. Añadir la función al buscador global si es importante.
7. Documentar un nuevo patrón en este archivo si introduce un componente reusable.

## 21. Filosofía visual

Axolotl debe sentirse:

- Oscuro.
- Técnico.
- Moderno.
- Compacto.
- Multipropósito.
- Neo/Prism en visualización de datos.
- Sobrio en estructura y navegación.

Evitar:

- Tarjetas blancas.
- Glows excesivos.
- Gradientes en todos los componentes.
- Colores hardcoded cuando existe token.
- Controles sin función.
- Inconsistencias de padding y radios.

---

**Fuente de verdad del diseño:** este archivo `design.md`.

Cuando se agregue un patrón visual o componente reusable nuevo, actualizar este documento en el mismo cambio.


## 22. Motion / animaciones

La plantilla usa dos librerías por CDN:

- **Lenis** para desplazamiento suave.
- **AOS** para entradas discretas al hacer scroll.

CDN actual de la plantilla:

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/aos@2.3.4/dist/aos.css">
<script src="https://cdn.jsdelivr.net/npm/lenis@1.3.11/dist/lenis.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/aos@2.3.4/dist/aos.js"></script>
```

### Scroll suave

El botón "volver arriba" usa:

```js
smoothScrollTo(0, { duration: 1.15 });
```

Para desplazarse a una sección:

```js
smoothScrollTo(document.querySelector('#mi-seccion'), {
  offset: -10,
  duration: 1
});
```

Si Lenis no carga, existe fallback a `scrollIntoView()` / `window.scrollTo()` nativos.

### Animaciones AOS

Las tarjetas principales se inicializan automáticamente con `fade-up`.

Las animaciones deben:

- Ser cortas.
- No mover agresivamente la interfaz.
- No bloquear interacción.
- Ejecutarse una sola vez por defecto.
- Respetar `prefers-reduced-motion`.

No usar animación únicamente como decoración en cada elemento pequeño.


### Compatibilidad con `file://`

Los navegadores tratan cada URL `file://` como un origen de seguridad especial. Por ello, algunas librerías, requests y recursos locales pueden producir mensajes como:

```text
'file:' URLs are treated as unique security origins
```

Axolotl no depende de Lenis para el botón "volver arriba". Si la página se abre con `file://`, `smoothScrollTo()` utiliza automáticamente un motor interno basado en `requestAnimationFrame`, de modo que el desplazamiento suave sigue funcionando sin CDN.

Lenis se activa automáticamente en `http://` o `https://` cuando su CDN está disponible.

Para desarrollo se recomienda servir la carpeta por HTTP local:

```bash
python -m http.server 8000
```

y abrir:

```text
http://localhost:8000
```

Esto evita las restricciones especiales de `file://` y se acerca más al entorno real de producción.



### Nota técnica: RAF + `scroll-behavior`

No mezclar una animación manual con `requestAnimationFrame` y `scroll-behavior: smooth` en cada frame. Si el documento conserva `scroll-behavior: smooth`, cada llamada a `window.scrollTo()` puede iniciar su propia interpolación nativa y producir un salto aparente.

El motor interno de Axolotl hace lo siguiente durante la animación:

```js
const root = document.documentElement;
const previous = root.style.scrollBehavior;

root.style.scrollBehavior = 'auto';

// animación requestAnimationFrame...
window.scrollTo(0, y);

// al terminar:
root.style.scrollBehavior = previous;
```

La interpolación usa `easeInOutCubic`, que acelera y desacelera de forma visible, en lugar de un `easeOutExpo` demasiado agresivo al comienzo.

El botón de volver arriba usa aproximadamente 1.3 segundos para que el movimiento sea perceptible incluso en páginas largas.


## 23. Skeletons de carga

Los datos que dependan de API, base de datos, archivos o imágenes deben mostrar un estado de carga consistente.

### Hosts de datos

Ejemplo:

```html
<article id="clientes" class="card" data-skeleton="summary">
  ...
</article>
```

Tipos incorporados:

- `stats`
- `summary`
- `market`
- `news`

API:

```js
setLoading('#clientes', true, 'summary');

// después de recibir la información:
setLoading('#clientes', false, 'summary');
```

Para promesas/requests:

```js
const result = await withSkeleton(
  '#clientes',
  'summary',
  async () => {
    return await cargarClientes();
  }
);
```

`withSkeleton()` siempre limpia el loading en `finally`, incluso si la petición falla.

### Skeletons de imagen

Estructura recomendada:

```html
<div class="skeleton-image-wrap">
  <div class="image-skeleton-layer skeleton"></div>
  <img
    data-skeleton-image
    src="imagen.jpg"
    alt="Descripción"
  >
</div>
```

Inicializar:

```js
initImageSkeletons();
```

El skeleton desaparece tanto en `load` como en `error`, evitando placeholders infinitos.

### Shimmer

El shimmer usa tokens de `primary` y `secondary`, por lo que responde al editor global de tema.

Con `prefers-reduced-motion: reduce`, la animación shimmer se deshabilita.

### Regla de layout

Un skeleton debe ocupar aproximadamente el mismo espacio que el contenido final. Esto evita layout shift al hidratar datos.

## 24. Sidebar con subcategorías

El sidebar soporta grupos expandibles.

```html
<div class="nav-group">
  <button class="nav-item nav-group-toggle" aria-expanded="false">
    <i class="fa-solid fa-chart-line"></i>
    <span class="nav-label">Analítica</span>
    <i class="fa-solid fa-chevron-right nav-chevron"></i>
  </button>

  <div class="nav-submenu">
    <button class="nav-subitem" data-target="market-section">
      <i class="fa-solid fa-table-list"></i>
      <span>Resumen general</span>
    </button>
  </div>
</div>
```

Reglas:

- El chevron rota al abrir el grupo.
- El padre se marca activo cuando una subcategoría está activa.
- Cada subcategoría lleva icono.
- Un padre expandible no cierra el drawer móvil; solo abre/cierra sus hijos.
- Una subcategoría ejecuta la navegación/acción y cierra el drawer en móvil.
- Si el sidebar de escritorio está colapsado y se pulsa un grupo, primero se expande y después se abre.
- En modo colapsado no se muestran submenús.

### Ventanas compactas

Entre 781 y 1180 px el sidebar ya no se fuerza a icon-only. Conserva una versión estrecha de 174–188 px con texto visible. El usuario todavía puede colapsarlo manualmente.

Por debajo de 780 px el sidebar se convierte en drawer off-canvas:

- ancho `clamp(238px, 76vw, 272px)`;
- navegación con texto completo;
- subcategorías habilitadas;
- `Axolotl AI` usa una tarjeta compacta y ya no queda un gran espacio vacío;
- backdrop para cerrar al tocar fuera;
- scroll interno solo cuando la altura lo requiere.

## 25. Laboratorio UI visible

La plantilla incluye una sección `#code-lab-section` para probar visualmente el design system.

Incluye:

- botones semánticos;
- modal reusable;
- modal de confirmación;
- skeleton de carga;
- toasts;
- pestañas HTML/CSS/JavaScript;
- bloques de código copiables.

### Copiar código

Todo bloque:

```html
<pre><code>...</code></pre>
```

es procesado por:

```js
enhanceCodeBlocks(document);
```

y recibe automáticamente un botón:

```text
Copiar
```

Después de copiar cambia temporalmente a:

```text
✓ Copiado
```

No es necesario escribir el botón manualmente.

### Añadir ejemplos nuevos

```html
<article class="demo-card">
  <div class="demo-card-head">
    <div>
      <strong>Mi componente</strong>
      <small>Descripción.</small>
    </div>
  </div>

  <div class="demo-preview">
    <!-- componente funcional -->
  </div>

  <pre><code>&lt;button class="btn-ui btn-primary"&gt;Demo&lt;/button&gt;</code></pre>
</article>
```

Si el ejemplo necesita JavaScript, inicializarlo dentro de `initCodeLab()`.

La finalidad del laboratorio es que un desarrollador pueda ver el componente real, interactuar con él y copiar la implementación sin salir del dashboard.

## 26. Syntax highlighting y detección automática

La plantilla usa **Highlight.js** mediante CDN.

```html
<link
  rel="stylesheet"
  href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/styles/atom-one-dark.min.css"
>
<script
  src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/highlight.min.js">
</script>
```

Todo bloque:

```html
<pre><code>...</code></pre>
```

es procesado por `enhanceCodeBlocks()`.

El sistema:

1. conserva el código original en `data-raw-code`;
2. intenta detectar automáticamente el lenguaje con Highlight.js;
3. muestra una etiqueta visual con el lenguaje detectado;
4. colorea la sintaxis;
5. agrega el botón **Copiar**;
6. al copiar cambia a **✓ Copiado**;
7. copia el código original, no el HTML coloreado por Highlight.js.



### Importante: build de navegador

No usar este archivo del paquete npm directamente en el navegador:

```text
highlight.js/lib/common.min.js
```

Ese archivo pertenece al build CommonJS/Node y puede producir:

```text
Uncaught ReferenceError: require is not defined
```

En una página HTML por CDN se debe usar el bundle preparado para navegador:

```html
<script
  src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/highlight.min.js">
</script>
```

Este bundle expone `window.hljs` y no necesita `require()`.

### Forzar un lenguaje

Para evitar ambigüedad:

```html
<pre><code data-language="javascript">
const app = 'Axolotl';
</code></pre>
```

Lenguajes comunes disponibles en el build CDN:

- HTML/XML
- CSS
- JavaScript
- TypeScript
- JSON
- Bash/Shell
- Python
- SQL
- Markdown
- PHP
- Java
- C/C++

### API manual

```js
const language = highlightCodeElement(
  document.querySelector('#miCodigo'),
  'javascript'
);
```

Para actualizar la etiqueta:

```js
updateCodeLanguageBadge(shell, language);
```

### Fallback

Si el CDN no carga, Axolotl usa `fallbackLanguageDetect()` para identificar aproximadamente:

- HTML
- JSON
- SQL
- Python
- Bash
- CSS
- JavaScript

El fallback mantiene la etiqueta y el botón Copiar, aunque el coloreado completo depende de Highlight.js.

### Laboratorio visible

`Laboratorio UI & Código` incluye un selector interactivo para probar:

- JavaScript
- Python
- SQL
- JSON
- HTML
- CSS
- Bash

Al cambiar el lenguaje, el bloque se vuelve a colorear y actualiza su etiqueta.

### Sobre el warning `file:` del navegador

Si el dashboard se abre con doble clic como:

```text
file:///C:/...
```

Chrome puede mostrar warnings de origen como:

```text
'file:' URLs are treated as unique security origins
```

Eso es una restricción del navegador al trabajar con archivos locales, no un error de sintaxis de Axolotl.

Para desarrollo conviene servir la carpeta:

```bash
python -m http.server 8000
```

y abrir:

```text
http://localhost:8000
```

El dashboard mantiene fallback para varias funciones cuando corre por `file://`, pero `localhost` reproduce mejor un entorno web real.

## 27. Paleta Axolotl Rosado

La colorimetría predeterminada está inspirada directamente en el arte de ajolote suministrado.

Los dos colores dominantes aproximados extraídos de la referencia son:

| Token visual | Hex | Uso |
|---|---|---|
| Branquias / acento | `#FC779B` | `--primary`, foco, navegación activa, gráficas |
| Cuerpo / blush | `#FCBBD5` | `--secondary`, gradientes, acentos suaves |
| Rosa profundo | `#E9638A` | `--primary-2` |
| Lavanda complementaria | `#D49AF6` | `--info`, final de gradientes Prism |
| Noche Axolotl | `#09070B` | fondo principal oscuro |

Tokens de referencia:

```css
--axolotl-gills: #FC779B;
--axolotl-body: #FCBBD5;
--axolotl-deep: #E9638A;
--axolotl-glow: #F5A2C1;
--axolotl-lavender: #D49AF6;
--axolotl-night: #09070B;
```

### Base oscura

La paleta no convierte la plantilla en un dashboard rosa claro. Las superficies continúan siendo oscuras con un matiz ciruela:

```css
--body-bg: #09070B;
--sidebar-bg: #0D0910;
--surface: #151018;
--surface-2: #1B131E;
--border: #352330;
--border-soft: #281B25;
--body-color: #FFF2F7;
--muted: #B58C9B;
```

### Neo / Prism

Las gráficas de línea usan principalmente:

```text
Branquias (#FC779B)
→ Cuerpo (#FCBBD5)
→ Lavanda (#D49AF6)
```

El glow debe permanecer discreto.

El doughnut conserva colores planos en sus segmentos y el efecto luminoso únicamente dentro del hueco, conforme a la regla existente del componente.

### Colores semánticos

Aunque la identidad principal es rosa, estados funcionales siguen siendo distinguibles:

```css
--success: #8EDBC7;
--danger:  #FF5F7F;
--warning: #F5B86F;
--info:    #D49AF6;
```

Esto evita que éxito, advertencia y error se vuelvan indistinguibles.

### Persistencia del editor de tema

Esta versión guarda personalizaciones en:

```text
axolotl-theme-v21
```

La clave está versionada para que un tema turquesa almacenado por versiones anteriores no reemplace accidentalmente la nueva colorimetría Axolotl.

## 28. Scrollbars y overlays

### Scrollbar principal

Abrir un modal, buscador global, drawer o ventana superpuesta **no debe ocultar la barra de desplazamiento principal**.

La clase histórica `.no-scroll` ya no usa `overflow: hidden`. Se conserva por compatibilidad con el JavaScript existente, pero su función visual es mantener el gutter estable:

```css
html {
  scrollbar-gutter: stable;
}

.no-scroll {
  overflow-y: scroll !important;
  overflow-x: hidden !important;
}
```

Esto evita:

- que desaparezca la perilla principal;
- que el contenido cambie de ancho;
- el salto horizontal al abrir/cerrar un modal.

Los componentes superpuestos siguen usando backdrop/obscurecimiento.

Para evitar que el scroll interno de un modal o drawer se propague accidentalmente a la página se usa:

```css
overscroll-behavior: contain;
```

en `.modal`, `.drawer` y `.search-results`.

### Scrollbar del sidebar

La barra del sidebar debe ser visualmente discreta.

En escritorio:

- permanece oculta mientras el mouse está fuera del sidebar;
- aparece al hacer hover sobre el sidebar;
- usa el degradado `primary → secondary`;
- no modifica el ancho del sidebar.

Ejemplo:

```css
.sidebar {
  scrollbar-color: transparent transparent;
}

.sidebar:hover {
  scrollbar-color: rgba(var(--primary-rgb), .72) transparent;
}

.sidebar::-webkit-scrollbar-thumb {
  background: transparent;
}

.sidebar:hover::-webkit-scrollbar-thumb {
  background: linear-gradient(
    180deg,
    rgba(var(--primary-rgb), .76),
    rgba(var(--secondary-rgb), .48)
  );
}
```

En dispositivos touch se oculta por completo porque no existe un hover persistente y el desplazamiento se hace por gesto.

## 29. Métricas del servidor

La plantilla incluye `#server-metrics-section` dentro de la página dedicada `/servidor`.

Las métricas técnicas no deben volver a mezclarse con el overview principal de `/dashboard`.

Componentes:

### CPU

- velocímetro semicircular SVG;
- arco animado;
- aguja con transición suave;
- lectura porcentual;
- temperatura;
- frecuencia;
- procesos;
- uptime;
- ocho mini círculos por núcleo.

Los núcleos usan `conic-gradient` y la variable:

```css
--core-value
```

Ejemplo:

```js
setCore(0, 54);
```

### RAM

Utiliza el mismo componente de velocímetro que CPU.

API:

```js
setGauge('ram', 68);
```

La animación modifica:

- `stroke-dashoffset` del arco;
- rotación de la aguja;
- texto porcentual.

### Red

Usa Chart.js con:

- descarga;
- subida;
- latencia;
- paquetes por segundo;
- tooltip interactivo.

La gráfica conserva el lenguaje Neo/Prism de Axolotl.

### Almacenamiento

Cada volumen usa:

```html
<div class="storage-track">
  <div class="storage-fill"></div>
</div>
```

El ancho cambia suavemente con `transition`.

Incluye:

- Sistema;
- Datos;
- Backups;
- I/O lectura;
- I/O escritura.

### Simulación

La demostración actual actualiza métricas aproximadamente cada 2.4 segundos.

```js
serverMetricsTimer = setInterval(
  tickServerMetrics,
  2400
);
```

En producción, reemplazar `tickServerMetrics()` por datos reales obtenidos mediante:

- endpoint REST;
- WebSocket;
- Server-Sent Events;
- agente del servidor;
- Prometheus/Grafana API;
- servicio propio de monitoreo.

### Pausar

El botón `Pausar` detiene únicamente la simulación visual y permite probar el estado estático del componente.

### Corrección de animación del velocímetro

El velocímetro no debe depender de `CSS transform` sobre un `<g>` SVG, porque algunos navegadores —especialmente al abrir el archivo mediante `file://`— pueden no aplicar el origen de transformación como se espera.

La aguja se anima mediante el atributo SVG:

```js
needle.setAttribute(
  'transform',
  `rotate(${angle} 105 105)`
);
```

El arco obtiene su longitud real:

```js
const length = arc.getTotalLength();
arc.style.strokeDasharray = length;
arc.style.strokeDashoffset =
  length * (1 - value / 100);
```

La transición completa usa `requestAnimationFrame` y `easeInOutCubic`, incluyendo la numeración.

El valor porcentual del velocímetro se coloca **debajo del SVG**, fuera de la trayectoria de la aguja, para que texto y aguja nunca se superpongan.

Los núcleos también interpolan su valor y su número mediante `requestAnimationFrame`. La propiedad `--core-value` se registra con `@property` para permitir transiciones suaves del `conic-gradient`.

La sección de métricas del servidor no utiliza skeleton de arranque en la demo, porque ocultaría la primera animación del velocímetro.

### Inicialización obligatoria

El monitor debe arrancarse explícitamente al final del documento:

```js
initServerMetrics();
```

La función no debe quedar solamente declarada. Sin esa llamada:

- el velocímetro conserva la aguja en su posición HTML inicial;
- los núcleos no animan;
- la gráfica de red queda vacía;
- el intervalo de actualización nunca comienza.

La secuencia de arranque recomendada es:

```js
loadTheme();
initCharts();
initSmoothMotion();
initSkeletonDemo();
enhanceCodeBlocks(document);
initCodeLab();
initServerMetrics();
```

En consola debe aparecer:

```text
Axolotl: monitor de servidor iniciado.
```

## 30. Reglas finales para convertir Axolotl en plantilla

Este documento es la fuente de verdad del sistema visual y funcional.

### Reglas obligatorias

1. **No hardcodear colores de identidad**
   - Usar `--primary`, `--secondary`, `--success`, `--danger`, `--warning`, `--info`.
   - Para transparencias usar sus equivalentes `--*-rgb`.
   - Mantener la paleta Axolotl como valor predeterminado.

2. **Toda interacción visible debe funcionar**
   - Si parece clicable, debe ejecutar una acción.
   - No dejar botones, dropdowns, tabs, iconos o enlaces decorativos sin comportamiento.
   - Si una función aún no existe, ocultar o deshabilitar el control.

3. **Mantener el dashboard multipropósito**
   - No volver a convertir la interfaz en una plantilla exclusivamente financiera.
   - Los módulos de mercado existentes funcionan como ejemplos/demo de visualización.
   - Los nuevos módulos deben usar nombres y patrones reutilizables.

4. **Usar componentes existentes antes de crear nuevos**
   - Cards: `.card`, `.section`.
   - Botones: `.btn-ui` + variante semántica.
   - Modales: `openUiModal()`.
   - Toasts: `toast()`.
   - Skeletons: `setLoading()` / `withSkeleton()`.
   - Código: `<pre><code>` + `enhanceCodeBlocks()`.
   - Scroll: `smoothScrollTo()`.

5. **Responsive obligatorio**
   - Desktop amplio.
   - Desktop compacto.
   - Tablet.
   - Mobile.
   - No introducir scroll horizontal global.
   - Sidebar:
     - desktop: expandible/colapsable;
     - tablet: estrecho con texto;
     - móvil: drawer off-canvas.

6. **Sidebar**
   - Las categorías pueden usar `nav-group`.
   - Subcategorías usan `nav-subitem`.
   - Los padres expandibles rotan su chevron.
   - La subcategoría activa marca también al padre.
   - Los iconos deben permanecer centrados y consistentes.

7. **Overlays**
   - Modales, drawer y buscador no deben ocultar la scrollbar principal.
   - Conservar `scrollbar-gutter: stable`.
   - Mantener `overscroll-behavior: contain` en superficies internas desplazables.

8. **Animación**
   - Debe ser suave y corta.
   - Respetar `prefers-reduced-motion`.
   - Para scroll usar el motor disponible mediante `smoothScrollTo()`.
   - Para valores animados usar `requestAnimationFrame` cuando sea necesario.
   - No abusar de glow o movimiento decorativo.

9. **Neo / Prism**
   - Permitido principalmente en:
     - gráficas de líneas;
     - sparklines;
     - estados activos;
     - elementos de datos.
   - Evitar glows fuertes en cards normales.
   - El donut mantiene segmentos planos.
   - Cualquier brillo del donut debe quedar dentro del hueco central.

10. **Skeletons**
    - Todo contenido dependiente de API, DB o imágenes debe tener estado loading.
    - El skeleton debe conservar aproximadamente el tamaño final para evitar layout shift.
    - Eliminarlo tanto en éxito como en error.

11. **Código**
    - Todo bloque `<pre><code>` debe:
      - detectar lenguaje;
      - colorear sintaxis;
      - mostrar etiqueta de lenguaje;
      - ofrecer botón Copiar;
      - cambiar temporalmente a `✓ Copiado`.
    - Para ejemplos ambiguos usar `data-language`.

12. **Accesibilidad**
    - Usar `<button>` para acciones.
    - `aria-label` en controles solo-icono.
    - `aria-expanded` en dropdowns/submenús.
    - Foco visible.
    - Escape debe cerrar overlays cuando corresponda.
    - Navegación usable con teclado.

13. **Datos reales**
    - Sustituir simulaciones por REST, WebSocket, SSE u otra fuente.
    - Mantener la misma API de renderizado del frontend cuando sea posible.
    - No acoplar los estilos al backend.

14. **Métricas del servidor**
    - CPU y RAM usan velocímetro animado.
    - Núcleos usan círculos independientes.
    - Red usa Chart.js.
    - Almacenamiento usa barras de progreso.
    - La demo actual es simulada; producción debe conectarse a una fuente real.

15. **Tema**
    - El editor de tema modifica tokens semánticos.
    - Nuevos componentes deben reaccionar a cambios globales.
    - Evitar hex directos salvo excepciones documentadas.

### Estructura mínima del paquete

```text
axolotl-template/
├── index.html
└── design.md
```

### Recomendación para producción

La versión actual es una plantilla frontend en un único archivo HTML.

Antes de producción puede separarse en:

```text
src/
├── css/
│   ├── tokens.css
│   ├── components.css
│   └── responsive.css
├── js/
│   ├── app.js
│   ├── charts.js
│   ├── server-metrics.js
│   └── ui.js
└── index.html
```

La separación es opcional; el archivo actual funciona como referencia completa y portable.

---

**Estado del proyecto:** plantilla base lista para reutilización y adaptación.



## 31. Integración Goravel

La implementación activa del design system vive en una aplicación Goravel 100% Go.

### Mapeo de archivos

```text
design.md                              fuente de verdad visual
public/css/app.css                     tokens y componentes visuales
public/js/app.js                       interactividad frontend
resources/views/dashboard.tmpl        página principal
resources/views/partials/head.tmpl     metadatos y estilos
resources/views/partials/scripts.tmpl  dependencias JS y app.js
resources/views/partials/sidebar.tmpl  navegación lateral
resources/views/partials/topbar.tmpl   header superior
resources/views/partials/stats.tmpl    métricas superiores
resources/views/partials/dashboard-grid.tmpl
resources/views/partials/server-metrics.tmpl
resources/views/partials/code-lab.tmpl
resources/views/partials/overlays.tmpl modales/drawer/buscador/toasts
```

### Regla de plantillas

El equivalente funcional a Blade es `html/template` integrado por Goravel/Gin.

Los componentes compartidos deben declararse como templates nombrados:

```gotemplate
{{ define "partials/mi-componente" }}
  ...
{{ end }}
```

y consumirse mediante:

```gotemplate
{{ template "partials/mi-componente" . }}
```

No agregar otro motor de plantillas mientras `html/template` cubra el caso de uso.

### Regla de assets

- CSS de aplicación: `public/css/`.
- JavaScript de aplicación: `public/js/`.
- Las vistas no deben volver a incluir bloques monolíticos grandes de `<style>` o `<script>`.
- Los assets públicos se sirven bajo `/assets/*`.
- Las dependencias CDN actuales son demostrativas; si una aplicación necesita operación offline o CSP estricta, migrarlas a assets locales.

### Backend y frontend

Los controllers entregan datos; las vistas presentan esos datos.

No colocar consultas ORM directamente en templates ni mezclar lógica de persistencia dentro de JavaScript de presentación.

Para datos dinámicos futuros:

```text
Route -> Middleware -> Controller -> ORM/servicio -> View/JSON
```

Crear services solamente cuando exista lógica de negocio real que justifique una capa adicional.

### Rutas y middleware

La plantilla base es multipágina:

```text
GET /              dashboard de entrada
GET /dashboard     overview general
GET /cuenta        perfil, preferencias y seguridad
GET /servidor      métricas técnicas
GET /componentes   laboratorio y catálogo UI
GET /health        healthcheck
/assets/*          estáticos
```

Las páginas comparten sidebar, topbar, overlays, scripts y design system. No concentrar funcionalidades independientes en una única vista larga cuando representan áreas distintas del producto.

El middleware global `SecurityHeaders` aplica headers defensivos sin romper los CDNs de la demo.

### Base de datos

La conexión predeterminada es:

```dotenv
DB_CONNECTION=sqlite
DB_DATABASE=storage/axolotl.sqlite
```

Los modelos pertenecen a `app/models/` y las migraciones a `database/migrations/`.

No hardcodear rutas absolutas de base de datos dentro de modelos o controllers.

### Dashboard multipropósito

La página inicial debe demostrar patrones reutilizables de varias áreas, no representar una sola industria.

Contenido recomendado de `/dashboard`:

```text
KPIs
→ rendimiento y tendencias
→ proyectos y tareas
→ actividad
→ equipo, agenda e integraciones
```

Las áreas técnicas y de configuración viven en rutas dedicadas:

```text
/servidor    → métricas técnicas
/cuenta      → perfil, preferencias y seguridad
/componentes → laboratorio UI
/archivos    → explorador backend del workspace demo
```

Los datos repetibles del demo salen de `DashboardController` y se renderizan con `html/template`. Un proyecto consumidor puede sustituirlos por ORM, APIs o servicios manteniendo los mismos partials.

Reglas:

- evitar cards que repitan la misma información;
- cada módulo debe demostrar un patrón diferente;
- no centrar la plantilla en finanzas, cripto, e-commerce u otra vertical;
- preservar IDs de sección usados por sidebar, búsqueda y scroll spy;
- mantener el logo SVG visible una sola vez en el sidebar;
- no repetir la marca como decoración en el footer del sidebar;
- usar partials separados cuando el módulo tenga responsabilidad visual propia.

### Regla para nuevos módulos

Al añadir un módulo:

1. crear el markup como partial si tiene responsabilidad visual propia;
2. reutilizar tokens existentes;
3. añadir JS solo si HTML/CSS no resuelven la interacción;
4. añadir controller/ruta únicamente si requiere backend;
5. añadir migración/modelo solo si persiste datos;
6. añadir prueba mínima para lógica no trivial;
7. actualizar `contexto/` y `tareas/`;
8. validar `go test ./...` y `go vet ./...`.

## 32. Autenticación y primer setup

Axolotl incluye un flujo de acceso reutilizable para proyectos multipropósito.

### Rutas públicas

```text
/           decide setup, login o dashboard
/setup      primer administrador
/login      inicio de sesión
/register   registro estándar
/health     healthcheck
/assets/*   archivos públicos
```

### Rutas privadas

```text
/dashboard
/cuenta
/servidor
/componentes
/archivos
```

Todas las páginas privadas usan `RequireAuth`.

### Primer administrador

- Al arrancar, si no existe un usuario activo con rol `admin`, se genera un PIN de 6 dígitos con `crypto/rand`.
- El PIN vive solo en memoria.
- Se escribe mediante el canal de log `setup`, que imprime a consola y guarda `storage/logs/setup.log`.
- `/setup` nunca recibe el PIN desde el backend; el usuario debe consultarlo en el entorno del servidor.
- El frontend usa seis inputs OTP independientes.
- El backend concatena y valida los seis dígitos.
- La comparación del PIN usa `subtle.ConstantTimeCompare`.
- Al crear correctamente el administrador, el PIN se invalida.
- Si ya existe administrador, `/setup` no permite crear otro.

### Login y registro

- Contraseñas con Argon2id.
- Sesiones de Goravel; el login regenera el ID de sesión.
- Formularios POST protegidos con CSRF.
- Setup, login y registro tienen rate limit por IP.
- `/register` solo puede crear `role=user`; el servidor ignora cualquier rol enviado por el navegador.
- El registro se puede desactivar con `AUTH_ALLOW_REGISTRATION=false`.
- Logout es POST con CSRF.
- En producción se debe usar HTTPS y `SESSION_SECURE=true`.

### UI de autenticación

Los estilos de setup/login/registro viven separados del dashboard:

```text
resources/views/auth.tmpl
public/css/auth.css
public/js/auth.js
```

El diseño mantiene la marca Axolotl, pero evita cargar Chart.js, AOS, Lenis y el resto del dashboard antes de autenticar.

## 33. Branding oficial Axolotl

La marca oficial vive en:

```text
public/brand/axolotl.svg
public/brand/axolotl.png
```

Rutas públicas:

```text
/assets/brand/axolotl.svg
/assets/brand/axolotl.png
```

### Reglas de uso

- `axolotl.svg` es la marca principal para interfaz, sidebar y componentes escalables.
- `axolotl.png` es el fallback raster para favicon/plataforma y usos donde SVG no sea apropiado.
- El botón principal del sidebar usa el SVG como logotipo y conserva su comportamiento de collapse/expand al hacer hover/click.
- No sustituir la marca principal por iconos genéricos de Font Awesome.
- No alterar los colores internos del SVG con filtros agresivos; el archivo ya contiene la paleta oficial.
- Mantener `object-fit: contain` y preservar proporciones.
- No deformar ancho/alto de manera independiente.
- Para nuevos favicons o metadata, reutilizar estos assets antes de crear duplicados.

### Head recomendado

```html
<link rel="icon" type="image/svg+xml" href="/assets/brand/axolotl.svg">
<link rel="icon" type="image/png" href="/assets/brand/axolotl.png">
<link rel="apple-touch-icon" href="/assets/brand/axolotl.png">
```

### Sidebar

El logo del sidebar debe seguir siendo también el control de expandir/colapsar. En desktop, hover revela el icono de acción; en touch/mobile, el estado del drawer controla qué icono se muestra.

## 34. Explorador de archivos

Axolotl incluye una página independiente:

```text
/archivos
```

Su objetivo es demostrar un módulo full-stack reusable para trabajar con archivos desde Go sin exponer libremente el filesystem del servidor.

### Backend

La raíz se configura mediante:

```dotenv
FILES_ROOT=storage/demo-files
FILES_MAX_EDIT_BYTES=2097152
```

El backend vive en:

```text
app/files/workspace.go
app/files/configured.go
app/http/controllers/file_controller.go
```

API protegida:

```text
GET  /api/files
GET  /api/files/read
POST /api/files/create
POST /api/files/save
POST /api/files/delete
```

### Reglas de seguridad

- todas las operaciones deben permanecer dentro de `FILES_ROOT`;
- nunca aceptar rutas absolutas del cliente;
- bloquear `../` y cualquier path traversal equivalente;
- no seguir ni permitir symlinks;
- no devolver rutas físicas absolutas del host;
- el root nunca puede eliminarse;
- nombres nuevos deben validarse en backend;
- la edición se limita a texto UTF-8;
- archivos binarios no se abren en el editor;
- respetar `FILES_MAX_EDIT_BYTES`;
- create/save/delete requieren sesión y CSRF;
- un error interno no debe revelar paths del sistema operativo.

### UI tipo explorador

La página usa:

- toolbar de Atrás / Subir / Actualizar;
- breadcrumbs;
- acceso rápido a carpetas raíz;
- vista de detalles con nombre, tipo, tamaño y fecha;
- selección visual;
- doble clic para abrir carpeta o archivo;
- creación de archivos y carpetas;
- eliminación con modal de confirmación.

El item `Archivos` del sidebar usa `data-route="files"` y debe quedar activo en `/archivos`.

### Editor tipo VS Code

Los archivos editables se abren dentro de la misma página.

Componentes:

- pestaña con nombre;
- indicador de cambios sin guardar;
- números de línea;
- textarea transparente superpuesta al bloque resaltado;
- Highlight.js mediante el CDN ya utilizado por Axolotl;
- detección de lenguaje según extensión;
- botón Guardar;
- atajo `Ctrl/Cmd + S`;
- Tab inserta dos espacios;
- barra inferior con ruta y estado.

### Cierre con cambios sin guardar

Cerrar o navegar fuera del archivo modificado debe preguntar antes de descartar cambios.

Dentro de la aplicación se reutiliza `openUiModal()` con tres opciones:

```text
Cancelar
Cerrar sin guardar
Guardar y cerrar
```

Además se mantiene `beforeunload` como protección del navegador cuando el usuario abandona la página completa.

### Assets

```text
resources/views/files.tmpl
resources/views/partials/files-main.tmpl
public/css/files.css
public/js/files.js
```

Los estilos y JavaScript del explorador se cargan únicamente cuando `ActiveRoute == "files"`.

