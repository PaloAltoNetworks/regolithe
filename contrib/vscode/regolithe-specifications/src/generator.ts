'use strict';

import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { exec } from 'child_process';

import { shouldConsiderDocument } from './utils';
import { regoGenFileName } from './const';

export function codegen(doc: vscode.TextDocument): void {

    if (!shouldConsiderDocument(doc)) {
        return;
    }

    const docDir = path.dirname(doc.fileName)
    const p = path.join(docDir, regoGenFileName)

    if (!fs.existsSync(p)) {
        return
    }

    const cmd = fs.readFileSync(p).toString();

    exec(`cd '${docDir}' && ${cmd}`)
}
