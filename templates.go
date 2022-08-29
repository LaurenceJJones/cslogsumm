package main

import (
	"text/template"

	log "github.com/sirupsen/logrus"
)

var (
	DefaultTemplates = map[string]string{
		"default": "I am default template",
	}
)

type CustomTemplate struct {
	Name     string `yaml:"name"`
	FilePath string `yaml:"file_path,omitempty"`
	Format   string `yaml:"format,omitempty"`
}

type TemplateEngine struct {
	CustomTemplates []CustomTemplate   `yaml:"custom_templates,omitempty"`
	Engine          *template.Template `yaml:"-"`
}

func DefaultTemplateEngine(Config *CslsConfig) *TemplateEngine {
	T := TemplateEngine{
		CustomTemplates: Config.CustomTemplates,
		Engine:          &template.Template{},
	}
	for key, val := range DefaultTemplates {
		log.Infof("Loading %s with %s as value", key, val)
		T.Engine.New(key).Parse(val)
	}
	_, err := T.Engine.New("Main").Parse(Config.Format)
	if err != nil {
		log.Debug(err.Error())
	}
	return &T
}
