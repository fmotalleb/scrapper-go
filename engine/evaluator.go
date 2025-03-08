package engine

import (
	"fmt"
	"regexp"
	"strings"
)

type operator func(string, string) (bool, error)

var operators = map[string]operator{
	"is": func(s1, s2 string) (bool, error) { return s1 == s2, nil },
	"match": func(s1, s2 string) (bool, error) {
		r, err := regexp.Compile(s1)
		if err != nil {
			return false, err
		}
		return r.MatchString(s2), nil
	},
	"contains": func(s1, s2 string) (bool, error) {
		fmt.Println(s1, s2)
		return strings.Contains(s1, s2), nil
	},
}

type Query struct {
	Field string
	Op    string
	Value string
}

func ParseQuery(query string) (*Query, error) {
	parts := strings.Fields(query)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid query format")
	}

	field := parts[0]
	op := parts[1]
	value := strings.Join(parts[2:], " ")
	value = strings.Trim(value, "\"") // Remove surrounding quotes

	if _, ok := operators[op]; !ok {
		return nil, fmt.Errorf("unsupported operator: %s", op)
	}

	return &Query{Field: field, Op: op, Value: value}, nil
}

func (q *Query) EvaluateQuery(data map[string]string) (bool, error) {
	val, exists := data[q.Field]
	if !exists {
		return false, fmt.Errorf("field %s not found", q.Field)
	}

	if op, ok := operators[q.Op]; ok {
		return op(val, q.Value)
	}

	return false, fmt.Errorf("unknown operation")
}
