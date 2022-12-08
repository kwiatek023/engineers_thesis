package utils

import (
	"github.com/Knetic/govaluate"
	"log"
	"math"
	"strconv"
)

// ParseStrToPositiveInt - parses str to positive int
func ParseStrToPositiveInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil || number <= 0 {
		log.Fatal("Could not parse to positive integer. Make sure parameters are positive integers.")
	}

	return number
}

// EvaluateExpression - evaluates value based on given expression
func EvaluateExpression(expr string, parameters map[string]interface{}) float64 {
	if expr == "" {
		return 0.0
	}
	functions := map[string]govaluate.ExpressionFunction{
		"log": func(args ...interface{}) (interface{}, error) {
			parameter := args[0].(float64)
			result := math.Log(parameter)
			return result, nil
		},
	}

	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	result, _ := expression.Evaluate(parameters)
	p := result.(float64)
	if p < 0 || p > 1 {
		log.Fatal("Probability should be in range [0,1].")
	}
	return p
}
