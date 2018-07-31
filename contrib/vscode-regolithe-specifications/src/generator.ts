'use strict';

import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { exec } from 'child_process';

import { shouldConsiderDocument } from './utils';


export class RegolitheGenerator {

    regoGenFileName: string;
    outputChannel: vscode.OutputChannel

    constructor(regoGenFileName: string, outputChannel: vscode.OutputChannel) {
        this.regoGenFileName = regoGenFileName;
        this.outputChannel = outputChannel;
    }

    public generate(doc: vscode.TextDocument): void {

        if (!shouldConsiderDocument(doc)
            && !doc.fileName.endsWith('_type.mapping')
            && !doc.fileName.endsWith('_api.info')
            && !doc.fileName.endsWith('_parameters')
        ) {
            return;
        }

        const docDir = path.dirname(doc.fileName)
        const p = path.join(docDir, this.regoGenFileName)

        if (!fs.existsSync(p)) {
            return
        }

        const cmd = fs.readFileSync(p).toString();

        exec(`cd '${docDir}' && ${cmd}`, (err: Error, stdout: string, stderr: string) => {
            if (err) {
                this.outputChannel.clear();
                this.outputChannel.appendLine("Error during generation:")
                this.outputChannel.append(stderr);
                this.outputChannel.show();
            } else {
                this.outputChannel.clear();
                this.outputChannel.hide();
            }
        })

        console.log('Regolithe: model generated for', doc.fileName);
    }
}
