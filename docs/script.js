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

  const sidebar = document.getElementById('sidebar');
  const menuToggle = document.getElementById('menuToggle');
  const navLinks = sidebar.querySelectorAll('nav a');

  var backdrop = document.createElement('div');
  backdrop.className = 'sidebar-backdrop';
  document.body.appendChild(backdrop);

  function toggleSidebar(open) {
    if (open === undefined) {
      sidebar.classList.toggle('open');
    } else if (open) {
      sidebar.classList.add('open');
    } else {
      sidebar.classList.remove('open');
    }
    var isOpen = sidebar.classList.contains('open');
    backdrop.classList.toggle('visible', isOpen);
    if (window.innerWidth <= 768) {
      document.body.style.overflow = isOpen ? 'hidden' : '';
    }
  }

  menuToggle.addEventListener('click', function () {
    toggleSidebar();
  });

  navLinks.forEach(function (link) {
    link.addEventListener('click', function () {
      toggleSidebar(false);
    });
  });

  backdrop.addEventListener('click', function () {
    toggleSidebar(false);
  });

  document.addEventListener('click', function (e) {
    if (window.innerWidth <= 768 &&
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
})();
