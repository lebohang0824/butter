(function () {
  const LS_KEY = 'butter-docs-theme';

  function getPreferredTheme() {
    const stored = localStorage.getItem(LS_KEY);
    if (stored) return stored;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function setTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem(LS_KEY, theme);
    const btn = document.getElementById('themeToggle');
    if (btn) btn.textContent = theme === 'dark' ? '☀️' : '🌙';
  }

  function toggleTheme() {
    const current = document.documentElement.getAttribute('data-theme') || 'light';
    setTheme(current === 'dark' ? 'light' : 'dark');
  }

  setTheme(getPreferredTheme());

  document.getElementById('themeToggle').addEventListener('click', toggleTheme);

  const SIDEBAR_KEY = 'butter-docs-sidebar';

  const sidebar = document.getElementById('sidebar');
  const menuToggle = document.getElementById('menuToggle');
  const navLinks = sidebar.querySelectorAll('nav a');

  var backdrop = document.createElement('div');
  backdrop.className = 'sidebar-backdrop';
  document.body.appendChild(backdrop);

  function isDesktop() {
    return window.innerWidth > 768;
  }

  function updateBackdrop() {
    var isOpen = sidebar.classList.contains('open');
    backdrop.classList.toggle('visible', isOpen);
    if (!isDesktop()) {
      document.body.style.overflow = isOpen ? 'hidden' : '';
    }
  }

  function setDesktopCollapsed(collapsed) {
    document.body.classList.toggle('sidebar-collapsed', collapsed);
    localStorage.setItem(SIDEBAR_KEY, collapsed ? '1' : '0');
  }

  function toggleSidebar(open) {
    if (isDesktop()) {
      var collapsed = open === undefined
        ? !document.body.classList.contains('sidebar-collapsed')
        : !open;
      setDesktopCollapsed(collapsed);
    } else {
      if (open === undefined) {
        sidebar.classList.toggle('open');
      } else if (open) {
        sidebar.classList.add('open');
      } else {
        sidebar.classList.remove('open');
      }
      updateBackdrop();
    }
  }

  if (isDesktop() && localStorage.getItem(SIDEBAR_KEY) === '1') {
    document.body.classList.add('sidebar-collapsed');
  }

  menuToggle.addEventListener('click', function () {
    toggleSidebar();
  });

  navLinks.forEach(function (link) {
    link.addEventListener('click', function () {
      if (!isDesktop()) toggleSidebar(false);
    });
  });

  backdrop.addEventListener('click', function () {
    toggleSidebar(false);
  });

  document.addEventListener('click', function (e) {
    if (!isDesktop() &&
        !sidebar.contains(e.target) &&
        !menuToggle.contains(e.target)) {
      toggleSidebar(false);
    }
  });

  var sections = document.querySelectorAll('section[id]');
  function updateActiveLink() {
    var scrollY = window.scrollY + 100;
    var currentId = '';
    sections.forEach(function (sec) {
      if (sec.offsetTop <= scrollY) {
        currentId = sec.id;
      }
    });
    navLinks.forEach(function (a) {
      a.classList.remove('active');
      if (a.getAttribute('href') === '#' + currentId) {
        a.classList.add('active');
      }
    });
  }

  window.addEventListener('scroll', updateActiveLink);
  updateActiveLink();

  function setupCopyButtons() {
    document.querySelectorAll('pre code').forEach(function (code) {
      var pre = code.parentElement;
      if (!pre || pre.querySelector('.code-copy')) return;

      var btn = document.createElement('button');
      btn.className = 'code-copy';
      btn.type = 'button';
      btn.textContent = 'Copy';
      btn.setAttribute('aria-label', 'Copy code to clipboard');
      pre.appendChild(btn);

      btn.addEventListener('click', function () {
        var text = code.textContent.replace(/\n+$/, '');
        function done() {
          btn.classList.add('copied');
          btn.textContent = 'Copied!';
          setTimeout(function () {
            btn.classList.remove('copied');
            btn.textContent = 'Copy';
          }, 1800);
        }
        if (navigator.clipboard && navigator.clipboard.writeText) {
          navigator.clipboard.writeText(text).then(done);
        } else {
          var ta = document.createElement('textarea');
          ta.value = text;
          ta.setAttribute('readonly', '');
          ta.style.position = 'fixed';
          ta.style.opacity = '0';
          document.body.appendChild(ta);
          ta.select();
          document.execCommand('copy');
          document.body.removeChild(ta);
          done();
        }
      });
    });
  }

  setupCopyButtons();
})();
