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
  input := "(YOUR TNSNAMES.Ora CONTENT HERE)"
  tnsFile, err := tnsparser.ParseTNSString(input)
  if err != nil {
      fmt.Println("Error:", err)
      return
  }
  fmt.Println("Parsed TNSFile:", tnsFile)
}
```
