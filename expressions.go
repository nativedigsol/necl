package necl

import (
	"fmt"
	"strconv"
	"strings"
)

// Transforms a string into an interface for the IfExpression function
func parseStringToInterfaceIfExpression(value string, attributes map[string]Attribute) (string, interface{}, error) {
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
func IfExpression(line string, attributes map[string]Attribute) (string, interface{}, error) {
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
	positiveType, positiveValue, err := parseStringToInterfaceIfExpression(positiveOutcome, attributes)
	if err != nil {
		return "", nil, err
	}
	negativeType, negativeValue, err := parseStringToInterfaceIfExpression(negativeOutcome, attributes)
	if err != nil {
		return "", nil, err
	}

	// Get condition
	lineNoOutcomes := line[:positiveOutcomeIndex-1]
	conditionRaw := strings.TrimSpace(strings.TrimPrefix(lineNoOutcomes, "if"))
	conditionType, condition, err := parseStringToInterfaceIfExpression(conditionRaw, attributes)
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

// Transforms a string into an interface for the ForExpression function
func parseStringToInterfaceForExpression(value string, attributes map[string]Attribute) (string, error) {
	// Comparison (will transform into a boolean by the end)
	if ContainsMany(value, []string{"==", "!=", "<", "<=", ">", ">="}) {
		return "comparison", nil
	}

	// Arithmetic operation
	if ContainsMany(value, []string{"+", "-", "*", "/"}) {
		return "arithmetic", nil
	}

	// String functions
	if ContainsMany(value, []string{"upper(", "lower(", "concat(", "contains(", "length("}) {
		return "func-string", nil
	}

	// Mathematical functions
	if ContainsMany(value, []string{"power(", "floor(", "remainder("}) {
		return "func-math", nil
	}

	// Logical functions
	if ContainsMany(value, []string{"and(", "or(", "nand(", "nor(", "xor(", "xnor("}) {
		return "func-logic", nil
	}

	// Attribute
	if (attributes[value].Value != nil) || ContainsMany(value, []string{"index", "value"}) {
		return "attribute", nil
	}

	// Unknown type
	err := fmt.Errorf("unknown type for %s", value)
	return "unknown", err
}

// forExpression will create a collection by projecting the items from another collection into it
func ForExpression(line string, attributes map[string]Attribute) ([]interface{}, error) {
	// Get outcome index
	outcomeIndex := 0
	for i := len(line) - 1; i >= 0; i-- {
		if string(line[i]) == ":" && outcomeIndex == 0 {
			outcomeIndex = i
		}
	}

	if outcomeIndex == 0 {
		err := fmt.Errorf("missing outcome in line %s", line)
		return nil, err
	}

	// Get outcome
	outcome := strings.TrimSpace(line[outcomeIndex+1:])

	// Parse outcome
	outcomeType, err := parseStringToInterfaceForExpression(outcome, attributes)
	if err != nil {
		return nil, err
	}

	// Get condition
	lineNoOutcome := line[:outcomeIndex-1]
	conditionRaw := strings.TrimSpace(strings.TrimPrefix(lineNoOutcome, "for"))

	// Check if the condition is an array or an attribute (must be array type)
	conditionArray := attributes[conditionRaw].Array
	if conditionArray == nil {
		if strings.HasPrefix(conditionRaw, `[`) && strings.HasSuffix(conditionRaw, `]`) {
			// Parse array
			conditionArray, err = parseArrayAttributes(conditionRaw, attributes)
			if err != nil {
				return nil, err
			}
		} else {
			err := fmt.Errorf("condition to a 'for' expression must be either a call to an array attribute, or a definition of an array: %s", line)
			return nil, err
		}
	}

	// Create the result array
	resultArray := []interface{}{}

	// Loop through elements of the array with the condition
	for index, value := range conditionArray {
		// Update index and value on the map of attributes
		attributes["index"] = Attribute{
			Name:  "index",
			Value: index,
		}
		attributes["value"] = Attribute{
			Name:  "value",
			Value: value,
		}

		// Calculate condition
		var newEntry interface{}
		switch outcomeType {
		case "attribute":
			newEntry = attributes[outcome].Value
			if newEntry == nil {
				err := fmt.Errorf("unknown attribute %s for expression %s", outcome, line)
				return nil, err
			}
		case "comparison":
			newEntry, err = PerformComparison(outcome, attributes)
			if err != nil {
				return nil, err
			}
		case "arithmetic":
			newEntry, err = PerformArithmeticOperation(outcome, attributes)
			if err != nil {
				return nil, err
			}
		case "func-string":
			_, newEntry, err = StringFunctions(outcome, attributes)
			if err != nil {
				return nil, err
			}
		case "func-math":
			newEntry, err = MathFunctions(outcome, attributes)
			if err != nil {
				return nil, err
			}
		case "func-logic":
			newEntry, err = LogicFunctions(outcome, attributes)
			if err != nil {
				return nil, err
			}
		default:
			err := fmt.Errorf("unknown attribute type %s for expression %s", outcome, line)
			return nil, err
		}
		resultArray = append(resultArray, newEntry)
	}

	return resultArray, nil
}
