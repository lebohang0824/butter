Prism.languages.butter = {
  'comment': {
    pattern: /(?:#|\/\/).*/,
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
  'endpoint-declaration': {
    pattern: /\b(endpoint)\s+([A-Za-z_]\w*)\s+"([^"]*)"/,
    inside: {
      'keyword': /endpoint/,
      'function': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'function'
      },
      'string': {
        pattern: /"[^"]*"/,
        alias: 'string'
      }
    }
  },
  'enforce-statement': {
    pattern: /\b(enforce)\s+("[^"]*")/,
    inside: {
      'keyword': /enforce/,
      'string': {
        pattern: /"[^"]*"/,
        alias: 'string'
      }
    }
  },
  'returns-line': {
    pattern: /^\s+(\d{3})\s+([A-Za-z_]\w*)\s*$/m,
    inside: {
      'number': /\d{3}/,
      'class-name': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'class-name'
      }
    }
  },
  'returns-string': {
    pattern: /^\s+(\d{3})\s+("[^"]*")\s*$/m,
    inside: {
      'number': /\d{3}/,
      'string': {
        pattern: /"[^"]*"/,
        alias: 'string'
      }
    }
  },
  'returns bare': {
    pattern: /^\s+(\d{3})\s*$/m,
    inside: {
      'number': /\d{3}/
    }
  },
  'keyword': /\b(?:app|description|version|feature|endpoint|rules|method|params|actions|responses|returns|enforce|string|integer|double|boolean|enum|array)\b/
};
