package goldmark_classname

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Config struct {
	ClassNameProvider func(node ast.Node) string
}

type ClassNameTransformer struct {
	Config Config
}

func (c *ClassNameTransformer) Transform(node *ast.Document, _ text.Reader, _ parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (status ast.WalkStatus, err error) {
		if entering {
			n.SetAttributeString("class", c.Config.ClassNameProvider(n))
		}
		return ast.WalkContinue, nil
	})
}

func NewTransformer(config Config) parser.ASTTransformer {
	return &ClassNameTransformer{Config: config}
}
