# integrity

[![GitHub Super-Linter](https://github.com/sans-blue-team/integrity/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)
[![pre-release](https://github.com/sans-blue-team/integrity/actions/workflows/pre-release.yml/badge.svg)](https://github.com/sans-blue-team/integrity/actions/workflows/pre-release.yml)

## Command line flags

- `-d`: Specify directory (defaults to current directory)

- `-c`: Specify suffix to VERSION file

  - For example, `-c SEC123-21-01` could create the file `VERSION-SEC123-21-01`

- `-v`: Verify that the files created earlier by this tool have not changed

- `-j`: Output results as JSON instead of writing a VERSION file

- `-p`: Only verify the "-part" version file

## Examples

Create VERSION-TEST.txt manifest file in current directory

```bash
integrity -c TEST
```

Create VERSION-TEST.txt manifest file in the `/tmp` directory

```bash
integrity -c TEST -d /tmp
```

Verify VERSION-TEST.txt manifest file in current working directory

```bash
integrity -c TEST -v
```

Verify VERSION-TEST.txt manifest file in the `/tmp` directory

```bash
integrity -c TEST -d /tmp -v
```

Output results as JSON

```bash
integrity -c TEST -j
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
