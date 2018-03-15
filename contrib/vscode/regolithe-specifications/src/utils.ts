'use strict';

import * as vscode from 'vscode';
import { languageId, specFileExt, absFileExt } from './const';

export const shouldConsiderDocument = (doc: vscode.TextDocument): boolean =>
    doc.languageId == languageId && (doc.fileName.endsWith(specFileExt) || doc.fileName.endsWith(absFileExt));
