package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

var r *html.Renderer = html.NewRenderer(html.RendererOptions{
	Flags:          html.CommonFlags | html.HrefTargetBlank,
	RenderNodeHook: customRendering,
})

func Parse(content []byte) string {
	parsed := markdown.ToHTML(content, nil, r)

	return string(parsed)
}

func ParseString(content string) string {
	return Parse([]byte(content))
}

func customRendering(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {

	switch node := node.(type) {
	case *ast.CodeBlock:
		return processCodeBlock(w, node)
	case *ast.Paragraph:
		return processParagraph(w, node)
	case *ast.Text:
		return processText(w, node)
	default:
		return ast.GoToNext, false
	}
}

func processCodeBlock(w io.Writer, cb *ast.CodeBlock) (ast.WalkStatus, bool) {
	var attrs []string
	info := cb.Info

	literal := cb.Literal
	literal = bytes.TrimLeft(literal, "\n\r")

	if len(info) > 0 {
		endOfLang := bytes.IndexAny(info, "\t ")
		if endOfLang < 0 {
			endOfLang = len(info)
		}
		s := `class="language-` + string(info[:endOfLang]) + `"`
		attrs = append(attrs, s)
	}
	attrs = append(attrs, html.BlockAttrs(cb)...)

	write(w, `<pre class="my-2 p-2">`)
	write(w, html.TagWithAttributes("<code", attrs))

	if string(info) != "java" {
		write(w, string(literal))
	} else {
		write(w, string(addCustomClasses(prefixLinesWithSpan(literal))))
	}

	write(w, "</code>")
	write(w, "</pre>")
	write(w, "\n")

	return ast.GoToNext, true
}

func processParagraph(w io.Writer, p *ast.Paragraph) (ast.WalkStatus, bool) {
	if string(p.Literal) == "" {
		return ast.GoToNext, true
	} else {
		return ast.GoToNext, false
	}
}

func processText(w io.Writer, p *ast.Text) (ast.WalkStatus, bool) {
	if string(p.Literal) == "" {
		return ast.GoToNext, true
	} else {
		return ast.GoToNext, false
	}
}

func write(w io.Writer, value string) {
	w.Write([]byte(value))
}

type Replacement struct {
	regex    *regexp.Regexp
	newValue string
}

var replacements = []Replacement{
	{
		regexp.MustCompile(`((?:private)|(?:public))`),
		`<span class="cb-access-modifier">$1</span>`,
	},
	{
		regexp.MustCompile(`(abstract)`),
		`<span class="cb-abstract">$1</span>`},
	{
		regexp.MustCompile(`(final)`),
		`<span class="cb-final">$1</span>`,
	},
	{
		regexp.MustCompile(`((?:(?:class)|(?:interface)|(?:record)|(?:enum))\s)([A-z0-9]*)`),
		`$1<span class="cb-type-name">$2</span>`,
	},
	{
		regexp.MustCompile(`((?:class)|(?:interface)|(?:record)|(?:enum))\s`),
		`<span class="cb-type">$1</span> `,
	},
	{
		regexp.MustCompile(`((?:int)|(?:long)|(?:double)|(?:char)|(?:byte)|(?:String)|(?:Integer)|(?:Long)|(?:Double)|(?:Character)|(?:BigInteger)|(?:Optional<.+>)|(?:Stream<.+>))\s`),
		`<span class="cb-type-name">$1</span> `,
	},
	{
		regexp.MustCompile(`(\s*)([A-z0-9]+)(\(.*\))`),
		`$1<span class="cb-method">$2</span>$3`,
	},
}

func addCustomClasses(value []byte) []byte {

	builder := bytes.NewBuffer([]byte{})
	html.Escape(builder, value)
	escapedValue := builder.Bytes()

	for _, v := range replacements {
		if v.regex.Match(escapedValue) {
			escapedValue = v.regex.ReplaceAll(escapedValue, []byte(v.newValue))
		}

	}

	return escapedValue
}

func prefixLinesWithSpan(value []byte) []byte {
	builder := bytes.NewBuffer([]byte{})

	splits := bytes.Split(value, []byte("\n"))
	splitsNo := len(splits) - 1

	for i, v := range splits {
		if i == splitsNo && bytes.Equal(v, []byte{}) {
			continue
		}

		builder.Write([]byte(fmt.Sprintf("<span class=\"line-number\">%02d  </span>%v\n", i+1, string(v))))
	}

	return builder.Bytes()
}
