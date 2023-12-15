# mt940
Parse mt940 message to struct

- Package don't handle any logic on fields level it just simply returns struct with all fields from mt940 message.
```go
type Statement struct {
	Header       string
	Fields       map[string]interface{}
	Transactions []map[string]interface{}
}
```
- eg: 
```go
fmt.Println(statement.Fields["F_20"])
```

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
	fmt.Println(statement.Fields["F_20"])
}
```

- Note: Currently in development. There are no any validations or result checks. So use it on your own risk.