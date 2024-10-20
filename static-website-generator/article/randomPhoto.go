package article

import (
	"bytes"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	goldmarkTailwindcss "markdown2tailwindcss/goldmark-tailwindcss"
	"os"
	"path"
	"sort"
)

type Article struct {
	ProjectDirname string
	Title          string
	ImgSrc         string
	Order          int
	Content        string
	Dependencies   []string
}

type ArticleSlice []Article

func (s ArticleSlice) Len() int           { return len(s) }
func (s ArticleSlice) Less(i, j int) bool { return s[i].Order < s[j].Order }
func (s ArticleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

var rootPath = path.Join("../")

var mdConverter = goldmark.New(
	goldmark.WithExtensions(
		meta.Meta,
	),
	goldmark.WithRenderer(
		renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(html.NewRenderer(), 1000))),
	),
	goldmark.WithParser(
		parser.NewParser(parser.WithBlockParsers(parser.DefaultBlockParsers()...),
			parser.WithInlineParsers(parser.DefaultInlineParsers()...),
			parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
			parser.WithASTTransformers(
				util.Prioritized(goldmarkTailwindcss.NewTransformer(), 100),
			),
		),
	),
)

func processDirectory(dir os.DirEntry) (Article, error) {
	var article = Article{
		ProjectDirname: dir.Name(),
		Order:          -1,
	}

	articleMd, err := os.ReadFile(path.Join(rootPath, dir.Name(), "README.md"))
	if err != nil {
		return Article{}, err
	}

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := mdConverter.Convert(articleMd, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)
	if metaData["Title"] != nil {
		article.Title = metaData["Title"].(string)
	}
	if metaData["ImgSrc"] != nil {
		article.ImgSrc = metaData["ImgSrc"].(string)
	}
	if metaData["Order"] != nil {
		article.Order = metaData["Order"].(int)
	}
	if metaData["Dependencies"] != nil {
		for _, dep := range metaData["Dependencies"].([]interface{}) {
			article.Dependencies = append(article.Dependencies, dep.(string))
		}
	}
	article.Content = buf.String()

	return article, nil
}

func GetArticlesData() ArticleSlice {
	var articles ArticleSlice

	root, err := os.ReadDir(rootPath)
	if err != nil {
		panic(err)
	}
	for _, dir := range root {
		if dir.IsDir() {
			a, err := processDirectory(dir)
			if err != nil {
				continue
			}
			articles = append(articles, a)
		}
	}

	sort.Sort(articles)

	return articles
}
