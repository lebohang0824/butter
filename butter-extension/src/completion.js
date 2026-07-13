const vscode = require('vscode');

const DOCS = {
  app: 'Declare the root application specification.\n\n```butter\napp MyApp\n\tdescription "..."\n\tversion "1.0.0"\n```',
  product: 'Declare the root product specification (alias for `app`).\n\n```butter\nproduct MyProduct\n\tdescription "..."\n\tversion "1.0.0"\n```',
  description: 'A human-readable description of the current block.\n\n```butter\ndescription "My application"\n```',
  version: 'Version identifier for the app, feature, endpoint, or listener.\n\n```butter\nversion "1.0.0"\n```',
  feature: 'Define a logical feature group within the app.\n\n```butter\nfeature MyFeature\n\tdescription "..."\n\tparams\n\t\tparam X\n\t\t\ttype string\n```',
  endpoint: 'Define a synchronous HTTP network gateway with route, method, params, responses, actions, and return mappings.\n\n```butter\nendpoint ProcessOrder\n\tdescription "Process a checkout order"\n\troute "/api/checkout/orders"\n\tmethod "POST"\n\tparams\n\t\tparam checkout_token\n\t\t\ttype string\n\t\t\trequired true\n```',
  listener: 'Define an asynchronous message consumer with topic, params, actions, and return state mappings.\n\n```butter\nlistener ProcessEvent\n\tdescription "Process incoming events"\n\ttopic "events.process"\n\tparams\n\t\tparam event_id\n\t\t\ttype string\n\t\t\trequired true\n```',
  topic: 'Define the message topic or queue name for the listener.\n\n```butter\ntopic "events.process"\n```',
  params: 'Begin a parameter definition block for the current feature, endpoint, or listener.\n\n```butter\nparams\n\tparam Name\n\t\ttype string\n```',
  param: 'Define a single configuration parameter.\n\n```butter\nparam ParamName\n\ttype string\n\trequired true\n```',
  type: 'Set the data type of a parameter or field.\n\nValues: `string`, `int`, `float`, `bool`, `enum[...]`',
  required: 'Mark the parameter as required (`true`) or optional (`false`).',
  default: 'Set a default value for the parameter, used when no value is provided.',
  validate: 'Add a numeric validation rule (e.g., `">0"`, `">=1"`, `"=<100"`). Requires `int` or `float` type.',
  length: 'Enforce a maximum length for string or numeric parameters.',
  actions: 'Begin an action definition block for the current feature, endpoint, or listener.\n\n```butter\nactions\n\taction "Do something"\n```',
  action: 'Define a behavioral action step with a quoted description.\n\n```butter\naction "Execute the pipeline"\n\tenforce "Must be valid"\n```',
  enforce: 'Specify an invariant enforcement rule for the action.\n\n```butter\nenforce "The value must be positive"\n```',
  route: 'Define the relative URL path pattern for the endpoint.\n\n```butter\nroute "/api/resource/{id}"\n```',
  method: 'The transport-layer HTTP verb for the endpoint.\n\n```butter\nmethod "POST"\n```',
  responses: 'Begin a response schema definition block. Each response defines an internal JSON payload format.\n\n```butter\nresponses\n\tresponse OrderSuccess\n\t\tfield order_id\n\t\t\ttype string\n```',
  response: 'Define a reusable, internal response payload schema.\n\n```butter\nresponse OrderSuccess\n\tfield order_id\n\t\ttype string\n\tfield total_amount\n\t\ttype float\n```',
  field: 'Declare a JSON key inside a response schema. Defines the name and type of a single output field. Defaults to `string` if no type is given.\n\n```butter\nresponse OrderSuccess\n\tfield order_id\n\tfield total_amount\n\t\ttype float\n```',
  returns: 'Begin a return mapping block that binds execution outcomes to HTTP status codes (endpoint) or message states (listener).\n\n```butter\nreturns\n\treturn 200 OrderSuccess | if "Transaction succeeded"\n\treturn 400 "Invalid input" | if "Validation fails"\n```',
  return: 'Map an HTTP status code to a response payload or string literal with an optional condition.\n\nSyntax: `return <StatusCode> [<ResponseName> | <"String">] [| if/unless <"Condition">]`\n\n```butter\nreturn 201 OrderSuccess | if "Transaction authorized"\nreturn 400 "Invalid token" | if "Token validation fails"\nreturn 500 | if "Provider times out"\n```',
  ack: 'Acknowledge successful message processing. The message is removed from the queue.\n\n```butter\nreturn ack | if "Processing succeeded"\n```',
  nack: 'Negative acknowledgment. The message is rejected and may be requeued.\n\n```butter\nreturn nack | if "Invalid message format"\n```',
  retry: 'Request message retry. The message will be reprocessed later.\n\n```butter\nreturn retry | if "Temporary failure"\n```',
  dlq: 'Send message to dead letter queue. The message is moved to a dead letter queue for manual inspection.\n\n```butter\nreturn dlq | if "Permanent failure"\n```',
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
  snippet('endpoint', 'endpoint ${1:EndpointName}\n\t${0}'),
  snippet('listener', 'listener ${1:ListenerName}\n\t${0}'),
]);

