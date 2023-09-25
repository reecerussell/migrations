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
	_ "github.com/reecerussell/migrations/providers/mysql"
)

const (
	defaultFileContext   = "."
	defaultConfigFile    = "migrations.yaml"
	defaultTransactional = false
	version              = "v0.2.5"
)

var (
	fileContext   string
	configFile    string
	target        string
	transactional bool
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go handleShutdown(cancel)

	upCommand := flag.NewFlagSet("up", flag.ExitOnError)
	upCommand.StringVar(&fileContext, "context", defaultFileContext, "The execution path of the migrations")
	upCommand.StringVar(&configFile, "file", defaultConfigFile, "The name of the migrations config file")
	upCommand.StringVar(&target, "target", "", "The migration to apply")
	upCommand.BoolVar(&transactional, "trans", defaultTransactional, "Determines wether if one migration in the list to apply fails, any applied in the list are rolled back.")

	downCommand := flag.NewFlagSet("down", flag.ExitOnError)
	downCommand.StringVar(&fileContext, "context", defaultFileContext, "The execution path of the migrations")
	downCommand.StringVar(&configFile, "file", defaultConfigFile, "The name of the migrations config file")
	downCommand.StringVar(&target, "target", "", "The migration to rollback")
	downCommand.BoolVar(&transactional, "trans", false, "Determines wether if one migration in the list to rollback fails, any rolled backed in the list are reapplied.")

	if len(os.Args) < 2 {
		help()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "up":
		upCommand.Parse(os.Args[2:])
		break
	case "down":
		downCommand.Parse(os.Args[2:])
		break
	case "version":
		fmt.Printf("Migrations %s\n", version)
		os.Exit(0)
		break
	default:
		help()
		os.Exit(0)
		break
	}

	fmt.Printf("Migrate transactionally: %v\n", transactional)
	fmt.Printf("Using context: %s\n", fileContext)

	configPath := path.Join(fileContext, configFile)
	config, err := migrations.LoadConfigFromFile(configPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Using config file: %s\n", configFile)

	p := providers.Get(config.Provider, config.Config)
	fmt.Printf("Using provider: %s\n", config.Provider)

	fr := migrations.NewFileReader(fileContext)

	if upCommand.Parsed() {
		err = migrations.Apply(ctx, config.Migrations, p, fr, target)
	}

	if downCommand.Parsed() {
		err = migrations.Rollback(ctx, config.Migrations, p, fr, target)
	}

	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}
}

func handleShutdown(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	fmt.Printf("\rAborting...\n")
	cancel()

	os.Exit(1)
}

func help() {
	fmt.Printf("Migrations %s\n---\n\nCommands:\n\n", version)

	// Up
	fmt.Printf("up\n---\n")
	fmt.Printf("description: Applies unapplied migrations.\n")
	fmt.Printf("usage: %s up --context example --file migrations.yaml --trans\n", os.Args[0])
	fmt.Printf("arguments:\n")
	fmt.Printf("\tcontext\tThe execution path of the migrations (default: %s)\n", defaultFileContext)
	fmt.Printf("\tfile\tThe name of the migrations config file (default: %s)\n", defaultConfigFile)
	fmt.Printf("\ttarget\tThe migration to apply. This will apply all migrations leading up to the target.\n")
	fmt.Printf("\ttrans\tDetermines wether to apply the migrations transactionally (default: %v)\n", defaultTransactional)

	fmt.Printf("\n")

	// Down
	fmt.Printf("down\n---\n")
	fmt.Printf("description: Rolls back applied migrations.\n")
	fmt.Printf("usage: %s down --context example --file migrations.yaml --trans\n", os.Args[0])
	fmt.Printf("arguments:\n")
	fmt.Printf("\tcontext: The execution path of the migrations (default: %s)\n", defaultFileContext)
	fmt.Printf("\tfile: The name of the migrations config file (default: %s)\n", defaultConfigFile)
	fmt.Printf("\ttarget: The migration to rollback. This will rollback all subsequent migrations, leading up to the target.\n")
	fmt.Printf("\ttrans: Determines wether to rollback the migrations transactionally (default: %v)\n", defaultTransactional)
}
