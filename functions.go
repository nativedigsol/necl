package necl

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Gets elements required for a string function
func getValuesForStringFunc(targetsRaw string, attributes map[string]Attribute) ([]string, error) {
	var targets []string
	previousComma := 0
	for commaIndex, findComma := range targetsRaw {
		if string(findComma) == "," {
			newElement := strings.TrimSpace(targetsRaw[previousComma:commaIndex])
			targets = append(targets, newElement)
			previousComma = int(commaIndex) + 1
		}

		// Last element
		if commaIndex == len(targetsRaw)-1 && targets != nil {
			newElement := strings.TrimSpace(targetsRaw[previousComma+1:])
			targets = append(targets, newElement)
		}
	}

	// Confirm there are 2 values
	if targets[1] == "" {
		err := fmt.Errorf("this function requires two values on: %s", targetsRaw)
		return nil, err
	}

	// Search for attributes and transform into string
	var targetsString []string
	for _, val := range targets {
		var stringVal string
		stringVal = string(val)

		// Search for attribute with name
		if attributes[val].Value != nil {
			stringVal = attributes[val].Value.(string)
		} else {
			stringVal = stringVal[1 : len(stringVal)-1]
		}

		// Remove quote if not attribute
		targetsString = append(targetsString, stringVal)
	}

	return targetsString, nil
}

// StringFunctions is a super set of all string functions
func StringFunctions(line string, attributes map[string]Attribute) (string, interface{}, error) {
	// Upper
	if strings.Contains(line, "upper(") {
		// Get func value
		target := strings.TrimSpace(line[6 : len(line)-1])

		// Search for attribute with name
		if attributes[target].Value != nil {
			target = attributes[target].Value.(string)
		} else {
			target = target[1 : len(target)-1]
		}

		// Perform operation
		return "upper", strings.ToUpper(target), nil
	}

	// Lower
	if strings.Contains(line, "lower(") {
		// Get func value
		target := strings.TrimSpace(line[6 : len(line)-1])

		// Search for attribute with name
		if attributes[target].Value != nil {
			target = attributes[target].Value.(string)
		} else {
			target = target[1 : len(target)-1]
		}

		// Perform operation
		return "lower", strings.ToLower(target), nil
	}

	// Concat
	if strings.Contains(line, "concat(") {
		// Get func values
		targetsRaw := strings.TrimSpace(line[7 : len(line)-1])

		// Get values
		targets, err := getValuesForStringFunc(targetsRaw, attributes)
		if err != nil {
			return "concat", "", err
		}

		// Perform the operation
		return "concat", strings.Join(targets, " "), nil
	}

	// Contains
	if strings.Contains(line, "contains(") {
		// Get func values
		targetsRaw := strings.TrimSpace(line[9 : len(line)-1])

		// Get values
		targets, err := getValuesForStringFunc(targetsRaw, attributes)
		if err != nil {
			return "contains", "", err
		}
		// Perform the operation
		return "contains", strings.Contains(targets[0], targets[1]), nil
	}

	// Length
	if strings.Contains(line, "length(") {
		// Get func value
		target := strings.TrimSpace(line[7 : len(line)-1])

		// Search for attribute with name
		if attributes[target].Value != nil {
			target = attributes[target].Value.(string)
		} else {
			target = target[1 : len(target)-1]
		}

		// Perform operation
		return "length", len(target), nil
	}

	err := fmt.Errorf("unknown function on %s", line)
	return "", "", err
}

// Gets elements required for a mathematical function
func getValuesForMathFunc(targetsRaw string, attributes map[string]Attribute) ([]int, error) {
	var targets []string
	previousComma := 0
	for commaIndex, findComma := range targetsRaw {
		if string(findComma) == "," {
			newElement := strings.TrimSpace(targetsRaw[previousComma:commaIndex])
			targets = append(targets, newElement)
			previousComma = int(commaIndex) + 1
		}

		// Last element
		if commaIndex == len(targetsRaw)-1 && targets != nil {
			newElement := strings.TrimSpace(targetsRaw[previousComma+1:])
			targets = append(targets, newElement)
		}
	}

	// Confirm there are 2 values
	if targets[1] == "" {
		err := fmt.Errorf("this function requires two values on: %s", targetsRaw)
		return nil, err
	}

	// Search for attributes and transform into int
	var targetsInt []int
	for _, val := range targets {
		var intVal int
		intVal, err := strconv.Atoi(val)
		if err != nil {
			// Search for attribute with name
			if attributes[val].Value == nil {
				err := fmt.Errorf("no attribute named %s was found", val)
				return nil, err
			}
			intVal = attributes[val].Value.(int)
		}

		targetsInt = append(targetsInt, intVal)
	}

	return targetsInt, nil
}

