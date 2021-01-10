package meta

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

var markdown = goldmark.New(
	goldmark.WithExtensions(
		New(),
	),
)

func TestMeta(t *testing.T) {
	source := `---
title: goldmark-meta
tags:
  - one
ignored: hi
---
# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct {
		Title string
		Tags  []string
	}{}
	if err := Unmarshal(context, &out); err != nil {
		t.Errorf("YAML unmarshal failed: %v", err.Error())
	}
	if out.Title != "goldmark-meta" {
		t.Errorf("Title must be 'goldmark-meta', but got '%s'", out.Title)
	}
	if len(out.Tags) != 1 {
		t.Error("Tags must be a slice that has 1 element")
	}
	if out.Tags[0] != "one" {
		t.Errorf("First tag must be 'one' but got '%s'", out.Tags[0])
	}
	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("Should render '<h1>Hello goldmark-meta</h1>', but got '%s'", buf.String())
	}
}

func TestNoMeta(t *testing.T) {
	source := `# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct {
		Title string
	}{}
	if err := Unmarshal(context, &out); err != nil {
		t.Error("YAML unmarshal failed")
	}
	if out.Title != "" {
		t.Errorf("Title must be empty, but got '%s'", out.Title)
	}
	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("Should render '<h1>Hello goldmark-meta</h1>', but got '%s'", buf.String())
	}
}

func TestEmptyMeta(t *testing.T) {
	source := `---
---
# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct {
		Title string
	}{}
	if err := Unmarshal(context, &out); err != nil {
		t.Error("YAML unmarshal failed")
	}
	if out.Title != "" {
		t.Errorf("Title must be empty, but got '%s'", out.Title)
	}
	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("Should render '<h1>Hello goldmark-meta</h1>', but got '%s'", buf.String())
	}
}

func TestMetaError(t *testing.T) {
	source := `---
bad:
  - : {
  }
    - markdown
    - goldmark
---
# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	out := struct{}{}
	if err := Unmarshal(context, &out); err == nil {
		t.Error("Should throw unmarshal error")
	}
}

func TestInvalidMetaBuffer(t *testing.T) {
	source := `# Hello goldmark-meta`

	var buf bytes.Buffer
	context := parser.NewContext()
	context.Set(contextKey, 0) // Not the expected bytes.Buffer
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(parser.NewContext())); err != nil {
		panic(err)
	}
	out := struct {
		Title string
	}{}
	if err := Unmarshal(context, &out); err == nil {
		t.Error("Should throw unavailable YAML buffer error")
	}
}
