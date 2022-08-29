package main

import (
	"os"

	"github.com/crowdsecurity/crowdsec/pkg/csconfig"
	"github.com/crowdsecurity/crowdsec/pkg/database"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type CslsConfig struct {
	FilePath        string                `yaml:"-"`
	Mode            string                `yaml:"mode"`
	DbConfig        *csconfig.DatabaseCfg `yaml:"db_config,omitempty"`
	CsConfigPath    string                `yaml:"cs_config_path,omitempty"`
	CsConfig        csconfig.Config       `yaml:"-"`
	DbClient        *database.Client      `yaml:"-"`
	Format          string                `yaml:"format,omitempty"`
	CustomTemplates []CustomTemplate      `yaml:"custom_templates,omitempty"`
	TemplateEngine  *TemplateEngine       `yaml:"-"`
	LogLevel        string                `yaml:"log_level,omitempty"`
	EmailSettings   *EmailSettings        `yaml:"email_settings,omitempty"`
	EmailClient     *Email                `yaml:"-"`
}

func NewDefaultConfig(FilePath string) *CslsConfig {
	cfg := CslsConfig{
		FilePath: FilePath,
	}
	if FilePath == "" {
		return &cfg
	}
	readYAML(FilePath, &cfg)
	if cfg.CsConfigPath != "" {
		readYAML(cfg.CsConfigPath, &cfg.CsConfig)
	}
	if cfg.CsConfig.DbConfig != nil {
		cfg.LoadDbIfExist(cfg.CsConfig.DbConfig)
	}
	if cfg.DbConfig != nil {
		cfg.LoadDbIfExist(cfg.DbConfig)
	}
	return &cfg
}

func (C *CslsConfig) LoadDbIfExist(cs *csconfig.DatabaseCfg) {
	if cs.DbPath != "" {
		_, error := os.Stat(cs.DbPath)
		if !errors.Is(error, os.ErrNotExist) {
			dbClient, _ := database.NewClient(cs)
			C.DbClient = dbClient
		}
	}
}

func readYAML[T CslsConfig | csconfig.Config](filePath string, cfg *T) error {
	var content []byte
	var err error
	if content, err = os.ReadFile(filePath); err != nil {
		return errors.Wrap(err, "while reading yaml file")
	}
	configData := os.ExpandEnv(string(content))
	yaml.UnmarshalStrict([]byte(configData), &cfg)
	return nil
}
