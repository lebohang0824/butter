const vscode = require('vscode');
const cp = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { ButterCompletionProvider, ButterHoverProvider, ButterDefinitionProvider } = require('./completion');

function activate(context) {
  const diagnostics = vscode.languages.createDiagnosticCollection('butter');
  const output = vscode.window.createOutputChannel('Butter');
  let formatting = false;

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

  function format(document) {
    return new Promise((resolve, reject) => {
      const compilerPath = vscode.workspace.getConfiguration('butter').get('compilerPath', 'butter');
      const text = document.getText();
      const tmpFile = path.join(os.tmpdir(), `butter-fmt-${Date.now()}.butter`);

      fs.writeFile(tmpFile, text, 'utf8', writeErr => {
        if (writeErr) { reject(writeErr); return; }

        cp.execFile(compilerPath, ['fmt', tmpFile], { timeout: 10000 }, fmtErr => {
          if (fmtErr) {
            fs.unlink(tmpFile, () => {});
            reject(fmtErr);
            return;
          }

          fs.readFile(tmpFile, 'utf8', (readErr, formatted) => {
            fs.unlink(tmpFile, () => {});
            if (readErr) { reject(readErr); return; }
            if (text === formatted) { resolve([]); return; }

            const range = new vscode.Range(
              document.positionAt(0),
              document.positionAt(text.length)
            );
            resolve([new vscode.TextEdit(range, formatted)]);
          });
        });
      });
    });
  }

  async function formatOnSave(doc) {
    if (doc.languageId !== 'butter' || doc.uri.scheme !== 'file') return;

    if (formatting) {
      lint(doc);
      return;
    }

    formatting = true;
    let needsLint = true;
    try {
      const edits = await format(doc);
      if (edits.length > 0) {
        const edit = new vscode.WorkspaceEdit();
        edit.replace(doc.uri, edits[0].range, edits[0].newText);
        await vscode.workspace.applyEdit(edit);
        await doc.save();
        needsLint = false;
      }
    } catch (e) {
      output.appendLine(`Format error: ${e.message}`);
    } finally {
      formatting = false;
    }

    if (needsLint) {
      lint(doc);
    }
  }

  context.subscriptions.push(
    vscode.workspace.onDidSaveTextDocument(formatOnSave),
    vscode.workspace.onDidOpenTextDocument(lint),
    vscode.commands.registerCommand('butter.lint', () => {
      const editor = vscode.window.activeTextEditor;
      if (editor) lint(editor.document);
    }),
    vscode.workspace.onDidChangeConfiguration(e => {
      if (e.affectsConfiguration('butter.compilerPath')) {
        vscode.workspace.textDocuments.forEach(lint);
      }
    }),
    vscode.languages.registerDocumentFormattingEditProvider('butter', {
      provideDocumentFormattingEdits(document) {
        return format(document);
      }
    }),
    vscode.commands.registerCommand('butter.fmt', () => {
      const editor = vscode.window.activeTextEditor;
      if (editor) {
        vscode.commands.executeCommand('editor.action.formatDocument');
      }
    }),
    vscode.languages.registerCompletionItemProvider('butter', new ButterCompletionProvider(), '|', ' ', '['),
    vscode.languages.registerHoverProvider('butter', new ButterHoverProvider()),
    vscode.languages.registerDefinitionProvider('butter', new ButterDefinitionProvider())
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