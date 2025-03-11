// Package query contains core functionality of internal scripting engine
package query

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

type operator func(string, string) (bool, error)

var operators = map[string]operator{
	"is": func(s1, s2 string) (bool, error) { return s1 == s2, nil },
	"match": func(s1, s2 string) (bool, error) {
		r, err := regexp.Compile(s1)
		if err != nil {
			slog.Error("invalid regex", slog.Any("err", err), slog.Any("pattern", s1))
			return false, fmt.Errorf("invalid regex: %v", err)
		}
		return r.MatchString(s2), nil
	},
	"contains": func(s1, s2 string) (bool, error) {
		slog.Debug("checking contains", slog.Any("s1", s1), slog.Any("s2", s2))
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
		slog.Error("invalid query format", slog.String("query", query))
		return nil, fmt.Errorf("invalid query format: %s", query)
	}

	field := parts[0]
	op := parts[1]
	value := strings.Join(parts[2:], " ")
	value = strings.Trim(value, "\"") // Remove surrounding quotes

	if _, ok := operators[op]; !ok {
		slog.Error("unsupported operator", slog.String("op", op))
		return nil, fmt.Errorf("unsupported operator: %s", op)
	}

	slog.Debug("parsed query", slog.String("field", field), slog.String("operator", op), slog.String("value", value))
	return &Query{Field: field, Op: op, Value: value}, nil
}

func (q *Query) EvaluateQuery(data map[string]string) (bool, error) {
	val, exists := data[q.Field]
	if !exists {
		slog.Info("field not found, evaluating as a fixed value", slog.String("field", q.Field))
		val = q.Field
	}

	// Check if the operator exists and evaluate
	op, ok := operators[q.Op]
	if !ok {
		slog.Error("unknown operation", slog.String("op", q.Op))
		return false, fmt.Errorf("unknown operation: %s", q.Op)
	}

	slog.Debug("evaluating query", slog.String("field", q.Field), slog.String("value", val), slog.Any("operator", q.Op), slog.Any("query_value", q.Value))

	result, err := op(val, q.Value)
	if err != nil {
		slog.Error("error evaluating operator", slog.String("op", q.Op), slog.Any("err", err))
		return false, err
	}

	slog.Debug("query evaluation result", slog.String("field", q.Field), slog.Any("result", result))
	return result, nil
}
