package mt940

import (
	"fmt"
	"regexp"
	"strings"
)

// TODO: Add unit tests

type tMessage struct {
	Header       string
	Fields       map[string]interface{}
	Transactions []map[string]interface{}
}

// Parse parses SWIFT message and returns tMessage struct
func Parse(data []byte) (tMessage, error) {
	return parseMessage(string(data))
}

// ParseMultimessage parses SWIFT multimessage and returns slice of tMessage structs
// separator is a string that separates messages e.g. "\r\n$\r\n"
func ParseMultimessage(data []byte, separator string) ([]tMessage, error) {
	messages := strings.Split(string(data), separator)
	var messagesParsed []tMessage
	for _, message := range messages {
		messageParsed, err := parseMessage(message)
		if err != nil {
			return messagesParsed, err
		}
		messagesParsed = append(messagesParsed, messageParsed)
	}
	return messagesParsed, nil
}

// parseMessage parses SWIFT message and returns tMessage struct
func parseMessage(data string) (tMessage, error) {
	message := tMessage{Fields: make(map[string]interface{}), Transactions: make([]map[string]interface{}, 0)}
	header, headerEnd, err := getHeaderLine(&data) // get header
	if err != nil {
		return message, err
	}
	message.Header = header
	messageEnd := strings.Index(data, "\r\n-}") // Find End of message
	if messageEnd == -1 {
		return message, fmt.Errorf("messageEndNotFound")
	}
	data = data[headerEnd:messageEnd]  // get pure message
	fields, err := getAllFields(&data) // get all fields
	if err != nil {
		return message, err
	}
	sortFields(fields, &message) // sort fields to message struct
	return message, nil
}

// getHeaderLine returns header line from SWIFT message
func getHeaderLine(data *string) (header string, eol int, err error) {
	sol := strings.Index(*data, "{1:") // find Start of header
	if sol == -1 {
		return "", -1, fmt.Errorf("headerStartNotFound")
	}
	eol = strings.Index((*data)[sol:], "{4:\r\n") + 3 // find End of header
	if eol == 2 {                                     // -1 + 3(lineAbove) = 2 -> headerEndNotFound
		return "", -1, fmt.Errorf("headerEndNotFound")
	}
	header = (*data)[sol:eol] // get header
	if header == "" {         // check if header is valid
		return "", -1, fmt.Errorf("headerNotFound")
	}
	return
}

// getAllFields returns all fields from SWIFT message
func getAllFields(data *string) (fields []map[string]string, err error) {
	mi := regexp.MustCompile(`\r\n:\w{2,3}:`).FindAllStringIndex(*data, -1) // find all Start message indexes
	for i := 0; i < len(mi); i++ {                                          // loop through all fields and add them to fields slice
		if i < len(mi)-1 {
			fields = append(fields, map[string]string{(*data)[mi[i][0]+3 : mi[i][1]-1]: (*data)[mi[i][1]:mi[i+1][0]]})
		} else {
			fields = append(fields, map[string]string{(*data)[mi[i][0]+3 : mi[i][1]-1]: (*data)[mi[i][1]:]})
		}
	}
	if len(fields) == 0 { // check if fields has been found
		err = fmt.Errorf("fieldsNotFound")
	}
	return
}

// sortFields sorts fields to message struct
func sortFields(fields []map[string]string, message *tMessage) {
	for i := 0; i < len(fields); i++ {
		for k, v := range fields[i] {
			// If field is 61-Transaction, add it to transactions slice
			if k == "61" { // check if field is transaction
				message.Transactions = append(message.Transactions, map[string]interface{}{"F_" + k: v})
				if i+1 <= len(fields) && fields[i+1]["86"] != "" { // check if field 86-TranDetails is present
					message.Transactions[len(message.Transactions)-1]["F_86"] = fields[i+1]["86"]
					i++
				}
				continue
			}
			message.Fields["F_"+k] = v // add field to fields map
		}
	}
}
