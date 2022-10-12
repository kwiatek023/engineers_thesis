package utils

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"math"
	"strconv"
)

func ParseStrToPositiveInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil || number <= 0 {
		// TODO handle error
		panic(err)
	}

	return number
}

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
		fmt.Println("Probability should be in range [0,1]")
		panic(nil)
		//TODO raise error
	}
	return p
}
