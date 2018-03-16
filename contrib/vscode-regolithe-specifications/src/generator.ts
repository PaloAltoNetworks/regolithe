'use strict';

import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { exec } from 'child_process';

import { shouldConsiderDocument } from './utils';


export class RegolitheGenerator {

    regoGenFileName: string;

    constructor(regoGenFileName: string) {
        this.regoGenFileName = regoGenFileName
    }

    public generate(doc: vscode.TextDocument): void {

        if (!shouldConsiderDocument(doc)) {
            return;
        }

        const docDir = path.dirname(doc.fileName)
        const p = path.join(docDir, this.regoGenFileName)

        if (!fs.existsSync(p)) {
            return
        }

        const cmd = fs.readFileSync(p).toString();

        exec(`cd '${docDir}' && ${cmd}`)

        console.log('Regolithe: model generated for', doc.fileName);
    }
}
