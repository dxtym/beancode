## Beancode

Working with BitTorrent? You need Beancode.

## Features

* Easy API based on stdlib
* Almost full unit test coverage
* Support int, str, array, and map

## Usage

To marshal into Bencode:
```go
from := struct {
	Foo []string `bencode:"foo"`
	Boo struct {
		Foo int `bencode:"foo"`
		Bar int `bencode:"bar"`
	} `bencode:"boo"`
}{
	Foo: []string{"boo", "foo"},
	Boo: struct {
		Foo int `bencode:"foo"`
		Bar int `bencode:"bar"`
	}{
		Foo: 100,
		Bar: 100,
	},
}

val, err := beancode.Marshal(from)
```

To unmarshal from Bencode:
```go
var to struct {
	Foo []string `bencode:"foo"`
	Boo struct {
		Foo int `bencode:"foo"`
		Bar int `bencode:"bar"`
	} `bencode:"boo"`
}
from := "d3:fool3:boo3:bare3:bood3:fooi100e3:bari100eee"

err := beancode.Unmarshal(from, &to)
```

## Install

```go
go get github.com/dxtym/beancode
```

## Benchmarks

Covered on AMD Ryzen 3 4300U with Radeon Graphics (8GB RAM).

* Marshal: 2535 ns/op, 454 B/op, 16 allocs/op
* Unmarshal: 2658 ns/op, 1499 B/op, 21 allocs/op

## Plans

* Have some documentation

## License

[MIT License](LICENSE)