
# amp

  Abstract Message Protocol codec and streaming parser for golang.

## Installation

```
$ go get github.com/nickpoorman/go-amp
```

## Example

TODO

## Protocol

  AMP is a simple versioned protocol for framed messages containing
  zero or more "arguments". Each argument is opaque binary, thus you
  may use JSON, BSON, msgpack and others on top of AMP. Multiple argument
  support is used to allow a hybrid of binary/non-binary message args without
  requiring higher level serialization libraries like msgpack or BSON.

  All multi-byte integers are big endian. The `version` and `argc` integers
  are stored in the first byte, followed by a sequence of zero or more
  `<length>` / `<data>` pairs, where `length` is a 32-bit unsigned integer.

```
      0        1 2 3 4     <length>    ...
+------------+----------+------------+
| <ver/argc> | <length> | <data>     | additional arguments
+------------+----------+------------+
```

## Implementations

 - [c](https://github.com/clibs/amp) (~10m ops/s)
 - [node](https://github.com/visionmedia/node-amp) (~1.5m ops/s)
 - [golang](https://github.com/nickpoorman/go-amp) this library --not yet benchmarked

# License

  MIT
