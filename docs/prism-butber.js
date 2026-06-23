Prism.languages.butter = {
  'comment': {
    pattern: /#.*/,
    greedy: true
  },
  'string': {
    pattern: /"(?:[^"\\]|\\.)*"/,
    greedy: true
  },
  'boolean': /\b(?:true|false)\b/,
  'app-declaration': {
    pattern: /\b(app)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /app/,
      'class-name': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'class-name'
      }
    }
  },
  'feature-declaration': {
    pattern: /\b(feature)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /feature/,
      'function': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'function'
      }
    }
  },
  'param-declaration': {
    pattern: /\b(param)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /param/,
      'variable': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'variable'
      }
    }
  },
  'conditional': /\b(?:if|unless|when|while)\b/,
  'keyword': /\b(?:app|description|version|feature|params|param|type|required|default|actions|action)\b/
};