// MathFunctions is a super set of all mathematical functions
func MathFunctions(line string, attributes map[string]Attribute) (interface{}, error) {
	// Power
	if strings.Contains(line, "power(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[6 : len(line)-1])

		// Get targets
		targets, err := getValuesForMathFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}

		return int(math.Pow(float64(targets[0]), float64(targets[1]))), nil
	}

	// Floor division
	if strings.Contains(line, "floor(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[6 : len(line)-1])

		// Get targets
		targets, err := getValuesForMathFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}

		return int(math.Floor(float64(targets[0]) / float64(targets[1]))), nil
	}

	// Remainder
	if strings.Contains(line, "remainder(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[10 : len(line)-1])

		// Get targets
		targets, err := getValuesForMathFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}

		return targets[0] % targets[1], nil
	}

	err := fmt.Errorf("unknown function on %s", line)
	return "", err
}

// Gets elements required for a logical function
func getValuesForLogicFunc(targetsRaw string, attributes map[string]Attribute) ([]bool, error) {
	var targets []string
	previousComma := 0
	for commaIndex, findComma := range targetsRaw {
		if string(findComma) == "," {
			newElement := strings.TrimSpace(targetsRaw[previousComma:commaIndex])
			targets = append(targets, newElement)
			previousComma = int(commaIndex) + 1
		}

		// Last element
		if commaIndex == len(targetsRaw)-1 && targets != nil {
			newElement := strings.TrimSpace(targetsRaw[previousComma+1:])
			targets = append(targets, newElement)
		}
	}

	// Confirm there are 2 values
	if targets[1] == "" {
		err := fmt.Errorf("this function requires two values on: %s", targetsRaw)
		return nil, err
	}

	// Search for attributes and transform into boolean
	var targetsBool []bool
	for _, val := range targets {
		var boolVal bool
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			// Search for attribute with name
			if attributes[val].Value == nil {
				err := fmt.Errorf("no attribute named %s was found", val)
				return nil, err
			}
			boolVal = attributes[val].Value.(bool)
		}

		targetsBool = append(targetsBool, boolVal)
	}

	return targetsBool, nil
}

// LogicFunctions is a super set of all mathematical functions
func LogicFunctions(line string, attributes map[string]Attribute) (interface{}, error) {

	// XOR
	if strings.Contains(line, "xor(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[4 : len(line)-1])

		// Get targets
		targets, err := getValuesForLogicFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}
		return targets[0] != targets[1], nil
	}

	// XNOR
	if strings.Contains(line, "xnor(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[5 : len(line)-1])

		// Get targets
		targets, err := getValuesForLogicFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}
		return !(targets[0] != targets[1]), nil
	}

	// NAND
	if strings.Contains(line, "nand(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[5 : len(line)-1])

		// Get targets
		targets, err := getValuesForLogicFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}
		return !(targets[0] == targets[1]), nil
	}

	// NOR
	if strings.Contains(line, "nor(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[4 : len(line)-1])

		// Get targets
		targets, err := getValuesForLogicFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}
		return !(targets[0] || targets[1]), nil
	}

	// AND
	if strings.Contains(line, "and(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[4 : len(line)-1])

		// Get targets
		targets, err := getValuesForLogicFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}
		return targets[0] == targets[1], nil
	}

	// OR
	if strings.Contains(line, "or(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[3 : len(line)-1])

		// Get targets
		targets, err := getValuesForLogicFunc(targetsRaw, attributes)
		if err != nil {
			return nil, err
		}
		return targets[0] || targets[1], nil
	}

	err := fmt.Errorf("unknown function on %s", line)
	return false, err
}
