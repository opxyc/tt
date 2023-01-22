package main

import (
	"github.com/opxyc/tt/pkg/repository/sqlite"
	"github.com/opxyc/tt/pkg/tt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
)

var ttService tt.TTService

func main() {
	app := &cli.App{Name: "tt", Commands: commands}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("failed to init")
	}

	// sqlite file is stored at $HOME/.tt/db
	sqliteDBPath := path.Join(homeDir, ".tt")
	err = os.MkdirAll(sqliteDBPath, os.ModePerm)
	if err != nil {
		log.Fatal("failed to init")
	}

	repository, err := sqlite.NewSqliteRespository(path.Join(sqliteDBPath, "db"))
	if err != nil {
		log.Fatalf("failed to init: %s", err)
	}

	ttService = tt.NewTTService(repository)
}
