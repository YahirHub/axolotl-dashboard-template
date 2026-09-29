(() => {
  const page = document.getElementById('server-tasks-page');
  if (!page) return;

  const csrf = document.querySelector('meta[name="csrf-token"]')?.content || '';
  const list = document.getElementById('serverTaskList');
  const empty = document.getElementById('serverTaskEmpty');
  const count = document.getElementById('serverTaskCount');
  const createBtn = document.getElementById('newServerTaskBtn');
  const refreshBtn = document.getElementById('refreshServerTasksBtn');

  async function api(url, options = {}) {
    const response = await fetch(url, {
      credentials: 'same-origin',
      ...options,
      headers: {
        Accept: 'application/json',
        ...(options.headers || {})
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

  function showToast(title, message) {
    if (typeof toast === 'function') toast(title, message);
  }

  function formatBytes(bytes) {
    if (!Number.isFinite(bytes) || bytes <= 0) return bytes === 0 ? '0 B' : '—';
    if (bytes < 1024) return bytes + ' B';
    return (bytes / 1024).toFixed(bytes >= 10 * 1024 ? 0 : 1) + ' KB';
  }

  function formatDate(value) {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '—';
    return new Intl.DateTimeFormat('es-MX', {
      dateStyle: 'short',
      timeStyle: 'short'
    }).format(date);
  }

  async function loadTasks() {
    try {
      const payload = await api('/api/server/tasks');
      renderTasks(Array.isArray(payload.tasks) ? payload.tasks : []);
    } catch (error) {
      showToast('No se pudieron cargar las tareas', error.message);
    }
  }

  function renderTasks(tasks) {
    list.innerHTML = '';
    count.textContent = String(tasks.length);
    empty.classList.toggle('hidden', tasks.length > 0);

    tasks.forEach(task => {
      const row = document.createElement('article');
      row.className = 'server-task-row';

      const file = document.createElement('div');
      file.className = 'server-task-file';
      file.innerHTML = '<i class="fa-regular fa-file-code"></i><div><strong></strong><span></span></div>';
      file.querySelector('strong').textContent = task.filename;
      file.querySelector('span').textContent = 'Plantilla local · no instalada';

      const kind = document.createElement('span');
      kind.className = 'server-task-kind';
      kind.textContent = task.kind === 'timer' ? 'timer' : 'service';

      const size = document.createElement('span');
      size.className = 'server-task-meta';
      size.textContent = formatBytes(task.size);

      const modified = document.createElement('span');
      modified.className = 'server-task-meta';
      modified.textContent = formatDate(task.modified_at);

      const actions = document.createElement('div');
      actions.className = 'server-task-actions';

      const view = document.createElement('button');
      view.type = 'button';
      view.className = 'server-task-icon-btn';
      view.title = 'Ver plantilla';
      view.setAttribute('aria-label', 'Ver ' + task.filename);
      view.innerHTML = '<i class="fa-regular fa-eye"></i>';
      view.addEventListener('click', () => viewTask(task.filename));

      const remove = document.createElement('button');
      remove.type = 'button';
      remove.className = 'server-task-icon-btn danger';
      remove.title = 'Eliminar plantilla';
      remove.setAttribute('aria-label', 'Eliminar ' + task.filename);
      remove.innerHTML = '<i class="fa-regular fa-trash-can"></i>';
      remove.addEventListener('click', () => confirmDelete(task.filename));

      actions.append(view, remove);
      row.append(file, kind, size, modified, actions);
      list.appendChild(row);
    });
  }

  async function viewTask(filename) {
    try {
      const task = await api('/api/server/tasks/read?filename=' + encodeURIComponent(filename));

      openUiModal({
        title: task.filename,
        subtitle: 'Plantilla systemd · no instalada',
        icon: 'fa-file-code',
        size: 'lg',
        body: '<div class="server-task-preview"><pre><code id="serverTaskPreviewCode" data-language="ini"></code></pre></div>',
        footer: '<button class="btn-ui btn-primary" id="serverTaskPreviewClose" type="button">Cerrar</button>'
      });

      const code = document.getElementById('serverTaskPreviewCode');
      if (code) {
        code.textContent = task.content || '';
        if (typeof enhanceCodeBlocks === 'function') {
          enhanceCodeBlocks(document.getElementById('modalContent'));
        }
      }
      document.getElementById('serverTaskPreviewClose')?.addEventListener('click', closeModal);
    } catch (error) {
      showToast('No se pudo abrir la plantilla', error.message);
    }
  }

  function openCreateModal() {
    openUiModal({
      title: 'Nueva plantilla systemd',
      subtitle: 'Solo se genera el archivo; no se ejecuta en el host.',
      icon: 'fa-list-check',
      size: 'lg',
      body: `
        <form id="serverTaskForm" autocomplete="off">
          <div class="server-task-form-grid">
            <div class="field">
              <label for="serverTaskName">Nombre</label>
              <input id="serverTaskName" maxlength="64" placeholder="mi-aplicacion" required>
            </div>
            <div class="field">
              <label for="serverTaskKind">Tipo</label>
              <select id="serverTaskKind">
                <option value="service">systemd .service</option>
                <option value="timer">systemd .timer</option>
              </select>
            </div>
            <div class="field full">
              <label for="serverTaskDescription">Descripción</label>
              <input id="serverTaskDescription" maxlength="160" placeholder="Servicio de ejemplo para Axolotl">
            </div>
          </div>

          <div class="server-task-kind-fields" id="serverServiceFields">
            <div class="server-task-form-grid" style="margin-top:10px">
              <div class="field full">
                <label for="serverTaskExec">ExecStart</label>
                <input id="serverTaskExec" maxlength="500" placeholder="/usr/bin/php /var/www/app/artisan queue:work">
              </div>
              <div class="field">
                <label for="serverTaskWorkingDir">WorkingDirectory</label>
                <input id="serverTaskWorkingDir" maxlength="240" placeholder="/var/www/app">
              </div>
              <div class="field">
                <label for="serverTaskUser">User</label>
                <input id="serverTaskUser" maxlength="64" placeholder="www-data">
              </div>
              <div class="field">
                <label for="serverTaskRestart">Restart</label>
                <select id="serverTaskRestart">
                  <option value="on-failure">on-failure</option>
                  <option value="always">always</option>
                  <option value="no">no</option>
                </select>
              </div>
              <div class="field">
                <label for="serverTaskAfter">After</label>
                <input id="serverTaskAfter" maxlength="160" value="network.target">
              </div>
            </div>
          </div>

          <div class="server-task-kind-fields hidden" id="serverTimerFields">
            <div class="server-task-form-grid" style="margin-top:10px">
              <div class="field">
                <label for="serverTaskCalendar">OnCalendar</label>
                <input id="serverTaskCalendar" maxlength="160" placeholder="*-*-* 03:00:00">
              </div>
              <div class="field">
                <label for="serverTaskTarget">Unit objetivo</label>
                <input id="serverTaskTarget" maxlength="80" placeholder="mi-aplicacion.service">
              </div>
              <div class="field full">
                <label class="switch-row" style="padding:0">
                  <span class="switch-copy"><strong>Persistent</strong><small>Ejecuta la tarea pendiente después de un apagado.</small></span>
                  <span class="switch"><input id="serverTaskPersistent" type="checkbox" checked><span></span></span>
                </label>
              </div>
            </div>
          </div>

          <div class="server-task-warning">
            <i class="fa-solid fa-shield-halved"></i>
            <span>Axolotl valida y guarda esta definición en <code>storage/server-tasks</code>. No llama a <code>systemctl</code>, no copia archivos a <code>/etc/systemd</code> y no ejecuta comandos.</span>
          </div>
        </form>
      `,
      footer: `
        <button class="btn-ui btn-ghost" id="serverTaskCancel" type="button">Cancelar</button>
        <button class="btn-ui btn-primary" id="serverTaskCreate" type="button">Crear plantilla</button>
      `
    });

    const kind = document.getElementById('serverTaskKind');
    const serviceFields = document.getElementById('serverServiceFields');
    const timerFields = document.getElementById('serverTimerFields');

    const syncKind = () => {
      const timer = kind?.value === 'timer';
      serviceFields?.classList.toggle('hidden', timer);
      timerFields?.classList.toggle('hidden', !timer);
    };

    kind?.addEventListener('change', syncKind);
    syncKind();

    document.getElementById('serverTaskCancel')?.addEventListener('click', closeModal);
    document.getElementById('serverTaskCreate')?.addEventListener('click', createTask);
    document.getElementById('serverTaskName')?.focus();
  }

  async function createTask() {
    const kind = document.getElementById('serverTaskKind')?.value || 'service';
    const values = {
      name: document.getElementById('serverTaskName')?.value || '',
      kind,
      description: document.getElementById('serverTaskDescription')?.value || '',
      exec_start: document.getElementById('serverTaskExec')?.value || '',
      working_directory: document.getElementById('serverTaskWorkingDir')?.value || '',
      user: document.getElementById('serverTaskUser')?.value || '',
      restart: document.getElementById('serverTaskRestart')?.value || 'on-failure',
      after: document.getElementById('serverTaskAfter')?.value || '',
      on_calendar: document.getElementById('serverTaskCalendar')?.value || '',
      target_unit: document.getElementById('serverTaskTarget')?.value || '',
      persistent: document.getElementById('serverTaskPersistent')?.checked ? 'true' : 'false'
    };

    try {
      const task = await postForm('/api/server/tasks/create', values);
      closeModal();
      await loadTasks();
      showToast('Plantilla creada', task.filename);
    } catch (error) {
      showToast('No se pudo crear la plantilla', error.message);
    }
  }

  function confirmDelete(filename) {
    openUiModal({
      title: 'Eliminar plantilla',
      subtitle: filename,
      icon: 'fa-trash-can',
      size: 'sm',
      body: '<div class="modal-callout">Se eliminará únicamente el archivo de plantilla local. Esta acción no interactúa con systemd.</div>',
      footer: `
        <button class="btn-ui btn-ghost" id="serverTaskDeleteCancel" type="button">Cancelar</button>
        <button class="btn-ui btn-danger" id="serverTaskDeleteConfirm" type="button">Eliminar</button>
      `
    });

    document.getElementById('serverTaskDeleteCancel')?.addEventListener('click', closeModal);
    document.getElementById('serverTaskDeleteConfirm')?.addEventListener('click', async () => {
      try {
        await postForm('/api/server/tasks/delete', { filename });
        closeModal();
        await loadTasks();
        showToast('Plantilla eliminada', filename);
      } catch (error) {
        showToast('No se pudo eliminar la plantilla', error.message);
      }
    });
  }

  createBtn?.addEventListener('click', openCreateModal);
  refreshBtn?.addEventListener('click', loadTasks);
  loadTasks();
})();