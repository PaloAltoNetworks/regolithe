'use strict';

import * as vscode from 'vscode';
import * as cp from 'child_process';

import { shouldConsiderDocument } from './utils';

export class RegolitheDocumentFormattingEditProvider implements vscode.DocumentFormattingEditProvider {

    formatCommandBinPath: string;

    constructor(toolPath: string) {

        this.formatCommandBinPath = toolPath;
    }

    public provideDocumentFormattingEdits(document: vscode.TextDocument, options: vscode.FormattingOptions, token: vscode.CancellationToken): Thenable<vscode.TextEdit[]> {

        if (!shouldConsiderDocument(document)) {
            return null
        }

        return this.format(document).then(edits => edits)
    }

    private format(doc: vscode.TextDocument): Thenable<vscode.TextEdit[]> {

        return new Promise<vscode.TextEdit[]>((resolve, reject) => {

            let stdout = '';
            let stderr = '';

            const p = cp.spawn(this.formatCommandBinPath, ['beautify'])
            p.stdout.setEncoding('utf8');
            p.stdout.on('data', data => stdout += data);
            p.stderr.on('data', data => stderr += data);
            p.on('error', err => {
                return reject();
            });

            p.on('close', code => {

                if (code !== 0) {
                    return reject(stderr);
                }

                const edit = new vscode.WorkspaceEdit()
                const range = new vscode.Range(new vscode.Position(0, 0), doc.lineAt(doc.lineCount - 1).range.end);
                const edits: vscode.TextEdit[] = [new vscode.TextEdit(range, stdout)]

                return resolve(edits);
            });

            p.stdin.end(doc.getText());
        });
    }
}
