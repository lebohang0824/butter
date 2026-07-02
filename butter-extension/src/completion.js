const vscode = require('vscode');

const DOCS = {
  app: 'Declare the root application specification.\n\n```butter\napp MyApp\n\tdescription "..."\n\tversion "1.0.0"\n```',
  product: 'Declare the root product specification (alias for `app`).\n\n```butter\nproduct MyProduct\n\tdescription "..."\n\tversion "1.0.0"\n```',
  description: 'A human-readable description of the current block.\n\n```butter\ndescription "My application"\n```',
  version: 'Version identifier for the app or feature.\n\n```butter\nversion "1.0.0"\n```',
  feature: 'Define a logical feature group within the app.\n\n```butter\nfeature MyFeature\n\tdescription "..."\n\tparams\n\t\tparam X\n\t\t\ttype string\n```',
  params: 'Begin a parameter definition block for the current feature.\n\n```butter\nparams\n\tparam Name\n\t\ttype string\n```',
  param: 'Define a single configuration parameter.\n\n```butter\nparam ParamName\n\ttype string\n\trequired true\n```',
  type: 'Set the data type of the parameter.\n\nValues: `string`, `int`, `float`, `bool`, `enum[...]`',
  required: 'Mark the parameter as required (`true`) or optional (`false`).',
  default: 'Set a default value for the parameter, used when no value is provided.',
  validate: 'Add a numeric validation rule (e.g., `">0"`, `">=1"`, `"=<100"`). Requires `int` or `float` type.',
  length: 'Enforce a maximum length for string or numeric parameters.',
  actions: 'Begin an action definition block for the current feature.\n\n```butter\nactions\n\taction "Do something"\n```',
  action: 'Define a behavioral action step with a quoted description.\n\n```butter\naction "Execute the pipeline"\n\tenforce "Must be valid"\n```',
  enforce: 'Specify an invariant enforcement rule for the action.\n\n```butter\nenforce "The value must be positive"\n```',
  if: 'Execute this action only when the condition is true.\n\n```butter\naction "Do something" | if "condition == true"\n```',
  unless: 'Execute this action only when the condition is false.\n\n```butter\naction "Do something" | unless "condition == false"\n```',
  when: 'Execute this action when the condition becomes true.\n\n```butter\naction "Do something" | when "event occurs"\n```',
  while: 'Repeatedly execute this action while the condition remains true.\n\n```butter\naction "Do something" | while "condition == true"\n```',
  string: 'Text string data type.\n\n```butter\nparam Name\n\ttype string\n```',
  int: 'Integer number data type.\n\n```butter\nparam Count\n\ttype int\n```',
  float: 'Floating-point number data type.\n\n```butter\nparam Ratio\n\ttype float\n```',
  bool: 'Boolean (true/false) data type.\n\n```butter\nparam Enabled\n\ttype bool\n```',
  boolean: 'Boolean (true/false) data type (alias for `bool`).',
  'enum[...]': 'Restricted to values from a predefined list.\n\n```butter\nparam Format\n\ttype enum["json", "yaml"]\n```',
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
  snippet('app', 'app ${1:AppName}\n\t${0}'),
  snippet('product', 'product ${1:ProductName}\n\t${0}'),
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  snippet('feature', 'feature ${1:FeatureName}\n\t${0}'),
]);

const APP_BODY = setDocs([
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  snippet('feature', 'feature ${1:FeatureName}\n\t${0}'),
]);

const FEATURE_BODY = setDocs([
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  item('params', vscode.CompletionItemKind.Keyword),
  item('actions', vscode.CompletionItemKind.Keyword),
]);

const PARAMS_BODY = setDocs([
  snippet('param', 'param ${1:ParamName}\n\t${0}'),
]);

const PARAM_BODY = setDocs([
  item('type', vscode.CompletionItemKind.Keyword),
  item('required', vscode.CompletionItemKind.Keyword),
  item('default', vscode.CompletionItemKind.Keyword),
  item('validate', vscode.CompletionItemKind.Keyword),
  item('length', vscode.CompletionItemKind.Keyword),
]);

const ACTIONS_BODY = setDocs([
  snippet('action', 'action "${1:statement}"\n\t${0}'),
]);

const ACTION_BODY = setDocs([
  item('enforce', vscode.CompletionItemKind.Keyword),
]);

