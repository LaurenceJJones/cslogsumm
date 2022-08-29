package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var ConfigFilePath string
var Config *CslsConfig

func prepender(filename string) string {
	const header = `---
id: %s
title: %s
---
`
	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	return fmt.Sprintf(header, base, strings.Replace(base, "_", " ", -1))
}

func linkHandler(name string) string {
	return fmt.Sprintf("/cslogsumm/%s", name)
}

func setLogger(Config *CslsConfig) {
	logFormatter := &log.TextFormatter{TimestampFormat: "02-01-2006 03:04:05 PM", FullTimestamp: true}
	log.SetFormatter(logFormatter)
	switch strings.ToLower(Config.LogLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
	log.Debugf("Logging is set to %s", Config.LogLevel)
}

func main() {

	var rootCmd = &cobra.Command{
		Use:               "cslogsumm",
		Short:             "cslogsumm allows you to generate reports",
		Long:              `cslogsumm allows you to generate html/text reports from crowdsec via cscli/db/api/`,
		DisableAutoGenTag: true,
		SilenceErrors:     false,
		SilenceUsage:      true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if ConfigFilePath != "" {
				Config = NewDefaultConfig(ConfigFilePath)
			}
			setLogger(Config)
			Config.TemplateEngine = DefaultTemplateEngine(Config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := Config.TemplateEngine.Engine.ExecuteTemplate(os.Stdout, "Main", Config)
			if err != nil {
				log.Debug(err.Error())
			}
		},
	}
	var cmdDocGen = &cobra.Command{
		Use:               "doc",
		Short:             "Generate the documentation in `./doc/`. Directory must exist.",
		Args:              cobra.ExactArgs(0),
		Hidden:            true,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := doc.GenMarkdownTreeCustom(rootCmd, "./doc/", prepender, linkHandler); err != nil {
				log.Fatalf("Failed to generate cobra doc: %s", err)
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", "", "path to config file")
	rootCmd.AddCommand(cmdDocGen)

	if err := rootCmd.Execute(); err != nil {
		exitCode := 1
		log.NewEntry(log.StandardLogger()).Log(log.FatalLevel, err)
		os.Exit(exitCode)
	}
}
