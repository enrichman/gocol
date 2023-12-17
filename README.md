# gocol

Go color your test coverage! ✨

```
go test -cover ./... | gocol
```

`gocol` will turn this:

<img alt="gocol 1" src="./assets/gocol_1.png" width="600px">

into this!

<img alt="gocol 2" src="./assets/gocol_2.png" width="600px">

See immediately how much your projects is covered!

## Installation

To install `gocol` just use `go install`
```
go install github.com/enrichman/gocol@v0.0.1
```

## Usage

Pipe the output of a `go test -cover` to `gocol`
```
go test -cover ./... | gocol
```

If you are using the verbose `-v` then the `PASS|FAIL|SKIP` lines will be coloured as well.

<img alt="gocol 4" src="./assets/gocol_4.png" width="600px">


# Colors and ranges 🌈

Currently only a fixed range of colors and percentage is available.

<img alt="gocol 3" src="./assets/gocol_3.png" width="600px">

# Feedback
If you like the project please star it on Github 🌟, and feel free to drop me a note, or [open an issue](https://github.com/enrichman/gocol/issues/new)!

[Twitter](https://twitter.com/enrichmann)

# License

[MIT](LICENSE)