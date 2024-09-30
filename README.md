# Quake Log Parser

The Quake Log Parser is a command line tool (CLI) designed to process Quake 3 Arena server logs. It extracts useful match data, generating a JSON report with players, kills, leaderboard and causes of death for each match.

This project was built for a CloudWalk interview process.

## Requirements

In order to run this parser, you should have:

- Golang with version 1.20 or higher
- Quake 3 Arena server logs (or use the provided in the `input` folder)

## Instructions

In order to run this Log Parser CLI, do the following steps:

First, build the project with the following command:

```bash
$ go build -o ./bin/quake-parser cmd/main.go
```

Once built, run the parser with the following command:

```bash
$ ./bin/quake-parser --type parallel --input_file input/quake.log
```

### Parser Options

This parser requires two arguments:
- `--type [parallel|sequential]`: Defines the type of the parser. Use `parallel` for concurrent parsing and `sequential` for linear parsing.
- `input_file [file]`. The path of file that needs to be parsed.

For more information, use the `help` command:
```
$ ./bin/quake-parser --help

Usage ./bin/quake-parser -type [parallel|sequential] -input_file [file]
-input_file string
        The path of file to be parsed
-type string
        The type of parsing (e.g., sequential, parallel)
```

## Testing

To run all the unit tests in the project, use the following command:

```
$ go test -v ./...
```

This command will automatically find and execute all test files in the project, providing verbose output.

## Additional Information

For a better user experience, the parser tracks each player's nickname throughout the match. If a player changes their nickname mid-game, the parser ensures that all kills are correctly attributed to their new name. The final report displays the player's last used nickname.

The parser also handles reconnections by generating a hash based on the player's last-used nickname. If a player reconnects **with the same name**, their kill count is restored. While this does not reflect the actual game mechanics (as kills are typically lost upon disconnecting), it provides a more accurate report for analysis.

