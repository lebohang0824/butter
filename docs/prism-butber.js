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
  'endpoint-declaration': {
    pattern: /\b(endpoint)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /endpoint/,
      'function': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'function'
      }
    }
  },
  'listener-declaration': {
    pattern: /\b(listener)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /listener/,
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
  'response-declaration': {
    pattern: /\b(response)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /response/,
      'class-name': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'class-name'
      }
    }
  },
  'field-declaration': {
    pattern: /\b(field)\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /field/,
      'variable': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'variable'
      }
    }
  },
  'return-statement': {
    pattern: /\b(return)\s+(\d{3})\s+([A-Za-z_]\w*)/,
    inside: {
      'keyword': /return/,
      'number': /\d{3}/,
      'class-name': {
        pattern: /[A-Za-z_]\w*/,
        alias: 'class-name'
      }
    }
  },
  'listener-return': {
    pattern: /\b(return)\s+(ack|nack|retry|dlq)\b(?:\s*\|\s*(if|unless)\b)?/,
    inside: {
      'keyword': /return/,
      'number': /ack|nack|retry|dlq/,
      'conditional': /if|unless/
    }
  },
  'return-string': {
    pattern: /\b(return)\s+(\d{3})\s+("[^"]*")/,
    inside: {
      'keyword': /return/,
      'number': /\d{3}/,
      'string': {
        pattern: /"[^"]*"/,
        alias: 'string'
      }
    }
  },
  'return bare': {
    pattern: /\b(return)\s+(\d{3})\b/,
    inside: {
      'keyword': /return/,
      'number': /\d{3}/
    }
  },
  'conditional': /\b(?:if|unless|when|while)\b/,
  'keyword': /\b(?:app|product|description|version|feature|endpoint|listener|params|param|type|required|default|validate|bool|boolean|length|enforce|actions|action|route|method|topic|responses|response|field|returns|return)\b/
};
