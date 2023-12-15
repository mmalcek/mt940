package mt940

import (
	"regexp"
	"strings"
)

type Statement struct {
	Header       string
	Fields       map[string]interface{}
	Transactions []map[string]interface{}
}

func Parse(file []byte) (Statement, error) {
	return parseFile(string(file))
}

func parseFile(file string) (Statement, error) {
	// TODO: add values basic validation
	re := regexp.MustCompile(`\r\n:\w{2,3}:`)
	mi := re.FindAllStringIndex(file, -1)
	sta := Statement{
		Fields:       make(map[string]interface{}),
		Transactions: make([]map[string]interface{}, 0),
	}
	sta.Header = file[:strings.Index(file, "\r\n")]
	for i := 0; i < len(mi); i++ {
		if i < len(mi)-1 {
			// if transaction
			if file[mi[i][0]+3:mi[i][1]-1] == "61" {
				tran := make(map[string]interface{})
				tran["F_"+file[mi[i][0]+3:mi[i][1]-1]] = file[mi[i][1]:mi[i+1][0]]
				// if include line 86
				if file[mi[i+1][0]+3:mi[i+1][1]-1] == "86" {
					tran["F_"+file[mi[i+1][0]+3:mi[i+1][1]-1]] = file[mi[i+1][1]:mi[i+2][0]]
					i++
				}
				sta.Transactions = append(sta.Transactions, tran)
				continue
			} else {
				sta.Fields["F_"+file[mi[i][0]+3:mi[i][1]-1]] = file[mi[i][1]:mi[i+1][0]]
			}
		} else {
			eol := strings.Index(file[mi[i][1]:], "\r\n")
			sta.Fields["F_"+file[mi[i][0]+3:mi[i][1]-1]] = file[mi[i][1] : mi[i][1]+eol]
		}
	}
	return sta, nil
}

/* TODO:
- Field 86 may appears alone - edit in parser
"As a starting point, field 86 is used to supply extra information related to the preceding statement line in file 61.
However, a separate field 86 can also be supplied.
If other channels than SWIFT NET FIN are used as receiving channel - the first instance of field 86 may contain general information about the account being part of a Group Cash Pool arrangement. This instance is not related to a statement line.
If SWIFT NET FIN is used as receiving channel, the same information may be given as a last Instance of field 86, which is not related to a statement line.
"

- check specification if header fields can repeat and update parser
e.g.
:28C: 00532/001 (First message of statement 532)
:28C: 00532/002 (Second message of statement 532)

- make optional "deepParser" that parse fields by common specification as per example specs in ./docs/MT940_Structured.pdf

*/
