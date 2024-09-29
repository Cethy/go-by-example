package plugin_article

import (
	plugin_fragment "go-by-example/static-website-generator/plugin-fragment"
	"slices"
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

func getIndexArticleListItem(link, title, imgSrc string) string {
	/*link := "/article.html"
	title := "foo"
	imgSrc := "/static/article1.jpg"*/

	return `<a href="` + link + `" class="p-4 md:w-1/2 hover:underline" title="` + title + `">
    <img src="` + imgSrc + `" class="mb-3 rounded-lg" alt="image article">
    <h3>` + title + `</h3>
</a>`
}

func getIndexArticleList() string {
	articles := GetArticlesData()
	list := ""
	for link, article := range articles {
		list = list + getIndexArticleListItem(link, article.Title, article.ImgSrc)
	}
	return list
}

var pluginFragmentId = "{index-articles}"

func PreBuildFile(fileContent, srcDir string) (string, error) {
	fragmentIds := plugin_fragment.GetAllFragmentIds(fileContent)

	if slices.Contains(fragmentIds, pluginFragmentId) {
		fileContent = strings.ReplaceAll(fileContent, pluginFragmentId, getIndexArticleList())
	}
	return fileContent, nil
}
