'use strict';

import * as vscode from 'vscode';
import * as cp from 'child_process';
import * as path from 'path';

import { tmFile, vmFile, pmFile, shouldConsiderDocument } from './utils';

export class RegolitheDocumentFormattingEditProvider {

    formatCommandBinPath: string;
    outputChannel: vscode.OutputChannel

    constructor(toolPath: string, outputChannel: vscode.OutputChannel) {
        this.formatCommandBinPath = toolPath;
        this.outputChannel = outputChannel;
    }

    public provideDocumentFormattingEdits(doc: vscode.TextDocument, options: vscode.FormattingOptions, token: vscode.CancellationToken): Thenable<vscode.TextEdit[]> {
        return this.format(doc)
    }

    public format(doc: vscode.TextDocument): Thenable<vscode.TextEdit[]> {

        if (!shouldConsiderDocument(doc)) {
            return null
        }

        return new Promise<vscode.TextEdit[]>((resolve, reject) => {

            let stdout = '';
            let stderr = '';

            const params = ['format'];
            if (doc.fileName.endsWith(tmFile)) {
                params.push("--mode", "typemapping");
            } else if (doc.fileName.endsWith(vmFile)) {
                params.push("--mode", "validationmapping");
            } else if (doc.fileName.endsWith(pmFile)) {
                params.push("--mode", "parametermapping");
            }

            const p = cp.spawn(this.formatCommandBinPath, params)

            p.stdout.setEncoding('utf8');
            p.stdout.on('data', data => stdout += data);
            p.stderr.on('data', data => stderr += data);
            p.on('error', err => {
                console.error('Regolithe: unable to format: ', err);
                return reject();
            });

            p.on('close', code => {

                if (code !== 0) {
                    this.outputChannel.clear()
                    this.outputChannel.appendLine("Error during formatting:")
                    this.outputChannel.append(stderr)
                    this.outputChannel.show();
                    return reject(stderr);
                } else {
                    this.outputChannel.clear();
                    this.outputChannel.hide();
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
