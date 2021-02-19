package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/providers"

	// providers
	_ "github.com/reecerussell/migrations/providers/mssql"
)

var (
	migrationContext  string
	migrationFilename string
	target            string
	rollback          bool
)

func main() {
	flag.StringVar(&migrationContext, "context", ".", "the migration context path")
	flag.StringVar(&migrationFilename, "file", "migrations.yaml", "the migration config file")
	flag.StringVar(&target, "target", "", "the target migration to apply or rollback")
	flag.BoolVar(&rollback, "rollback", false, "determines wether the migrations will be rolled back")
	flag.Parse()

	configPath := path.Join(migrationContext, migrationFilename)
	config, err := migrations.LoadConfigFromFile(configPath)
	if err != nil {
		panic(err)
	}

	p := providers.Get(config.Provider)
	fmt.Printf("Using provider: %s\n", config.Provider)

	ctx, cancel := context.WithCancel(context.Background())
	ctx = migrations.NewContext(ctx, migrationContext)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-stop

		fmt.Printf("\rAborting...\n")
		cancel()

		os.Exit(1)
	}()

	if rollback {
		err = migrations.Rollback(ctx, config.Migrations, p, target)
	} else {
		err = migrations.Apply(ctx, config.Migrations, p, target)
	}

	if err != nil {
		fmt.Printf("An error occured: %v\n", err)
		os.Exit(1)
	}
}
