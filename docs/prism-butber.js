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
    pattern: /\b(app|product)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /app|product/,
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
  'keyword': /\b(?:app|product|description|version|feature|params|param|type|required|default|validate|bool|boolean|length|actions|action)\b/
};
