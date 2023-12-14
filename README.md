# mt940
Parse mt940 message to struct

- Package don't handle any logic on fields level it just simply returns struct with all fields from mt940 message.
```go
type Statement struct {
	Header       string
	Fields       map[string]string
	Transactions []map[string]string
}
```
eg: Fields -> [25]:C343201CZK2382150,25


## usage
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
	statement, err := mt940.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(statement.Header)
}
```

- Note: Currently in development. There are no any validations or result checks. So use it on your own risk.