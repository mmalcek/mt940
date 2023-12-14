package mt940

import (
	"regexp"
	"strings"
)

type Statement struct {
	Header       string
	Fields       map[string]string
	Transactions []map[string]string
}

func Parse(file []byte) (Statement, error) {
	return parseFile(string(file))
}

func parseFile(file string) (Statement, error) {
	// TODO: add values basic validation
	re := regexp.MustCompile(`\r\n:\w{2,3}:`)
	mi := re.FindAllStringIndex(file, -1)
	sta := Statement{
		Fields:       make(map[string]string),
		Transactions: make([]map[string]string, 0),
	}
	sta.Header = file[:strings.Index(file, "\r\n")]
	for i := 0; i < len(mi); i++ {
		if i < len(mi)-1 {
			// if transaction
			if file[mi[i][0]+3:mi[i][1]-1] == "61" {
				tran := make(map[string]string)
				tran[file[mi[i][0]+3:mi[i][1]-1]] = file[mi[i][1] : mi[i+1][0]-1]
				// if include line 86
				if file[mi[i+1][0]+3:mi[i+1][1]-1] == "86" {
					tran[file[mi[i+1][0]+3:mi[i+1][1]-1]] = file[mi[i+1][1] : mi[i+2][0]-1]
					i++
				}
				sta.Transactions = append(sta.Transactions, tran)
				continue
			} else {
				sta.Fields[file[mi[i][0]+3:mi[i][1]-1]] = file[mi[i][1] : mi[i+1][0]-1]
			}
		} else {
			eol := strings.Index(file[mi[i][1]:], "\r\n")
			sta.Fields[file[mi[i][0]+3:mi[i][1]-1]] = file[mi[i][1] : mi[i][1]+eol]
		}
	}
	return sta, nil
}
