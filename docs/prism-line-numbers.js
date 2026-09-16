(function () {
  function numberBlocks() {
    document.querySelectorAll('pre').forEach(function (pre) {
      if (pre.classList.contains('code-block')) return;

      var code = null;
      for (var i = 0; i < pre.children.length; i++) {
        if (pre.children[i].tagName === 'CODE') {
          code = pre.children[i];
          break;
        }
      }
      if (!code) return;

      if (code.className && String(code.className).includes('language-bash')) return;

      var text = code.textContent.replace(/\n$/, '');
      var lines = text.split('\n');

      pre.classList.add('code-block');

      var gutter = document.createElement('div');
      gutter.className = 'code-gutter';
      gutter.setAttribute('aria-hidden', 'true');

      for (var n = 1; n <= lines.length; n++) {
        var span = document.createElement('span');
        span.textContent = n;
        gutter.appendChild(span);
      }

      pre.insertBefore(gutter, pre.firstChild);
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', numberBlocks);
  } else {
    numberBlocks();
  }
})();