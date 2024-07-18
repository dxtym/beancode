## Beancode

Working with BitTorrent? You need Beancode.

## Features

* Easy API based on stdlib
* Almost full unit test coverage
* Support int, str, array, and map

## Usage

To marshal into Bencode:
```go
in := struct {
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

val, err := beancode.Marshal(in)
```

To unmarshal from Bencode:
```go
var out struct {
	Foo []string `bencode:"foo"`
	Boo struct {
		Foo int `bencode:"foo"`
		Bar int `bencode:"bar"`
	} `bencode:"boo"`
}
in := "d3:fool3:boo3:bare3:bood3:fooi100e3:bari100eee"

err := beancode.Unmarshal(in, &out)
```

## Install

```go
go get github.com/dxtym/beancode
```

## Benchmarks

Covered on AMD Ryzen 3 4300U with Radeon Graphics (8GB RAM).

* Marshal: 1787 ns/op, 457 B/op, 19 allocs/op
* Unmarshal: 1397 ns/op, 1435 B/op, 19 allocs/op

## Plans

* Add struct decode feature
* Get rid of external dependencies
* Work automated testing
* Have some documentation

## License

[MIT License](LICENSE)