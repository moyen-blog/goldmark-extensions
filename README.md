# Goldmark Meta

Goldmark Meta is an extension for [Goldmark](http://github.com/yuin/goldmark) that allows you to define document metadata in YAML format. Additionally, Goldmark Meta optionally generates a snipped of the markdown documents text.

Originally forked from [yuin/goldmark-meta](http://github.com/yuin/goldmark-meta), this project has undergone large, breaking changes and now stands on its own.

## Overview

Goldmark Meta provides the following functions to extend Goldmark.

1. Scrape document metadata from [YAML front matter](https://jekyllrb.com/docs/front-matter/) defined at the top of a markdown document.
2. Aggregate a configurable-length snippet of the markdown documents text content.

## Installation

```
go get github.com/moyen-blog/goldmark-meta
```

## YAML Front Matter

### Markdown Syntax

YAML front matter is metadata defined at the top of a markdown document. It can not contain markdown itself. The front matter block must be surrounded by lines containing the YAML metadata separator `---`. If metadata is defined, the separator must be the first line of the document. The YAML metadata block must end with a YAML metadata separator.

All valid YAML is supported. The underlying parsing of YAML is done by [go-yaml/yaml](https://github.com/go-yaml/yaml).

Example:

```
---
title: goldmark-meta
tags:
  - one
---

# Heading 1
```

### Accessing Metadata

```go
import (
    "bytes"
    "fmt"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    "github.com/yuin/goldmark/parser"
    meta "github.com/moyen-blog/goldmark-meta"
)

func main() {
    markdown := goldmark.New(
        goldmark.WithExtensions(
            meta.Meta,
        ),
    )
    source := `---
title: goldmark-meta
summary: Add YAML metadata to the document
tags:
  - one
---

# Hello goldmark-meta
`

    var buf bytes.Buffer
    context := parser.NewContext()
    if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
        panic(err)
    }
    // Note: Struct fields must be public in order to correctly populate the data
    out := struct {
        Title string
        Tags []string
    }{}
    if err := meta.Unmarshal(context, &out); err != nil {
        panic(err)
    }
    fmt.Println(out.Title, out.Tags)
}
```

## Markdown Text Snippet

### Markdown Syntax

There is no specific syntax needed for generating a text snippet. The snippet is made up of text from paragraphs within the markdown document only.

Example:

Example:

```
# Heading 1

Here's some text in a paragraph.

# Heading 2

More text here.
```

### Accessing Text Snippet

TODO
