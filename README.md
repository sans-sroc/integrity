# Integrity

[![GitHub Super-Linter](https://github.com/sans-blue-team/integrity/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)
[![pre-release](https://github.com/sans-blue-team/integrity/actions/workflows/pre-release.yml/badge.svg)](https://github.com/sans-blue-team/integrity/actions/workflows/pre-release.yml)

## Overview

File validation at it's finest.

## Help

```help
NAME:
   integrity - integrity

USAGE:
   integrity [global options] command [command options] [arguments...]

AUTHORS:
   Ryan Nicholson <rnicholson@sans.org>
   Don Williams <dwilliams@sans.org>
   Erik Kristensen <ekristensen@sans.org>

COMMANDS:
   create    create integrity files
   validate  validate integrity files
   version   print version

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### Create

```help
NAME:
   integrity create - create integrity files

USAGE:
   integrity create [command options] [arguments...]

OPTIONS:
   --name value, -n value       The name that will be given to the ISO volume during USB creation. [$NAME]
   --user value, -u value       allow setting what user created the file (default: "ekristen") [$USER]
   --log-level value, -l value  Log Level (default: "info") [$LOG_LEVEL]
   --directory value, -d value  The directory that will be the current working directory for the tool when it runs (default: ".") [$DIRECTORY]
   --help, -h                   show help (default: false)

```

### Validate

```help
NAME:
   integrity validate - validate integrity files

USAGE:
   integrity validate [command options] [arguments...]

OPTIONS:
   --output-format value, --format value  Chose which format to output the validation results (default is none) (valid options: none, json) (default: "none") [$OUTPUT_FORMAT]
   --output value, -o value               When output-format is specified, this controls where it goes, (defaults to stdout) (default: "-") [$OUTPUT]
   --log-level value, -l value            Log Level (default: "info") [$LOG_LEVEL]
   --directory value, -d value            The directory that will be the current working directory for the tool when it runs (default: ".") [$DIRECTORY]
   --help, -h                             show help (default: false)
```

#### Validate Output

The validate output options change the behavior of the too slightly.

If the `--output-format` is set to `json` and the `--log-level` has not been set to `none` it will write all logs to `STDERR` while the JSON format is written to `STDOUT`, this is to allow the capture of the `json` separately from the log output.

## Examples

### Simple Create

```bash
integrity create -n 572.00.0
```

### Create w/ Specified Directory

```bash
integrity create -n 572.00.0 -d /tmp
```

### Simple Validate

**Note:** this assumes the create was run in the current working directory and `sans-integrity.yml` already exists.

```bash
integrity validate
```

### Validate w/ Specified Directory

```bash
integrity validate -d /tmp
```

### Validate w/ JSON Output

**Note:** this is really only useful for programmatic validation purposes.

```bash
integrity validate --output-format json 
```

### Validate w/ JSON Output to File

**Note:** this is really only useful for programmatic validation purposes.

```bash
integrity validate --output-format json --output results.json
```

## Building

### Building Additional Versions

You can find three builds in the Releases section of GitHub, but to generate more, install go on your machine and set the environment variables `GOOS` and `GOARCH` for your target system:

| GOOS | GOARCH |
|------|--------|
| android | arm |
| darwin | 386 |
| darwin | amd64 |
| darwin | arm |
| darwin | arm64 |
| dragonfly | amd64 |
| freebsd | 386 |
| freebsd | amd64 |
| freebsd | arm |
| linux | 386 |
| linux | amd64 |
| linux | arm |
| linux | arm64 |
| linux | ppc64 |
| linux | ppc64le |
| linux | mips |
| linux | mipsle |
| linux | mips64 |
| linux | mips64le |
| netbsd | 386 |
| netbsd | amd64 |
| netbsd | arm |
| openbsd | 386 |
| openbsd | amd64 |
| openbsd | arm |
| plan9 | 386 |
| plan9 | amd64 |
| solaris | amd64 |
| windows | 386 |
| windows | amd64 |

This list was gathered from [DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04).

You can find a complete list of `GOOS` and `GOARCH` values [here](https://github.com/golang/go/blob/master/src/go/build/syslist.go).

You can build like so on a Linux/macOS system:

```bash
env GOOS=solaris GOARCH=amd64 go build -o binaries/integrity-solaris-amd64 integrity.go
```

## Build with Docker instead

You can also build with Docker if you're uneasy installing go on your machine (or just prefer to use Docker):

```bash
docker run -it --rm -v ${PWD}:/usr/src/integrity -w /usr/src/integrity \
    golang env GOOS=solaris GOARCH=amd64 go build \
    -o binaries/integrity-solaris-amd64 integrity.go
```

## Development

There are go modules included on this project, so you will need to make sure you run `go mod vendor` to bring them to your local directory if you are using the Makefile as the makefile prefers the use of the vendor directory.

If you are simply running `go run main.go` the modules will be pulled from vendor or your go root depending on where it finds it. Golang will also automatically pull the mods down when you run if there are changes.

During iterative updates if the modules change you will find yourself needing to run `go mod vendor` or at least `go mod download` to ensure you have the updated modules locally.

### Ignore Files

There are two types of ignore files. `IgnoreOnCreate` and `IgnoreAlways`, both are defined in [pkg/common/constants.go](pkg/common/constants.go). Files that should go in the `IgnoreOnCreate` are things like `.DS_Store`, whereas files that should go into `IgnoreAlways` is the `sans-integrity.yml` and all GPG files with a match of `.*\\.gpg$`.

Changing the ignore files will require a new release of the tool.

**Note:** to aid developers, the option `-i` is present that allows you to pass a custom ignore strictly for development purposes while testing and developing on the tool. This is useful when needing to ensure the validation tool is properly picking up files that are not in the file.
