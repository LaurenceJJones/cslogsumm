package main

import (
	"os"

	"github.com/crowdsecurity/crowdsec/pkg/csconfig"
	"github.com/crowdsecurity/crowdsec/pkg/database"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type CslsConfig struct {
	FilePath     string                `yaml:"-"`
	Mode         string                `yaml:"mode"`
	DbConfig     *csconfig.DatabaseCfg `yaml:"db_config,omitempty"`
	CsConfigPath string                `yaml:"cs_config_path,omitempty"`
	CsConfig     csconfig.Config       `yaml:"-"`
	DbClient     *database.Client      `yaml:"-"`
	Format       string                `yaml:"format,omitempty"`
}

func NewDefaultConfig(FilePath string) *CslsConfig {
	cfg := CslsConfig{
		FilePath: FilePath,
	}
	readYAML(FilePath, &cfg)
	if cfg.CsConfigPath != "" {
		readYAML(cfg.CsConfigPath, &cfg.CsConfig)
	}
	if cfg.DbConfig != nil {
		if dbClient, err := database.NewClient(cfg.DbConfig); err == nil {
			cfg.DbClient = dbClient
		}
	}
	return &cfg
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
