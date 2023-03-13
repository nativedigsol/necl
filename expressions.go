package necl

import (
	"fmt"
	"strconv"
	"strings"
)

// Transforms a string into an interface
func parseStringToInterface(value string, attributes map[string]Attribute) (string, interface{}, error) {
	// Attribute
	if attributes[value].Value != nil {
		return attributes[value].Type, attributes[value].Value, nil
	}

	// String
	if (strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) || (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) {
		return "string", value[1 : len(value)-1], nil
	}

	// Comparison
	if ContainsMany(value, []string{"==", "!=", "<", "<=", ">", ">="}) {
		return "comparison", value, nil
	}

	// Functions
	// String functions
	if strings.Contains(value, "contains(") {
		return "func-str", value, nil
	}

	// Logical functions
	if ContainsMany(value, []string{"and(", "or(", "nand(", "nor(", "xor(", "xnor("}) {
		return "func-logic", value, nil
	}

	// Boolean
	if strings.Contains(strings.ToLower(value), "true") || strings.Contains(strings.ToLower(value), "false") {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return "boolean", b, err
		}
		return "boolean", b, nil
	}

	// Number
	_, err := strconv.ParseFloat(value, 32)
	if err == nil {
		if strings.Contains(value, ".") || strings.Contains(value, ",") {
			val, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return "", nil, err
			}
			return "number", val, nil
		} else {
			val, err := strconv.Atoi(value)
			if err != nil {
				return "", nil, err
			}
			return "number", val, nil
		}
	}

	// Unknown type
	err = fmt.Errorf("unknown type for %s", value)
	return "unknown", nil, err
}

// ifExpression will calculate the value of an attribute with an "if" expression
func ifExpression(line string, attributes map[string]Attribute) (string, interface{}, error) {
	// Get outcome indexes
	positiveOutcomeIndex := 0
	negativeOutcomeIndex := 0

	for i := len(line) - 1; i >= 0; i-- {
		// Negative outcome
		if string(line[i]) == ":" && negativeOutcomeIndex == 0 {
			negativeOutcomeIndex = i
		}

		// Positive outcome
		// Only start looking for it when the negative outcome index is found
		if string(line[i]) == "?" && positiveOutcomeIndex == 0 && negativeOutcomeIndex != 0 {
			positiveOutcomeIndex = i
		}
	}

	if positiveOutcomeIndex == 0 || negativeOutcomeIndex == 0 {
		err := fmt.Errorf("missing outcome in line %s", line)
		return "", nil, err
	}

	// Get outcomes
	positiveOutcome := strings.TrimSpace(line[positiveOutcomeIndex+1 : negativeOutcomeIndex-1])
	negativeOutcome := strings.TrimSpace(line[negativeOutcomeIndex+1:])

	// Parse outcomes
	positiveType, positiveValue, err := parseStringToInterface(positiveOutcome, attributes)
	if err != nil {
		return "", nil, err
	}
	negativeType, negativeValue, err := parseStringToInterface(negativeOutcome, attributes)
	if err != nil {
		return "", nil, err
	}

	// Get condition
	lineNoOutcomes := line[:positiveOutcomeIndex-1]
	conditionRaw := strings.TrimSpace(strings.TrimPrefix(lineNoOutcomes, "if"))
	conditionType, condition, err := parseStringToInterface(conditionRaw, attributes)
	if err != nil {
		return "", nil, err
	}

	// Invalid types for if condition
	if ContainsMany(conditionType, []string{"string", "array", "arithmetic", "func-math", "number"}) {
		err := fmt.Errorf("invalid type %s for condition %s", conditionType, condition)
		return "", nil, err
	}

	// func-strings has some invalid functions
	if ContainsMany(conditionRaw, []string{"upper(", "lower(", "concat(", "length("}) {
		err := fmt.Errorf("invalid function %s for condition %s. only contains() is a valid funciton for if condition", conditionType, condition)
		return "", nil, err
	}

	// Calculate condition
	var conditionResult bool
	var resultValue interface{}
	var resultType string
	switch conditionType {
	case "boolean":
		conditionResult = condition.(bool)
	case "comparison":
		conditionResult, err = PerformComparison(condition.(string), attributes)
		if err != nil {
			return "", nil, err
		}
	case "func-str":
		_, conditionResultInterface, err := StringFunctions(condition.(string), attributes)
		if err != nil {
			return "", nil, err
		}
		conditionResult = conditionResultInterface.(bool)
	case "func-logic":
		conditionResultInterface, err := LogicFunctions(condition.(string), attributes)
		if err != nil {
			return "", nil, err
		}
		conditionResult = conditionResultInterface.(bool)
	default:
		err := fmt.Errorf("unknown condition: %s", condition)
		return "", nil, err
	}

	if conditionResult {
		resultType = positiveType
		resultValue = positiveValue
	} else {
		resultType = negativeType
		resultValue = negativeValue
	}

	return resultType, resultValue, nil
}

// ParseExpression is a super set for all expressions
func ParseExpression(line string, attributes map[string]Attribute) (string, interface{}, error) {
	if strings.Contains(line, "if") {
		resultType, resultValue, err := ifExpression(line, attributes)
		if err != nil {
			return "", nil, err
		}
		return resultType, resultValue, nil
	}

	err := fmt.Errorf("unknown condition on line: %s", line)
	return "", nil, err
}
