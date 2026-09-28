
(() => {
  const explorer = document.getElementById('fileExplorer');
  if (!explorer) return;

  const csrf = document.querySelector('meta[name="csrf-token"]')?.content || '';
  const rootLabel = explorer.dataset.rootLabel || 'Demo';

  const els = {
    back: document.getElementById('fileBackBtn'),
    up: document.getElementById('fileUpBtn'),
    refresh: document.getElementById('fileRefreshBtn'),
    breadcrumbs: document.getElementById('fileBreadcrumbs'),
    quickFolders: document.getElementById('fileQuickFolders'),
    newFolder: document.getElementById('newFolderBtn'),
    newFile: document.getElementById('newFileBtn'),
    delete: document.getElementById('deleteFileBtn'),
    list: document.getElementById('fileList'),
    listView: document.getElementById('fileListView'),
    empty: document.getElementById('fileEmpty'),
    editorView: document.getElementById('fileEditorView'),
    editorName: document.getElementById('editorFileName'),
    editorLanguage: document.getElementById('editorLanguage'),
    editorInput: document.getElementById('editorInput'),
    editorWrap: document.getElementById('editorCodeWrap'),
    editorHighlight: document.getElementById('editorHighlight'),
    editorHighlightCode: document.getElementById('editorHighlightCode'),
    editorLines: document.getElementById('editorLineNumbers'),
    editorSave: document.getElementById('editorSaveBtn'),
    editorClose: document.getElementById('editorCloseBtn'),
    editorDirty: document.getElementById('editorDirtyDot'),
    editorPath: document.getElementById('editorPath'),
    editorStatus: document.getElementById('editorStatus'),
    count: document.getElementById('fileStatusCount'),
    selected: document.getElementById('fileStatusSelected'),
    rootButton: document.querySelector('[data-file-root]')
  };

  const state = {
    path: '',
    loaded: false,
    history: [],
    entries: [],
    selected: null,
    rootFolders: [],
    currentFile: null,
    originalContent: '',
    dirty: false
  };

  function showToast(title, message) {
    if (typeof toast === 'function') toast(title, message);
  }

  async function api(url, options) {
    const requestOptions = options || {};
    const response = await fetch(url, {
      credentials: 'same-origin',
      ...requestOptions,
      headers: {
        Accept: 'application/json',
        ...(requestOptions.headers || {})
      }
    });

    const payload = await response.json().catch(() => ({}));
    if (!response.ok) {
      throw new Error(payload.error || ('Error HTTP ' + response.status));
    }
    return payload;
  }

  async function postForm(url, values) {
    const body = new URLSearchParams();
    Object.entries(values).forEach(([key, value]) => body.set(key, String(value ?? '')));

    return api(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
        'X-CSRF-TOKEN': csrf
      },
      body
    });
  }

  function parentPath(path) {
    if (!path) return '';
    const parts = path.split('/').filter(Boolean);
    parts.pop();
    return parts.join('/');
  }

  function joinPath(parent, name) {
    return [parent, name].filter(Boolean).join('/');
  }

  function formatBytes(bytes) {
    if (!Number.isFinite(bytes) || bytes <= 0) return bytes === 0 ? '0 B' : '—';
    const units = ['B', 'KB', 'MB', 'GB'];
    let value = bytes;
    let index = 0;
    while (value >= 1024 && index < units.length - 1) {
      value /= 1024;
      index += 1;
    }
    return (value >= 10 || index === 0 ? value.toFixed(0) : value.toFixed(1)) + ' ' + units[index];
  }

  function formatDate(value) {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '—';
    return new Intl.DateTimeFormat('es-MX', {
      dateStyle: 'short',
      timeStyle: 'short'
    }).format(date);
  }

  function fileIcon(entry) {
    if (entry.type === 'folder') return 'fa-regular fa-folder';
    const ext = (entry.ext || '').toLowerCase();
    if (['go', 'js', 'ts', 'tsx', 'html', 'css', 'json', 'py', 'sql', 'sh', 'yaml', 'yml', 'php', 'java', 'c', 'h', 'cpp', 'hpp', 'cs', 'rs', 'rb', 'toml'].includes(ext)) {
      return 'fa-regular fa-file-code';
    }
    if (['md', 'txt', 'log', 'csv'].includes(ext)) return 'fa-regular fa-file-lines';
    return 'fa-regular fa-file';
  }

  function typeLabel(entry) {
    if (entry.type === 'folder') return 'Carpeta';
    const ext = (entry.ext || '').toUpperCase();
    return ext ? ('Archivo ' + ext) : 'Archivo';
  }

  function escapeHtml(value) {
    return String(value ?? '')
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  }

  function escapeAttr(value) {
    return escapeHtml(value);
  }

  function updateToolbar() {
    els.back.disabled = state.history.length === 0;
    els.up.disabled = state.path === '';
    els.delete.disabled = !state.selected;
    els.count.textContent = state.entries.length + ' ' + (state.entries.length === 1 ? 'elemento' : 'elementos');
    els.selected.textContent = state.selected
      ? (state.selected.name + ' · ' + (state.selected.type === 'folder' ? 'Carpeta' : formatBytes(state.selected.size)))
      : 'Ningún elemento seleccionado';
  }

  function renderBreadcrumbs() {
    const parts = state.path.split('/').filter(Boolean);
    const nodes = [];

    nodes.push(
      '<button class="file-crumb ' + (parts.length === 0 ? 'current' : '') + '" type="button" data-path="">' +
      '<i class="fa-regular fa-folder-open"></i><span>' + escapeHtml(rootLabel) + '</span></button>'
    );

    let current = '';
    parts.forEach((part, index) => {
      current = joinPath(current, part);
      nodes.push('<span class="file-crumb-sep"><i class="fa-solid fa-chevron-right"></i></span>');
      nodes.push(
        '<button class="file-crumb ' + (index === parts.length - 1 ? 'current' : '') +
        '" type="button" data-path="' + escapeAttr(current) + '">' + escapeHtml(part) + '</button>'
      );
    });

    els.breadcrumbs.innerHTML = nodes.join('');
    els.breadcrumbs.querySelectorAll('[data-path]').forEach(button => {
      button.addEventListener('click', () => requestEditorClose(() => loadDirectory(button.dataset.path, true)));
    });
  }

  function renderQuickFolders() {
    if (!state.rootFolders.length) {
      els.quickFolders.innerHTML = '<div class="file-nav-title">Sin carpetas</div>';
      return;
    }

    els.quickFolders.innerHTML = state.rootFolders.map(entry =>
      '<button class="file-nav-item" type="button" data-quick-path="' + escapeAttr(entry.path) + '">' +
      '<i class="fa-regular fa-folder"></i><span>' + escapeHtml(entry.name) + '</span></button>'
    ).join('');

    els.quickFolders.querySelectorAll('[data-quick-path]').forEach(button => {
      button.addEventListener('click', () => requestEditorClose(() => loadDirectory(button.dataset.quickPath, true)));
    });
  }

  function renderList() {
    els.list.innerHTML = '';
    els.empty.classList.toggle('hidden', state.entries.length !== 0);

    state.entries.forEach(entry => {
      const row = document.createElement('button');
      row.type = 'button';
      row.className = 'file-row';
      row.dataset.path = entry.path;
      row.dataset.type = entry.type;
      row.setAttribute('role', 'option');
      row.setAttribute('aria-selected', 'false');

      row.innerHTML =
        '<span class="file-name-cell">' +
          '<i class="' + fileIcon(entry) + '"></i>' +
          '<span>' + escapeHtml(entry.name) + '</span>' +
        '</span>' +
        '<span class="file-meta-cell">' + typeLabel(entry) + '</span>' +
        '<span class="file-meta-cell">' + (entry.type === 'folder' ? '—' : formatBytes(entry.size)) + '</span>' +
        '<span class="file-meta-cell">' + formatDate(entry.modified) + '</span>';

      row.addEventListener('click', () => selectEntry(entry, row));
      row.addEventListener('dblclick', () => openEntry(entry));
      els.list.appendChild(row);
    });

    updateToolbar();
  }

  function selectEntry(entry, row) {
    state.selected = entry;
    els.list.querySelectorAll('.file-row').forEach(item => {
      const selected = item === row;
      item.classList.toggle('selected', selected);
      item.setAttribute('aria-selected', selected ? 'true' : 'false');
    });
    updateToolbar();
  }

  async function loadDirectory(path, pushHistory) {
    const nextPath = path || '';
    const shouldPush = pushHistory !== false;

    try {
      const payload = await api('/api/files?path=' + encodeURIComponent(nextPath));

      if (shouldPush && state.loaded && nextPath !== state.path) {
        state.history.push(state.path);
      }

      state.path = payload.path || nextPath;
      state.entries = Array.isArray(payload.entries) ? payload.entries : [];
      state.selected = null;
      state.loaded = true;

      if (state.path === '') {
        state.rootFolders = state.entries.filter(entry => entry.type === 'folder');
        renderQuickFolders();
      }

      renderBreadcrumbs();
      renderList();
      els.rootButton?.classList.toggle('active', state.path === '');
      els.listView.classList.remove('hidden');
      els.editorView.classList.add('hidden');
    } catch (error) {
      showToast('No se pudo abrir la carpeta', error.message);
    }
  }

  async function openEntry(entry) {
    if (entry.type === 'folder') {
      requestEditorClose(() => loadDirectory(entry.path, true));
      return;
    }
    requestEditorClose(() => openFile(entry));
  }

  async function openFile(entry) {
    try {
      const file = await api('/api/files/read?path=' + encodeURIComponent(entry.path));
      state.currentFile = file;
      state.originalContent = file.content;
      state.dirty = false;

      els.editorName.textContent = file.name;
      els.editorLanguage.textContent = file.language || 'plaintext';
      els.editorPath.textContent = rootLabel + '/' + file.path;
      els.editorInput.value = file.content;
      els.editorStatus.textContent = 'UTF-8 · ' + formatBytes(file.size) + ' · Listo';

      updateEditorVisuals();
      setDirty(false);
      els.listView.classList.add('hidden');
      els.editorView.classList.remove('hidden');
      els.editorInput.focus();
    } catch (error) {
      showToast('No se pudo abrir el archivo', error.message);
    }
  }

  function languageForHighlight(language) {
    if (language === 'html') return 'xml';
    if (language === 'plaintext') return '';
    return language || '';
  }

  function updateEditorVisuals() {
    const editorContent = els.editorInput.value;
    const requestedLanguage = state.currentFile?.language || 'plaintext';
    const language = languageForHighlight(requestedLanguage);
    let highlighted = false;

    if (window.hljs) {
      try {
        if (language && window.hljs.getLanguage(language)) {
          const result = window.hljs.highlight(editorContent, {
            language,
            ignoreIllegals: true
          });
          els.editorHighlightCode.innerHTML = result.value + '\n';
          els.editorHighlightCode.className = 'hljs language-' + language;
          highlighted = true;
        } else if (editorContent && typeof window.hljs.highlightAuto === 'function') {
          const result = window.hljs.highlightAuto(editorContent);
          els.editorHighlightCode.innerHTML = result.value + '\n';
          els.editorHighlightCode.className = result.language
            ? 'hljs language-' + result.language
            : 'hljs';
          highlighted = true;
        }
      } catch (error) {
        console.warn('Axolotl: no se pudo colorear el archivo.', error);
      }
    }

    if (!highlighted) {
      els.editorHighlightCode.textContent = editorContent + '\n';
      els.editorHighlightCode.className = '';
    }

    els.editorWrap?.classList.toggle('highlight-ready', highlighted);

    const lineCount = Math.max(1, editorContent.split('\n').length);
    els.editorLines.textContent = Array.from({ length: lineCount }, (_, index) => index + 1).join('\n');
    syncEditorScroll();
  }

  function syncEditorScroll() {
    els.editorHighlight.scrollTop = els.editorInput.scrollTop;
    els.editorHighlight.scrollLeft = els.editorInput.scrollLeft;
    els.editorLines.scrollTop = els.editorInput.scrollTop;
  }

  function setDirty(value) {
    state.dirty = Boolean(value);
    els.editorDirty.classList.toggle('show', state.dirty);
    els.editorSave.disabled = !state.dirty;
    els.editorStatus.textContent = state.dirty
      ? 'UTF-8 · Cambios sin guardar'
      : ('UTF-8 · ' + formatBytes(state.currentFile?.size || 0) + ' · Guardado');
  }

  async function saveEditor() {
    if (!state.currentFile) return false;

    try {
      await postForm('/api/files/save', {
        path: state.currentFile.path,
        content: els.editorInput.value
      });

      state.originalContent = els.editorInput.value;
      state.currentFile.content = state.originalContent;
      state.currentFile.size = new Blob([state.originalContent]).size;
      setDirty(false);
      showToast('Archivo guardado', state.currentFile.name);
      return true;
    } catch (error) {
      showToast('No se pudo guardar', error.message);
      return false;
    }
  }

  function closeEditorWithoutPrompt() {
    state.currentFile = null;
    state.originalContent = '';
    state.dirty = false;
    els.editorInput.value = '';
    els.editorHighlightCode.textContent = '';
    els.editorView.classList.add('hidden');
    els.listView.classList.remove('hidden');
    setDirty(false);
  }

  function requestEditorClose(next) {
    const action = typeof next === 'function' ? next : closeEditorWithoutPrompt;

    if (!state.currentFile) {
      action();
      return;
    }

    if (!state.dirty) {
      if (action === closeEditorWithoutPrompt) {
        action();
      } else {
        closeEditorWithoutPrompt();
        action();
      }
      return;
    }

    if (typeof openUiModal !== 'function') {
      if (window.confirm('Hay cambios sin guardar. ¿Cerrar sin guardar?')) {
        closeEditorWithoutPrompt();
        action();
      }
      return;
    }

    openUiModal({
      title: 'Cambios sin guardar',
      subtitle: state.currentFile.name,
      icon: 'fa-triangle-exclamation',
      size: 'sm',
      body: '<div class="modal-callout">Este archivo tiene cambios que todavía no se han guardado. Puedes guardarlos antes de cerrar o descartarlos.</div>',
      footer:
        '<button class="btn-ui btn-ghost" type="button" id="editorCancelClose">Cancelar</button>' +
        '<button class="btn-ui btn-danger" type="button" id="editorDiscardClose">Cerrar sin guardar</button>' +
        '<button class="btn-ui btn-primary" type="button" id="editorSaveClose">Guardar y cerrar</button>'
    });

    document.getElementById('editorCancelClose')?.addEventListener('click', closeModal);
    document.getElementById('editorDiscardClose')?.addEventListener('click', () => {
      closeModal();
      if (action === closeEditorWithoutPrompt) {
        action();
      } else {
        closeEditorWithoutPrompt();
        action();
      }
    });
    document.getElementById('editorSaveClose')?.addEventListener('click', async () => {
      if (await saveEditor()) {
        closeModal();
        if (action === closeEditorWithoutPrompt) {
          action();
        } else {
          closeEditorWithoutPrompt();
          action();
        }
      }
    });
  }

  function openCreateModal(kind) {
    const isFolder = kind === 'folder';
    const placeholder = isFolder ? 'Nueva carpeta' : 'archivo.txt';

    openUiModal({
      title: isFolder ? 'Nueva carpeta' : 'Nuevo archivo',
      subtitle: state.path ? (rootLabel + '/' + state.path) : rootLabel,
      icon: isFolder ? 'fa-folder-plus' : 'fa-file-circle-plus',
      size: 'sm',
      body:
        '<div class="field">' +
          '<label>Nombre</label>' +
          '<input id="newItemName" maxlength="120" placeholder="' + placeholder + '" autocomplete="off">' +
        '</div>',
      footer:
        '<button class="btn-ui btn-ghost" type="button" id="newItemCancel">Cancelar</button>' +
        '<button class="btn-ui btn-primary" type="button" id="newItemCreate">Crear</button>'
    });

    const input = document.getElementById('newItemName');
    input?.focus();
    document.getElementById('newItemCancel')?.addEventListener('click', closeModal);

    const create = async () => {
      const name = input?.value.trim() || '';
      if (!name) {
        showToast('Nombre requerido', 'Escribe un nombre para continuar.');
        return;
      }
      try {
        const created = await postForm('/api/files/create', {
          parent: state.path,
          name,
          type: kind
        });
        closeModal();
        await loadDirectory(state.path, false);
        showToast(isFolder ? 'Carpeta creada' : 'Archivo creado', name);
        if (!isFolder) {
          await openFile(created);
        }
      } catch (error) {
        showToast('No se pudo crear', error.message);
      }
    };

    document.getElementById('newItemCreate')?.addEventListener('click', create);
    input?.addEventListener('keydown', event => {
      if (event.key === 'Enter') {
        event.preventDefault();
        create();
      }
    });
  }

  function confirmDeleteSelected() {
    if (!state.selected) return;
    const item = state.selected;
    const folderNotice = item.type === 'folder'
      ? '<p class="text-muted" style="margin:8px 0 0;font-size:10px">La carpeta y todo su contenido serán eliminados.</p>'
      : '';

    openUiModal({
      title: 'Eliminar elemento',
      subtitle: item.name,
      icon: 'fa-trash-can',
      size: 'sm',
      body: '<div class="modal-callout">Esta acción no se puede deshacer.' + folderNotice + '</div>',
      footer:
        '<button class="btn-ui btn-ghost" type="button" id="deleteCancel">Cancelar</button>' +
        '<button class="btn-ui btn-danger" type="button" id="deleteConfirm">Eliminar</button>'
    });

    document.getElementById('deleteCancel')?.addEventListener('click', closeModal);
    document.getElementById('deleteConfirm')?.addEventListener('click', async () => {
      try {
        await postForm('/api/files/delete', { path: item.path });
        closeModal();

        if (state.currentFile?.path === item.path) {
          closeEditorWithoutPrompt();
        }

        await loadDirectory(state.path, false);
        showToast('Elemento eliminado', item.name);
      } catch (error) {
        showToast('No se pudo eliminar', error.message);
      }
    });
  }

  els.rootButton?.addEventListener('click', () => requestEditorClose(() => loadDirectory('', true)));
  els.refresh.addEventListener('click', () => requestEditorClose(() => loadDirectory(state.path, false)));

  els.back.addEventListener('click', () => {
    if (!state.history.length) return;
    requestEditorClose(() => {
      const previous = state.history.pop();
      loadDirectory(previous, false);
    });
  });

  els.up.addEventListener('click', () => {
    if (!state.path) return;
    requestEditorClose(() => loadDirectory(parentPath(state.path), true));
  });

  els.newFolder.addEventListener('click', () => openCreateModal('folder'));
  els.newFile.addEventListener('click', () => openCreateModal('file'));
  els.delete.addEventListener('click', confirmDeleteSelected);
  els.editorSave.addEventListener('click', saveEditor);
  els.editorClose.addEventListener('click', () => requestEditorClose(closeEditorWithoutPrompt));

  els.editorInput.addEventListener('input', () => {
    setDirty(els.editorInput.value !== state.originalContent);
    updateEditorVisuals();
  });

  els.editorInput.addEventListener('scroll', syncEditorScroll);

  els.editorInput.addEventListener('keydown', event => {
    if (event.key === 'Tab') {
      event.preventDefault();
      const start = els.editorInput.selectionStart;
      const end = els.editorInput.selectionEnd;
      const value = els.editorInput.value;
      els.editorInput.value = value.slice(0, start) + '  ' + value.slice(end);
      els.editorInput.selectionStart = els.editorInput.selectionEnd = start + 2;
      els.editorInput.dispatchEvent(new Event('input', { bubbles: true }));
    }

    if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 's') {
      event.preventDefault();
      saveEditor();
    }
  });

  document.addEventListener('keydown', event => {
    if (event.key === 'Delete' && state.selected && !state.currentFile && !event.target.closest('input,textarea')) {
      confirmDeleteSelected();
    }
    if (event.key === 'Enter' && state.selected && !state.currentFile && !event.target.closest('input,textarea')) {
      openEntry(state.selected);
    }
  });

  window.addEventListener('beforeunload', event => {
    if (!state.dirty) return;
    event.preventDefault();
    event.returnValue = '';
  });

  loadDirectory('', false);
})();