const APP_BODY = setDocs([
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  snippet('feature', 'feature ${1:FeatureName}\n\t${0}'),
  snippet('endpoint', 'endpoint ${1:EndpointName}\n\t${0}'),
  snippet('listener', 'listener ${1:ListenerName}\n\t${0}'),
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
  item('route', vscode.CompletionItemKind.Keyword),
  item('method', vscode.CompletionItemKind.Keyword),
  item('params', vscode.CompletionItemKind.Keyword),
  item('responses', vscode.CompletionItemKind.Keyword),
  item('actions', vscode.CompletionItemKind.Keyword),
  item('returns', vscode.CompletionItemKind.Keyword),
]);

const LISTENER_BODY = setDocs([
  item('description', vscode.CompletionItemKind.Keyword),
  item('version', vscode.CompletionItemKind.Keyword),
  item('topic', vscode.CompletionItemKind.Keyword),
  item('params', vscode.CompletionItemKind.Keyword),
  item('actions', vscode.CompletionItemKind.Keyword),
  item('returns', vscode.CompletionItemKind.Keyword),
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

const RESPONSES_BODY = setDocs([
  snippet('response', 'response ${1:ResponseName}\n\t${0}'),
]);

const RESPONSE_BODY = setDocs([
  snippet('field', 'field ${1:FieldName}'),
]);

const FIELD_BODY = setDocs([
  item('type', vscode.CompletionItemKind.Keyword),
  snippet('field', 'field ${1:FieldName}'),
]);

const RETURNS_BODY = setDocs([
  snippet('return', 'return ${1:200} ${2:ResponseName} | if "${3:condition}"'),
]);

const LISTENER_RETURNS_BODY = setDocs([
  snippet('return', 'return ${1:ack} | if "${2:condition}"'),
]);

const MESSAGE_STATES = setDocs([
  item('ack', vscode.CompletionItemKind.Constant),
  item('nack', vscode.CompletionItemKind.Constant),
  item('retry', vscode.CompletionItemKind.Constant),
  item('dlq', vscode.CompletionItemKind.Constant),
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
  if (pk === 'endpoint') {
    if (currentIndent === 1) return 'endpoint-body';
    if (gk === 'params') return 'params-body';
    if (gk === 'actions') return 'actions-body';
    if (gk === 'responses') return 'responses-body';
    if (gk === 'returns') return 'returns-body';
    return 'endpoint-body';
  }
  if (pk === 'listener') {
    if (currentIndent === 1) return 'listener-body';
    if (gk === 'params') return 'params-body';
    if (gk === 'actions') return 'actions-body';
    if (gk === 'returns') return 'listener-returns-body';
    return 'listener-body';
  }
  if (pk === 'params') return 'params-body';
  if (pk === 'param') return 'param-body';
  if (pk === 'actions') return 'actions-body';
  if (pk === 'action') return 'action-body';
  if (pk === 'responses') return 'responses-body';
  if (pk === 'response') return 'response-body';
  if (pk === 'field') return 'field-body';
  if (pk === 'returns') return 'returns-body';

  if (gk === 'params') return 'params-body';
  if (gk === 'actions') return 'actions-body';
  if (gk === 'param') return 'param-body';
  if (gk === 'action') return 'action-body';
  if (gk === 'responses') return 'responses-body';
  if (gk === 'response') return 'response-body';
  if (gk === 'field') return 'field-body';
  if (gk === 'returns') return 'returns-body';

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
      case 'endpoint-body': return ENDPOINT_BODY;
      case 'listener-body': return LISTENER_BODY;
      case 'params-body': return PARAMS_BODY;
      case 'param-body': return PARAM_BODY;
      case 'actions-body': return ACTIONS_BODY;
      case 'action-body': return ACTION_BODY;
      case 'responses-body': return RESPONSES_BODY;
      case 'response-body': return RESPONSE_BODY;
      case 'field-body': return FIELD_BODY;
      case 'returns-body': return RETURNS_BODY;
      case 'listener-returns-body': return LISTENER_RETURNS_BODY;
      default: return TOP_LEVEL;
    }
  }
}

class ButterHoverProvider {
  provideHover(document, position) {
    const lineText = document.lineAt(position.line).text;

    const returnMatch = lineText.match(/\breturn\s+\d{3}\s+(\w+)/);
    if (returnMatch) {
      const respName = returnMatch[1];
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

    const returnMatch = lineText.match(/\breturn\s+\d{3}\s+(\w+)/);
    if (returnMatch) {
      const respName = returnMatch[1];
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
  const re = /^\s*response\s+(\w+)\s*$/;
  for (let i = 0; i < document.lineCount; i++) {
    const m = document.lineAt(i).text.match(re);
    if (m && m[1] === name) {
      const fields = [];
      let j = i + 1;
      while (j < document.lineCount) {
        const fl = document.lineAt(j).text;
        const fm = fl.match(/^\s+field\s+(\w+)\s*$/);
        if (fm) {
          let type = 'string';
          if (j + 1 < document.lineCount) {
            const tl = document.lineAt(j + 1).text;
            const tm = tl.match(/^\s+type\s+(\S+)\s*$/);
            if (tm) { type = tm[1]; j++; }
          }
          fields.push({ name: fm[1], type });
          j++;
        } else if (fl.match(/^\s*response\s+/) || fl.match(/^\s*returns\s*$/) || fl.match(/^\s*actions\s*$/) || (fl.match(/^\S/) && fl.trim() !== '')) {
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
