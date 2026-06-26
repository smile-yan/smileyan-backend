package utils

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownToHTML 将 Markdown 转换为 HTML（与前端 Editor.vue 渲染效果一致）
func MarkdownToHTML(markdown []byte) []byte {
	// 创建 goldmark 实例，配置与前端 markdown-it 一致
	// 启用 Raw HTML 以支持内嵌 SVG 等 HTML 元素
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Linkify,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithRawHtml(),
		),
		goldmark.WithRendererOptions(
			html.WithRawHtml(),
		),
	)

	// 渲染为 HTML
	var buf bytes.Buffer
	err := md.Convert(markdown, &buf)
	if err != nil {
		return []byte("")
	}

	// 后处理：与前端 markdown-it 渲染结果一致
	result := postProcessHTML(buf.String())

	return []byte(result)
}

// postProcessHTML 将 goldmark 输出的 HTML 转换为与前端 markdown-it 一致的格式
func postProcessHTML(html string) string {
	// 1. 代码块处理 - 使用正则表达式一次性替换
	// goldmark: <pre><code class="language-go">code</code></pre>
	// 前端: <pre class="language-go"><code>code</code></pre>
	html = replaceCodeBlock(html)

	// 2. 先处理图片，添加标题（支持 Markdown 图片语法：![标题](url)）
	// 必须在标题处理之前执行，避免图片被包含在标题标签中
	reImg := regexp.MustCompile(`<img\s+src="([^"]+)"\s+alt="([^"]+)"\s*/?>`)
	html = reImg.ReplaceAllString(html, `<figure class="article-figure"><img class="article-image" src="$1" alt="$2"><figcaption class="article-figcaption">$2</figcaption></figure>`)

	// 3. 标题添加 class
	html = strings.ReplaceAll(html, "<h1>", `<h1 class="article-h1">`)
	html = strings.ReplaceAll(html, "<h2>", `<h2 class="article-h2">`)
	html = strings.ReplaceAll(html, "<h3>", `<h3 class="article-h3">`)
	html = strings.ReplaceAll(html, "<h4>", `<h4 class="article-h4">`)
	html = strings.ReplaceAll(html, "<h5>", `<h5 class="article-h5">`)
	html = strings.ReplaceAll(html, "<h6>", `<h6 class="article-h6">`)

	// 4. 其他元素添加 class
	html = strings.ReplaceAll(html, "<p>", `<p class="article-p">`)
	html = strings.ReplaceAll(html, "<blockquote>", `<blockquote class="article-blockquote">`)
	html = strings.ReplaceAll(html, "<ul>", `<ul class="article-list">`)
	html = strings.ReplaceAll(html, "<ol>", `<ol class="article-list">`)
	html = strings.ReplaceAll(html, "<li>", `<li class="article-list-item">`)
	html = strings.ReplaceAll(html, "<a ", `<a class="article-link" target="_blank" rel="noopener noreferrer" `)

	// 先移除表格的 inline style，再添加 class
	html = removeTableInlineStyles(html)

	html = strings.ReplaceAll(html, "<table>", `<table class="article-table">`)
	html = strings.ReplaceAll(html, "<thead>", `<thead class="article-thead">`)
	html = strings.ReplaceAll(html, "<tbody>", `<tbody class="article-tbody">`)
	html = strings.ReplaceAll(html, "<tr>", `<tr class="article-tr">`)
	html = strings.ReplaceAll(html, "<th>", `<th class="article-th">`)
	html = strings.ReplaceAll(html, "<td>", `<td class="article-td">`)
	html = strings.ReplaceAll(html, "<hr>", `<hr class="article-hr" />`)
	html = strings.ReplaceAll(html, "<strong>", `<strong class="article-strong">`)
	html = strings.ReplaceAll(html, "<em>", `<em class="article-em">`)

	return html
}

// replaceCodeBlock 处理代码块，使用正则表达式
func replaceCodeBlock(html string) string {
	// 处理有语言的代码块: <pre><code class="language-python">code</code></pre>
	// 替换为: <pre class="language-python"><code>code</code></pre>
	// 使用 [\s\S]*? 来匹配包括换行符在内的任意字符
	re := regexp.MustCompile(`(?s)<pre><code class="language-([^"]+)">([\s\S]*?)</code></pre>`)
	html = re.ReplaceAllString(html, `<pre class="language-$1 line-numbers"><code>$2</code></pre>`)

	// 处理没有语言的代码块: <pre><code>code</code></pre>
	// 添加 line-numbers class
	re2 := regexp.MustCompile(`(?s)<pre><code>([\s\S]*?)</code></pre>`)
	html = re2.ReplaceAllString(html, `<pre class="line-numbers"><code>$1</code></pre>`)

	return html
}

// removeTableInlineStyles 移除表格元素的 inline style 并修复结构
func removeTableInlineStyles(html string) string {
	// goldmark 表格输出: <table><tbody><tr><th>...</th></tr><tr><td>...</td></tr></tbody></table>
	// 正确结构: <table><thead><tr><th>...</th></tr></thead><tbody><tr><td>...</td></tr></tbody></table>

	// 1. 移除表格元素的所有属性（包括 style）
	reClean := regexp.MustCompile(`<(thead|tbody|tr|th|td)(?:\s+(?:style|class| colspan| rowspan|[a-z-]+)="[^"]*")*>`)
	html = reClean.ReplaceAllString(html, "<$1>")
	html = regexp.MustCompile(`<table(?:\s+(?:style|class|[a-z-]+)="[^"]*")*>`).ReplaceAllString(html, "<table>")

	// 使用正则表达式匹配整个表格
	re := regexp.MustCompile(`(?s)<table[^>]*>(.*?)</table>`)
	html = re.ReplaceAllStringFunc(html, func(match string) string {
		// 匹配表头行（包含 th）和数据行（包含 td）
		re2 := regexp.MustCompile(`(?s)(<tbody[^>]*>)\s*(<tr[^>]*>.*?<th[^>]*>.*?</th>.*?</tr>)\s*(.*?</tbody>)`)
		result := re2.ReplaceAllString(match, "<table class=\"article-table\"><thead class=\"article-thead\">$2</thead><tbody class=\"article-tbody\">$3")

		// 添加 class
		result = strings.ReplaceAll(result, "<thead>", `<thead class="article-thead">`)
		result = strings.ReplaceAll(result, "<tbody>", `<tbody class="article-tbody">`)
		result = strings.ReplaceAll(result, "<tr>", `<tr class="article-tr">`)
		result = strings.ReplaceAll(result, "<th>", `<th class="article-th">`)
		result = strings.ReplaceAll(result, "<td>", `<td class="article-td">`)

		return result
	})

	// 确保 table 标签有 class
	html = strings.ReplaceAll(html, "<table>", `<table class="article-table">`)

	return html
}