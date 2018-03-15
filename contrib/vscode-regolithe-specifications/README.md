# Regolithe Specifications for VSCode

This package provides support for editing Regolithe Specifications.

## What is Regolithe

For now it is internal. You mostly don't need this extension.

Yet :)

## Features

### Schema validation

The VSCode Regolithe extension provides auto completion, validation and documentation for Regolithe specification files using the integrated json schema.

![schema](https://imgur.com/wDQcpt1.gif)

### Snippets

To quicky edit the specifications, the extension adds the following snippets:

- `spec<tab>`: creates a skeleton for a new specification file
- `attr<tab>`: adds a new attribute
- `rel<tab>`: adds a new relation

![snippets](https://imgur.com/tEWJE6D.gif)

### Linting

In order to keep your specification files clean, the extension will sort and lint the files on save, keeping things in order.

![lint](https://imgur.com/jsEzzor.gif)

### Auto codegen

You can create a file named `.regolithe-gen-cmd` in the root folder of the specification that contains the command to run to generate your specification.

For instance:

```shell
elegen folder -d . -o ..
```

![codegen](https://imgur.com/f6Tn1Yr.gif)
