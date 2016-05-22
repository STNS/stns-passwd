# stns-passwd
This command is used to generate a password to be used in STNS

## Description

## Usage
```bash
$ stns-passwd -s <username> -m <hash_method> -c <stretching_count> <password_strings>
```

## options
```bash
Usage of stns-passwd:
  -c int
        The number of times to stretching
  -m string
        Specifies the hash function(sha256/sha512) (default "sha256")
  -s string
        String used to salt. The SNS is the user name
  -version
        Print version information and quit.
```
## Install
### Donwnload
* [Github Release](https://github.com/STNS/stns-passwd/releases)

### go get
```bash
$ go get -d github.com/STNS/stns-passwd
```

## Contribution

1. Fork ([https://github.com/STNS/stns-passwd/fork](https://github.com/STNS/stns-passwd/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author
* [pyama86](http://github.com/pyama86)
