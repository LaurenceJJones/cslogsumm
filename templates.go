package main

import (
	"text/template"

	log "github.com/sirupsen/logrus"
)

var (
	DefaultTemplates = map[string]string{
		"default": `I am default template
I also like to be multi lined`,
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
		log.Debugf("Loading %s with %s as value", key, val)
		_, err := T.Engine.New(key).Parse(val)
		if err != nil {
			log.Debugf("Error whilst paring %s with %s", key, err.Error())
		}
	}
	_, err := T.Engine.New("Main").Parse(Config.Format)
	if err != nil {
		log.Debugf("Error whilst paring %s with %s", Config.Format, err.Error())
	}
	log.Debug(T.Engine.DefinedTemplates())
	return &T
}
