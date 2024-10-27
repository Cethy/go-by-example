package file

import (
	"os"
	"regexp"
	"ssh-multitodolist/data"
	"strings"
)

/**
IMPORT/EXPORT data format:
```markdown
# list 1

- [ ] Foo
- [x] Bar
- [] Baz

---

# list 2

---

- [ ] unnamed list's item
- [x] Bar2
- [ ] Baz2
```
*/

func readList(raw string) []data.ListItem {
	pattern := "\\- \\[(?P<Checked> ?x?)\\] (?P<Value>[A-z0-9].*)"
	r, _ := regexp.Compile(pattern)
	all := r.FindAllStringSubmatch(raw, -1)
	var listItems []data.ListItem
	for _, item := range all {
		listItems = append(listItems, data.ListItem{
			Value:   item[r.SubexpIndex("Value")],
			Checked: item[r.SubexpIndex("Checked")] == "x",
		})
	}

	return listItems
}

func readData(sourcePath string) ([]data.NamedList, error) {
	var namedLists []data.NamedList

	rawContent, err := os.ReadFile(sourcePath)
	if os.IsNotExist(err) {
		rawContent = []byte(`# `)
	} else if err != nil {
		return namedLists, err
	}

	// split lists
	rawLists := strings.Split(string(rawContent), "---")
	patternTitle := "(?m)\\# (?P<Title>.*)$"
	r, _ := regexp.Compile(patternTitle)
	for _, rawList := range rawLists {
		titleMatch := r.FindStringSubmatch(rawList)

		name := ""
		if len(titleMatch) > 0 {
			name = titleMatch[r.SubexpIndex("Title")]
		}

		namedLists = append(namedLists, data.NamedList{
			Name:  name,
			Items: readList(rawList),
		})
	}

	return namedLists, nil
}

func writeData(namedLists []data.NamedList, targetPath string) error {
	content := ""
	for i, namedList := range namedLists {
		if i > 0 {
			content += "\n\n---\n\n"
		}
		if namedList.Name != "" {
			content += "# " + namedList.Name + "\n\n"
		}
		for _, listItem := range namedList.Items {
			Checked := " "
			if listItem.Checked {
				Checked = "x"
			}
			content = content + "- [" + Checked + "] " + listItem.Value + "\n"
		}
	}

	return os.WriteFile(targetPath, []byte(content), 0644)
}