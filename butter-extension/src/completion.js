const vscode = require('vscode');

const DOCS = {
  app: 'Declare the root application specification.\n\n```butter\napp MyApp\n  description "..."\n  version "1.0.0"\n```',
  description: 'A human-readable description of the current block.\n\n```butter\ndescription "My application"\n```',
  version: 'Version identifier for the app, feature, or endpoint.\n\n```butter\nversion "1.0.0"\n```',
  feature: 'Define a logical feature group within the app.\n\n```butter\nfeature MyFeature\n  description "..."\n  params\n    name string\n```',
  endpoint: 'Define a synchronous HTTP network gateway with route, method, params, responses, actions, and returns mappings.\n\n```butter\nendpoint ProcessOrder "/api/checkout/orders"\n  description "Process a checkout order"\n  method "POST"\n  params\n    checkout_token string\n```',
  rules: 'Begin an app-level rules block of quoted strings.\n\n```butter\nrules\n  "Use MVC pattern"\n```',
  params: 'Begin a parameter definition block.\n\n```butter\nparams\n  name string\n```',
  actions: 'Begin an action definition block.\n\n```butter\nactions\n  "Do something"\n```',
  enforce: 'Enforce a constraint directly below its parent action string.\n\n```butter\n  "Validate name"\n    enforce "Sanitize name before storing"\n```',
  responses: 'Begin a response schema definition block.\n\n```butter\nresponses\n  OrderSuccess\n    order_id string\n```',
  method: 'The HTTP method for the endpoint, as a quoted string.\n\n```butter\nmethod "POST"\n```',
  returns: 'Begin a returns block mapping HTTP status codes to response payloads.\n\n```butter\nreturns\n  201 OrderSuccess\n  500 "Server error"\n```',
  string: 'Text string data type.\n\n```butter\nname string\n```',
  integer: 'Integer number data type.\n\n```butter\ncount integer\n```',
  double: 'Floating-point number data type.\n\n```butter\nratio double\n```',
  boolean: 'Boolean (true/false) data type.\n\n```butter\ncompleted boolean\n```',
  'enum[...]': 'Restricted to values from a predefined list.\n\n```butter\npriority enum["low", "high"]\n```',
  'array[...]': 'An array of a primitive type.\n\n```butter\nmembers_id array[integer]\n```',
  true: 'Boolean true value.',
  false: 'Boolean false value.',
};

function snippet(label, insert) {
  const item = new vscode.CompletionItem(label, vscode.CompletionItemKind.Snippet);
  item.insertText = new vscode.SnippetString(insert);
  return item;
}

function item(label, kind) {
  const i = new vscode.CompletionItem(label, kind);
  i.documentation = DOCS[label] ? new vscode.MarkdownString(DOCS[label]) : undefined;
  return i;
}

function setDocs(items) {
  for (const i of items) {
    if (DOCS[i.label]) {
      i.documentation = new vscode.MarkdownString(DOCS[i.label]);
    }
  }
  return items;
}

const TOP_LEVEL = setDocs([
  snippet('app', 'app ${1:AppName}\n  ${0}'),
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  snippet('feature', 'feature ${1:FeatureName}\n  ${0}'),
  snippet('endpoint', 'endpoint ${1:EndpointName} "${2:/path}"\n  ${0}'),
  item('rules', vscode.CompletionItemKind.Keyword),
]);

const APP_BODY = setDocs([
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  item('rules', vscode.CompletionItemKind.Keyword),
  snippet('feature', 'feature ${1:FeatureName}\n  ${0}'),
  snippet('endpoint', 'endpoint ${1:EndpointName} "${2:/path}"\n  ${0}'),
]);

const FEATURE_BODY = setDocs([
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  item('params', vscode.CompletionItemKind.Keyword),
  item('actions', vscode.CompletionItemKind.Keyword),
]);

const ENDPOINT_BODY = setDocs([
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  item('method', vscode.CompletionItemKind.Keyword),
  item('params', vscode.CompletionItemKind.Keyword),
  item('responses', vscode.CompletionItemKind.Keyword),
  item('actions', vscode.CompletionItemKind.Keyword),
  item('returns', vscode.CompletionItemKind.Keyword),
]);

