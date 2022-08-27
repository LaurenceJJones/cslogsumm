package main

import (
	"os"

	"github.com/crowdsecurity/crowdsec/pkg/csconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type CslsConfig struct {
	FilePath     string                `yaml:"-"`
	Mode         string                `yaml:"mode"`
	DbConfig     *csconfig.DatabaseCfg `yaml:"db_config,omitempty"`
	CsConfigPath string                `yaml:"cs_config_path,omitempty"`
	CsConfig     csconfig.Config       `yaml:"-"`
}

func NewDefaultConfig(FilePath string) *CslsConfig {
	cfg := CslsConfig{
		FilePath: FilePath,
	}
	readYAML(FilePath, &cfg)
	if cfg.CsConfigPath != "" {
		readYAML(cfg.CsConfigPath, &cfg.CsConfig)
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
