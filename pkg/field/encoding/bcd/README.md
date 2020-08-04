# bcd
[![Build Status](https://travis-ci.org/albenik/bcd.svg?branch=master)](https://travis-ci.org/albenik/bcd)
[![GoDoc](https://godoc.org/github.com/albenik/bcd?status.svg)](https://godoc.org/github.com/albenik/bcd)

Go implementation of BCD conversion functions

Usage

```go
package main

import (
	"fmt"
	"github.com/albenik/bcd"
)

func main() {
	fmt.Printf("Uint32: %d", bcd.ToUint32([]byte{0x11, 0x22, 0x33, 0x44}))
	fmt.Printf("BCD: %x", bcd.FromUint32(11223344))
}
```

## Documentation
For more documentation see [package documentation](https://godoc.org/github.com/albenik/bcd)