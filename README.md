# wkt-s2-decoder


Library that converts wkt into [S2](https://github.com/golang/geo) datatypes, relies mostly on https://github.com/IvanZagoskin/wkt

## Install

```bash
go get -u github.com/purtato/wkt-s2-decoder
```

## Example
```go
package main

import (
	"bytes"
	"fmt"
	decoder "github.com/purtato/wkt-s2-decoder"
)

func main() {
	d := decoder.New()
	poly, _ := d.ParseLinestring(bytes.NewReader([]byte("LINESTRING (30 10, 10 30, 40 40)")))
	fmt.Printf("%+v", poly)
}

```
