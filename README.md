# mt940
Parse mt940 message to struct
- Package don't handle any logic on fields level it just simply returns struct with all fields from mt940 message as strings.

## this package is developed for [BaFi](https://github.com/mmalcek/bafi) project
BaFi allows additional parsing on fields level and create formated output using go templates
https://github.com/mmalcek/bafi

### struct
```go
type tMessage struct {
	Header       string // Message header {1:......{4:
	Fields       map[string]interface{} // :20:......, :25:......, ...
	Transactions []map[string]interface{} // []{:61:......, :86:......}
}
```

### example
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mmalcek/mt940"
)

func main() {
	file, err := os.ReadFile("test.sta")
	if err != nil {
		log.Fatal(err)
	}
	message, err := mt940.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message.Fields["F_20"])
}
```
- Parse multiple messages in one file (eg: Multicash)
```go
messages, err := mt940.ParseMultimessage(file, "\r\n$\r\n") // message separator = "\r\n$\r\n"
```
