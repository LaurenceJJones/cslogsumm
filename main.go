package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

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
		},
		Run: func(cmd *cobra.Command, args []string) {
			template, _ := template.New("").Parse(Config.Format)
			b := new(strings.Builder)
			template.Execute(b, Config.DbClient)
			log.Info(b.String())
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
