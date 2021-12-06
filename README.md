# Celli aka Cellementatry cli

`Celli` is a cli util to work with [Cellementary](https://github.com/MonkeyBuisness/cellementary-extension) extension notebooks.

It allows you to write notebooks content inside the `markdown` file and the convert it to the notebook format that extension can render in the VS Code editor.

## Example usage

Imagine you have a markdown document named `example.md` contains `java-code` inserts.

Then command
```console
$ celli convert t2b example.md > example.javabook
```
will convert **example.md** file to the **example.javabook** file that you can open in the VS Code editor.

Then command
```console
$ celli convert b2t example.javabook > example.md
```
will convert **example.javabook** file to the **example.md** (will do the reverse conversion).

To see more usage options run
```console
$ celli --help
``` 

## Installation

Go to the [Releases](https://github.com/MonkeyBuisness/celli/releases) page and download the latest release to your machine.

### Install from source (Not recommended)
```sh
# Go 1.16+
go install -o github.com/MonkeyBuisness/celli@latest

# Go version < 1.16
go get -u github.com/MonkeyBuisness/celli@latest
```

## Serializable comments

As the main idea of this extension is allow you to create notebooks without VS Code editor, it's very important to provide opportunity to have full control on the notebook creation process.

For this case util allows to insert special comments  `<!-- -->` (we will call it `serializable comments`) with metadata that will be serialized in the special way to the notebook content.

The supported serializable comments:

1. ```html
    <!-- notebook:{
        "any": "metadata here",
        ...
    } -->
    ```
    will be transformed to the
    ```json
    {
        "metadata": {
            "any": "metadata here"
        },
        "cells": []
    }
    ```
    during the convertaion process.
    In other words, the `<!-- notebook:{} -->` comment uses for setting notebook metadata.

2. ```html
    # Hello
    
    <!-- br: -->

    # World
    ```
    will be transformed to the
    ```json
    {
        "metadata": {
           
        },
        "cells": [
            {
                "languageId": "markdown",
                "kind": 1,
                "content": "# Hello"
            },
            {
                "languageId": "markdown",
                "kind": 1,
                "content": "# World"
            }
        ]
    }
    ```
    during the convertaion process.
    In other words, the `<!-- br: -->` comment uses to split markup text into cells.
3. ```html
    <!-- autors:[
        {
            "name": "John Smith",
            "avatar": "https://example.com/john-smith.png",
            "link": "https://example.com/john-smith",
            "about": "Java developer..."
        }
    ] -->
    ```
    will be transformed to the
    ```json
    {
        "metadata": {
           
        },
        "cells": [
            {
                "languageId": "markdown",
                "kind": 1,
                "content": "<authors info table>"
            }
        ]
    }
    ```
    during the convertaion process.
    In other words, the `<!-- authors:[] -->` comment uses to add authors info to the notebook document.
4. ```html
    <!-- code:{
        "lang": "java",
        "meta": {
            "is-executable": "false"
        },
        "content": "package example;"
    } -->
    ```
    will be transformed to the
    ```json
    {
        "metadata": {
           
        },
        "cells": [
            {
                "languageId": "java",
                "kind": 2,
                "content": "package main;",
                "metadata": {
                     "is-executable": "false"
                },
                "uri": ""
            }
        ]
    }
    ```
    during the convertaion process.
    In other words, the `<!-- code:{} -->` comment uses to add code cell the notebook document.
    > if **uri** field is provided, then **content** field of the cell will be overwritten with the content of the provided URI. The uri may contain path to the local file (`file:///home/examples/Main.java`) or link to the remote file (`https://www.github.com/test-repo/main/blob/Main.java`).

See more examples [here](https://github.com/MonkeyBuisness/celli/tree/master/example).
