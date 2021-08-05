# GoldMark Extensions

GoldMark Extensions is a collection of extensions for [GoldMark](http://github.com/yuin/goldmark).

Originally forked from [yuin/goldmark-meta](http://github.com/yuin/goldmark-meta), this project has undergone breaking changes and now stands on its own.

## Overview

GoldMark Extensions provides the following functions to extend GoldMark.

1. Scrape document metadata from [YAML front matter](https://jekyllrb.com/docs/front-matter/) defined at the top of a Markdown document.
2. Aggregate a configurable-length snippet of the Markdown documents text content.
3. Add attribute `class="task"` to tasklist items.

## Installation

```
go get github.com/moyen-blog/goldmark-extensions
```

## YAML Front Matter

Include the GoldMark extension with `extensions.MetadataExtension`.

### Markdown Syntax

YAML front matter is metadata defined at the top of a Markdown document. It can not contain Markdown itself. The front matter block must be surrounded by lines containing the YAML metadata separator `---`. If metadata is defined, the separator must be the first line of the document. The YAML metadata block must end with a YAML metadata separator.

All valid YAML is supported. The underlying parsing of YAML is done by [go-yaml/yaml](https://github.com/go-yaml/yaml).

Example:

```
---
title: goldmark-extensions
tags:
  - one
---

# Heading 1
```

### Accessing Metadata

The metadata can be parsed via the `meta.Unmarshal()` function.

```go
import (
    "bytes"
    "fmt"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/parser"
    extensions "github.com/moyen-blog/goldmark-extensions"
)

func main() {
    markdown := goldmark.New(
        goldmark.WithExtensions(
            extensions.MetadataExtension,
        ),
    )
    source := `---
title: goldmark-extensions
summary: Add YAML metadata to the document
tags:
  - one
---

# Hello goldmark-extensions
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
    if err := extensions.Unmarshal(context, &out); err != nil {
        panic(err)
    }
    fmt.Println(out.Title, out.Tags)
}
```

## Markdown Text Snippet

Include the GoldMark extension with `extensions.SnippetExtension(max)` where `max` is the maximum length of the generated snippet.

### Markdown Syntax

There is no specific syntax needed for generating a text snippet. The snippet is made up of text from paragraphs within the Markdown document only.

Example:

```
# Heading 1

Here's some text in a paragraph.

# Heading 2

More text here.
```

### Accessing Text Snippet

The generated snippet is accessed via the `extensions.Snippet()` function.

```go
import (
    "bytes"
    "fmt"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/parser"
    extensions "github.com/moyen-blog/goldmark-extensions"
)

func main() {
    var markdown = goldmark.New(
        goldmark.WithExtensions(
            extensions.SnippetExtension(100), // Maximum length of snippet
        ),
    )
    source := `# Hello
Paragraph text here.

Another one.
And continued here.`

    var buf bytes.Buffer
    context := parser.NewContext()
    if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
        panic(err)
    }
    s, err := extensions.Snippet(context)
    if err != nil {
        panic(err)
    }
    fmt.Println(s)
}
```


## Tasklist Styling

Include the GoldMark extension with `extensions.TasklistExtension`. This extension relies on the built-in [GoldMark tasklist extension](https://github.com/yuin/goldmark#built-in-extensions).

### Markdown Syntax

Standard tasklist Markdown syntax applies. With the class applied, further styling can be applied specifically to tasklist items while leaving standard list items unaffected e.g. `li.task`.

Example:

```
# Markdown Tasklist

Here's my list:
- [ ] Unchecked task
- Plaintext list item
- [ ] Checked task
```

In the previous example, only the first and third items will have `class="task"` inline styling applied.
