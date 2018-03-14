'use strict';
import * as vscode from 'vscode';
import { RegolitheDocumentFormattingEditProvider } from './formatter';
import * as path from 'path';

export function activate(ctx: vscode.ExtensionContext) {

    const regoPath = path.join(ctx.extensionPath, 'bin', 'rego')

    ctx.subscriptions.push(
        vscode.languages.registerDocumentFormattingEditProvider(
            'json',
            new RegolitheDocumentFormattingEditProvider(regoPath),
        ),
    );
}

// this method is called when your extension is deactivated
export function deactivate() {
}
