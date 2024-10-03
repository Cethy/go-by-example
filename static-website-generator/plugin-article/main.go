package plugin_article

import (
	"bytes"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	goldmarkTailwindcss "go-by-example/libs/goldmark-tailwindcss"
	"go-by-example/static-website-generator/generator"
	pluginfragment "go-by-example/static-website-generator/plugin-fragment"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

var pluginFragmentArticleFilename = "article"
var pluginFragmentListItemFilename = "article-list-item"
var pluginFragmentListItemId = "{index-articles}"

type Article struct {
	Filename string
	Title    string
	ImgSrc   string
	Order    int
	Content  string
}

type ArticleSlice []Article

func (s ArticleSlice) Len() int           { return len(s) }
func (s ArticleSlice) Less(i, j int) bool { return s[i].Order < s[j].Order }
func (s ArticleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func getArticlesData() ArticleSlice {
	var articles ArticleSlice

	mdConverter := goldmark.New(
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

	rootPath := path.Join("./")
	root, err := os.ReadDir(rootPath)
	if err != nil {
		panic(err)
	}
	for _, dir := range root {
		if dir.IsDir() {
			articleMd, err := os.ReadFile(path.Join(rootPath, dir.Name(), "article.md"))
			if err != nil {
				continue
			}

			var buf bytes.Buffer
			context := parser.NewContext()
			if err := mdConverter.Convert(articleMd, &buf, parser.WithContext(context)); err != nil {
				panic(err)
			}
			metaData := meta.Get(context)

			articles = append(articles, Article{
				Filename: dir.Name() + ".html",
				Title:    metaData["Title"].(string),
				ImgSrc:   metaData["ImgSrc"].(string),
				Order:    metaData["Order"].(int),
				Content:  buf.String(),
			})
		}
	}

	sort.Sort(articles)

	return articles
}

func getIndexArticleListItem(articleData Article, srcDir string) string {
	article := pluginfragment.GetFragmentContent(pluginFragmentListItemFilename, srcDir)

	article = strings.ReplaceAll(article, "{link}", articleData.Filename)
	article = strings.ReplaceAll(article, "{title}", articleData.Title)
	article = strings.ReplaceAll(article, "{imgSrc}", articleData.ImgSrc)

	return article
}

func getIndexArticleList(srcDir string) string {
	articles := getArticlesData()
	list := ""
	for _, article := range articles {
		list = list + getIndexArticleListItem(article, srcDir)
	}
	return list
}

var articleFilesBuilt []string

func buildArticleFile(srcDir, targetPath string, article Article) {
	if slices.Contains(articleFilesBuilt, targetPath) {
		return
	}

	articleRaw := pluginfragment.GetFragmentContent(pluginFragmentArticleFilename, srcDir)

	articleRaw = strings.ReplaceAll(articleRaw, "{title}", article.Title)
	articleRaw = strings.ReplaceAll(articleRaw, "{imgSrc}", article.ImgSrc)
	articleRaw = strings.ReplaceAll(articleRaw, "{content}", article.Content)
	fileContent, err := generator.PreProcess(articleRaw)
	if err != nil {
		panic(err)
	}

	articleFilesBuilt = append(articleFilesBuilt, targetPath)
	generator.BuildFile(targetPath, fileContent)
}

var pluginIsRunning = false

func PreBuildFile(fileContent string, config generator.Config) (string, error) {
	if pluginIsRunning {
		return fileContent, nil
	}
	pluginIsRunning = true

	fileContent = strings.ReplaceAll(fileContent, pluginFragmentListItemId, getIndexArticleList(config.SrcDir))

	articles := getArticlesData()
	for _, article := range articles {
		buildArticleFile(config.SrcDir, filepath.Join(config.OutputDir, article.Filename), article)
	}

	pluginIsRunning = false

	return fileContent, nil
}
