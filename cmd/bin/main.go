package main

import (
	"os"
	"product-service/cmd"
	"product-service/internal/adapter"
	"product-service/internal/infrastructure"

	"flag"

	"github.com/rs/zerolog/log"
)

func main() {
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError) // create a new flag set for server command
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)     // create a new flag set for seed command

	if len(os.Args) < 2 { // check if no command provided
		log.Info().Msg("No command provided, defaulting to 'server'")
		cmd.RunServer(serverCmd, os.Args[1:])
		os.Exit(0)
	}

	switch os.Args[1] { // check the command provided
	case "server":
		cmd.RunServer(serverCmd, os.Args[2:])
	case "seed":
		cmd.RunSeed(seedCmd, os.Args[2:])
	default:
		log.Info().Msg("Invalid command provided, defaulting to 'server' with provided flags")
		if os.Args[1][0] == '-' { // check if the first argument is a flag
			cmd.RunServer(serverCmd, os.Args[1:])
			os.Exit(0)
		}

		cmd.RunServer(serverCmd, os.Args[2:]) // default to server if invalid command and flags are provided
	}
}

// init function to initialize the configuration and adapter
// before the main function is executed
func init() {
	infrastructure.Configuration(
		infrastructure.WithPath("./"),
		infrastructure.WithFilename(".env"),
	).Initialize()

	adapter.Adapters = &adapter.Adapter{}
}
