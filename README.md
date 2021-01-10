# goldmark-meta

goldmark-meta is an extension for the [goldmark](http://github.com/yuin/goldmark) that allows you to define document metadata in YAML format.

## Changes

There are significant changes from the upstream [yuin/goldmark-meta](http://github.com/yuin/goldmark-meta). The extension is largely simplified and has undergone the following changes.
- Rendering a table is no longer supported
- YAML metadata is accessed via a typical `Unmarshal` interface
- YAML metadata must use a separator line of minimum three `-` characters

## Usage

### Installation

```
go get github.com/moyen-blog/goldmark-meta
```

### Markdown syntax

YAML metadata block is a leaf block that can not have any markdown element as a child.

YAML metadata must start with a **YAML metadata separator**. This separator must be at first line of the document.

A **YAML metadata separator** is a line that only `-` is repeated three or more times.

YAML metadata must end with a **YAML metadata separator**.

You can define objects as a 1st level item. At deeper level, you can define any kind of YAML element.

Example:

```
---
title: goldmark-meta
summary: Add YAML metadata to the document
tags:
  - markdown
  - goldmark
---

# Heading 1
```

### Access the metadata

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
  - markdown
  - goldmark
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
    }{}
    if err := meta.Unmarshal(context, &out); err != nil {
        panic(err)
    }
    fmt.Print(out.Title)
}
```
