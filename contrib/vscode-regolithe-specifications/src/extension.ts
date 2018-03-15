'use strict';
import * as vscode from 'vscode';
import { RegolitheDocumentFormattingEditProvider } from './formatter';
import * as path from 'path';

import { codegen } from './generator';
import { languageId } from './const';

export function activate(ctx: vscode.ExtensionContext) {

    const regoPath = path.join(ctx.extensionPath, 'bin', 'rego')

    ctx.subscriptions.push(
        vscode.languages.registerDocumentFormattingEditProvider(
            languageId,
            new RegolitheDocumentFormattingEditProvider(regoPath),
        ),
    );

    vscode.workspace.onDidSaveTextDocument(codegen)
}

// this method is called when your extension is deactivated
export function deactivate() {
}
