package vdf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func ParseFile(filePath string, caseInsensitive bool) (*KeyValue, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return parse(bufio.NewScanner(f), caseInsensitive)
}

func ParseReader(reader io.Reader, caseInsensitive bool) (*KeyValue, error) {
	return parse(bufio.NewScanner(reader), caseInsensitive)
}

func ParseBytes(b []byte, caseInsensitive bool) (*KeyValue, error) {
	r := bytes.NewReader(b)
	return parse(bufio.NewScanner(r), caseInsensitive)
}

// Based on https://github.com/rossengeorgiev/vdf-parser/blob/master/vdf.js
func parse(scanner *bufio.Scanner, ci bool) (*KeyValue, error) {
	regex, err := regexp.Compile(`"[^"\\]*(?:\\.[^"\\]*)*"`)
	if err != nil {
		return nil, err
	}

	obj := make(map[string]interface{})
	stack := []map[string]interface{}{obj}

	expectBracket := false
	for i := 0; scanner.Scan(); i++ {

		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, `//`) {
			continue
		}

		if line[0] == '{' {
			expectBracket = false
			continue
		}

		if expectBracket {
			return nil, errors.New(fmt.Sprintf("bracket expected; invalid syntax on line %d", i+1))
		}

		if line[0] == '}' {
			stack = stack[:len(stack)-1]
			continue
		}

		for {
			m := regex.FindAllString(line, -1)

			if len(m) == 0 {
				return nil, errors.New(fmt.Sprintf("invalid syntax on line %d", i+1))
			}

			key := strings.Trim(m[0], `"`)
			if ci {
				key = strings.ToLower(key)
			}
			continuing := !strings.HasSuffix(line, `"`) || strings.HasSuffix(line, `\"`)

			if len(m) == 1 && !continuing {
				if _, ok := stack[len(stack)-1][key]; !ok {
					stack[len(stack)-1][key] = make(map[string]interface{})
				}
				stack = append(stack, stack[len(stack)-1][key].(map[string]interface{}))
				expectBracket = true
			} else if continuing {
				scanner.Scan()
				line += "\n" + scanner.Text()
				continue
			} else {
				stack[len(stack)-1][key] = strings.Trim(m[1], `"`)
			}

			break
		}

	}

	if len(stack) != 1 {
		return nil, errors.New("open parentheses somewhere")
	}

	return NewKeyValue(obj, ci)
}
