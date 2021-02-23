package main

import (
	"log"
	"os"

	"github.com/dictybase-playground/query-benchmarks/internal/app"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "query-benchmarks"
	app.Usage = "cli for timing various stock list queries"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text",
			Value: "json",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
	}
	app.Commands = []cli.Command{QueryCmd()}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error in running command %s", err)
	}
}

func serviceFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "stock-grpc-host",
			EnvVar:   "STOCK_API_SERVICE_HOST",
			Usage:    "stock grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "stock-grpc-port",
			EnvVar:   "STOCK_API_SERVICE_PORT",
			Usage:    "stock grpc port",
			Required: true,
		},
		cli.StringFlag{
			Name:     "annotation-grpc-host",
			EnvVar:   "ANNOTATION_API_SERVICE_HOST",
			Usage:    "annotation grpc host",
			Required: true,
		},
		cli.StringFlag{
			Name:     "annotation-grpc-port",
			EnvVar:   "ANNOTATION_API_SERVICE_PORT",
			Usage:    "annotation grpc port",
			Required: true,
		},
	}
}

func queryFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "query",
			Usage:    "stock list query to test",
			Required: true,
		},
		cli.IntFlag{
			Name:     "quantity",
			Usage:    "number of strains to loop through",
			Required: true,
		},
	}
}

func QueryCmd() cli.Command {
	flags := queryFlags()
	flags = append(flags, serviceFlags()...)
	return cli.Command{
		Name:   "query",
		Usage:  "time the specified query",
		Action: app.TimeQuery,
		Flags:  flags,
	}
}
