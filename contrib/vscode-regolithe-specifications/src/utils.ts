'use strict';
import * as vscode from 'vscode';

export const languageId = 'yaml';
export const specFileExt = '.spec';
export const absFileExt = '.abs';


export const shouldConsiderDocument =
    (doc: vscode.TextDocument): boolean =>
        doc.languageId == languageId && (
            doc.fileName.endsWith(specFileExt)
            || doc.fileName.endsWith(absFileExt)
        );