const PARAMS_BODY = setDocs([
  snippet('bare_param', '${1:name} ${2:string}'),
]);

const ACTIONS_BODY = setDocs([
  snippet('bare_action', '"${1:action statement}"'),
]);

const ENFORCE_BODY = setDocs([
  snippet('enforce_line', 'enforce "${1:constraint}"'),
]);

const RESPONSES_BODY = setDocs([
  snippet('response_header', '${1:ResponseName}'),
]);

const RESPONSE_BODY = setDocs([
  snippet('bare_field', '${1:name} ${2:string}'),
]);

const RETURNS_BODY = setDocs([
  snippet('returns_line', '${1:201} ${2:ResponseName}'),
  snippet('returns_string', '${1:500} "${2:Server error}"'),
]);

const METHODS = setDocs([
  snippet('"POST"', '"POST"'),
  snippet('"GET"', '"GET"'),
  snippet('"PUT"', '"PUT"'),
  snippet('"DELETE"', '"DELETE"'),
  snippet('"PATCH"', '"PATCH"'),
]);

const TYPES = setDocs([
  item('string', vscode.CompletionItemKind.TypeParameter),
  item('integer', vscode.CompletionItemKind.TypeParameter),
  item('double', vscode.CompletionItemKind.TypeParameter),
  item('boolean', vscode.CompletionItemKind.TypeParameter),
  snippet('enum[...]', 'enum[${1:values}]'),
  snippet('array[...]', 'array[${1:TypeName}]'),
]);

const BOOLS = setDocs([
  item('true', vscode.CompletionItemKind.Constant),
  item('false', vscode.CompletionItemKind.Constant),
]);

function getParentChain(document, lineNum) {
  const currentLine = document.lineAt(Math.min(lineNum, document.lineCount - 1)).text;
  const currentIndent = currentLine.search(/\S/);
  const chain = [];
  let minIndent = currentIndent >= 0 ? currentIndent : Infinity;

  for (let i = lineNum - 1; i >= 0; i--) {
    const line = document.lineAt(i).text;
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith('#') || trimmed.startsWith('//')) continue;
    const indent = line.search(/\S/);
    if (indent < minIndent) {
      const kw = trimmed.split(/\s+/)[0];
      chain.unshift({ keyword: kw, indent, line: i });
      minIndent = indent;
    }
  }
  return chain;
}

function getContext(document, lineNum) {
  const line = document.lineAt(Math.min(lineNum, document.lineCount - 1)).text;
  const trimmed = line.trim();
  const firstWord = trimmed.split(/\s+/)[0];
  const words = trimmed.split(/\s+/);
  let currentIndent = line.search(/\S/);

  if (currentIndent === -1) {
    for (let i = lineNum - 1; i >= 0; i--) {
      const prev = document.lineAt(i).text;
      const pi = prev.search(/\S/);
      if (pi >= 0) { currentIndent = pi; break; }
    }
    if (currentIndent === -1) currentIndent = 0;
  }

  const chain = getParentChain(document, lineNum);
  const parent = chain.length > 0 ? chain[chain.length - 1] : null;
  const grandparent = chain.length > 1 ? chain[chain.length - 2] : null;

  const context = {
    isEmpty: trimmed === '' || trimmed.startsWith('#') || trimmed.startsWith('//'),
    firstWord,
    words,
    wordCount: words.length,
    currentIndent,
    parent,
    grandparent,
    chain,
  };

  if (context.isEmpty || firstWord === '') {
    context.completionType = 'keyword';
  } else if (words.length <= 1) {
    context.completionType = 'keyword';
  } else {
    context.completionType = 'none';
  }

  return context;
}

function suggestionKind(context) {
  const { parent, grandparent, currentIndent } = context;

  const pk = parent ? parent.keyword : null;
  const gk = grandparent ? grandparent.keyword : null;

  if (currentIndent === 0) {
    if (pk === 'app') return 'app-body';
    return 'top-level';
  }

  if (pk === 'app') return 'app-body';
  if (pk === 'rules' || gk === 'rules') return 'rules-body';
  if (gk === 'params') return 'params-body';
  if (gk === 'actions') return 'actions-body';
  if (gk === 'enforce') return 'enforce-body';
  if (gk === 'returns') return 'returns-body';
  if (gk === 'responses') return 'responses-body';
  if (pk === 'params' || pk === 'actions') return 'params-body';
  if (pk === 'method') return 'methods';
  if (pk === 'returns') return 'returns-body';
  if (gk === 'method') return 'methods';
  if (pk === 'feature') return 'feature-body';
  if (pk === 'endpoint') return 'endpoint-body';

  return 'top-level';
}

