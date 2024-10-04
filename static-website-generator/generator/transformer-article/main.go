package transformer_article

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
	transformerfragment "go-by-example/static-website-generator/generator/transformer-fragment"
	"os"
	"path"
	"slices"
	"sort"
	"strings"
)

type Config struct {
	fragmentArticleFilename    string
	fragmentListItemFilename   string
	fragmentListItemId         string
	fragmentGithubLinkFilename string
}

func GetDefaultConfig() Config {
	return Config{
		fragmentArticleFilename:    "article",
		fragmentListItemFilename:   "article-list-item",
		fragmentListItemId:         "{index-articles}",
		fragmentGithubLinkFilename: "article-github-link",
	}
}

type ArticleTransformer interface {
	transform(fileContent string) (string, error)
}

type articleTransformer struct {
	Config              Config
	fragmentTransformer transformerfragment.FragmentTransformer
	generator           generator.Generator
	pluginIsRunning     bool
	articleFilesBuilt   []string
}

type Article struct {
	ProjectDirname string
	Title          string
	ImgSrc         string
	Order          int
	Content        string
}

type ArticleSlice []Article

func (s ArticleSlice) Len() int           { return len(s) }
func (s ArticleSlice) Less(i, j int) bool { return s[i].Order < s[j].Order }
func (s ArticleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (p *articleTransformer) getArticlesData() ArticleSlice {
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
				ProjectDirname: dir.Name(),
				//ProjectDirname: "projects/" + dir.Name() + ".html",
				Title:   metaData["Title"].(string),
				ImgSrc:  metaData["ImgSrc"].(string),
				Order:   metaData["Order"].(int),
				Content: buf.String(),
			})
		}
	}

	sort.Sort(articles)

	return articles
}

func (p *articleTransformer) getGithubLink(articleData Article) string {
	fragment := p.fragmentTransformer.GetFragmentContent(p.Config.fragmentGithubLinkFilename)
	fragment = strings.ReplaceAll(fragment, "{projectDir}", articleData.ProjectDirname)

	return fragment
}

func (p *articleTransformer) getTargetArticleRelativePath(articleData Article) string {
	// @todo add to config
	return "projects/" + articleData.ProjectDirname + ".html"
}

func (p *articleTransformer) getIndexArticleListItem(articleData Article) string {
	article := p.fragmentTransformer.GetFragmentContent(p.Config.fragmentListItemFilename)

	article = strings.ReplaceAll(article, "{link}", p.getTargetArticleRelativePath(articleData))
	article = strings.ReplaceAll(article, "{title}", articleData.Title)
	article = strings.ReplaceAll(article, "{imgSrc}", articleData.ImgSrc)

	return article
}

func (p *articleTransformer) getIndexArticleList() string {
	articles := p.getArticlesData()
	list := ""
	for _, article := range articles {
		list = list + p.getIndexArticleListItem(article)
	}
	return list
}

func (p *articleTransformer) buildArticleFile(article Article) {
	relativeFilePath := p.getTargetArticleRelativePath(article)

	if slices.Contains(p.articleFilesBuilt, relativeFilePath) {
		return
	}

	articleRaw := p.fragmentTransformer.GetFragmentContent(p.Config.fragmentArticleFilename)

	articleRaw = strings.ReplaceAll(articleRaw, "{title}", article.Title)
	articleRaw = strings.ReplaceAll(articleRaw, "{imgSrc}", article.ImgSrc)
	articleRaw = strings.ReplaceAll(articleRaw, "{content}", article.Content)
	articleRaw = strings.ReplaceAll(articleRaw, "{github-link}", p.getGithubLink(article))
	fileContent, err := p.generator.InvokeTransformers(articleRaw)
	if err != nil {
		panic(err)
	}

	p.generator.BuildFile(relativeFilePath, fileContent)
	p.articleFilesBuilt = append(p.articleFilesBuilt, relativeFilePath)
}

func (p *articleTransformer) Transform(fileContent string) (string, error) {
	if p.pluginIsRunning {
		return fileContent, nil
	}
	p.pluginIsRunning = true

	fileContent = strings.ReplaceAll(fileContent, p.Config.fragmentListItemId, p.getIndexArticleList())

	articles := p.getArticlesData()
	for _, article := range articles {
		p.buildArticleFile(article)
	}

	p.pluginIsRunning = false

	return fileContent, nil
}

func NewTransformer(config Config, generator generator.Generator, fragmentTransformer transformerfragment.FragmentTransformer) generator.Transformer {
	return &articleTransformer{
		Config:              config,
		generator:           generator,
		fragmentTransformer: fragmentTransformer,
		pluginIsRunning:     false,
	}
}
