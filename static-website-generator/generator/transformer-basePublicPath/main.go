package transformer_basePublicPath

import (
	"go-by-example/static-website-generator/generator"
	"strings"
)

type Config struct {
	BasePublicPath string
}

type transformer struct {
	Config Config
}

func (p *transformer) Transform(fileContent string) (string, error) {
	fileContent = strings.ReplaceAll(fileContent, "{basePublicPath}", p.Config.BasePublicPath)
	return fileContent, nil
}

func NewTransformer(config Config) generator.Transformer {
	return &transformer{Config: config}
}
