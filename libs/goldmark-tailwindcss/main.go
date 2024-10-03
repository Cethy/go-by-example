package goldmark_tailwindcss

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	goldmarkClassname "go-by-example/libs/goldmark-classname"
)

func getHeadingClassname(node ast.Node) string {
	n := node.(*ast.Heading)

	if n.Level == 1 {
		return "text-3xl font-bold"
	} else if n.Level == 2 {
		return "text-2xl font-bold"
	} else if n.Level == 3 {
		return "text-xl font-bold"
	} else if n.Level == 4 {
		return "text-l font-bold"
	} else if n.Level == 5 {
		return "text-l"
	} else {
		return "text-white/75"
	}
}

func getListClassname(node ast.Node) string {
	n := node.(*ast.List)
	className := "list-disc list-outside md:list-inside"
	if n.Parent().Kind() == ast.KindListItem {
		className = className + " pl-4"
	}
	return className
}

func provider(n ast.Node) string {
	if n.Kind() == ast.KindHeading {
		return getHeadingClassname(n)
	}
	if n.Kind() == ast.KindParagraph {
		return ""
	}
	if n.Kind() == ast.KindList {
		return getListClassname(n)
	}
	if n.Kind() == ast.KindLink {
		return "text-blue-600 hover:underline focus:border-blue-400 active:bg-blue-600 visited:text-purple-600"
	}
	return ""
}

func NewTransformer() parser.ASTTransformer {
	transformerConfig := goldmarkClassname.Config{
		ClassNameProvider: provider,
	}
	return goldmarkClassname.NewTransformer(transformerConfig)
}
