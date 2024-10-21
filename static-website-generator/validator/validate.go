package validator

import (
	"errors"
	"log"
	"os"
	"path"
	"static-website-generator/article"
	"static-website-generator/article/unsplash"
	"strconv"
	"strings"
)

type Error struct {
	article    article.Article
	violations []string
}

func tryFixImgSrc(projectDirName string) error {
	articlePath := path.Join("..", projectDirName, "README.md")

	accessToken, err := os.ReadFile(".unsplashAccessToken")
	if err != nil {
		log.Println("You need to a create a .unsplashAccessToken file")
		panic(err)
	}

	imgUrl, err := unsplash.GetRandomPhotoUrl(string(accessToken))
	if err != nil {
		return err
	}

	articleMd, err := os.ReadFile(articlePath)

	lines := strings.Split(string(articleMd), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if strings.TrimSpace(line) == "---" {
			newSlice := make([]string, len(lines)+1)
			copy(newSlice[:i], lines[:i])
			newSlice[i] = "ImgSrc: " + imgUrl
			copy(newSlice[i+1:], lines[i:])

			lines = newSlice
			break
		}
	}

	return os.WriteFile(articlePath, []byte(strings.Join(lines, "\n")), 0644)
}

func Validate() error {
	var myErrors []Error
	articles := article.GetArticlesData()

	for _, a := range articles {
		var violations []string
		if a.Title == "" {
			violations = append(violations, "Title is required")
		}
		if a.Order == -1 {
			violations = append(violations, "Order is required")
		}
		if a.ImgSrc == "" {
			violations = append(violations, "ImgSrc is missing")
		}

		myErrors = append(myErrors, Error{a, violations})
	}

	var mustFixCount int
	for _, e := range myErrors {
		if len(e.violations) == 0 {
			log.Println("[OK]", e.article.ProjectDirname)
			continue
		}
		log.Println("[KO]", e.article.ProjectDirname, "missing metadata:")
		for i, v := range e.violations {
			mustFixCount++

			indic := "|"
			if i == len(e.violations)-1 {
				indic = "└"
			}

			fixedTxt := ""
			// only try to fix img if it's the only error
			if v == "ImgSrc is missing" && len(e.violations) == 1 {
				fixedTxt = "- FIXED!"
				mustFixCount--

				err := tryFixImgSrc(e.article.ProjectDirname)
				if err != nil {
					fixedTxt = "- " + err.Error()
				}
			}

			log.Printf(" %s-> %s %s", indic, v, fixedTxt)
		}
	}

	if mustFixCount > 0 {
		log.Println("You have " + strconv.Itoa(mustFixCount) + " violations to fix.")
		return errors.New("You have " + strconv.Itoa(mustFixCount) + " violations to fix.")
	}
	log.Println("All is well ...")
	return nil
}
