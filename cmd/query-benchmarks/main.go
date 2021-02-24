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
	app.Commands = []cli.Command{RegCmd(), GWDICmd(), AnnoCmd()}
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
		cli.IntFlag{
			Name:     "qty, q",
			Usage:    "number of strains to loop through",
			Required: true,
		},
	}
}

func RegCmd() cli.Command {
	flags := queryFlags()
	flags = append(flags, serviceFlags()...)
	return cli.Command{
		Name:   "reg",
		Usage:  "times the query for getting non-gwdi strain list",
		Action: app.TimeRegQuery,
		Flags:  flags,
	}
}

func GWDICmd() cli.Command {
	flags := queryFlags()
	flags = append(flags, serviceFlags()...)
	return cli.Command{
		Name:   "gwdi",
		Usage:  "times the query for getting gwdi strain list",
		Action: app.TimeGWDIQuery,
		Flags:  flags,
	}
}

func AnnoCmd() cli.Command {
	flags := queryFlags()
	flags = append(flags, serviceFlags()...)
	return cli.Command{
		Name:   "anno",
		Usage:  "times the query for getting inventory list first",
		Action: app.TimeAnnoQuery,
		Flags:  flags,
	}
}
