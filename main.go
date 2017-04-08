package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "Google Container Registry"
	//TODO SC: Define usage
	app.Usage = "gcr publish plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{

		//
		// commit args
		//

		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},

		//
		// gcr plugin params
		// Preserve original docker-gcr interface.
		//

		cli.StringFlag{
			Name:   "registry",
			Usage:  "gcr plugin target registry",
			EnvVar: "PLUGIN_REGISTRY",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "gcr plugin gcloud auth token",
			EnvVar: "PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "gcr plugin target repo",
			EnvVar: "PLUGIN_REPO",
		},
		cli.StringSliceFlag{
			Name:   "tag",
			Usage:  "gcr plugin docker image tags",
			EnvVar: "PLUGIN_TAG",
		},
		cli.StringFlag{
			Name:   "file",
			Usage:  "gcr plugin docker file name (Default 'Dockerfile')",
			EnvVar: "PLUGIN_FILE",
		},
		cli.StringFlag{
			Name:   "context",
			Usage:  "gcr plugin docker execution context (path)",
			EnvVar: "PLUGIN_CONTEXT",
		},
		cli.StringFlag{
			Name:   "storage_driver",
			Usage:  "gcr plugin docker storage driver",
			EnvVar: "PLUGIN_STORAGE_DRIVER",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {
	plugin := Plugin{
		Commit: Commit{
			Sha: c.String("commit.sha"),
		},
		Config: Config{
			// plugin-specific parameters
			Registry: c.String("registry"),
			Token:    strings.TrimSpace(c.String("token")),
			Repo:     c.String("repo"),
			Tag:      StrSlice{c.StringSlice("tag")},
			File:     c.String("file"),
			Context:  c.String("context"),
			Storage:  c.String("storage_driver"),
		},
	}

	// Repository name should have gcr prefix
	if len(plugin.Config.Registry) == 0 {
		plugin.Config.Registry = "gcr.io"
	}
	// Set the Dockerfile name
	if len(plugin.Config.File) == 0 {
		plugin.Config.File = "Dockerfile"
	}
	// Set the Context value
	if len(plugin.Config.Context) == 0 {
		plugin.Config.Context = "."
	}
	// Set the Tag value
	if plugin.Config.Tag.Len() == 0 {
		plugin.Config.Tag = StrSlice{[]string{"latest"}}
	}
	// Concat the Registry URL and the Repository name if necessary
	if strings.Count(plugin.Config.Repo, "/") == 1 {
		plugin.Config.Repo = fmt.Sprintf("%s/%s", plugin.Config.Registry, plugin.Config.Repo)
	}

	if err := plugin.Exec(); err != nil {
		fmt.Printf("Unable to execute plugin: %v\n", err)
		os.Exit(1)
	}
}
