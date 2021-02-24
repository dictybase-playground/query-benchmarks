# query-benchmarks

Timing different stock list queries.

```bash
NAME:
   query-benchmarks - cli for timing various stock list queries

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   reg      times the query for getting non-gwdi strain list
   gwdi     times the query for getting gwdi strain list
   anno     times the query for getting inventory list first
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-format value  format of the logging out, either of json or text (default: "json")
   --log-level value   log level for the application (default: "error")
   --help, -h          show help
   --version, -v       print the version
```

The `reg` and `gwdi` commands will also generate and save a `.png` image with a
scatterplot of unavailable strains per query.
