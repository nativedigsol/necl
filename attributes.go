package necl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Discover the type of an attribute based on the NECL spec
func discoverAttributeType(value string) (string, error) {
	// String
	// Multiline strings are not checked here, it is checked by findAttribute beforehand and "compiled" by getMultilineString
	if (strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) || (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) {
		return "string", nil
	}

	// Array
	if strings.HasPrefix(value, `[`) && strings.HasSuffix(value, `]`) {
		return "array", nil
	}

	// Comparison (will transform into a boolean by the end)
	if ContainsMany(value, []string{"==", "!=", "<", "<=", ">", ">="}) {
		return "comparison", nil
	}

	// Arithmetic operation
	if ContainsMany(value, []string{"+", "-", "*", "/"}) {
		return "arithmetic", nil
	}

	// Functions
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

	// Boolean
	if strings.Contains(strings.ToLower(value), "true") || strings.Contains(strings.ToLower(value), "false") {
		return "boolean", nil
	}

	// Number
	_, err := strconv.ParseFloat(value, 32)
	if err == nil {
		return "number", nil
	}

	// No valid type was found
	err = fmt.Errorf("no valid type was found for attribute %s", value)
	return "", err
}

// getMultilineString looks if a certain line is an multiline string
func getMultilineString(data []string, startLine int) interface{} {
	// Get all lines of the multistring
	var stringLines []string
	for i, line := range data {
		if i >= startLine {
			// First line, so need to remove the variable name
			if i == startLine {
				j := strings.Index(line, "=")
				line = strings.TrimSpace(line[j+1:])
			}

			// Add lines to the array (already trimmed)
			removeBackslash := strings.TrimSuffix(strings.TrimSpace(line), `\`)
			removeWhitespace := (strings.TrimSpace(removeBackslash))
			newLine := removeWhitespace[1 : len(removeWhitespace)-1]
			stringLines = append(stringLines, newLine)

			// If the line doesn't have a `\`, it means that it is the last line of the multiline string
			if !strings.HasSuffix(strings.TrimSpace(line), `\`) {
				break
			}
		}
	}

	// join the stringLines
	return strings.Join(stringLines, " ")
}

// getAttribute transforms a string in a interface{} with the correct type
func getAttribute(attributeValueRaw string, isArray bool, currentAttributes map[string]Attribute) (string, interface{}, []interface{}, error) {
	var attributeValue interface{}
	var arrayElements []interface{}

	attributeType, err := discoverAttributeType(attributeValueRaw)
	if err != nil {
		return "", nil, nil, err
	}

	switch attributeType {
	case "string":
		attributeValue = attributeValueRaw[1 : len(attributeValueRaw)-1]
	case "number":
		if strings.Contains(attributeValueRaw, ".") || strings.Contains(attributeValueRaw, ",") {
			attributeValue, _ = strconv.ParseFloat(attributeValueRaw, 32)
		} else {
			attributeValue, _ = strconv.Atoi(attributeValueRaw)
		}
	case "boolean":
		attributeValue, _ = strconv.ParseBool(attributeValueRaw)
	case "comparison":
		attributeType = "boolean"
		attributeValue, err = PerformComparison(attributeValueRaw, currentAttributes)
		if err != nil {
			return "", nil, nil, err
		}
	case "arithmetic":
		attributeType = "number"
		attributeValue, err = PerformArithmeticOperation(attributeValueRaw, currentAttributes)
		if err != nil {
			return "", nil, nil, err
		}
	case "func-string":
		var functionDone string
		functionDone, attributeValue, err = StringFunctions(attributeValueRaw)
		if err != nil {
			return "", nil, nil, err
		}
		if functionDone == "contains" {
			attributeType = "boolean"
		} else if functionDone == "length" {
			attributeType = "number"
		} else {
			attributeType = "string"
		}
	case "func-math":
		attributeType = "number"
		attributeValue, err = MathFunctions(attributeValueRaw)
		if err != nil {
			return "", nil, nil, err
		}
	case "func-logic":
		attributeType = "boolean"
		attributeValue, err = LogicFunctions(attributeValueRaw, currentAttributes)
		if err != nil {
			return "", nil, nil, err
		}
	case "array":
		if isArray {
			err := errors.New("an attribute with array type can't have nested arrays")
			return "", nil, nil, err
		} else {
			var err error
			arrayElements, err = parseArrayAttributes(attributeValueRaw, currentAttributes)
			if err != nil {
				return "", nil, nil, err
			}
		}
	default:
		err := fmt.Errorf("unknown attribute type: %s", attributeType)
		return "", nil, nil, err
	}

	return attributeType, attributeValue, arrayElements, nil
}

// parseArrayAttributes parse all attributes in an array
func parseArrayAttributes(array string, currentAttributes map[string]Attribute) ([]interface{}, error) {
	// Remove array "[]"
	// The "if" is there so it doesn't remove actual values of a multiline array
	arrayRaw := array
	if strings.HasPrefix(strings.TrimSpace(arrayRaw), "[") && strings.HasSuffix(strings.TrimSpace(arrayRaw), "]") {
		arrayRaw = strings.TrimSpace(array[1 : len(array)-1])
	}

	// Find all elements
	// Elements are separated by a comma ","
	var arrayElementsRaw []string
	previousComma := 0
	for commaIndex, findComma := range arrayRaw {
		if string(findComma) == "," {
			newElement := strings.TrimSpace(arrayRaw[previousComma:commaIndex])
			arrayElementsRaw = append(arrayElementsRaw, newElement)
			previousComma = int(commaIndex) + 1
		}

		// Last element
		if commaIndex == len(arrayRaw)-1 && arrayElementsRaw != nil {
			newElement := strings.TrimSpace(arrayRaw[previousComma+1:])
			arrayElementsRaw = append(arrayElementsRaw, newElement)
		}
	}

	// If no commas, check if the array only has one value
	if arrayRaw != "" && arrayElementsRaw == nil {
		arrayElementsRaw = append(arrayElementsRaw, arrayRaw)
	}

	// Discover types of all elements
	var arrayElements []interface{}
	for _, item := range arrayElementsRaw {
		_, attributeValue, _, err := getAttribute(item, true, currentAttributes)
		if err != nil {
			return nil, err
		}

		arrayElements = append(arrayElements, attributeValue)
	}

	return arrayElements, nil
}

// isMultilineArray checks if a certain attribute is a multiline array
func isMultilineArray(data []string, startLine int) bool {
	// Check if is attribute
	if !strings.Contains(data[startLine], "=") {
		return false
	}

	// Remove attribute name
	i := strings.Index(data[startLine], "=")
	rawAttribute := strings.TrimSpace(data[startLine][i+1:])

	// Check if it's a multiline array
	if strings.HasPrefix(rawAttribute, `[`) && !strings.HasSuffix(rawAttribute, `]`) {
		return true
	}

	return false
}

// getMultilineArray parse all attributes in an multiline array
func getMultilineArray(data []string, startLine int, currentAttributes map[string]Attribute) ([]interface{}, error) {
	// Loop through data to find the end of the array and all the elements
	arrayElementsRaw := []string{}
	endOfArray := false
	for i, line := range data {
		if i >= startLine && !endOfArray {
			// Check end of array
			if strings.Contains(strings.TrimSpace(line), "]") {
				endOfArray = true
			}

			// First line of the array, need to remove the attribute name
			if i == startLine {
				j := strings.Index(data[startLine], "=")
				line = strings.TrimSpace(data[startLine][j+1:])
			}

			// Add all elements of the multiline array to a simple string
			arrayElementsRaw = append(arrayElementsRaw, strings.TrimSpace(line))
		}
	}
	// Transform the array of elements into a string
	stringElements := strings.Join(arrayElementsRaw, " ")

	// Parse elements of the array
	arrayElements, err := parseArrayAttributes(stringElements, currentAttributes)
	if err != nil {
		return nil, err
	}

	return arrayElements, nil
}

// findAttribute looks for an attribute in a single line
func findAttribute(data []string, line int, currentAttributes map[string]Attribute) (bool, Attribute, error) {
	// Look for '='
	if !strings.Contains(data[line], "=") {
		return false, Attribute{}, nil
	}

	// Find position of the '='
	i := strings.Index(data[line], "=")

	// Get attribute name
	attributeName := strings.TrimSpace(data[line][:i])

	// Name cannot be empty
	if attributeName == "" {
		err := errors.New("attribute name cannot be empty")
		return true, Attribute{
			Name:  "",
			Type:  "",
			Value: "",
			Array: []interface{}{},
		}, err
	}

	// Discover attribute value
	rawAttribute := strings.TrimSpace(data[line][i+1:])
	var attributeValue interface{}
	var attributeType string
	arrayValues := []interface{}{}
	var err error

	// If attribute is a multiline string, get the full string
	if strings.HasSuffix(rawAttribute, `\`) {
		attributeValue = getMultilineString(data, line)
		attributeType = "string"
	} else if isMultilineArray(data, line) { // Multiline array
		arrayValues, err = getMultilineArray(data, line, currentAttributes)
		if err != nil {
			return false, Attribute{}, err
		}
	} else {
		attributeType, attributeValue, arrayValues, err = getAttribute(rawAttribute, false, currentAttributes)
		if err != nil {
			return false, Attribute{}, err
		}
	}

	return true, Attribute{
		Name:  attributeName,
		Type:  attributeType,
		Value: attributeValue,
		Array: arrayValues,
	}, nil
}

// findAttributesInsideBlock looks for attributes definitions in a block
func findAttributes(data []string) (map[string]Attribute, error) {
	attributes := make(map[string]Attribute)

	for i := range data {
		found, newAttr, err := findAttribute(data, i, attributes)
		if err != nil {
			return nil, err
		}

		if found && (newAttr.Name != "") {
			attributes[newAttr.Name] = Attribute{
				Name:  newAttr.Name,
				Type:  newAttr.Type,
				Value: newAttr.Value,
				Array: newAttr.Array,
			}
		}
	}

	return attributes, nil
}
