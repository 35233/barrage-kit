package stt

import (
	"bytes"
	"sort"
	"strings"
)

func Decode(input string) interface{} {
	var toDoStr = input

	if input[len(input)-1] != '/' {
		toDoStr = input + "/"
	}

	var result interface{}
	tmpKey, tmpValue := bytes.Buffer{}, bytes.Buffer{}

	if strings.Contains(toDoStr, "@=") {
		result = make(map[string]interface{})
	} else {
		result = make([]interface{}, 0, 1)
	}

	for i, ln := 0, len(toDoStr); i < ln; i++ {
		char := toDoStr[i]
		if '/' == char {
			switch res := result.(type) {
			case []interface{}:
				result = append(res, tmpValue.String())
			case map[string]interface{}:
				res[tmpKey.String()] = tmpValue.String()
			}
			tmpKey, tmpValue = bytes.Buffer{}, bytes.Buffer{}
		} else if '@' == char {
			i += 1
			switch toDoStr[i] {
			case 'A':
				tmpValue.WriteByte('@')
				break
			case 'S':
				tmpValue.WriteByte('/')
				break
			case '=':
				tmpKey, tmpValue = tmpValue, tmpKey
			}
		} else {
			tmpValue.WriteByte(char)
		}
	}
	return result
}

func replaceValue(input string) string {
	buff := bytes.Buffer{}
	for i, ln := 0, len(input); i < ln; i++ {
		if input[i] == '@' {
			buff.WriteByte('@')
			buff.WriteByte('A')
		} else if input[i] == '/' {
			buff.WriteByte('@')
			buff.WriteByte('S')
		} else {
			buff.WriteByte(input[i])
		}
	}
	return buff.String()
}
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func Encode(input interface{}) string {
	if str, ok := input.(string); ok {
		return str
	}
	arrInput, isArray := input.([]interface{})
	mapInput, isMap := input.(map[string]interface{})
	buff := bytes.Buffer{}
	if isArray {
		for _, item := range arrInput {
			buff.WriteString(replaceValue(Encode(item)))
			buff.WriteString("/")
		}
		return buff.String()
	}
	if isMap {
		for _, k := range sortedKeys(mapInput) {
			buff.WriteString(replaceValue(Encode(k)))
			buff.WriteString("@=")
			buff.WriteString(replaceValue(Encode(mapInput[k])))
			buff.WriteString("/")
		}
		return buff.String()
	}
	return ""
}
