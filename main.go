package main

import (
	"fmt"
	"os"
	"strings"

	notecli "github.com/MonkeyBuisness/celli/notebook/cli"
	"github.com/MonkeyBuisness/celli/notebook/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	appVersion string
	appName    string
)

func main() {
	var (
		prettyBookFlag bool
	)

	app := &cli.App{
		Name:    appName,
		Usage:   "work with cellementaty notebooks in the easiest way",
		Version: appVersion,
		Commands: []*cli.Command{
			{
				Name:        "version",
				Description: "prints app version",
				Action: func(c *cli.Context) error {
					fmt.Printf("version: %s\n", appVersion)
					return nil
				},
			},
			{
				Name:     "new",
				Aliases:  []string{"n", "create"},
				Category: "template",
				Description: fmt.Sprintf("Supported template types: %s",
					strings.Join(types.SupportedBookTypes(), ",")),
				Usage:       "new <type of the notebook template to create>",
				Subcommands: createNewSubcommands(),
			},
			{
				Name:        "convert",
				Aliases:     []string{"c", "transform"},
				Category:    "template",
				Description: "converts existing notebook file to the template or existing template file to the notebook",
				Usage:       "convert book2tpl | tpl2book <path to the file>",
				Subcommands: []*cli.Command{
					{
						Name:    "book2tpl",
						Aliases: []string{"b2t"},
						Usage:   "book2tpl <path to the notebook file> > destination.md",
						Action: func(c *cli.Context) error {
							notebookPath := c.Args().First()
							return notecli.ConvertToTemplate(notebookPath)
						},
					},
					{
						Name:    "tpl2book",
						Aliases: []string{"t2b"},
						Usage:   "tpl2book <path to the template file> > destination.notebook",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:        "pretty",
								Aliases:     []string{"p"},
								Value:       false,
								Usage:       "pretty JSON output for notebook document",
								Destination: &prettyBookFlag,
							},
						},
						Action: func(c *cli.Context) error {
							templatePath := c.Args().First()
							return notecli.ConvertToNotebook(templatePath, prettyBookFlag)
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func createNewSubcommands() []*cli.Command {
	notebookTypes := types.SupportedBookTypes()
	cmds := make([]*cli.Command, len(notebookTypes))

	var templateDest string
	for i := range notebookTypes {
		cmds[i] = &cli.Command{
			Name: notebookTypes[i],
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:        "output",
					Aliases:     []string{"o", "dest", "dst"},
					Usage:       "output file or folder to save template data",
					DefaultText: "template.md in the current directory",
					Destination: &templateDest,
					Value:       "./",
				},
			},
			Action: func(c *cli.Context) error {
				return notecli.CreateTemplate(notebookTypes[i], templateDest)
			},
		}
	}

	return cmds
}
