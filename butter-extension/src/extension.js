const vscode = require('vscode');
const cp = require('child_process');

function activate(context) {
  const diagnostics = vscode.languages.createDiagnosticCollection('butter');
  const output = vscode.window.createOutputChannel('Butter');

  context.subscriptions.push(diagnostics, output);

  function lint(doc) {
    if (doc.languageId !== 'butter' || doc.uri.scheme !== 'file') return;

    const compilerPath = vscode.workspace.getConfiguration('butter').get('compilerPath', 'butter');

    diagnostics.delete(doc.uri);

    cp.execFile(compilerPath, ['compile', '--check', doc.fileName], { timeout: 10000 }, (err, _stdout, stderr) => {
      if (err) {
        if (err.code === 'ENOENT') {
          output.appendLine(`Butter compiler not found at '${compilerPath}'`);
          return;
        }

        const diags = parseErrors(stderr, doc);
        if (diags.length > 0) {
          diagnostics.set(doc.uri, diags);
        }
      } else {
        diagnostics.delete(doc.uri);
      }
    });
  }

  context.subscriptions.push(
    vscode.workspace.onDidSaveTextDocument(lint),
    vscode.workspace.onDidOpenTextDocument(lint),
    vscode.commands.registerCommand('butter.lint', () => {
      const editor = vscode.window.activeTextEditor;
      if (editor) lint(editor.document);
    }),
    vscode.workspace.onDidChangeConfiguration(e => {
      if (e.affectsConfiguration('butter.compilerPath')) {
        vscode.workspace.textDocuments.forEach(lint);
      }
    })
  );

  vscode.workspace.textDocuments.forEach(lint);
}

function deactivate() {}

function parseErrors(stderr, doc) {
  const diags = [];
  const re = /line (\d+): (.+)/g;
  let m;
  while ((m = re.exec(stderr)) !== null) {
    const line = Math.max(0, parseInt(m[1], 10) - 1);
    const text = m[2].trim();
    const range = new vscode.Range(line, 0, line, doc.lineAt(line).text.length);
    diags.push(new vscode.Diagnostic(range, text, vscode.DiagnosticSeverity.Error));
  }
  return diags;
}

module.exports = { activate, deactivate };
