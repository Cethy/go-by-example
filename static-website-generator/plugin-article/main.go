package plugin_article

import (
	"go-by-example/static-website-generator/generator"
	pluginfragment "go-by-example/static-website-generator/plugin-fragment"
	"path/filepath"
	"slices"
	"strings"
)

var pluginFragmentArticleFilename = "article"
var pluginFragmentListItemFilename = "article-list-item"
var pluginFragmentListItemId = "{index-articles}"

type Article struct {
	Title  string
	ImgSrc string
}

func getArticlesData() map[string]Article {
	articles := make(map[string]Article)
	articles["hello-world.html"] = Article{
		Title:  "Hello World",
		ImgSrc: "/static/article1.jpg",
	}
	articles["http-server.html"] = Article{
		Title:  "HTTP Server",
		ImgSrc: "/static/article2.jpg",
	}

	return articles
}

func getIndexArticleListItem(link, title, imgSrc, srcDir string) string {
	article := pluginfragment.GetFragmentContent(pluginFragmentListItemFilename, srcDir)

	article = strings.ReplaceAll(article, "{link}", link)
	article = strings.ReplaceAll(article, "{title}", title)
	article = strings.ReplaceAll(article, "{imgSrc}", imgSrc)

	return article
}

func getIndexArticleList(srcDir string) string {
	articles := getArticlesData()
	list := ""
	for link, article := range articles {
		list = list + getIndexArticleListItem(link, article.Title, article.ImgSrc, srcDir)
	}
	return list
}

var articleFilesBuilt = []string{}

func buildArticleFile(srcDir, targetPath string) {
	if slices.Contains(articleFilesBuilt, targetPath) {
		return
	}

	articleRaw := pluginfragment.GetFragmentContent(pluginFragmentArticleFilename, srcDir)

	fileContent, err := generator.PreProcess(articleRaw)
	if err != nil {
		panic(err)
	}

	articleFilesBuilt = append(articleFilesBuilt, targetPath)
	generator.BuildFile(targetPath, fileContent)
}

var pluginIsRunning bool = false

func PreBuildFile(fileContent string, config generator.Config) (string, error) {
	if pluginIsRunning {
		return fileContent, nil
	}
	pluginIsRunning = true

	fileContent = strings.ReplaceAll(fileContent, pluginFragmentListItemId, getIndexArticleList(config.SrcDir))

	articles := getArticlesData()
	for link := range articles {
		buildArticleFile(config.SrcDir, filepath.Join(config.OutputDir, link))
	}

	pluginIsRunning = false

	return fileContent, nil
}
