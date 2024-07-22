## Beancode

Working with BitTorrent? You need Beancode.

## Features

* Easy API based on stdlib
* Almost full unit test coverage: 86%
* Support int, string, slice, map, and struct

## Usage

Let's define a common struct for encoding and decoding processes:
```go
type Boo struct {
	Foo int `beancode:"foo"`
	Bar string `Beancode:"bar"`
}
```

To marshal into Bencode:
```go
from := Boo{Foo: 42, Bar: "qux"}
val, err := beancode.Marshal(from)
```

To unmarshal from Bencode:
```go
var to Boo
from := "d3:fooi42e3:bar3:quxe"

err := beancode.Unmarshal(from, &to)
```

Voila, as simple as that!

## Install

```go
go get github.com/dxtym/beancode
```

## Benchmarks

Covered on AMD Ryzen 3 4300U with Radeon Graphics (8GB RAM).

* Marshal: 3652 ns/op, 361 B/op, 17 allocs/op
* Unmarshal: 3672 ns/op, 1515 B/op, 22 allocs/op

## Plans

* Implement decode to struct
* Enhance benchmark performance
* Have some documentation

## License

[MIT License](LICENSE)