package plugin_article

import (
	pluginfragment "go-by-example/static-website-generator/plugin-fragment"
	"strings"
)

type Article struct {
	Title  string
	ImgSrc string
}

func GetArticlesData() map[string]Article {
	articles := make(map[string]Article)
	articles["hello-world"] = Article{
		Title:  "Hello World",
		ImgSrc: "/static/article1.jpg",
	}
	articles["http-server"] = Article{
		Title:  "HTTP Server",
		ImgSrc: "/static/article2.jpg",
	}

	return articles
}

func getIndexArticleListItem(link, title, imgSrc, srcDir string) string {
	article := pluginfragment.GetFragmentContent("article-list-item", srcDir)

	article = strings.ReplaceAll(article, "{link}", link)
	article = strings.ReplaceAll(article, "{title}", title)
	article = strings.ReplaceAll(article, "{imgSrc}", imgSrc)

	return article
}

func getIndexArticleList(srcDir string) string {
	articles := GetArticlesData()
	list := ""
	for link, article := range articles {
		list = list + getIndexArticleListItem(link, article.Title, article.ImgSrc, srcDir)
	}
	return list
}

var pluginFragmentId = "{index-articles}"

func PreBuildFile(fileContent, srcDir string) (string, error) {
	fileContent = strings.ReplaceAll(fileContent, pluginFragmentId, getIndexArticleList(srcDir))

	return fileContent, nil
}
