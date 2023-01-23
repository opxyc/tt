package main

import "github.com/urfave/cli/v2"

var commands = []*cli.Command{
	{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start a new activity",
		Action: func(cCtx *cli.Context) error {
			return start()
		},
	},
	{
		Name:    "pause",
		Aliases: []string{"p"},
		Usage:   "pause current activity",
		Action:  pause,
	},
	{
		Name:    "resume",
		Aliases: []string{"r"},
		Usage:   "resume an activity",
		Action:  resume,
	},
	{
		Name:    "stop",
		Aliases: []string{"x"},
		Usage:   "stop an activity",
		Action:  stop,
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list activities",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "date",
				Aliases:     []string{"dt"},
				DefaultText: "today",
				Value:       "today",
				Usage:       "YYYY-mm-dd (ex: 2023-01-21)",
			},
			&cli.BoolFlag{
				Name:        "csv",
				DefaultText: "true",
			},
		},
		Action: list,
	},
	{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete an activity",
		Action:  delete,
	},
}
