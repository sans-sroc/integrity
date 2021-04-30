# integrity

[![GitHub Super-Linter](https://github.com/sans-blue-team/integrity/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)
[![pre-release](https://github.com/sans-blue-team/integrity/actions/workflows/pre-release.yml/badge.svg)](https://github.com/sans-blue-team/integrity/actions/workflows/pre-release.yml)

## Help

```
NAME:
   integrity - integrity

USAGE:
   integrity [global options] command [command options] [arguments...]

VERSION:
   v2.0.0

AUTHORS:
   Ryan Nicholson <rnicholson@sans.org>
   Don Williams <dwilliams@sans.org>

COMMANDS:
   create    create integrity files
   validate  validate integrity files
   version   print version
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Create

```
NAME:
   integrity create - create integrity files

USAGE:
   integrity create [command options] [arguments...]

OPTIONS:
   --log-level value, -l value           Log Level (default: "info") [$LOG_LEVEL]
   --directory value, -d value           Target Directory (default: ".") [$DIRECTORY]
   --courseware-version value, -c value  Coursware Version Identifier [$COURSEWARE_VERSION]
   --json, -j                            Output in JSON (default: false)
   --json-pretty                         Output JSON in Pretty Print Format (default: true)
   --user value                          allow setting what user created the file (default: "ekristen")
   --help, -h                            show help (default: false)

```

### Validate

```
NAME:
   integrity validate - validate integrity files

USAGE:
   integrity validate [command options] [arguments...]

OPTIONS:
   --parts, -p                           Validate the VERSION-part.txt file (default: false)
   --first, -f                           Validate the VERSION-first.txt file (default: false)
   --log-level value, -l value           Log Level (default: "info") [$LOG_LEVEL]
   --directory value, -d value           Target Directory (default: ".") [$DIRECTORY]
   --courseware-version value, -c value  Coursware Version Identifier [$COURSEWARE_VERSION]
   --json, -j                            Output in JSON (default: false)
   --json-pretty                         Output JSON in Pretty Print Format (default: true)
   --user value                          allow setting what user created the file (default: "ekristen")
   --help, -h                            show help (default: false)
```

## Examples

Create VERSION-TEST.txt manifest file in current directory

```bash
integrity create -c TEST
```

Create VERSION-TEST.txt manifest file in the `/tmp` directory

```bash
integrity create -c TEST -d /tmp
```

Verify VERSION-TEST.txt manifest file in current working directory

```bash
integrity validate -c TEST
```

Verify VERSION-TEST.txt manifest file in the `/tmp` directory

```bash
integrity validate -c TEST -d /tmp
```

Output results as JSON (no VERSION file is created)

```bash
integrity create -c TEST -j
```

Verify VERSION-TEST.txt manifest in the `/tmp` directory and output the results in JSON

```bash
integrity validate -c TEST -d /tmp -j
```

## Building additional versions

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
