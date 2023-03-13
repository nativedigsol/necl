package necl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// performComparison will make a comparison check against 2 values and return a boolean as an interface
func PerformComparison(lineRaw string, currentAttributes map[string]Attribute) (bool, error) {
	// Discover the comparison
	comparisons := []string{"==", "!=", "<", "<=", ">", ">="}
	comparison := ""
	for _, c := range comparisons {
		if strings.Contains(lineRaw, c) {
			comparison = c
		}
	}
	if comparison == "" {
		err := fmt.Errorf("unknown comparator on line: %s", lineRaw)
		return false, err
	}

	// Get values
	value1Str := ""
	value2Str := ""
	double := false
	i := strings.Index(lineRaw, comparison)

	for _, c := range []string{"==", "!=", "<=", ">="} {
		if comparison == c {
			double = true
		}
	}
	value1Str = strings.TrimSpace(lineRaw[:i-1])
	if double {
		value2Str = strings.TrimSpace(lineRaw[i+2:])
	} else {
		value2Str = strings.TrimSpace(lineRaw[i+1:])
	}

	// Check if values are a number or an attribute's name
	var value1 interface{}
	var value2 interface{}

	value1, err := strconv.Atoi(value1Str)
	if err != nil {
		// Look for an attribute with this name
		if currentAttributes[value1Str].Value == nil {
			err := fmt.Errorf("no attribute named %s was found on line %s", value1Str, lineRaw)
			return false, err
		}
		value1 = currentAttributes[value1Str].Value
	}

	value2, err = strconv.Atoi(value2Str)
	if err != nil {
		// Look for an attribute with this name
		if currentAttributes[value2Str].Value == nil {
			err := fmt.Errorf("no attribute named %s was found on line %s", value2Str, lineRaw)
			return false, err
		}
		value2 = currentAttributes[value2Str].Value
	}

	// Transform interfaces in an int
	if strings.Contains(fmt.Sprintf("%v", value1), ".") || strings.Contains(fmt.Sprintf("%v", value1), ",") {
		err := errors.New("comparison operations can only be done to integer values")
		return false, err
	}
	if strings.Contains(fmt.Sprintf("%v", value2), ".") || strings.Contains(fmt.Sprintf("%v", value2), ",") {
		err := errors.New("comparison operations can only be done to integer values")
		return false, err
	}

	v1 := value1.(int)
	v2 := value2.(int)

	// Make comparison
	switch comparison {
	case "==":
		return v1 == v2, nil
	case "!=":
		return v1 != v2, nil
	case ">":
		return v1 > v2, nil
	case ">=":
		return v1 >= v2, nil
	case "<":
		return v1 < v2, nil
	case "<=":
		return v1 <= v2, nil
	}

	err = fmt.Errorf("unknown error")
	return false, err
}

// performArithmeticOperation performs an arithmetic operation with integers
func PerformArithmeticOperation(lineRaw string, currentAttributes map[string]Attribute) (int, error) {
	// Discover the operation
	operations := []string{"+", "-", "*", "/"}
	operation := ""
	for _, o := range operations {
		if strings.Contains(lineRaw, o) {
			operation = o
		}
	}
	if operation == "" {
		err := fmt.Errorf("unknown operator on line: %s", lineRaw)
		return 0, err
	}

	// Get values
	value1Str := ""
	value2Str := ""
	i := strings.Index(lineRaw, operation)
	value1Str = strings.TrimSpace(lineRaw[:i])
	value2Str = strings.TrimSpace(lineRaw[i+1:])

	// Check if values are a number or an attribute's name
	var value1Interface interface{}
	var value2Interface interface{}

	value1Interface, err := strconv.Atoi(value1Str)
	if err != nil {
		// Look for an attribute with this name
		if currentAttributes[value1Str].Value == nil {
			err := fmt.Errorf("no attribute named %s was found on line %s", value1Str, lineRaw)
			return 0, err
		}
		value1Interface = currentAttributes[value1Str].Value
	}

	value2Interface, err = strconv.Atoi(value2Str)
	if err != nil {
		// Look for an attribute with this name
		if currentAttributes[value2Str].Value == nil {
			err := fmt.Errorf("no attribute named %s was found on line %s", value2Str, lineRaw)
			return 0, err
		}
		value2Interface = currentAttributes[value2Str].Value
	}

	value1 := value1Interface.(int)
	value2 := value2Interface.(int)

	// Perform the operation
	result := 0
	switch operation {
	case "+":
		result = value1 + value2
	case "-":
		result = value1 - value2
	case "*":
		result = value1 * value2
	case "/":
		result = value1 / value2
	default:
		err := fmt.Errorf("unknown operation %s", operation)
		return 0, err
	}

	return result, nil
}
