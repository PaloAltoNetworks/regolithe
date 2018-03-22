'use strict';
import * as vscode from 'vscode';
import * as path from 'path';

import { RegolitheDocumentFormattingEditProvider } from './formatter';
import { RegolitheGenerator } from './generator';


export function activate(ctx: vscode.ExtensionContext) {

    const regoPath = path.join(ctx.extensionPath, 'bin', 'rego')
    const formatter = new RegolitheDocumentFormattingEditProvider(regoPath);

    const regoGenFileName = '.regolithe-gen-cmd';
    const generator = new RegolitheGenerator(regoGenFileName);

    ctx.subscriptions.push(
        vscode.workspace.onWillSaveTextDocument(
            (e: vscode.TextDocumentWillSaveEvent): void => {
                e.waitUntil(formatter.format(e.document).then(edits => edits))
            }
        ),
    );

    ctx.subscriptions.push(
        vscode.workspace.onDidSaveTextDocument(
            (doc: vscode.TextDocument) =>
                generator.generate(doc)
        ),
    );
}

// this method is called when your extension is deactivated
export function deactivate() {
}