const TYPES = setDocs([
  item('string', vscode.CompletionItemKind.TypeParameter),
  item('int', vscode.CompletionItemKind.TypeParameter),
  item('float', vscode.CompletionItemKind.TypeParameter),
  item('bool', vscode.CompletionItemKind.TypeParameter),
  item('boolean', vscode.CompletionItemKind.TypeParameter),
  snippet('enum[...]', 'enum[${1:values}]'),
]);

const BOOLS = setDocs([
  item('true', vscode.CompletionItemKind.Constant),
  item('false', vscode.CompletionItemKind.Constant),
]);

const CONDITIONALS = setDocs([
  item('if', vscode.CompletionItemKind.Keyword),
  item('unless', vscode.CompletionItemKind.Keyword),
  item('when', vscode.CompletionItemKind.Keyword),
  item('while', vscode.CompletionItemKind.Keyword),
]);

function getParentChain(document, lineNum) {
  const currentLine = document.lineAt(Math.min(lineNum, document.lineCount - 1)).text;
  const currentIndent = currentLine.search(/\S/);
  const chain = [];
  let minIndent = currentIndent >= 0 ? currentIndent : Infinity;

  for (let i = lineNum - 1; i >= 0; i--) {
    const line = document.lineAt(i).text;
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith('#')) continue;
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
    isEmpty: trimmed === '' || trimmed.startsWith('#'),
    firstWord,
    words,
    wordCount: words.length,
    currentIndent,
    parent,
    grandparent,
    chain,
  };

  if (firstWord === 'type' && words.length === 1) {
    context.completionType = 'type-value';
  } else if (firstWord === 'required' && words.length === 1) {
    context.completionType = 'bool-value';
  } else if (firstWord === 'default' && words.length === 1) {
    context.completionType = 'default-value';
  } else if (firstWord === '|' || trimmed.startsWith('|')) {
    context.completionType = 'pipe-condition';
  } else if (context.isEmpty || firstWord === '') {
    context.completionType = 'keyword';
  } else if (words.length <= 1) {
    context.completionType = 'keyword';
  } else {
    context.completionType = 'none';
  }

  return context;
}

function suggestionKind(context) {
  const { parent, grandparent, currentIndent, isEmpty } = context;

  const pk = parent ? parent.keyword : null;
  const gk = grandparent ? grandparent.keyword : null;

  if (currentIndent === 0) {
    if (pk === 'app' || pk === 'product') return 'app-body';
    return 'top-level';
  }

  if (pk === 'app' || pk === 'product') return 'app-body';
  if (pk === 'feature') {
    if (currentIndent === 1) return 'feature-body';
    if (gk === 'params') return 'params-body';
    if (gk === 'actions') return 'actions-body';
    return 'feature-body';
  }
  if (pk === 'params') return 'params-body';
  if (pk === 'param') return 'param-body';
  if (pk === 'actions') return 'actions-body';
  if (pk === 'action') return 'action-body';

  if (gk === 'params') return 'params-body';
  if (gk === 'actions') return 'actions-body';
  if (gk === 'param') return 'param-body';
  if (gk === 'action') return 'action-body';

  return 'keyword';
}

class ButterCompletionProvider {
  provideCompletionItems(document, position) {
    const ctx = getContext(document, position.line);

    if (ctx.completionType === 'type-value') {
      return TYPES;
    }
    if (ctx.completionType === 'bool-value') {
      return BOOLS;
    }
    if (ctx.completionType === 'default-value') {
      return [...BOOLS, item('"..."', vscode.CompletionItemKind.Value)];
    }
    if (ctx.completionType === 'pipe-condition') {
      return CONDITIONALS;
    }

    const kind = suggestionKind(ctx);

    switch (kind) {
      case 'top-level': return TOP_LEVEL;
      case 'app-body': return APP_BODY;
      case 'feature-body': return FEATURE_BODY;
      case 'params-body': return PARAMS_BODY;
      case 'param-body': return PARAM_BODY;
      case 'actions-body': return ACTIONS_BODY;
      case 'action-body': return ACTION_BODY;
      default: return TOP_LEVEL;
    }
  }
}

class ButterHoverProvider {
  provideHover(document, position) {
    const wordRange = document.getWordRangeAtPosition(position);
    if (!wordRange) return null;
    const word = document.getText(wordRange);
    if (DOCS[word]) {
      return new vscode.Hover(new vscode.MarkdownString(DOCS[word]));
    }
    return null;
  }
}

module.exports = { ButterCompletionProvider, ButterHoverProvider };
