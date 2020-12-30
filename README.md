# integrity

## Command line flags

- `-d`: Specify directory (defaults to current directory)

- `-c`: Specify suffix to VERSION file

    - For example, `-c SEC123-21-01` could create the file `VERSION-SEC123-21-01`

- `-v`: Verify that the files created earlier by this tool have not changed

## Examples

Create VERSION-TEST manifest file in current directory

```
integrity -c TEST
```

Create VERSION-TEST manifest file in the `/tmp` directory

```
integrity -c TEST -d /tmp
```

Verify VERSION-TEST manifest file in current working directory

```
integrity -c TEST -v
```

Verify VERSION-TEST manifest file in the `/tmp` directory

```
integrity -c TEST -d /tmp -v
```

## Building additional versions

The `build.sh` script will generate three versions of the tool at this time (all located in the `binaries` directory):

- macOS (amd64): Tested on macOS Catalina 10.15.7

- Linux (amd64): Tested on Ubuntu Desktop 20.04

- Windows (amd64): Tested on Windows 10 1809

To generate more, install go on your machine and set the environment variables `GOOS` and `GOARCH` for your target system:

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
