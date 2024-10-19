package transformer_article

import (
	"go-by-example/static-website-generator/article"
	"go-by-example/static-website-generator/generator"
	transformerfragment "go-by-example/static-website-generator/generator/transformer-fragment"
	"path"
	"path/filepath"
	"slices"
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

var rootPath = path.Join("./")

func (p *articleTransformer) getGithubLink(articleData article.Article) string {
	fragment := p.fragmentTransformer.GetFragmentContent(p.Config.fragmentGithubLinkFilename)
	fragment = strings.ReplaceAll(fragment, "{projectDir}", articleData.ProjectDirname)

	return fragment
}

func (p *articleTransformer) getTargetArticleRelativePath(articleData article.Article) string {
	// @todo add to config
	return "projects/" + articleData.ProjectDirname + ".html"
}

func getImgSrc(original string) string {
	if original[:4] != "http" {
		// if it's a local file, prefix basePublicPath
		return "{basePublicPath}" + original
	}
	return original
}

func (p *articleTransformer) getIndexArticleListItem(data article.Article) string {
	html := p.fragmentTransformer.GetFragmentContent(p.Config.fragmentListItemFilename)

	html = strings.ReplaceAll(html, "{link}", p.getTargetArticleRelativePath(data))
	html = strings.ReplaceAll(html, "{title}", data.Title)

	html = strings.ReplaceAll(html, "{imgSrc}", getImgSrc(data.ImgSrc))

	return html
}

func (p *articleTransformer) getIndexArticleList() string {
	articles := article.GetArticlesData()
	list := ""
	for _, data := range articles {
		list = list + p.getIndexArticleListItem(data)
	}
	return list
}

func (p *articleTransformer) buildArticleFile(data article.Article) {
	relativeFilePath := p.getTargetArticleRelativePath(data)

	if slices.Contains(p.articleFilesBuilt, relativeFilePath) {
		return
	}

	articleRaw := p.fragmentTransformer.GetFragmentContent(p.Config.fragmentArticleFilename)

	articleRaw = strings.ReplaceAll(articleRaw, "{title}", data.Title)
	articleRaw = strings.ReplaceAll(articleRaw, "{imgSrc}", getImgSrc(data.ImgSrc))
	articleRaw = strings.ReplaceAll(articleRaw, "{content}", data.Content)
	articleRaw = strings.ReplaceAll(articleRaw, "{github-link}", p.getGithubLink(data))

	// copy & update dependencies
	if len(data.Dependencies) > 0 {
		for _, dep := range data.Dependencies {
			relativeTargetPath := filepath.Join("projects/", data.ProjectDirname, dep)
			p.generator.CopyFile(filepath.Join(rootPath, data.ProjectDirname, dep), relativeTargetPath)

			articleRaw = strings.ReplaceAll(articleRaw, dep, filepath.Join(data.ProjectDirname, dep))
		}
	}

	p.generator.BuildFile(relativeFilePath, articleRaw)
	p.articleFilesBuilt = append(p.articleFilesBuilt, relativeFilePath)
}

func (p *articleTransformer) Transform(fileContent string) (string, error) {
	if p.pluginIsRunning {
		return fileContent, nil
	}
	p.pluginIsRunning = true

	fileContent = strings.ReplaceAll(fileContent, p.Config.fragmentListItemId, p.getIndexArticleList())

	articles := article.GetArticlesData()
	for _, data := range articles {
		p.buildArticleFile(data)
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
