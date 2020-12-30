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
