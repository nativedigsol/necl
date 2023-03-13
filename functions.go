package necl

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Gets elements required for a string function
func getValuesForStringFunc(targetsRaw string) ([]string, error) {
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

	// Remove quote sign from elements
	if strings.HasPrefix(targets[0], `"`) {
		targets[0] = strings.Trim(targets[0], `"`)
	} else if strings.HasPrefix(targets[0], `'`) {
		targets[0] = strings.Trim(targets[0], `'`)
	}
	if strings.HasPrefix(targets[1], `"`) {
		targets[1] = strings.Trim(targets[1], `"`)
	} else if strings.HasPrefix(targets[1], `'`) {
		targets[1] = strings.Trim(targets[1], `'`)
	}

	return targets, nil
}

// StringFunctions is a super set of all string functions
func StringFunctions(line string) (string, interface{}, error) {
	// Upper
	if strings.Contains(line, "upper(") {
		// Get func value
		target := strings.TrimSpace(line[6 : len(line)-1])

		// Remove quote signs
		target = strings.Trim(target, `"`)
		target = strings.Trim(target, `'`)

		// Perform operation
		return "upper", strings.ToUpper(target), nil
	}

	// Lower
	if strings.Contains(line, "lower(") {
		// Get func value
		target := strings.TrimSpace(line[6 : len(line)-1])

		// Remove quote signs
		target = strings.Trim(target, `"`)
		target = strings.Trim(target, `'`)

		// Perform operation
		return "lower", strings.ToLower(target), nil
	}

	// Concat
	if strings.Contains(line, "concat(") {
		// Get func values
		targetsRaw := strings.TrimSpace(line[7 : len(line)-1])

		// Get values
		targets, err := getValuesForStringFunc(targetsRaw)
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
		targets, err := getValuesForStringFunc(targetsRaw)
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

		// Remove quote signs
		target = strings.Trim(target, `"`)
		target = strings.Trim(target, `'`)

		// Perform operation
		return "length", len(target), nil
	}

	err := fmt.Errorf("unknown function on %s", line)
	return "", "", err
}

// Gets elements required for a mathematical function
func getValuesForMathFunc(targetsRaw string) ([]int, error) {
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

	// Transform string to int
	target0, err := strconv.Atoi(targets[0])
	if err != nil {
		return nil, err
	}
	target1, err := strconv.Atoi(targets[1])
	if err != nil {
		return nil, err
	}

	return []int{target0, target1}, nil
}

// MathFunctions is a super set of all mathematical functions
func MathFunctions(line string) (interface{}, error) {
	// Power
	if strings.Contains(line, "power(") {
		// Get func value
		targetsRaw := strings.TrimSpace(line[6 : len(line)-1])

		// Get targets
		targets, err := getValuesForMathFunc(targetsRaw)
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
		targets, err := getValuesForMathFunc(targetsRaw)
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
		targets, err := getValuesForMathFunc(targetsRaw)
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

			// Confirm attribute is a boolean
			if attributes[val].Type != "boolean" {
				err := fmt.Errorf("logical functions can only have boolean attributes as parameters: %s", val)
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
