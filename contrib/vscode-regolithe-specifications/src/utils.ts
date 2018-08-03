'use strict';
import * as vscode from 'vscode';

export const languageId = 'yaml';
export const specFileExt = '.spec';
export const absFileExt = '.abs';
export const tmFile = '_type.mapping';
export const vmFile = '_validation.mapping';
export const pmFile = '_parameter.mapping';

export const shouldConsiderDocument =
    (doc: vscode.TextDocument): boolean =>
        doc.languageId == languageId && (
            doc.fileName.endsWith(specFileExt)
            || doc.fileName.endsWith(absFileExt)
            || doc.fileName.endsWith(tmFile)
            || doc.fileName.endsWith(vmFile)
            || doc.fileName.endsWith(pmFile)
        );