class ButterCompletionProvider {
  provideCompletionItems(document, position) {
    const ctx = getContext(document, position.line);
    const kind = suggestionKind(ctx);

    switch (kind) {
      case 'top-level': return TOP_LEVEL;
      case 'app-body': return APP_BODY;
      case 'rules-body': return [snippet('rule_string', '"${1:rule statement}"')];
      case 'feature-body': return FEATURE_BODY;
      case 'endpoint-body': return ENDPOINT_BODY;
      case 'params-body': return [...PARAMS_BODY, ...TYPES];
      case 'actions-body': return ACTIONS_BODY;
      case 'enforce-body': return ENFORCE_BODY;
      case 'responses-body': return RESPONSES_BODY;
      case 'response-body': return RESPONSE_BODY;
      case 'returns-body': return RETURNS_BODY;
      case 'methods': return METHODS;
      default: return TOP_LEVEL;
    }
  }
}

class ButterHoverProvider {
  provideHover(document, position) {
    const lineText = document.lineAt(position.line).text;

    const returnsMatch = lineText.match(/^\s+(\d{3})\s+(\w+)/);
    if (returnsMatch) {
      const respName = returnsMatch[2];
      const start = lineText.indexOf(respName);
      const end = start + respName.length;
      if (position.character >= start && position.character <= end) {
        const decl = findResponseDecl(document, respName);
        if (decl) {
          const md = new vscode.MarkdownString(`**${respName}** — response schema\n\n`);
          for (const f of decl.fields) {
            md.appendMarkdown(`\`${f.name}\` — \`${f.type}\`\n\n`);
          }
          return new vscode.Hover(md);
        }
      }
    }

    const wordRange = document.getWordRangeAtPosition(position);
    if (!wordRange) return null;
    const word = document.getText(wordRange);
    if (DOCS[word]) {
      return new vscode.Hover(new vscode.MarkdownString(DOCS[word]));
    }
    return null;
  }
}

class ButterDefinitionProvider {
  provideDefinition(document, position) {
    const lineText = document.lineAt(position.line).text;

    const returnsMatch = lineText.match(/^\s+(\d{3})\s+(\w+)/);
    if (returnsMatch) {
      const respName = returnsMatch[2];
      const start = lineText.indexOf(respName);
      const end = start + respName.length;
      if (position.character >= start && position.character <= end) {
        const decl = findResponseDecl(document, respName);
        if (decl) {
          return new vscode.Location(document.uri, new vscode.Position(decl.line, 0));
        }
      }
    }

    return null;
  }
}

function findResponseDecl(document, name) {
  for (let i = 0; i < document.lineCount; i++) {
    const line = document.lineAt(i).text;
    const m = line.match(/^\s+([A-Z][A-Za-z0-9_]*)\s*$/);
    if (m && m[1] === name) {
      const fields = [];
      let j = i + 1;
      while (j < document.lineCount) {
        const fl = document.lineAt(j).text;
        const fm = fl.match(/^\s{2,}([A-Za-z_]\w*)\s+(string|integer|double|boolean|enum\[.*?\]|array\[.*?\])\s*$/);
        if (fm) {
          fields.push({ name: fm[1], type: fm[2] });
          j++;
        } else if (fl.match(/^\s*returns\s*$/) || fl.match(/^\s*actions\s*$/) || (fl.match(/^\S/) && fl.trim() !== '') || fl.match(/^\s+[A-Z][A-Za-z0-9_]*\s*$/)) {
          break;
        } else {
          j++;
        }
      }
      return { line: i, fields };
    }
  }
  return null;
}

module.exports = { ButterCompletionProvider, ButterHoverProvider, ButterDefinitionProvider };
