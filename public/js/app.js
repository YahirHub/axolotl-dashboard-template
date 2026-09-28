const $=(s,r=document)=>r.querySelector(s);
const $$=(s,r=document)=>[...r.querySelectorAll(s)];

function cssVar(name){
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
}
function hexToRgb(hex){
  const clean=String(hex||'').replace('#','').trim();
  if(!/^[0-9a-fA-F]{6}$/.test(clean)) return null;
  return [parseInt(clean.slice(0,2),16),parseInt(clean.slice(2,4),16),parseInt(clean.slice(4,6),16)];
}
function setSemanticColor(name,value){
  const root=document.documentElement;
  root.style.setProperty(`--${name}`,value);
  const rgb=hexToRgb(value);
  if(rgb) root.style.setProperty(`--${name}-rgb`,rgb.join(','));
}
const defaultTheme={
  primary:'#FC779B',
  secondary:'#FCBBD5',
  success:'#8EDBC7',
  danger:'#FF5F7F',
  warning:'#F5B86F',
  info:'#D49AF6'
};
function getCurrentTheme(){
  return Object.fromEntries(Object.keys(defaultTheme).map(k=>[k,cssVar(`--${k}`)]));
}
function saveTheme(){
  try{localStorage.setItem('axolotl-theme-v21',JSON.stringify(getCurrentTheme()))}catch(_){}
}
function loadTheme(){
  try{
    const saved=JSON.parse(localStorage.getItem('axolotl-theme-v21')||'null');
    if(saved) Object.entries(saved).forEach(([k,v])=>{if(defaultTheme[k]&&/^#[0-9a-fA-F]{6}$/.test(v))setSemanticColor(k,v)});
  }catch(_){}
}


const state={
  currency:'MXN',
  user:{name:'Usuario Axolotl',email:'usuario@axolotl.demo',role:'Administrador',alias:'AX'},
  rates:{MXN:1,USD:0.055,JPY:8.25},
  symbols:{MXN:'$',USD:'US$',JPY:'¥'},
  compact:false,
  animations:true,
  risk:'Moderado',
  selectedAsset:'WEB',
  notifications:[
    {id:1,title:'Despliegue completado',text:'Portal web v2.4.0 quedó disponible para validación.',time:'Hace 8 min',unread:true,icon:'fa-rocket'},
    {id:2,title:'Tarea próxima a vencer',text:'Validar flujo de onboarding vence hoy a las 16:30.',time:'Hace 24 min',unread:true,icon:'fa-list-check'},
    {id:3,title:'Capacidad elevada',text:'La cola de procesamiento alcanzó 82% de utilización.',time:'Hace 1 h',unread:true,icon:'fa-server'},
    {id:4,title:'Resumen diario listo',text:'El resumen operativo ya está disponible.',time:'Ayer',unread:false,icon:'fa-file-lines'}
  ]
};

const assets={
 WEB:{name:'Portal web',mxn:420000,change:12.4,forecast:8.2,momentum:'En ritmo',vol:'Media',read:'Saludable',series:[62,64,65,67,66,69,71,70,73,75,76,78,80,81,83,84,86]},
 CRM:{name:'CRM comercial',mxn:315000,change:7.2,forecast:5.6,momentum:'En ritmo',vol:'Baja',read:'Estable',series:[54,55,57,58,60,59,62,63,65,66,68,69,71,72,74,75,77]},
 APP:{name:'App móvil',mxn:580000,change:-3.1,forecast:4.1,momentum:'Atención',vol:'Alta',read:'Revisión',series:[72,73,71,70,69,68,70,67,66,65,64,66,65,67,68,69,71]},
 DATA:{name:'Data Hub',mxn:265000,change:4.9,forecast:6.4,momentum:'Estable',vol:'Media',read:'Saludable',series:[48,50,52,51,54,56,57,59,58,61,63,64,66,68,69,71,73]},
 OPS:{name:'Automatización Ops',mxn:190000,change:9.3,forecast:7.1,momentum:'En ritmo',vol:'Baja',read:'Saludable',series:[45,47,49,52,53,55,57,59,61,64,66,68,70,72,75,77,79]}
};

const periodData={
 '1H':{value:76,label:'Estable',caption:'La operación se mantiene estable durante la última hora.',series:[70,72,71,73,74,72,75,74,76,77,75,76,76]},
 '24H':{value:84,label:'Óptimo',caption:'La operación mantiene buen ritmo con capacidad disponible.',series:[68,70,72,71,74,77,75,79,78,80,81,83,86]},
 '7D':{value:79,label:'Saludable',caption:'La semana mantiene entregas constantes y pocas incidencias.',series:[64,66,68,69,70,72,74,73,76,77,78,80,79]},
 '30D':{value:73,label:'Saludable',caption:'El mes conserva un desempeño estable con margen de mejora.',series:[66,67,69,68,70,71,72,73,71,74,75,73,76]}
};

const news=[
 {title:'Despliegue completado',time:'Hace 18 min',body:'Portal web versión 2.4.0 quedó disponible. Este feed demuestra cómo presentar eventos provenientes de backend, auditoría o webhooks.'},
 {title:'Nuevo comentario',time:'Hace 42 min',body:'Se agregó una observación al CRM comercial. En producción este componente puede alimentarse desde la base de datos o un sistema de colaboración.'},
 {title:'Alerta de capacidad',time:'Hace 1 h',body:'La cola de procesamiento de la App móvil alcanzó 82%. El contenido es demostrativo y puede sustituirse por telemetría real.'}
];

let sentChart,predChart,portChart;
const charts=[];

const hoverLinePlugin={
  id:'hoverLine',
  afterDatasetsDraw(chart){
    const active=chart.tooltip && chart.tooltip.getActiveElements ? chart.tooltip.getActiveElements() : [];
    if(!active || !active.length || !chart.chartArea) return;
    const x=active[0].element.x;
    const {ctx,chartArea:{top,bottom}}=chart;
    ctx.save();
    ctx.beginPath();
    ctx.moveTo(x,top);
    ctx.lineTo(x,bottom);
    ctx.lineWidth=1;
    ctx.strokeStyle='rgba(89,228,217,.34)';
    ctx.setLineDash([3,3]);
    ctx.stroke();
    ctx.restore();
  }
};

const prismGlowPlugin={
  id:'prismGlow',
  beforeDatasetDraw(chart,args){
    if(!chart.chartArea) return;
    const dataset=chart.data.datasets[args.index];

    // Only datasets with an explicit prismGlowColor get an outer glow.
    // This intentionally excludes doughnut/pie charts.
    if(!dataset.prismGlowColor) return;

    const ctx=chart.ctx;
    ctx.save();
    dataset.__prismShadowSaved=true;
    ctx.shadowColor=dataset.prismGlowColor;
    ctx.shadowBlur=dataset.prismGlowBlur ?? 8;
    ctx.shadowOffsetX=0;
    ctx.shadowOffsetY=0;
  },
  afterDatasetDraw(chart,args){
    const dataset=chart.data.datasets[args.index];
    if(dataset && dataset.__prismShadowSaved){
      chart.ctx.restore();
      dataset.__prismShadowSaved=false;
    }
  }
};

Chart.register(hoverLinePlugin,prismGlowPlugin);

function prismLineGradient(chart,primary=cssVar('--primary'),secondary=cssVar('--secondary'),tertiary=cssVar('--info')){
  const area=chart.chartArea;
  if(!area) return primary;
  const g=chart.ctx.createLinearGradient(area.left,0,area.right,0);
  g.addColorStop(0,primary);
  g.addColorStop(.52,secondary);
  g.addColorStop(1,tertiary);
  return g;
}

function prismAreaGradient(chart,primary=cssVar('--primary-rgb'),secondary=cssVar('--secondary-rgb')){
  const area=chart.chartArea;
  if(!area) return 'rgba(89,228,217,.04)';
  const g=chart.ctx.createLinearGradient(0,area.top,0,area.bottom);
  g.addColorStop(0,`rgba(${primary},.22)`);
  g.addColorStop(.42,`rgba(${secondary},.085)`);
  g.addColorStop(1,'rgba(3,16,21,0)');
  return g;
}


function formatMoney(mxn, compact=false){
  const rate=state.rates[state.currency];
  const value=Number(mxn)*rate;
  const locale=state.currency==='JPY'?'ja-JP':state.currency==='USD'?'en-US':'es-MX';
  const currency=state.currency;
  if(compact && Math.abs(value)>=1e6){
    return new Intl.NumberFormat(locale,{style:'currency',currency,notation:'compact',maximumFractionDigits:2}).format(value);
  }
  const max=value<100?2:0;
  return new Intl.NumberFormat(locale,{style:'currency',currency,maximumFractionDigits:max}).format(value);
}

function updateCurrency(){
  $$('.money').forEach(el=>{
    const mxn=Number(el.dataset.mxn);
    if(!Number.isFinite(mxn)) return;
    const useCompact=el.closest('.donut-center')!==null;
    const extra=[...el.children].map(x=>x.outerHTML).join('');
    if(el.classList.contains('market-price')){
      const span=el.querySelector('span');
      const change=span?span.outerHTML:'';
      el.innerHTML=formatMoney(mxn)+change;
    }else{
      el.innerHTML=formatMoney(mxn,useCompact)+extra;
    }
  });
  const label={MXN:'MXN / $',USD:'USD / $',JPY:'JPY / ¥'}[state.currency];
  $('#currencyLabel').textContent=label;
  $$('#currencyMenu .drop-item').forEach(btn=>{
    const sel=btn.dataset.currency===state.currency;
    btn.classList.toggle('selected',sel);
    $('.check',btn).textContent=sel?'✓':'';
  });
}

function makeSpark(id,data,color='#FC779B',secondary='#FCBBD5',tertiary='#D49AF6'){
  const c=new Chart($(id),{
    type:'line',
    data:{
      labels:data.map((_,i)=>i),
      datasets:[{
        data,
        borderColor:ctx=>prismLineGradient(ctx.chart,color,secondary,tertiary),
        backgroundColor:ctx=>prismAreaGradient(ctx.chart),
        fill:true,
        borderWidth:1.65,
        pointRadius:0,
        pointHoverRadius:3,
        pointHoverBackgroundColor:'#FFF2F7',
        pointHoverBorderColor:'#FC779B',
        pointHoverBorderWidth:1,
        tension:.38,
        prismGlowColor:color,
        prismGlowBlur:7
      }]
    },
    options:{
      responsive:true,maintainAspectRatio:false,
      interaction:{mode:'index',intersect:false},
      plugins:{
        legend:{display:false},
        tooltip:{
          enabled:true,
          displayColors:false,
          backgroundColor:'rgba(5,19,24,.96)',
          borderColor:'rgba(89,228,217,.28)',
          borderWidth:1,
          titleColor:'#B18D9C',
          bodyColor:'#FFF2F7',
          callbacks:{
            title:items=>`Punto ${items[0].dataIndex+1}`,
            label:ctx=>`Valor: ${Number(ctx.parsed.y).toFixed(2)}`
          }
        }
      },
      hover:{mode:'index',intersect:false},
      scales:{x:{display:false},y:{display:false}},
      animation:{duration:350}
    }
  });
  charts.push(c);
}

function initCharts(){
  makeSpark('#sp1',[4,4.2,4.1,4.6,4.4,4.9,4.7,5.3,5.1,5.6,5.4,6.0],cssVar('--primary'),cssVar('--secondary'),cssVar('--info'));
  makeSpark('#sp2',[3.7,4.0,3.9,4.3,4.7,4.5,4.9,5.3,5.1,5.7,5.7,6.3],cssVar('--success'),cssVar('--primary'),cssVar('--info'));
  makeSpark('#sp3',[4.0,4.2,4.1,4.5,4.3,4.8,4.6,5.0,4.9,5.4,5.3,5.8],cssVar('--secondary'),cssVar('--info'),cssVar('--primary'));
  makeSpark('#a1',[3,3.4,3.2,3.9,3.5,4.2,4.0,4.7,4.5,5.2],cssVar('--warning'),cssVar('--secondary'),cssVar('--primary'));
  makeSpark('#a2',[3.0,3.2,3.5,3.3,3.8,4.2,4.0,4.4,4.8,5.3],'#C9ACFF',cssVar('--info'),cssVar('--primary'));
  makeSpark('#a3',[5.4,5.2,5.0,5.1,4.7,4.8,4.3,4.1,4.0,3.6],cssVar('--danger'),cssVar('--secondary'),cssVar('--info'));
  makeSpark('#a4',[3.1,3.0,3.5,3.4,3.8,4.0,3.9,4.5,4.4,5.0],cssVar('--secondary'),cssVar('--info'),cssVar('--primary'));

  sentChart=new Chart($('#sentChart'),{
    type:'line',
    data:{labels:['00:00','','','','06:00','','','','12:00','','','','18:00'],
      datasets:[{
        data:periodData['24H'].series,
        borderColor:ctx=>prismLineGradient(ctx.chart,cssVar('--primary'),cssVar('--secondary'),cssVar('--info')),
        backgroundColor:ctx=>prismAreaGradient(ctx.chart),
        borderWidth:2,
        pointRadius:0,
        pointHoverRadius:4,
        pointHoverBackgroundColor:'#FFF2F7',
        pointHoverBorderColor:'#FC779B',
        pointHoverBorderWidth:1,
        tension:.38,
        fill:true,
        prismGlowColor:cssVar('--primary'),
        prismGlowBlur:10
      }]},
    options:{
      responsive:true,maintainAspectRatio:false,
      interaction:{mode:'index',intersect:false},
      plugins:{
        legend:{display:false},
        tooltip:{
          enabled:true,
          displayColors:false,
          backgroundColor:'rgba(5,19,24,.97)',
          borderColor:'rgba(89,228,217,.24)',
          borderWidth:1,
          titleColor:'#B18D9C',
          bodyColor:'#FFF2F7',
          callbacks:{
            title:items=>`Rendimiento · ${sentimentHoverLabel(items[0].dataIndex)}`,
            label:ctx=>`Índice operativo: ${ctx.parsed.y}/100`,
            afterLabel:ctx=>sentimentReading(ctx.parsed.y)
          }
        }
      },
      hover:{mode:'index',intersect:false},
      onHover:(event,elements,chart)=>{
        if(!elements.length) return;
        const i=elements[0].index;
        const value=chart.data.datasets[0].data[i];
        $('#sentHoverTime').textContent=sentimentHoverLabel(i);
        $('#sentHoverValue').textContent=`${value} / 100`;
        const reading=sentimentReading(value);
        const el=$('#sentHoverRead');
        el.textContent=reading;
        el.className=value>=60?'up':value<45?'down':'';
      },
      scales:{x:{grid:{display:false},border:{display:false},ticks:{font:{size:9},color:'#8E6F7B',maxTicksLimit:4}},
      y:{position:'right',min:0,max:100,grid:{color:'rgba(120,170,170,.11)'},border:{display:false},ticks:{font:{size:9},color:'#8E6F7B',stepSize:25}}}
    }
  });

  predChart=new Chart($('#predChart'),{
    type:'line',
    data:{labels:['00:00','','','','06:00','','','','12:00','','','','18:00','','','','24:00'],
      datasets:[{
        data:assets.WEB.series,
        borderColor:ctx=>prismLineGradient(ctx.chart,cssVar('--primary'),cssVar('--secondary'),cssVar('--info')),
        backgroundColor:ctx=>prismAreaGradient(ctx.chart),
        borderWidth:2,
        pointRadius:0,
        pointHoverRadius:4,
        pointHoverBackgroundColor:'#FFF2F7',
        pointHoverBorderColor:'#FCBBD5',
        tension:.34,
        fill:true,
        prismGlowColor:cssVar('--secondary'),
        prismGlowBlur:10,
        segment:{borderDash:ctx=>ctx.p0DataIndex>=11?[5,5]:undefined}
      }]},
    options:{
      responsive:true,maintainAspectRatio:false,
      interaction:{mode:'index',intersect:false},
      plugins:{
        legend:{display:false},
        tooltip:{
          enabled:true,
          displayColors:false,
          backgroundColor:'rgba(5,19,24,.97)',
          borderColor:'rgba(89,228,217,.24)',
          borderWidth:1,
          titleColor:'#B18D9C',
          bodyColor:'#FFF2F7',
          callbacks:{
            title:items=>`Tendencia · ${predictionHoverLabel(items[0].dataIndex)}`,
            label:ctx=>`Índice: ${Number(ctx.parsed.y).toFixed(2)}`,
            afterLabel:ctx=>ctx.dataIndex>=11?'Zona proyectada':'Histórico ilustrativo'
          }
        }
      },
      hover:{mode:'index',intersect:false},
      onHover:(event,elements,chart)=>{
        if(!elements.length) return;
        const i=elements[0].index;
        const value=chart.data.datasets[0].data[i];
        const prev=i>0?chart.data.datasets[0].data[i-1]:value;
        $('#predHoverTime').textContent=predictionHoverLabel(i);
        $('#predHoverValue').textContent=Number(value).toFixed(2);
        const move=value-prev;
        const el=$('#predHoverMove');
        el.textContent=(move>=0?'▲ +':'▼ ')+Math.abs(move).toFixed(2);
        el.className=move>=0?'up':'down';
      },
      scales:{x:{grid:{display:false},border:{display:false},ticks:{font:{size:9},color:'#8E6F7B',maxTicksLimit:5}},
      y:{grid:{color:'rgba(120,170,170,.11)'},border:{display:false},ticks:{font:{size:9},color:'#8E6F7B'}}}
    }
  });

  portChart=new Chart($('#portChart'),{
    type:'doughnut',
    data:{labels:['Producto','Ventas','Operaciones','Analítica','Soporte'],
      datasets:[{
        data:[34,22,19,15,10],
        backgroundColor:[
          cssVar('--warning'),
          '#B58AF6',
          cssVar('--primary'),
          cssVar('--secondary'),
          cssVar('--muted-strong')
        ],
        borderWidth:0,
        hoverOffset:2
      }]},
    options:{
      responsive:true,
      maintainAspectRatio:true,
      aspectRatio:1,
      cutout:'72%',
      layout:{padding:2},
      plugins:{legend:{display:false}},
      onClick:(evt,els)=>{if(els[0]) selectPortfolio(els[0].index)}
    }
  });
}

function sentimentReading(v){
  if(v>=82) return 'Óptimo';
  if(v>=70) return 'Saludable';
  if(v>=55) return 'Atención';
  return 'Crítico';
}
function sentimentHoverLabel(i){
  const labels=['00:00','01:30','03:00','04:30','06:00','07:30','09:00','10:30','12:00','13:30','15:00','16:30','18:00'];
  return labels[i]||`Punto ${i+1}`;
}
function predictionHoverLabel(i){
  const total=16;
  const hour=Math.round((24*i)/total);
  return String(hour).padStart(2,'0')+':00';
}

function setPeriod(period){
  const d=periodData[period];
  $('#gaugeValue').textContent=d.value;
  $('#gaugeLabel').textContent=d.label;
  $('#gauge').style.background=`conic-gradient(var(--teal) 0 ${d.value}%, rgba(255,255,255,.07) ${d.value}% 100%)`;
  $('#aiCaption').innerHTML=`<i class="fa-solid fa-wand-magic-sparkles"></i>${d.caption}`;
  sentChart.data.datasets[0].data=d.series;
  sentChart.update();
  $$('#periodButtons button').forEach(b=>b.classList.toggle('active',b.dataset.period===period));
}

function selectAsset(symbol){
  state.selectedAsset=symbol;
  const a=assets[symbol];
  $$('.asset').forEach(el=>el.classList.toggle('selected',el.dataset.symbol===symbol));
  $('#predictionTitle').textContent=`Tendencia · ${a.name}`;
  $('#predCurrent').textContent=formatMoney(a.mxn);
  $('#predForecast').textContent=(a.forecast>=0?'+':'')+a.forecast.toFixed(2)+'%';
  $('#insightMomentum').textContent=a.momentum;
  $('#insightVol').textContent=a.vol;
  $('#insightRead').textContent=a.read;
  predChart.data.datasets[0].data=a.series;
  predChart.update();
  toast(`${a.name} seleccionado`,`El resumen y la tendencia se actualizaron.`);
}

function selectPortfolio(index){
  const rows=$$('#portfolioLegend .legend-row');
  rows.forEach((r,i)=>r.classList.toggle('selected',i===index));
  const labels=['Producto','Ventas','Operaciones','Analítica','Soporte'];
  const pct=[34,22,19,15,10][index];
  $('#portfolioCenter').innerHTML=`${labels[index]}<strong>${pct.toFixed(1)}%</strong><span>de capacidad</span>`;
}

function toggleDropdown(menu,button){
  const isHidden=menu.classList.contains('hidden');
  $$('.dropdown').forEach(d=>d.classList.add('hidden'));
  $$('#currencyBtn,#userBtn').forEach(b=>b.setAttribute('aria-expanded','false'));
  if(isHidden){menu.classList.remove('hidden');button.setAttribute('aria-expanded','true')}
}

function openDrawer(){
  renderNotifications();
  $('#notificationDrawer').classList.add('show');
  $('#backdrop').classList.add('show');
  $('#notificationDrawer').setAttribute('aria-hidden','false');
  document.body.classList.add('no-scroll');
}
function closeDrawer(){
  $('#notificationDrawer').classList.remove('show');
  $('#backdrop').classList.remove('show');
  $('#notificationDrawer').setAttribute('aria-hidden','true');
  document.body.classList.remove('no-scroll');
}
function renderNotifications(){
  const wrap=$('#notificationList');wrap.innerHTML='';
  state.notifications.forEach(n=>{
    const el=document.createElement('button');
    el.type='button';el.className='notification '+(n.unread?'unread':'');
    el.style.width='100%';el.style.color='inherit';el.style.textAlign='left';el.style.cursor='pointer';
    el.innerHTML=`<div class="notification-top"><div class="notification-icon"><i class="fa-solid ${n.icon}"></i></div><div><strong>${n.title}</strong><p>${n.text}</p><small>${n.time}</small></div></div>`;
    el.addEventListener('click',()=>{n.unread=false;renderNotifications();updateBadge();toast('Notificación leída',n.title)});
    wrap.appendChild(el);
  });
}
function updateBadge(){
  const count=state.notifications.filter(n=>n.unread).length;
  $('#notifBadge').textContent=count;
  $('#notifBadge').classList.toggle('hidden',count===0);
}

function openModal(title,subtitle,content,wide=false){
  const box=$('#modalBox');
  $('#modalTitle').textContent=title;
  $('#modalSubtitle').textContent=subtitle||'';
  $('#modalContent').innerHTML=content;
  box.classList.remove('modal-sm','modal-lg','modal-xl');
  box.classList.toggle('wide',wide);
  $('#modalLayer').classList.add('show');
  $('#modalLayer').setAttribute('aria-hidden','false');
  document.body.classList.add('no-scroll');
  enhanceCodeBlocks($('#modalContent'));
}

/* Reusable modal API for future modules */
function openUiModal(options={}){
  const {
    title='Modal',
    subtitle='',
    body='',
    size='md',
    icon='fa-window-maximize',
    footer=''
  }=options;
  const iconHtml=`<div class="modal-icon"><i class="fa-solid ${icon}"></i></div>`;
  openModal(title,subtitle,`${iconHtml}<div style="margin-top:12px">${body}</div>${footer?`<div class="modal-footer">${footer}</div>`:''}`,false);
  const box=$('#modalBox');
  if(['sm','lg','xl'].includes(size)) box.classList.add(`modal-${size}`);
}

async function copyText(value){
  if(navigator.clipboard && window.isSecureContext){
    await navigator.clipboard.writeText(value);
    return;
  }
  const ta=document.createElement('textarea');
  ta.value=value;
  ta.style.position='fixed';
  ta.style.opacity='0';
  document.body.appendChild(ta);
  ta.focus();
  ta.select();
  document.execCommand('copy');
  ta.remove();
}


function checkHighlightRuntime(){
  if(window.hljs){
    console.info('Axolotl: Highlight.js browser bundle cargado correctamente.');
    return true;
  }
  console.warn('Axolotl: Highlight.js no está disponible; se usará el detector de lenguaje interno sin coloreado avanzado.');
  return false;
}

function normalizeLanguageLabel(lang){
  const map={
    xml:'HTML',
    html:'HTML',
    javascript:'JavaScript',
    js:'JavaScript',
    typescript:'TypeScript',
    ts:'TypeScript',
    css:'CSS',
    scss:'SCSS',
    json:'JSON',
    bash:'Bash',
    shell:'Shell',
    python:'Python',
    sql:'SQL',
    plaintext:'Texto',
    text:'Texto',
    markdown:'Markdown',
    php:'PHP',
    java:'Java',
    csharp:'C#',
    cpp:'C++'
  };
  return map[String(lang||'').toLowerCase()] || (lang ? String(lang).toUpperCase() : 'Código');
}

function fallbackLanguageDetect(text){
  const src=String(text||'').trim();
  if(!src)return 'plaintext';
  if(/^<!doctype html>|^<html|<\/?[a-z][\s\S]*>/i.test(src))return 'xml';
  if(/^\s*[\{\[][\s\S]*[\}\]]\s*$/.test(src)){
    try{JSON.parse(src);return 'json'}catch(_){}
  }
  if(/\b(SELECT|INSERT|UPDATE|DELETE|CREATE TABLE|FROM|WHERE|JOIN)\b/i.test(src))return 'sql';
  if(/^\s*(def |import |from .+ import |print\(|class .+:)/m.test(src))return 'python';
  if(/^\s*(#!\/bin\/|npm |node |python |cd |echo |git )/m.test(src))return 'bash';
  if(/[.#][\w-]+\s*\{[\s\S]*:[^;]+;/.test(src))return 'css';
  if(/\b(const|let|var|function|=>|console\.log|document\.|addEventListener)\b/.test(src))return 'javascript';
  return 'plaintext';
}

function highlightCodeElement(code,forceLanguage=null){
  if(!code)return 'plaintext';

  /* Preserve the original source so the copy button always copies raw code. */
  if(!code.dataset.rawCode) code.dataset.rawCode=code.textContent;

  const raw=code.dataset.rawCode;
  const requested=forceLanguage || code.dataset.language || '';

  /* Remove previous Highlight.js bookkeeping before re-highlighting. */
  delete code.dataset.highlighted;
  code.className='';

  let detected=requested || fallbackLanguageDetect(raw);

  if(window.hljs){
    try{
      let result;
      if(requested && window.hljs.getLanguage(requested==='html'?'xml':requested)){
        const lang=requested==='html'?'xml':requested;
        result=window.hljs.highlight(raw,{language:lang,ignoreIllegals:true});
      }else{
        result=window.hljs.highlightAuto(raw);
      }
      code.innerHTML=result.value;
      detected=result.language || detected;
      code.classList.add('hljs');
      if(detected) code.classList.add(`language-${detected}`);
      code.dataset.highlighted='yes';
    }catch(err){
      code.textContent=raw;
      console.warn('Highlight.js no pudo colorear este bloque.',err);
    }
  }else{
    code.textContent=raw;
  }

  return detected;
}

function updateCodeLanguageBadge(shell,language){
  let badge=shell.querySelector(':scope > .code-language');
  if(!badge){
    badge=document.createElement('span');
    badge.className='code-language';
    shell.appendChild(badge);
  }
  badge.innerHTML=`<i class="fa-solid fa-code"></i><span>${normalizeLanguageLabel(language)}</span>`;
}

function enhanceCodeBlocks(root=document){
  $$('pre code',root).forEach(code=>{
    const pre=code.closest('pre');
    if(!pre)return;

    let shell=pre.parentElement;
    if(!shell.classList.contains('code-shell')){
      shell=document.createElement('div');
      shell.className='code-shell';
      pre.parentNode.insertBefore(shell,pre);
      shell.appendChild(pre);
    }

    const detected=highlightCodeElement(code,code.dataset.language||null);
    updateCodeLanguageBadge(shell,detected);

    if(pre.dataset.copyReady==='1') return;
    pre.dataset.copyReady='1';

    const btn=document.createElement('button');
    btn.type='button';
    btn.className='code-copy';
    btn.innerHTML='<i class="fa-regular fa-copy"></i><span>Copiar</span>';
    btn.addEventListener('click',async()=>{
      try{
        await copyText(code.dataset.rawCode || code.textContent);
        btn.classList.add('copied');
        btn.innerHTML='<i class="fa-solid fa-check"></i><span>Copiado</span>';
        setTimeout(()=>{
          btn.classList.remove('copied');
          btn.innerHTML='<i class="fa-regular fa-copy"></i><span>Copiar</span>';
        },1500);
      }catch(_){
        toast('No se pudo copiar','Selecciona el código manualmente.');
      }
    });
    shell.appendChild(btn);
  });
}
function closeModal(){
  $('#modalLayer').classList.remove('show');
  $('#modalLayer').setAttribute('aria-hidden','true');
  document.body.classList.remove('no-scroll');
}
function toast(title,text){
  const el=document.createElement('div');el.className='toast';
  el.innerHTML=`<i class="fa-solid fa-circle-check"></i><div><strong>${title}</strong><span>${text}</span></div>`;
  $('#toastWrap').appendChild(el);
  setTimeout(()=>el.remove(),3200);
}

function assetDetail(symbol){
  const a=assets[symbol];
  if(!a)return;
  openModal(`${a.name} · ${symbol}`,'Detalle de proyecto demostrativo.',
    `<div class="modal-grid">
      <div class="field"><label>Presupuesto</label><div class="card section" style="font-size:20px;font-weight:800">${formatMoney(a.mxn)}</div></div>
      <div class="field"><label>Variación</label><div class="card section" style="font-size:20px;font-weight:800;color:${a.change>=0?'var(--green)':'var(--red)'}">${a.change>=0?'+':''}${a.change.toFixed(2)}%</div></div>
      <div class="field"><label>Estado</label><div class="card section">${a.momentum}</div></div>
      <div class="field"><label>Complejidad</label><div class="card section">${a.vol}</div></div>
    </div>
    <div class="modal-actions">
      <button class="secondary" type="button" id="modalAnalyzeBtn">Ver lectura</button>
      <button class="primary" type="button" id="modalSelectBtn">Usar en tendencia</button>
    </div>`);
  $('#modalAnalyzeBtn').addEventListener('click',()=>toast('Lectura del proyecto',`Estado actual: ${a.read}.`));
  $('#modalSelectBtn').addEventListener('click',()=>{selectAsset(symbol);closeModal();smoothScrollTo($('#prediction-section'),{offset:-40,duration:.9})});
}

function allMarkets(){
  const rows=[
    ['Conversión','18.4%','+2.1 pts','Saludable'],
    ['Tiempo de respuesta','2m 18s','-14%','Mejorando'],
    ['Órdenes procesadas','1,842','+9.8%','En ritmo'],
    ['Backlog','37','+5','Atención']
  ].map(r=>`<tr><td><b>${r[0]}</b></td><td>${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td></tr>`).join('');
  openModal('Reporte operativo','Ejemplo de tabla reutilizable para indicadores.',`<table class="data-table"><thead><tr><th>Indicador</th><th>Valor</th><th>Cambio</th><th>Estado</th></tr></thead><tbody>${rows}</tbody></table>`,true);
}

function tradeModal(symbol){
  const a=assets[symbol];
  if(!a)return;
  openModal(`Nueva acción · ${a.name}`,'Ejemplo de formulario contextual reutilizable.',
    `<form id="tradeForm">
      <div class="modal-grid">
        <div class="field"><label>Responsable</label><input id="taskOwner" value="Usuario Axolotl"></div>
        <div class="field"><label>Prioridad</label><select id="taskPriority"><option>Media</option><option>Alta</option><option>Baja</option></select></div>
      </div>
      <div class="field" style="margin-top:10px"><label>Descripción</label><input id="taskDescription" value="Dar seguimiento a ${a.name}"></div>
      <div class="modal-actions"><button class="secondary" type="button" id="cancelTrade">Cancelar</button><button class="primary" type="submit">Crear acción</button></div>
    </form>`);
  $('#cancelTrade').addEventListener('click',closeModal);
  $('#tradeForm').addEventListener('submit',e=>{
    e.preventDefault();
    closeModal();
    toast('Acción creada',`Se agregó una acción para ${a.name}.`);
  });
}


function refreshThemeCharts(){
  if(typeof sentChart!=='undefined' && sentChart) sentChart.update('none');
  if(typeof predChart!=='undefined' && predChart) predChart.update('none');
  charts.forEach(c=>c.update('none'));
  if(typeof portChart!=='undefined' && portChart){
    portChart.data.datasets[0].backgroundColor=[
      cssVar('--warning'),'#B58AF6',cssVar('--primary'),cssVar('--secondary'),cssVar('--muted-strong')
    ];
    portChart.update('none');
  }
}

function themeModal(){
  const theme=getCurrentTheme();
  const labels={
    primary:'Primary',
    secondary:'Secondary',
    success:'Success',
    danger:'Danger',
    warning:'Warning',
    info:'Info'
  };
  openModal('Personalizar tema','Tokens semánticos globales, estilo Bootstrap.',
    `<div class="theme-grid">
      ${Object.keys(labels).map(k=>`<label class="theme-color-field">
        <input class="theme-token-input" type="color" data-token="${k}" value="${theme[k]}">
        <span><strong>${labels[k]}</strong><small>--${k}</small></span>
      </label>`).join('')}
    </div>
    <div class="theme-preview">
      <span class="theme-chip bg-primary">Primary · Branquias</span>
      <span class="theme-chip bg-secondary">Secondary · Cuerpo</span>
      <span class="theme-chip bg-success">Success</span>
      <span class="theme-chip bg-danger">Danger</span>
      <span class="theme-chip bg-warning">Warning</span>
      <span class="theme-chip bg-info">Info · Lavanda</span>
    </div>
    <div class="modal-callout" style="margin-top:10px">
      <b>Paleta Axolotl Rosado:</b> el rosa intenso proviene de las branquias de la referencia
      (<code>#FC779B</code>) y el rosa claro del cuerpo (<code>#FCBBD5</code>).
      Las superficies conservan un negro ciruela para mantener el modo oscuro.
    </div>
    <div class="modal-callout" style="margin-top:12px">
      Los componentes usan tokens semánticos. Cambiar <b>--primary</b>, <b>--secondary</b>, etc. recolorea la interfaz y las gráficas Neo/Prism.
    </div>
    <pre><code>&lt;div class="bg-primary text-muted rounded p-3"&gt;Contenido&lt;/div&gt;
&lt;button class="btn-ui btn-primary"&gt;Guardar&lt;/button&gt;
&lt;span class="text-success"&gt;Correcto&lt;/span&gt;</code></pre>
    <div class="modal-footer">
      <button class="btn-ui btn-ghost" type="button" id="resetThemeBtn">Restablecer</button>
      <button class="btn-ui btn-primary" type="button" id="saveThemeBtn">Guardar tema</button>
    </div>`,true);

  $$('.theme-token-input',$('#modalContent')).forEach(input=>{
    input.addEventListener('input',()=>{
      setSemanticColor(input.dataset.token,input.value);
      refreshThemeCharts();
    });
  });
  $('#resetThemeBtn').addEventListener('click',()=>{
    Object.entries(defaultTheme).forEach(([k,v])=>setSemanticColor(k,v));
    $$('.theme-token-input',$('#modalContent')).forEach(input=>input.value=defaultTheme[input.dataset.token]);
    refreshThemeCharts();
  });
  $('#saveThemeBtn').addEventListener('click',()=>{
    saveTheme();
    closeModal();
    toast('Tema guardado','Los tokens de color se conservarán en este navegador.');
  });
}

function componentsModal(){
  openModal('Componentes UI','Ejemplos reutilizables que siguen el design system.',
    `<div class="modal-callout">
      Usa estos componentes como base. Todos dependen de tokens semánticos y conservan la interfaz Axolotl.
    </div>

    <div class="card section" style="margin-top:12px">
      <div class="card-title">Modal estándar</div>
      <p class="text-muted" style="font-size:10px;line-height:1.55">Componente para confirmaciones, formularios, detalles y flujos secundarios.</p>
      <button class="btn-ui btn-primary" type="button" id="demoModalBtn">Abrir ejemplo</button>
      <pre><code>openUiModal({
  title: 'Nuevo elemento',
  subtitle: 'Descripción breve',
  icon: 'fa-plus',
  size: 'sm',
  body: '&lt;div class="modal-callout"&gt;Contenido&lt;/div&gt;',
  footer: '&lt;button class="btn-ui btn-primary"&gt;Guardar&lt;/button&gt;'
});</code></pre>
    </div>

    <div class="card section" style="margin-top:12px">
      <div class="card-title">Compatibilidad local</div>
      <p class="text-muted" style="font-size:10px;line-height:1.55">
        El smooth scroll incluye un motor RAF interno para funcionar incluso al abrir <code>index.html</code> mediante <code>file://</code>. Para desarrollo completo sigue siendo recomendable servir la carpeta mediante localhost.
      </p>
      <pre><code>python -m http.server 8000
# Después abre:
# http://localhost:8000</code></pre>
    </div>

    <div class="card section" style="margin-top:12px">
      <div class="card-title">Skeletons de carga</div>
      <p class="text-muted" style="font-size:10px;line-height:1.55">Úsalos mientras llegan datos de API/base de datos o mientras termina de cargar una imagen.</p>
      <div class="skeleton-news-row" style="margin-top:10px">
        <div class="skeleton" style="width:58px;height:46px"></div>
        <div>
          <div class="skeleton skeleton-line" style="width:92%"></div>
          <div class="skeleton skeleton-line" style="width:64%;margin-top:7px"></div>
          <div class="skeleton skeleton-line sm" style="width:28%;margin-top:8px"></div>
        </div>
      </div>
      <pre><code>await withSkeleton('#market-section', 'market', async () =&gt; {
  const data = await loadFromDatabase();
  renderMarket(data);
});

&lt;div class="skeleton-image-wrap"&gt;
  &lt;div class="image-skeleton-layer skeleton"&gt;&lt;/div&gt;
  &lt;img data-skeleton-image src="imagen.jpg" alt=""&gt;
&lt;/div&gt;</code></pre>
    </div>

    <div class="card section" style="margin-top:12px">
      <div class="card-title">Utilidades semánticas</div>
      <div class="theme-preview">
        <span class="theme-chip bg-primary">bg-primary</span>
        <span class="theme-chip bg-secondary">bg-secondary</span>
        <span class="theme-chip bg-success">bg-success</span>
        <span class="theme-chip bg-danger">bg-danger</span>
      </div>
      <pre><code>&lt;button class="btn-ui btn-primary"&gt;Primary&lt;/button&gt;
&lt;div class="bg-primary-soft border-primary rounded p-3"&gt;Panel&lt;/div&gt;
&lt;p class="text-danger"&gt;Mensaje de error&lt;/p&gt;</code></pre>
    </div>`,true);

  $('#demoModalBtn').addEventListener('click',()=>{
    openUiModal({
      title:'Nuevo elemento',
      subtitle:'Ejemplo del modal base Axolotl',
      icon:'fa-layer-group',
      size:'sm',
      body:`<div class="modal-callout">Este modal usa la misma superficie, borde, radio y tokens semánticos del dashboard.</div>
            <div class="field" style="margin-top:12px"><label>Nombre</label><input value="Elemento de ejemplo"></div>`,
      footer:`<button class="btn-ui btn-ghost" type="button" id="demoCancel">Cancelar</button>
              <button class="btn-ui btn-primary" type="button" id="demoSave">Guardar</button>`
    });
    $('#demoCancel').addEventListener('click',closeModal);
    $('#demoSave').addEventListener('click',()=>{closeModal();toast('Componente','Acción simulada correctamente.')});
  });
}

function settingsModal(){
  openModal('Preferencias','Configura el comportamiento visual de esta plantilla.',
    `<div class="switch-row"><div class="switch-copy"><strong>Modo compacto</strong><small>Reduce separaciones y altura de tarjetas.</small></div><label class="switch"><input id="compactToggle" type="checkbox" ${state.compact?'checked':''}><span class="slider"></span></label></div>
     <div class="switch-row"><div class="switch-copy"><strong>Animaciones de gráficas</strong><small>Activa o desactiva transiciones.</small></div><label class="switch"><input id="animToggle" type="checkbox" ${state.animations?'checked':''}><span class="slider"></span></label></div>
     <div class="modal-callout" style="margin-top:12px">Los módulos del dashboard son independientes: puedes reemplazar proyectos, tareas, agenda o métricas sin cambiar el layout general.</div>
     <div class="modal-actions"><button class="secondary" type="button" id="resetSettings">Restablecer</button><button class="primary" type="button" id="saveSettings">Guardar cambios</button></div>`);
  $('#resetSettings').addEventListener('click',()=>{$('#compactToggle').checked=false;$('#animToggle').checked=true});
  $('#saveSettings').addEventListener('click',()=>{
    state.compact=$('#compactToggle').checked;
    state.animations=$('#animToggle').checked;
    document.body.classList.toggle('compact',state.compact);
    Chart.defaults.animation=state.animations?{duration:500}:false;
    closeModal();toast('Preferencias guardadas','Los cambios visuales se aplicaron.');
  });
}

function profileModal(){
  openModal('Mi perfil','Datos demostrativos editables.',
    `<form id="profileForm"><div class="modal-grid">
      <div class="field"><label>Nombre de usuario</label><input id="profileName" value="${state.user.name}"></div>
      <div class="field"><label>Correo</label><input id="profileEmail" type="email" value="${state.user.email}"></div>
      <div class="field"><label>Alias</label><input id="profileAlias" maxlength="3" value="${state.user.alias}"></div>
      <div class="field"><label>Rol</label><select id="profileRole"><option ${state.user.role==='Administrador'?'selected':''}>Administrador</option><option ${state.user.role==='Editor'?'selected':''}>Editor</option><option ${state.user.role==='Analista'?'selected':''}>Analista</option><option ${state.user.role==='Usuario'?'selected':''}>Usuario</option></select></div>
      <div class="field"><label>Zona horaria</label><select><option>America/Mexico_City</option><option>America/Cancun</option><option>UTC</option></select></div>
    </div><div class="modal-actions"><button class="primary" type="submit">Guardar perfil</button></div></form>`);
  $('#profileForm').addEventListener('submit',e=>{
    e.preventDefault();
    state.user.name=$('#profileName').value.trim()||'Usuario Axolotl';
    state.user.email=$('#profileEmail').value.trim()||'usuario@axolotl.demo';
    state.user.alias=($('#profileAlias').value.trim()||'AX').slice(0,3).toUpperCase();
    state.user.role=$('#profileRole').value;
    syncUserUI();
    closeModal();
    toast('Perfil actualizado','Nombre, alias y rol se actualizaron en la plantilla.');
  });
}

function securityModal(){
  openModal('Seguridad','Controles demostrativos del menú de usuario.',
   `<div class="switch-row"><div class="switch-copy"><strong>Autenticación en dos pasos</strong><small>Simula el estado 2FA.</small></div><label class="switch"><input id="twoFA" type="checkbox" checked><span class="slider"></span></label></div>
    <div class="switch-row"><div class="switch-copy"><strong>Alertas de inicio de sesión</strong><small>Notificar nuevos accesos.</small></div><label class="switch"><input id="loginAlerts" type="checkbox" checked><span class="slider"></span></label></div>
    <div class="modal-actions"><button class="primary" type="button" id="saveSecurity">Guardar</button></div>`);
  $('#saveSecurity').addEventListener('click',()=>{closeModal();toast('Seguridad actualizada','Configuración guardada en la demo.')});
}

function alertsModal(){
  openDrawer();
}

function createAlertModal(){
  closeDrawer();
  openModal('Crear alerta','Configura una regla reutilizable para cualquier módulo.',
    `<form id="alertForm"><div class="modal-grid">
      <div class="field"><label>Módulo</label><select id="alertAsset">${Object.entries(assets).map(([s,a])=>`<option value="${s}">${a.name}</option>`).join('')}</select></div>
      <div class="field"><label>Condición</label><select id="alertCond"><option>Progreso menor a</option><option>Variación mayor a</option><option>Estado cambia a atención</option></select></div>
      <div class="field"><label>Umbral</label><input id="alertValue" type="number" min="0" step="1" value="70"></div>
      <div class="field"><label>Canal</label><select><option>Dashboard</option><option>Correo (demo)</option></select></div>
    </div><div class="modal-actions"><button class="primary" type="submit">Crear alerta</button></div></form>`);
  $('#alertForm').addEventListener('submit',e=>{
    e.preventDefault();
    const symbol=$('#alertAsset').value,cond=$('#alertCond').value,val=Number($('#alertValue').value);
    if(!Number.isFinite(val)||val<0){toast('Valor inválido','Revisa el umbral de la alerta.');return}
    const name=assets[symbol]?.name||symbol;
    state.notifications.unshift({id:Date.now(),title:`Alerta creada · ${name}`,text:`${cond}: ${val}.`,time:'Ahora',unread:true,icon:'fa-bell'});
    updateBadge();closeModal();toast('Alerta creada',`${name}: ${cond} ${val}.`);
  });
}

function allSignals(){
  const rows=[
    ['T-184','Validar flujo de onboarding','Mariana','72%','Hoy · 16:30'],
    ['T-207','Revisar segmentación de clientes','Diego','48%','Mañana · 10:00'],
    ['T-231','Publicar build interno','Sofía','61%','30 Sep'],
    ['T-249','Documentar modelo de eventos','Alex','34%','2 Oct']
  ].map(r=>`<tr><td><b>${r[0]}</b></td><td>${r[1]}</td><td>${r[2]}</td><td>${r[3]}</td><td>${r[4]}</td></tr>`).join('');
  openModal('Todas las tareas','Ejemplo de tabla de trabajo reutilizable.',`<table class="data-table"><thead><tr><th>ID</th><th>Tarea</th><th>Responsable</th><th>Avance</th><th>Vence</th></tr></thead><tbody>${rows}</tbody></table>`,true);
}

function newsModal(index){
  const n=news[index];
  openModal(n.title,n.time,`<div style="font-size:12px;line-height:1.7;color:#a8bfbb">${n.body}</div><div class="modal-actions"><button class="primary" type="button" id="newsDone">Cerrar</button></div>`);
  $('#newsDone').addEventListener('click',closeModal);
}

function allNews(){
  openModal('Toda la actividad','Eventos recientes de demostración.',
    news.map((n,i)=>`<button class="notification news-open" data-news="${i}" type="button" style="width:100%;color:inherit;text-align:left"><strong>${n.title}</strong><p>${n.body}</p><small>${n.time}</small></button>`).join(''),true);
  $$('.news-open',$('#modalContent')).forEach(b=>b.addEventListener('click',()=>newsModal(Number(b.dataset.news))));
}

function managePortfolio(){
  const labels=['Producto','Ventas','Operaciones','Analítica','Soporte'];
  const defaults=[34,22,19,15,10];
  openModal('Administrar capacidad','Edita la distribución demostrativa. El total debe ser 100%.',
    `<form id="portfolioForm">${labels.map((s,i)=>`<div class="field" style="margin-bottom:8px"><label>${s}</label><input class="port-input" data-index="${i}" type="number" min="0" max="100" step="0.1" value="${defaults[i]}"></div>`).join('')}
    <div id="portfolioTotal" style="font-size:11px;color:var(--muted)">Total: 100.0%</div>
    <div class="modal-actions"><button class="primary" type="submit">Aplicar distribución</button></div></form>`);
  const updateTotal=()=>{
    const t=$$('.port-input').reduce((s,x)=>s+(Number(x.value)||0),0);
    $('#portfolioTotal').textContent=`Total: ${t.toFixed(1)}%`;
    $('#portfolioTotal').style.color=Math.abs(t-100)<.01?'var(--green)':'var(--red)';
  };
  $$('.port-input').forEach(i=>i.addEventListener('input',updateTotal));
  $('#portfolioForm').addEventListener('submit',e=>{
    e.preventDefault();const vals=$$('.port-input').map(x=>Number(x.value)||0);const total=vals.reduce((a,b)=>a+b,0);
    if(Math.abs(total-100)>.01){toast('Distribución inválida','El total debe ser exactamente 100%.');return}
    portChart.data.datasets[0].data=vals;portChart.update();
    $$('#portfolioLegend .legend-row strong').forEach((el,i)=>el.textContent=vals[i].toFixed(1)+'%');
    closeModal();toast('Capacidad actualizada','La gráfica fue recalculada.');
  });
}

function aiStatus(){
  openModal('Estado de la aplicación','Resumen técnico de la plantilla.',
   `<div class="modal-grid">
     <div class="card section"><div style="color:var(--muted);font-size:10px">Aplicación</div><strong style="color:var(--green)">Operativa</strong></div>
     <div class="card section"><div style="color:var(--muted);font-size:10px">Base de datos</div><strong>SQLite conectada</strong></div>
     <div class="card section"><div style="color:var(--muted);font-size:10px">Backend</div><strong>Goravel</strong></div>
     <div class="card section"><div style="color:var(--muted);font-size:10px">Healthcheck</div><strong>/health</strong></div>
   </div>`);
}

function riskDetail(){
  openModal('Disponibilidad del sistema','Indicador general de salud de la aplicación.',
   `<div class="modal-grid">
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Disponibilidad</div><strong style="font-size:18px">99.97%</strong></div>
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Incidentes abiertos</div><strong style="font-size:18px">0</strong></div>
    </div>
    <div class="modal-actions"><button class="primary" type="button" id="riskSettingsBtn">Ver métricas</button></div>`);
  $('#riskSettingsBtn').addEventListener('click',()=>{closeModal();smoothScrollTo($('#server-metrics-section'),{offset:-30})});
}

function logout(){
  openModal('Cerrar sesión','Acción de demostración.',
    `<div style="font-size:12px;color:#C4A3B0;line-height:1.6">¿Quieres simular el cierre de sesión? No se perderá información del archivo.</div>
     <div class="modal-actions"><button class="secondary" type="button" id="logoutCancel">Cancelar</button><button class="primary" type="button" id="logoutConfirm">Cerrar sesión</button></div>`);
  $('#logoutCancel').addEventListener('click',closeModal);
  $('#logoutConfirm').addEventListener('click',()=>{closeModal();toast('Sesión cerrada','Demo: aquí redirigirías a tu pantalla de acceso.')});
}


/* Multipurpose search */
const searchLayer=$('#searchLayer');
const searchInput=$('#globalSearchInput');
const searchResults=$('#searchResults');
let searchFilter='all';
let searchActiveIndex=0;
let searchVisibleItems=[];

const searchIndex=[
  {type:'section',title:'Dashboard',desc:'Volver al centro de operaciones.',icon:'fa-house',action:()=>smoothScrollTo($('#dashboard'))},
  {type:'section',title:'Proyectos',desc:'Abrir proyectos destacados.',icon:'fa-diagram-project',action:()=>smoothScrollTo($('#summary-section'))},
  {type:'section',title:'Tareas',desc:'Abrir tareas prioritarias.',icon:'fa-list-check',action:()=>smoothScrollTo($('#signals-section'))},
  {type:'section',title:'Actividad',desc:'Abrir el feed de actividad.',icon:'fa-clock-rotate-left',action:()=>smoothScrollTo($('#news-section'))},
  {type:'section',title:'Rendimiento',desc:'Ver analítica operativa.',icon:'fa-chart-line',action:()=>smoothScrollTo($('#sentiment-section'))},
  {type:'section',title:'Distribución',desc:'Ver distribución de capacidad.',icon:'fa-chart-pie',action:()=>smoothScrollTo($('#portfolio-section'))},
  {type:'section',title:'Equipo',desc:'Ver carga y disponibilidad.',icon:'fa-users',action:()=>smoothScrollTo($('#team-section'))},
  {type:'section',title:'Agenda',desc:'Abrir próximos hitos.',icon:'fa-calendar',action:()=>smoothScrollTo($('#schedule-section'))},
  {type:'section',title:'Integraciones',desc:'Ver servicios conectables.',icon:'fa-plug',action:()=>smoothScrollTo($('#integrations-section'))},
  {type:'section',title:'Métricas del servidor',desc:'CPU, RAM, red y almacenamiento.',icon:'fa-server',action:()=>smoothScrollTo($('#server-metrics-section'))},
  {type:'section',title:'Laboratorio UI',desc:'Probar componentes y copiar ejemplos de código.',icon:'fa-code',action:()=>smoothScrollTo($('#code-lab-section'))},

  {type:'action',title:'Crear alerta',desc:'Configurar una regla de monitoreo.',icon:'fa-bell',action:createAlertModal},
  {type:'action',title:'Abrir notificaciones',desc:'Ver el centro de notificaciones.',icon:'fa-inbox',action:openDrawer},
  {type:'action',title:'Preferencias',desc:'Configurar apariencia y comportamiento.',icon:'fa-sliders',action:settingsModal},
  {type:'action',title:'Personalizar tema',desc:'Editar la paleta Axolotl.',icon:'fa-palette',action:themeModal},
  {type:'action',title:'Componentes UI',desc:'Ver componentes reutilizables.',icon:'fa-layer-group',action:componentsModal},
  {type:'action',title:'Mi perfil',desc:'Editar nombre, correo y alias.',icon:'fa-user',action:profileModal},
  {type:'action',title:'Seguridad',desc:'Abrir opciones de seguridad.',icon:'fa-shield-halved',action:securityModal},
  {type:'action',title:'Administrar capacidad',desc:'Editar la distribución general.',icon:'fa-chart-pie',action:managePortfolio},
  {type:'action',title:'Estado de la aplicación',desc:'Ver backend, base de datos y healthcheck.',icon:'fa-microchip',action:aiStatus},
  {type:'action',title:'Modo compacto',desc:'Abrir ajustes visuales.',icon:'fa-compress',action:settingsModal},

  ...Object.entries(assets).map(([symbol,a])=>({
    type:'asset',
    title:`${symbol} · ${a.name}`,
    desc:`Presupuesto ${formatMoney(a.mxn)} · variación ${a.change>=0?'+':''}${a.change.toFixed(2)}%`,
    icon:'fa-diagram-project',
    action:()=>assetDetail(symbol)
  })),

  ...news.map((n,i)=>({
    type:'news',
    title:n.title,
    desc:n.time,
    icon:'fa-clock-rotate-left',
    action:()=>newsModal(i)
  }))
];

function searchTypeLabel(type){
  return {section:'Sección',asset:'Proyecto',action:'Acción',news:'Actividad'}[type]||'Resultado';
}

function normalizeSearch(v){
  return String(v||'').toLocaleLowerCase('es').normalize('NFD').replace(/[\u0300-\u036f]/g,'');
}

function openSearch(prefill=''){
  $$('.dropdown').forEach(d=>d.classList.add('hidden'));
  searchLayer.classList.add('show');
  searchLayer.setAttribute('aria-hidden','false');
  document.body.classList.add('no-scroll');
  searchInput.value=prefill;
  searchActiveIndex=0;
  renderSearchResults();
  requestAnimationFrame(()=>searchInput.focus());
}

function closeSearch(){
  searchLayer.classList.remove('show');
  searchLayer.setAttribute('aria-hidden','true');
  document.body.classList.remove('no-scroll');
}

function filteredSearchItems(){
  const q=normalizeSearch(searchInput.value.trim());
  let items=searchIndex.filter(item=>searchFilter==='all'||item.type===searchFilter);
  if(q){
    items=items.filter(item=>{
      const hay=normalizeSearch(`${item.title} ${item.desc} ${searchTypeLabel(item.type)}`);
      return hay.includes(q);
    });
  }
  return items.slice(0,14);
}

function renderSearchResults(){
  const q=searchInput.value.trim();
  searchVisibleItems=filteredSearchItems();
  searchResults.innerHTML='';

  if(!q){
    const recent=document.createElement('div');
    recent.className='search-group-title';
    recent.textContent='Accesos rápidos';
    searchResults.appendChild(recent);
  }

  if(!searchVisibleItems.length){
    const answer=document.createElement('div');
    answer.className='search-answer';
    const safe=q.replace(/[<>&"]/g,m=>({'<':'&lt;','>':'&gt;','&':'&amp;','"':'&quot;'}[m]));
    answer.innerHTML=`<div class="search-answer-head"><i class="fa-solid fa-wand-magic-sparkles"></i> Respuesta simulada</div>
      <p>No encontré una coincidencia exacta para <b>${safe||'tu búsqueda'}</b>. En una implementación real, aquí puedes conectar búsqueda global, API, base de datos o un asistente. Como demo, puedes buscar “proyecto”, “tareas”, “equipo”, “servidor” o “integraciones”.</p>`;
    searchResults.appendChild(answer);
    return;
  }

  searchVisibleItems.forEach((item,index)=>{
    const btn=document.createElement('button');
    btn.type='button';
    btn.className='search-result'+(index===searchActiveIndex?' active':'');
    btn.innerHTML=`<span class="search-result-icon"><i class="fa-solid ${item.icon}"></i></span>
      <span class="search-result-copy"><strong>${item.title}</strong><small>${item.desc}</small></span>
      <span class="search-result-meta">${searchTypeLabel(item.type)} <i class="fa-solid fa-arrow-right"></i></span>`;
    btn.addEventListener('mouseenter',()=>{searchActiveIndex=index;updateSearchActive()});
    btn.addEventListener('click',()=>runSearchItem(index));
    searchResults.appendChild(btn);
  });

  if(q){
    const simulated=document.createElement('div');
    simulated.className='search-answer';
    simulated.innerHTML=`<div class="search-answer-head"><i class="fa-solid fa-wand-magic-sparkles"></i> Resumen de búsqueda</div>
      <p>Se encontraron <b>${searchVisibleItems.length}</b> coincidencias para “${q.replace(/[<>&"]/g,'')}”. Selecciona un resultado para ejecutar su acción.</p>`;
    searchResults.appendChild(simulated);
  }
}

function updateSearchActive(){
  $$('.search-result',searchResults).forEach((el,i)=>el.classList.toggle('active',i===searchActiveIndex));
  const active=$$('.search-result',searchResults)[searchActiveIndex];
  if(active) active.scrollIntoView({block:'nearest'});
}

function runSearchItem(index){
  const item=searchVisibleItems[index];
  if(!item) return;
  closeSearch();
  setTimeout(()=>item.action(),30);
  toast('Búsqueda',`Abriendo: ${item.title}`);
}

$('#searchBtn').addEventListener('click',()=>openSearch());
searchInput.addEventListener('input',()=>{searchActiveIndex=0;renderSearchResults()});
searchInput.addEventListener('keydown',e=>{
  if(e.key==='ArrowDown'){
    e.preventDefault();
    searchActiveIndex=Math.min(searchActiveIndex+1,Math.max(0,searchVisibleItems.length-1));
    updateSearchActive();
  }else if(e.key==='ArrowUp'){
    e.preventDefault();
    searchActiveIndex=Math.max(searchActiveIndex-1,0);
    updateSearchActive();
  }else if(e.key==='Enter'){
    e.preventDefault();
    if(searchVisibleItems.length) runSearchItem(searchActiveIndex);
  }else if(e.key==='Escape'){
    e.preventDefault();
    closeSearch();
  }
});
$('#searchFilters').addEventListener('click',e=>{
  const btn=e.target.closest('[data-filter]');
  if(!btn)return;
  searchFilter=btn.dataset.filter;
  $$('#searchFilters .search-filter').forEach(x=>x.classList.toggle('active',x===btn));
  searchActiveIndex=0;
  renderSearchResults();
});
searchLayer.addEventListener('click',e=>{if(e.target===searchLayer)closeSearch()});
document.addEventListener('keydown',e=>{
  const editable=e.target.matches('input,textarea,select,[contenteditable="true"]');
  if((e.ctrlKey||e.metaKey)&&e.key.toLowerCase()==='k'){
    e.preventDefault();
    if(searchLayer.classList.contains('show'))closeSearch();else openSearch();
  }else if(e.key==='/'&&!editable&&!searchLayer.classList.contains('show')){
    e.preventDefault();
    openSearch();
  }
});

function syncUserUI(){
  $('#headerUserAvatar').textContent=state.user.alias.slice(0,2).toUpperCase();
  $('#headerUserName').textContent=state.user.name;
  $('#headerUserRole').textContent=state.user.role;
}



/* =========================================================
   MOTION SYSTEM — Lenis + AOS + file:// safe RAF fallback
   ========================================================= */
let axLenis=null;
let axScrollAnimation=null;

function easeInOutCubic(t){
  return t<.5
    ? 4*t*t*t
    : 1-Math.pow(-2*t+2,3)/2;
}

function cancelNativeSmooth(){
  if(axScrollAnimation){
    cancelAnimationFrame(axScrollAnimation);
    axScrollAnimation=null;
  }
  document.documentElement.style.scrollBehavior='';
}

/*
  Guaranteed smooth-scroll fallback.
  It does not depend on CDN libraries, CORS, modules or the page origin,
  so it also works when index.html is opened directly through file://.
*/
function rafSmoothScrollTo(top,duration=1150){
  cancelNativeSmooth();

  const start=window.scrollY || document.documentElement.scrollTop || 0;
  const max=Math.max(0,document.documentElement.scrollHeight-window.innerHeight);
  const target=Math.max(0,Math.min(max,Number(top)||0));

  if(Math.abs(target-start)<2) return;

  /*
    Important:
    The page has CSS `scroll-behavior:smooth` as a baseline fallback.
    While animating frame-by-frame we MUST temporarily force `auto`,
    otherwise every RAF call starts another native smooth-scroll animation
    and some browsers appear to jump instantly.
  */
  const root=document.documentElement;
  const previousInlineBehavior=root.style.scrollBehavior;
  root.style.scrollBehavior='auto';

  const startTime=performance.now();
  const total=Math.max(500,Number(duration)||1150);

  const finish=()=>{
    axScrollAnimation=null;
    root.style.scrollBehavior=previousInlineBehavior;
  };

  const frame=now=>{
    const elapsed=now-startTime;
    const progress=Math.min(1,elapsed/total);
    const eased=easeInOutCubic(progress);
    const y=start+(target-start)*eased;

    window.scrollTo(0,y);

    if(progress<1){
      axScrollAnimation=requestAnimationFrame(frame);
    }else{
      window.scrollTo(0,target);
      finish();
    }
  };

  axScrollAnimation=requestAnimationFrame(frame);
}

function elementScrollTop(el,offset=-10){
  if(!el)return 0;
  const rect=el.getBoundingClientRect();
  return (window.scrollY || document.documentElement.scrollTop || 0)+rect.top+offset;
}

function initSmoothMotion(){
  const reduce=window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  const runningFromFile=location.protocol==='file:';

  /*
    Lenis is used when it is available over a normal HTTP/HTTPS origin.
    On file:// we intentionally use the inline RAF engine instead so the
    dashboard is fully functional when opened directly from Desktop.
  */
  if(!reduce && !runningFromFile && window.Lenis){
    try{
      axLenis=new Lenis({
        duration:1.08,
        smoothWheel:true,
        wheelMultiplier:.92,
        touchMultiplier:1.0,
        easing:t=>easeInOutCubic(t)
      });

      const raf=time=>{
        if(axLenis) axLenis.raf(time);
        requestAnimationFrame(raf);
      };
      requestAnimationFrame(raf);
    }catch(err){
      axLenis=null;
      console.warn('Lenis no pudo inicializarse. Se usará el motor RAF interno.',err);
    }
  }else if(runningFromFile){
    console.info('Axolotl: ejecución file:// detectada. El smooth scroll usa el motor RAF interno; para evitar restricciones de origen del navegador usa localhost.');
  }

  /* Add entrance animation declaratively to primary dashboard cards. */
  const animated=[
    ...$$('.stats .card'),
    ...$$('.left-col > .card'),
    ...$$('.right-col > .card')
  ];
  animated.forEach((el,index)=>{
    if(!el.dataset.aos){
      el.dataset.aos='fade-up';
      el.dataset.aosDuration='420';
      el.dataset.aosDelay=String(Math.min((index%4)*45,135));
      el.dataset.aosOnce='true';
    }
  });

  if(window.AOS){
    try{
      AOS.init({
        offset:24,
        duration:420,
        easing:'ease-out-cubic',
        once:true,
        mirror:false,
        anchorPlacement:'top-bottom',
        disable:reduce
      });
    }catch(err){
      console.warn('AOS no pudo inicializarse.',err);
    }
  }
}

function smoothScrollTo(target,options={}){
  const durationMs=Math.max(650,(options.duration??1.0)*1100);
  const offset=options.offset??-10;

  if(target===0){
    if(axLenis){
      axLenis.scrollTo(0,{duration:options.duration??1.05});
    }else{
      rafSmoothScrollTo(0,durationMs);
    }
    return;
  }

  const el=typeof target==='string'?$(target):target;
  if(!el)return;

  if(axLenis){
    axLenis.scrollTo(el,{offset,duration:options.duration??1.0});
  }else{
    rafSmoothScrollTo(elementScrollTop(el,offset),durationMs);
  }
}

/* =========================================================
   SKELETON SYSTEM
   ========================================================= */
function skeletonMarkup(type){
  if(type==='stats'){
    return `<div class="skeleton-overlay skeleton-grid-4">
      ${Array.from({length:4},()=>`<div class="skeleton-stat">
        <div class="skeleton skeleton-line" style="width:44%"></div>
        <div class="skeleton skeleton-line lg" style="width:68%;margin-top:13px"></div>
        <div class="skeleton skeleton-line sm" style="width:30%;margin-top:10px"></div>
        <div class="skeleton skeleton-line sm" style="width:88%;margin-top:20px"></div>
      </div>`).join('')}
    </div>`;
  }

  if(type==='market'){
    return `<div class="skeleton-overlay skeleton-card-pad">
      <div class="skeleton-card-head">
        <div class="skeleton skeleton-line" style="width:46%"></div>
        <div class="skeleton skeleton-circle" style="width:15px;height:15px"></div>
      </div>
      <div class="skeleton-stack">
        ${Array.from({length:5},()=>`<div class="skeleton-row">
          <div class="skeleton skeleton-circle" style="width:30px;height:30px"></div>
          <div>
            <div class="skeleton skeleton-line" style="width:70%"></div>
            <div class="skeleton skeleton-line sm" style="width:42%;margin-top:7px"></div>
          </div>
          <div>
            <div class="skeleton skeleton-line" style="width:100%"></div>
            <div class="skeleton skeleton-line sm" style="width:62%;margin:7px 0 0 auto"></div>
          </div>
        </div>`).join('')}
      </div>
    </div>`;
  }

  if(type==='news'){
    return `<div class="skeleton-overlay skeleton-card-pad">
      <div class="skeleton-card-head">
        <div class="skeleton skeleton-line" style="width:34%"></div>
        <div class="skeleton skeleton-circle" style="width:15px;height:15px"></div>
      </div>
      ${Array.from({length:3},()=>`<div class="skeleton-news-row">
        <div class="skeleton" style="width:58px;height:46px"></div>
        <div>
          <div class="skeleton skeleton-line" style="width:94%"></div>
          <div class="skeleton skeleton-line" style="width:70%;margin-top:7px"></div>
          <div class="skeleton skeleton-line sm" style="width:28%;margin-top:8px"></div>
        </div>
      </div>`).join('')}
    </div>`;
  }

  /* Summary/cards default */
  return `<div class="skeleton-overlay skeleton-card-pad">
    <div class="skeleton-card-head">
      <div class="skeleton skeleton-line" style="width:32%"></div>
      <div class="skeleton skeleton-line sm" style="width:14%"></div>
    </div>
    <div class="skeleton-grid-4">
      ${Array.from({length:4},()=>`<div style="padding:10px;border:1px solid rgba(255,255,255,.035);border-radius:11px">
        <div class="skeleton skeleton-circle" style="width:30px;height:30px"></div>
        <div class="skeleton skeleton-line sm" style="width:52%;margin-top:10px"></div>
        <div class="skeleton skeleton-line" style="width:72%;margin-top:8px"></div>
        <div class="skeleton skeleton-line sm" style="width:35%;margin-top:7px"></div>
        <div class="skeleton skeleton-line sm" style="width:95%;margin-top:12px"></div>
      </div>`).join('')}
    </div>
  </div>`;
}

function showSkeleton(host,type){
  if(!host || host.querySelector(':scope > .skeleton-overlay')) return;
  host.classList.add('skeleton-host');
  host.insertAdjacentHTML('beforeend',skeletonMarkup(type));
}
function hideSkeleton(host){
  if(!host)return;
  const overlay=host.querySelector(':scope > .skeleton-overlay');
  if(!overlay)return;
  overlay.classList.add('is-done');
  setTimeout(()=>overlay.remove(),340);
}
function setLoading(hostOrSelector,loading=true,type='summary'){
  const host=typeof hostOrSelector==='string'?$(hostOrSelector):hostOrSelector;
  if(!host)return;
  loading?showSkeleton(host,type):hideSkeleton(host);
}

/* Convenience wrapper for a real DB/API request:
   await withSkeleton('#market-section','market', async()=> fetch(...)) */
async function withSkeleton(hostOrSelector,type,task){
  const host=typeof hostOrSelector==='string'?$(hostOrSelector):hostOrSelector;
  setLoading(host,true,type);
  try{
    return await task();
  }finally{
    setLoading(host,false,type);
  }
}

/* Image loading skeleton:
   <div class="skeleton-image-wrap">
     <div class="image-skeleton-layer skeleton"></div>
     <img data-skeleton-image src="...">
   </div>
*/
function initImageSkeletons(root=document){
  $$('img[data-skeleton-image]',root).forEach(img=>{
    let wrap=img.closest('.skeleton-image-wrap');
    if(!wrap){
      wrap=document.createElement('div');
      wrap.className='skeleton-image-wrap';
      img.parentNode.insertBefore(wrap,img);
      wrap.appendChild(img);
    }
    if(!wrap.querySelector('.image-skeleton-layer')){
      const layer=document.createElement('div');
      layer.className='image-skeleton-layer skeleton';
      wrap.insertBefore(layer,img);
    }
    const done=()=>wrap.classList.add('is-loaded');
    if(img.complete) done();
    else{
      img.addEventListener('load',done,{once:true});
      img.addEventListener('error',done,{once:true});
    }
  });
}

function initSkeletonDemo(){
  const reduce=window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  $$('[data-skeleton]').forEach(host=>showSkeleton(host,host.dataset.skeleton));

  /* Demo delay only. Replace this block with hideSkeleton() after real requests resolve. */
  const delay=reduce?120:680;
  setTimeout(()=>{
    $$('[data-skeleton]').forEach((host,index)=>{
      setTimeout(()=>hideSkeleton(host),index*80);
    });
    if(window.AOS) AOS.refresh();
  },delay);

  initImageSkeletons();
}



/* A manual wheel/touch gesture always wins over the programmatic RAF scroll. */
window.addEventListener('wheel',()=>cancelNativeSmooth(),{passive:true});
window.addEventListener('touchstart',()=>cancelNativeSmooth(),{passive:true});


/* =========================================================
   VISIBLE CODE LAB DEMOS
   ========================================================= */
function initCodeLab(){
  const lab=$('#code-lab-section');
  if(!lab)return;

  const syntaxSamples={
    javascript:`const dashboard = {
  name: 'Axolotl',
  theme: 'neo-prism',
  modules: ['charts', 'search', 'sidebar']
};

console.log(dashboard.name);`,
    python:`class Dashboard:
    def __init__(self, name):
        self.name = name

dashboard = Dashboard("Axolotl")
print(dashboard.name)`,
    sql:`SELECT
  user_id,
  COUNT(*) AS total_actions
FROM activity_log
WHERE created_at >= CURRENT_DATE
GROUP BY user_id
ORDER BY total_actions DESC;`,
    json:`{
  "name": "Axolotl",
  "theme": "neo-prism",
  "responsive": true,
  "modules": ["charts", "search", "sidebar"]
}`,
    html:`<article class="card section">
  <div class="card-title">Axolotl</div>
  <button class="btn-ui btn-primary">Guardar</button>
</article>`,
    css:`.card {
  background: var(--surface);
  border: 1px solid var(--border-soft);
  border-radius: var(--radius);
}

.text-primary {
  color: var(--primary);
}`,
    bash:`python -m http.server 8000
echo "Axolotl disponible en localhost"
git status`
  };

  /* Any semantic demo button produces visible feedback. */
  $$('[data-demo-toast]',lab).forEach(btn=>{
    btn.addEventListener('click',()=>{
      toast(`${btn.dataset.demoToast} action`,'Componente interactivo funcionando.');
    });
  });

  $('#codeLabModalDemo')?.addEventListener('click',()=>{
    openUiModal({
      title:'Nuevo registro',
      subtitle:'Modal reutilizable desde el Laboratorio UI',
      icon:'fa-plus',
      size:'sm',
      body:`<div class="modal-callout">Este componente hereda automáticamente colores, radios y superficies globales.</div>
            <div class="field" style="margin-top:12px">
              <label>Nombre</label>
              <input id="labModalName" value="Elemento de prueba">
            </div>`,
      footer:`<button class="btn-ui btn-ghost" type="button" id="labModalCancel">Cancelar</button>
              <button class="btn-ui btn-primary" type="button" id="labModalSave">Guardar</button>`
    });
    $('#labModalCancel')?.addEventListener('click',closeModal);
    $('#labModalSave')?.addEventListener('click',()=>{
      const name=$('#labModalName')?.value||'Elemento';
      closeModal();
      toast('Guardado',`${name} se guardó en la simulación.`);
    });
  });

  $('#codeLabConfirmDemo')?.addEventListener('click',()=>{
    openUiModal({
      title:'¿Confirmar acción?',
      subtitle:'Ejemplo de confirmación',
      icon:'fa-circle-question',
      size:'sm',
      body:`<div class="modal-callout">En una implementación real aquí puedes confirmar eliminaciones, publicaciones o cambios sensibles.</div>`,
      footer:`<button class="btn-ui btn-ghost" type="button" id="labConfirmNo">Cancelar</button>
              <button class="btn-ui btn-danger" type="button" id="labConfirmYes">Confirmar</button>`
    });
    $('#labConfirmNo')?.addEventListener('click',closeModal);
    $('#labConfirmYes')?.addEventListener('click',()=>{
      closeModal();
      toast('Acción confirmada','La confirmación de prueba se ejecutó.');
    });
  });

  $('#codeLabSkeletonDemo')?.addEventListener('click',async()=>{
    const target=$('#codeLabSkeletonTarget');
    if(!target)return;
    await withSkeleton(target,'summary',()=>new Promise(resolve=>setTimeout(resolve,1450)));
    toast('Datos cargados','El skeleton desapareció después de la carga simulada.');
  });

  $('#codeLabToastSuccess')?.addEventListener('click',()=>toast('Correcto','La operación terminó correctamente.'));
  $('#codeLabToastWarning')?.addEventListener('click',()=>toast('Revisión necesaria','Este es un aviso de demostración.'));
  $('#codeLabToastDanger')?.addEventListener('click',()=>toast('Error simulado','Ejemplo de feedback para una acción fallida.'));

  $$('[data-code-tabs]',lab).forEach(tabBar=>{
    $$('.code-tab',tabBar).forEach(tab=>{
      tab.addEventListener('click',()=>{
        $$('.code-tab',tabBar).forEach(x=>x.classList.toggle('active',x===tab));
        const card=tabBar.closest('.demo-card');
        $$('.code-pane',card).forEach(pane=>pane.classList.toggle('active',pane.dataset.pane===tab.dataset.codePane));
      });
    });
  });

  const syntaxSelect=$('#syntaxDemoSelect');
  const syntaxCode=$('#syntaxDemoCode');
  if(syntaxSelect && syntaxCode){
    syntaxSelect.addEventListener('change',()=>{
      const lang=syntaxSelect.value;
      syntaxCode.dataset.rawCode=syntaxSamples[lang];
      syntaxCode.dataset.language=lang;
      const detected=highlightCodeElement(syntaxCode,lang);
      const shell=syntaxCode.closest('.code-shell');
      if(shell) updateCodeLanguageBadge(shell,detected);
    });
  }

  /* Static code blocks are present at page load, so highlighting + copy buttons are visible immediately. */
  enhanceCodeBlocks(lab);
}



/* =========================================================
   SERVER METRICS DEMO
   ========================================================= */
let serverMetricsTimer=null;
let serverMetricsPaused=false;
let networkChart=null;

const gaugeDisplayed={cpu:0,ram:0};
const gaugeAnimationFrames={cpu:null,ram:null};
const coreDisplayed=Array(8).fill(0);
const coreAnimationFrames=Array(8).fill(null);

function clamp(n,min=0,max=100){
  return Math.min(max,Math.max(min,n));
}
function jitter(value,range,min=0,max=100){
  return clamp(value+(Math.random()-.5)*range,min,max);
}

function rawSetGauge(prefix,value){
  const v=clamp(value);
  const arc=$(`#${prefix}GaugeArc`);
  const needle=$(`#${prefix}GaugeNeedle`);
  const label=$(`#${prefix}GaugeValue`);
  const top=$(`#${prefix}TopValue`);

  if(arc){
    let length=251.33;
    try{length=arc.getTotalLength()||length}catch(_){}
    arc.style.strokeDasharray=String(length);
    arc.style.strokeDashoffset=String(length*(1-v/100));
  }

  /*
    SVG transform attribute is used instead of CSS transform.
    This is reliable in Chrome even when the file is opened through file://.
  */
  if(needle){
    const angle=-90+(v*1.8);
    needle.setAttribute('transform',`rotate(${angle} 105 105)`);
  }

  const rounded=Math.round(v);
  if(label) label.textContent=`${rounded}%`;
  if(top) top.textContent=`${rounded}%`;
  gaugeDisplayed[prefix]=v;
}

function animateGauge(prefix,target,duration=900){
  const to=clamp(target);
  const from=Number.isFinite(gaugeDisplayed[prefix])?gaugeDisplayed[prefix]:0;

  if(gaugeAnimationFrames[prefix]){
    cancelAnimationFrame(gaugeAnimationFrames[prefix]);
    gaugeAnimationFrames[prefix]=null;
  }

  const start=performance.now();

  const frame=now=>{
    const t=Math.min(1,(now-start)/duration);
    const eased=easeInOutCubic(t);
    rawSetGauge(prefix,from+(to-from)*eased);

    if(t<1){
      gaugeAnimationFrames[prefix]=requestAnimationFrame(frame);
    }else{
      rawSetGauge(prefix,to);
      gaugeAnimationFrames[prefix]=null;
    }
  };

  gaugeAnimationFrames[prefix]=requestAnimationFrame(frame);
}

function rawSetCore(index,value){
  const core=$(`.core-circle[data-core="${index}"]`);
  if(!core)return;
  const v=clamp(value);
  core.style.setProperty('--core-value',v.toFixed(2));
  const text=$('span',core);
  if(text) text.textContent=`${Math.round(v)}%`;
  coreDisplayed[index]=v;
}

function animateCore(index,target,duration=820){
  const to=clamp(target);
  const from=Number.isFinite(coreDisplayed[index])?coreDisplayed[index]:0;

  if(coreAnimationFrames[index]){
    cancelAnimationFrame(coreAnimationFrames[index]);
    coreAnimationFrames[index]=null;
  }

  const start=performance.now();
  const frame=now=>{
    const t=Math.min(1,(now-start)/duration);
    const eased=easeInOutCubic(t);
    rawSetCore(index,from+(to-from)*eased);

    if(t<1){
      coreAnimationFrames[index]=requestAnimationFrame(frame);
    }else{
      rawSetCore(index,to);
      coreAnimationFrames[index]=null;
    }
  };
  coreAnimationFrames[index]=requestAnimationFrame(frame);
}

function initNetworkChart(){
  const canvas=$('#networkChart');
  if(!canvas || networkChart)return;

  networkChart=new Chart(canvas,{
    type:'line',
    data:{
      labels:Array.from({length:24},(_,i)=>i),
      datasets:[
        {
          label:'Descarga',
          data:Array.from({length:24},()=>50+Math.random()*38),
          borderColor:ctx=>prismLineGradient(ctx.chart,cssVar('--primary'),cssVar('--secondary'),cssVar('--info')),
          backgroundColor:ctx=>prismAreaGradient(ctx.chart,cssVar('--primary-rgb'),cssVar('--secondary-rgb')),
          fill:true,
          borderWidth:1.9,
          pointRadius:0,
          pointHoverRadius:3,
          tension:.36,
          prismGlowColor:cssVar('--primary'),
          prismGlowBlur:6
        },
        {
          label:'Subida',
          data:Array.from({length:24},()=>12+Math.random()*18),
          borderColor:cssVar('--secondary'),
          backgroundColor:'transparent',
          borderWidth:1.45,
          pointRadius:0,
          pointHoverRadius:3,
          tension:.36,
          prismGlowColor:cssVar('--secondary'),
          prismGlowBlur:4
        }
      ]
    },
    options:{
      responsive:true,
      maintainAspectRatio:false,
      animation:{duration:760,easing:'easeInOutCubic'},
      interaction:{mode:'index',intersect:false},
      plugins:{
        legend:{
          display:true,
          align:'end',
          labels:{
            usePointStyle:true,
            pointStyle:'circle',
            boxWidth:6,
            boxHeight:6,
            color:cssVar('--muted'),
            font:{size:8}
          }
        },
        tooltip:{
          enabled:true,
          displayColors:true,
          backgroundColor:'rgba(18,13,20,.96)',
          borderColor:'rgba(252,119,155,.18)',
          borderWidth:1,
          callbacks:{
            label:ctx=>`${ctx.dataset.label}: ${ctx.parsed.y.toFixed(1)} Mbps`
          }
        }
      },
      scales:{
        x:{display:false},
        y:{
          beginAtZero:true,
          suggestedMax:130,
          grid:{color:'rgba(255,255,255,.035)'},
          border:{display:false},
          ticks:{
            color:cssVar('--muted'),
            font:{size:8},
            maxTicksLimit:5,
            callback:v=>`${v} Mbps`
          }
        }
      }
    }
  });

  requestAnimationFrame(()=>networkChart.resize());
}

const serverState={
  cpu:42,
  ram:68,
  cores:[38,51,44,29,63,47,34,41],
  download:84.2,
  upload:21.6,
  latency:18,
  packets:18.4,
  diskC:72,
  diskD:51,
  diskB:34
};

function renderServerTextMetrics(){
  $('#cpuTemp').textContent=`${Math.round(46+serverState.cpu*.28)}°C`;
  $('#cpuFreq').textContent=`${(3.4+serverState.cpu*.018).toFixed(1)} GHz`;
  $('#cpuProcesses').textContent=String(Math.round(165+serverState.cpu*.42));

  const usedRam=(32*serverState.ram/100);
  const ramSub=$('#ramGaugeValue')?.nextElementSibling;
  if(ramSub) ramSub.textContent=`${usedRam.toFixed(1)} / 32 GB`;
  $('#ramAvailable').textContent=`${(32-usedRam).toFixed(1)} GB`;
  $('#ramSwap').textContent=`${(0.6+serverState.ram*.011).toFixed(1)} GB`;

  $('#netDownload').textContent=`${serverState.download.toFixed(1)} Mbps`;
  $('#netUpload').textContent=`${serverState.upload.toFixed(1)} Mbps`;
  $('#netLatency').textContent=`${Math.round(serverState.latency)} ms`;
  $('#netPackets').textContent=`${serverState.packets.toFixed(1)}k/s`;

  $('#diskCFill').style.width=`${serverState.diskC}%`;
  $('#diskDFill').style.width=`${serverState.diskD}%`;
  $('#diskBFill').style.width=`${serverState.diskB}%`;
  $('#diskCText').textContent=`${Math.round(serverState.diskC)}% · ${Math.round(512*serverState.diskC/100)} / 512 GB`;
  $('#diskDText').textContent=`${Math.round(serverState.diskD)}% · ${(2*serverState.diskD/100).toFixed(1)} / 2 TB`;
  $('#diskBText').textContent=`${Math.round(serverState.diskB)}% · ${Math.round(2000*serverState.diskB/100)} GB / 2 TB`;
  $('#diskRead').textContent=`${Math.round(330+serverState.cpu*1.95)} MB/s`;
  $('#diskWrite').textContent=`${Math.round(125+serverState.cpu*1.48)} MB/s`;
}

function animateServerMetrics(duration=900){
  animateGauge('cpu',serverState.cpu,duration);
  animateGauge('ram',serverState.ram,duration);
  serverState.cores.forEach((v,i)=>animateCore(i,v,Math.max(620,duration-70)));
  renderServerTextMetrics();
}

function tickServerMetrics(){
  if(serverMetricsPaused)return;

  serverState.cpu=jitter(serverState.cpu,16,18,86);
  serverState.ram=jitter(serverState.ram,5,48,84);
  serverState.cores=serverState.cores.map((v,i)=>{
    const target=serverState.cpu+(i%2===0?-7:7);
    return clamp(v+(target-v)*.36+(Math.random()-.5)*16,8,96);
  });
  serverState.download=jitter(serverState.download,28,22,122);
  serverState.upload=jitter(serverState.upload,11,6,48);
  serverState.latency=jitter(serverState.latency,7,8,38);
  serverState.packets=jitter(serverState.packets,5,9,32);

  animateServerMetrics(1100);

  if(networkChart){
    const down=networkChart.data.datasets[0].data;
    const up=networkChart.data.datasets[1].data;
    down.push(serverState.download); down.shift();
    up.push(serverState.upload); up.shift();
    networkChart.update();
  }
}

function initServerMetrics(){
  if(!$('#server-metrics-section'))return;

  console.info('Axolotl: monitor de servidor iniciado.');
  initNetworkChart();

  /* Start at zero while the section is already visible. */
  rawSetGauge('cpu',0);
  rawSetGauge('ram',0);
  serverState.cores.forEach((_,i)=>rawSetCore(i,0));
  renderServerTextMetrics();

  /*
    Delay only enough for the first layout/paint. The user can now clearly
    see the needles, arcs and core rings sweep into their real values.
  */
  requestAnimationFrame(()=>{
    requestAnimationFrame(()=>{
      setTimeout(()=>animateServerMetrics(1350),180);
    });
  });

  if(serverMetricsTimer)clearInterval(serverMetricsTimer);
  serverMetricsTimer=setInterval(tickServerMetrics,2600);

  $('#serverPauseBtn')?.addEventListener('click',()=>{
    serverMetricsPaused=!serverMetricsPaused;
    const btn=$('#serverPauseBtn');
    const status=$('#serverStatus');

    if(serverMetricsPaused){
      btn.innerHTML='<i class="fa-solid fa-play"></i> Reanudar';
      status.textContent='Pausado';
      status.style.color='var(--warning)';
      status.style.borderColor='rgba(var(--warning-rgb),.18)';
      status.style.background='rgba(var(--warning-rgb),.055)';
    }else{
      btn.innerHTML='<i class="fa-solid fa-pause"></i> Pausar';
      status.textContent='Online';
      status.removeAttribute('style');
      tickServerMetrics();
    }
  });

  window.addEventListener('resize',()=>{
    if(networkChart) networkChart.resize();
  });
}



const operationDetails={
  conversion:{title:'Conversión',value:'18.4%',detail:'Objetivo mensual 20%',note:'Ejemplo de KPI comercial con comparación contra objetivo.'},
  response:{title:'Tiempo de respuesta',value:'2m 18s',detail:'Mejora de 14%',note:'Ejemplo de KPI de soporte o SLA.'},
  orders:{title:'Órdenes procesadas',value:'1,842',detail:'+9.8% vs. ayer',note:'Ejemplo de throughput para operaciones o logística.'},
  backlog:{title:'Backlog',value:'37',detail:'5 pendientes nuevos',note:'Ejemplo de inventario de trabajo pendiente.'}
};

function operationDetail(key){
  const item=operationDetails[key];
  if(!item)return;
  openModal(item.title,'Indicador operativo demostrativo.',
    `<div class="modal-grid">
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Presupuesto</div><strong style="font-size:20px">${item.value}</strong></div>
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Contexto</div><strong>${item.detail}</strong></div>
    </div>
    <div class="modal-callout" style="margin-top:12px">${item.note}</div>`);
}

const taskDetails={
  'T-184':{title:'Validar flujo de onboarding',owner:'Mariana',due:'Hoy · 16:30',progress:72},
  'T-207':{title:'Revisar segmentación de clientes',owner:'Diego',due:'Mañana · 10:00',progress:48},
  'T-231':{title:'Publicar build interno',owner:'Sofía',due:'30 Sep',progress:61},
  'T-249':{title:'Documentar modelo de eventos',owner:'Alex',due:'2 Oct',progress:34}
};

function taskDetail(id){
  const item=taskDetails[id];
  if(!item)return;
  openModal(`${id} · ${item.title}`,'Detalle de tarea demostrativa.',
    `<div class="modal-grid">
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Responsable</div><strong>${item.owner}</strong></div>
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Vence</div><strong>${item.due}</strong></div>
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Avance</div><strong>${item.progress}%</strong></div>
      <div class="card section"><div style="font-size:9px;color:var(--muted)">Estado</div><strong style="color:var(--primary)">En proceso</strong></div>
    </div>`);
}

function runTemplateAction(action){
  if(action==='new-record'){
    openUiModal({
      title:'Nuevo registro',
      subtitle:'Formulario base conectado al sistema de modales.',
      icon:'fa-plus',
      size:'sm',
      body:'<div class="field"><label>Nombre</label><input placeholder="Nuevo elemento"></div><div class="field" style="margin-top:10px"><label>Descripción</label><input placeholder="Descripción breve"></div>',
      footer:'<button class="btn-ui btn-ghost" type="button" id="quickCancel">Cancelar</button><button class="btn-ui btn-primary" type="button" id="quickSave">Guardar</button>'
    });
    $('#quickCancel')?.addEventListener('click',closeModal);
    $('#quickSave')?.addEventListener('click',()=>{closeModal();toast('Registro creado','Acción demostrativa completada.')});
    return;
  }
  if(action==='export'){toast('Exportación preparada','Aquí puedes conectar CSV, XLSX o PDF.');return}
  if(action==='report'){toast('Reporte','Aquí puedes conectar un controller o job de reportes.');return}
  if(action==='invite'){profileModal();return}
  if(action==='team'){toast('Equipo','Conecta esta acción a tu módulo de usuarios o permisos.');return}
  if(action==='calendar' || action==='schedule'){toast('Agenda','Conecta esta acción a eventos, reservas o recordatorios.');return}
  if(action==='integrations' || action==='integration'){toast('Integraciones','Conecta aquí webhooks, correo, storage u otros servicios.');return}
}

/* Events */
const appEl=$('#app');
const sidebarToggle=$('#sidebarToggle');
const mobileMenuBtn=$('#mobileMenuBtn');
const mobileSidebarBackdrop=$('#mobileSidebarBackdrop');

function isMobile(){return window.matchMedia('(max-width:780px)').matches}
function setSidebarCollapsed(collapsed){
  if(isMobile()) return;
  appEl.classList.toggle('sidebar-collapsed',collapsed);
  sidebarToggle.setAttribute('aria-label',collapsed?'Mostrar sidebar':'Ocultar sidebar');
  sidebarToggle.title=collapsed?'Mostrar sidebar':'Ocultar sidebar';
  if(collapsed){
    $$('.nav-group').forEach(group=>{
      const toggle=$('.nav-group-toggle',group);
      if(toggle) toggle.setAttribute('aria-expanded','false');
    });
  }
}
sidebarToggle.addEventListener('click',()=>{
  if(isMobile()){
    appEl.classList.remove('mobile-sidebar-open');
    document.body.classList.remove('no-scroll');
    return;
  }
  setSidebarCollapsed(!appEl.classList.contains('sidebar-collapsed'));
});
mobileMenuBtn.addEventListener('click',()=>{
  appEl.classList.add('mobile-sidebar-open');
  document.body.classList.add('no-scroll');
});
mobileSidebarBackdrop.addEventListener('click',()=>{
  appEl.classList.remove('mobile-sidebar-open');
  document.body.classList.remove('no-scroll');
});
window.addEventListener('resize',()=>{
  if(!isMobile()){
    appEl.classList.remove('mobile-sidebar-open');
    document.body.classList.remove('no-scroll');
  }
});

$('#currencyBtn').addEventListener('click',e=>{e.stopPropagation();toggleDropdown($('#currencyMenu'),$('#currencyBtn'))});
$('#userBtn').addEventListener('click',e=>{e.stopPropagation();toggleDropdown($('#userMenu'),$('#userBtn'))});
document.addEventListener('click',e=>{if(!e.target.closest('.top-actions')){$$('.dropdown').forEach(d=>d.classList.add('hidden'));$('#currencyBtn').setAttribute('aria-expanded','false');$('#userBtn').setAttribute('aria-expanded','false')}});

$$('#currencyMenu .drop-item').forEach(btn=>btn.addEventListener('click',()=>{
  state.currency=btn.dataset.currency;updateCurrency();selectAsset(state.selectedAsset);
  $('#currencyMenu').classList.add('hidden');toast('Moneda actualizada',`Ahora se muestran valores en ${state.currency}.`);
}));
$$('#userMenu .drop-item').forEach(btn=>btn.addEventListener('click',()=>{
  $('#userMenu').classList.add('hidden');
  const a=btn.dataset.userAction;
  if(a==='profile')profileModal();
  else if(a==='settings')settingsModal();
  else if(a==='theme')themeModal();
  else if(a==='components')componentsModal();
  else if(a==='security')securityModal();
  else if(a==='logout')logout();
}));

$('#notificationBtn').addEventListener('click',openDrawer);
$('#closeDrawerBtn').addEventListener('click',closeDrawer);
$('#backdrop').addEventListener('click',closeDrawer);
$('#markAllReadBtn').addEventListener('click',()=>{state.notifications.forEach(n=>n.unread=false);renderNotifications();updateBadge();toast('Todo leído','Se limpiaron las notificaciones pendientes.')});
$('#createAlertBtn').addEventListener('click',createAlertModal);

$('#modalCloseBtn').addEventListener('click',closeModal);
$('#modalLayer').addEventListener('click',e=>{if(e.target===$('#modalLayer'))closeModal()});
document.addEventListener('keydown',e=>{if(e.key==='Escape'){closeModal();closeDrawer();$$('.dropdown').forEach(d=>d.classList.add('hidden'))}});

function closeMobileSidebar(){
  if(!isMobile()) return;
  appEl.classList.remove('mobile-sidebar-open');
  document.body.classList.remove('no-scroll');
}

function setNavGroup(group,open){
  if(!group)return;
  group.classList.toggle('open',open);
  const toggle=$('.nav-group-toggle',group);
  if(toggle) toggle.setAttribute('aria-expanded',open?'true':'false');
}

function toggleNavGroup(group){
  if(!group)return;

  if(!isMobile() && appEl.classList.contains('sidebar-collapsed')){
    setSidebarCollapsed(false);
    requestAnimationFrame(()=>setNavGroup(group,true));
    return;
  }
  setNavGroup(group,!group.classList.contains('open'));
}

function runSidebarAction(action){
  if(action==='alerts' || action==='notifications'){openDrawer();return}
  if(action==='create-alert'){createAlertModal();return}
  if(action==='settings'){settingsModal();return}
  if(action==='theme'){themeModal();return}
  if(action==='components'){componentsModal();return}
  if(action==='security'){securityModal();return}
  if(action==='portfolio'){managePortfolio();return}
  if(action==='ai-status'){aiStatus();return}
}

function activateSidebarTarget(targetId){
  $$('.nav-item,.nav-subitem').forEach(x=>x.classList.remove('active'));

  const direct=$(`.nav-item[data-target="${targetId}"]`);
  if(direct){direct.classList.add('active');return}

  const child=$(`.nav-subitem[data-target="${targetId}"]`);
  if(child){
    child.classList.add('active');
    const group=child.closest('.nav-group');
    if(group){
      const parent=$('.nav-group-toggle',group);
      if(parent) parent.classList.add('active');
      setNavGroup(group,true);
    }
  }
}

$$('.nav-group-toggle').forEach(btn=>btn.addEventListener('click',()=>{
  toggleNavGroup(btn.closest('.nav-group'));
}));

$$('.nav-item[data-target],.nav-subitem[data-target]').forEach(btn=>btn.addEventListener('click',()=>{
  const target=$('#'+btn.dataset.target);
  if(target){
    activateSidebarTarget(btn.dataset.target);
    smoothScrollTo(target,{offset:-10});
  }
  closeMobileSidebar();
}));

$$('.nav-item[data-action],.nav-subitem[data-action]').forEach(btn=>btn.addEventListener('click',()=>{
  $$('.nav-item,.nav-subitem').forEach(x=>x.classList.remove('active'));
  btn.classList.add('active');
  const group=btn.closest('.nav-group');
  if(group){
    const parent=$('.nav-group-toggle',group);
    if(parent) parent.classList.add('active');
    setNavGroup(group,true);
  }
  closeMobileSidebar();
  runSidebarAction(btn.dataset.action);
}));

$$('#periodButtons button').forEach(btn=>btn.addEventListener('click',()=>setPeriod(btn.dataset.period)));
$$('.asset').forEach(btn=>btn.addEventListener('click',()=>selectAsset(btn.dataset.symbol)));
$$('.operation-row').forEach(btn=>btn.addEventListener('click',()=>operationDetail(btn.dataset.operation)));
$$('.task-card').forEach(btn=>btn.addEventListener('click',()=>taskDetail(btn.dataset.task)));
$$('.news-item').forEach(btn=>btn.addEventListener('click',()=>newsModal(Number(btn.dataset.news))));
$$('#portfolioLegend .legend-row').forEach(btn=>btn.addEventListener('click',()=>selectPortfolio(Number(btn.dataset.index))));

$('#summaryDetailsBtn').addEventListener('click',()=>{
  const a=assets[state.selectedAsset];
  openModal(
    `Detalle del proyecto · ${state.selectedAsset}`,
    'Vista ampliada del elemento seleccionado en la plantilla multipropósito.',
    `<div class="modal-grid">
      <div class="card section">
        <div style="color:var(--muted);font-size:10px">Proyecto seleccionado</div>
        <strong style="display:block;margin-top:5px">${a.name} (${state.selectedAsset})</strong>
      </div>
      <div class="card section">
        <div style="color:var(--muted);font-size:10px">Presupuesto</div>
        <strong style="display:block;margin-top:5px">${formatMoney(a.mxn)}</strong>
      </div>
      <div class="card section">
        <div style="color:var(--muted);font-size:10px">Variación</div>
        <strong style="display:block;margin-top:5px;color:${a.change>=0?'var(--green)':'var(--red)'}">${a.change>=0?'+':''}${a.change.toFixed(2)}%</strong>
      </div>
      <div class="card section">
        <div style="color:var(--muted);font-size:10px">Lectura</div>
        <strong style="display:block;margin-top:5px">${a.read}</strong>
      </div>
    </div>
    <div style="margin-top:12px" class="card section">
      <div style="display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:10px">
        <div><span style="font-size:9px;color:var(--muted)">Estado</span><strong style="display:block;margin-top:4px">${a.momentum}</strong></div>
        <div><span style="font-size:9px;color:var(--muted)">Complejidad</span><strong style="display:block;margin-top:4px">${a.vol}</strong></div>
        <div><span style="font-size:9px;color:var(--muted)">Objetivo semanal</span><strong style="display:block;margin-top:4px;color:var(--teal)">${a.forecast>=0?'+':''}${a.forecast.toFixed(2)}%</strong></div>
      </div>
    </div>
    <div class="modal-actions">
      <button class="secondary" type="button" id="detailsCloseBtn">Cerrar</button>
      <button class="primary" type="button" id="detailsUseBtn">Ver tendencia</button>
    </div>`
  );
  $('#detailsCloseBtn').addEventListener('click',closeModal);
  $('#detailsUseBtn').addEventListener('click',()=>{
    selectAsset(state.selectedAsset);
    closeModal();
    smoothScrollTo($('#prediction-section'),{offset:-40,duration:.9});
  });
});
$('#allMarketsBtn').addEventListener('click',allMarkets);
$('#allSignalsBtn').addEventListener('click',allSignals);
$('#allNewsBtn').addEventListener('click',allNews);
$('#managePortfolioBtn').addEventListener('click',managePortfolio);
$('#aiStatusBtn')?.addEventListener('click',aiStatus);
$('#riskRing')?.addEventListener('click',riskDetail);
$$('.template-action').forEach(btn=>btn.addEventListener('click',()=>runTemplateAction(btn.dataset.templateAction)));

/* Tooltips */
const tt=$('#tooltip');
$$('[data-tooltip]').forEach(el=>{
  const show=()=>{
    tt.textContent=el.dataset.tooltip;tt.classList.add('show');
    const r=el.getBoundingClientRect();const w=tt.offsetWidth;const left=Math.min(window.innerWidth-w-10,Math.max(10,r.left+r.width/2-w/2));
    tt.style.left=left+'px';tt.style.top=(r.bottom+8)+'px';
  };
  el.addEventListener('mouseenter',show);el.addEventListener('focus',show);
  el.addEventListener('mouseleave',()=>tt.classList.remove('show'));el.addEventListener('blur',()=>tt.classList.remove('show'));
});

/* Demo scroll spy */
const tracked=[
  'dashboard','summary-section','signals-section','news-section',
  'sentiment-section','prediction-section','portfolio-section',
  'team-section','schedule-section','integrations-section',
  'server-metrics-section','code-lab-section'
].map(id=>$('#'+id)).filter(Boolean);
const spy=new IntersectionObserver(entries=>{
  const visible=entries.filter(e=>e.isIntersecting).sort((a,b)=>b.intersectionRatio-a.intersectionRatio)[0];
  if(!visible)return;
  activateSidebarTarget(visible.target.id);
},{threshold:[.25,.5,.75]});
tracked.forEach(el=>spy.observe(el));


const scrollTopBtn=$('#scrollTopBtn');
window.addEventListener('scroll',()=>scrollTopBtn.classList.toggle('show',window.scrollY>420),{passive:true});
scrollTopBtn.addEventListener('click',()=>{
  scrollTopBtn.classList.add('scrolling');
  smoothScrollTo(0,{duration:1.25});
  setTimeout(()=>scrollTopBtn.classList.remove('scrolling'),1450);
});

loadTheme();
initCharts();
initSmoothMotion();
initSkeletonDemo();
checkHighlightRuntime();
enhanceCodeBlocks(document);
initCodeLab();
initServerMetrics();
renderNotifications();
updateBadge();
updateCurrency();
syncUserUI();
