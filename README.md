# tns_parser

You ever have 12 different oracle databases configured to use TNS and everyone has root access to the server so every config file is different?

### Install

```bash
go get github.com/adam-baker/tns_parser
```

### Usage

```go
package main

import (
  "fmt"
  "github.com/adam-baker/tns_parser"
)

func main() {
  tns := tns_parser.NewTNSParser("/path/to/tnsnames.ora")
  fmt.Println(tns.Get("DB_NAME"))
}
```
