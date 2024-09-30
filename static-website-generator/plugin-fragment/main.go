package plugin_fragment

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var fragments = make(map[string]string)

func getFragmentContent(fragmentId, srcDir string) string {
	fragmentId = filepath.Join(srcDir, "../fragments", fragmentId+".html")

	if fragments[fragmentId] == "" {
		readFile, err := os.ReadFile(fragmentId)
		if err != nil {
			log.Println("WARNING", "cannot find file for fragment", fragmentId)
			return ""
		}

		fragments[fragmentId] = string(readFile)
	}
	return fragments[fragmentId]
}

func GetAllFragmentIds(fileContent string) []string {
	pattern := "(\\{[A-z-_]+\\})"
	r, _ := regexp.Compile(pattern)
	return r.FindAllString(fileContent, -1)
}

func PreBuildFile(fileContent, srcDir string) (string, error) {
	// insert fragments
	requiredFragments := GetAllFragmentIds(fileContent)
	for _, fragmentId := range requiredFragments {
		fragment := getFragmentContent(fragmentId[1:len(fragmentId)-1], srcDir)
		fileContent = strings.ReplaceAll(fileContent, fragmentId, fragment)
	}

	return fileContent, nil
}
