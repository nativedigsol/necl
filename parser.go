package necl

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Discover the type of an attribute based on the NECL spec
func discoverAttributeType(value string) (string, error) {
	// String
	// One quote
	if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
		return "string", nil
	}

	// Double quote
	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		return "string", nil
	}

	// Array
	if strings.HasPrefix(value, `[`) && strings.HasSuffix(value, `]`) {
		return "array", nil
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
	err = errors.New("no valid type was found")
	return "", err
}

// parseArrayAttributes parse all attributes in an array
func parseArrayAttributes(array string) ([]interface{}, error) {
	// Remove array "[]"
	arrayRaw := strings.TrimSpace(array[1 : len(array)-1])

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
	var attributeValue interface{}
	var arrayElements []interface{}
	for _, item := range arrayElementsRaw {
		attributeType, err := discoverAttributeType(item)
		Check(err)

		switch attributeType {
		case "string":
			attributeValue = item[1 : len(item)-1]
			arrayElements = append(arrayElements, attributeValue)
		case "number":
			if strings.Contains(item, ".") || strings.Contains(item, ",") {
				attributeValue, _ = strconv.ParseFloat(item, 32)
			} else {
				attributeValue, _ = strconv.Atoi(item)
			}
			arrayElements = append(arrayElements, attributeValue)
		case "boolean":
			attributeValue, err = strconv.ParseBool(item)
			arrayElements = append(arrayElements, attributeValue)
			Check(err)
		case "array":
			err := errors.New("an attribute with array type can't have nested arrays")
			return []interface{}{}, err
		default:
			return []interface{}{}, err
		}
	}

	return arrayElements, nil
}

// findAttribute looks for an attribute in a single line
func findAttribute(line string) (bool, Attribute, error) {
	// Look for '='
	if !strings.Contains(line, "=") {
		return false, Attribute{}, nil
	}

	// Ignore if line is a comment
	if strings.HasPrefix(line, "//") {
		return false, Attribute{}, nil
	}

	// Find position of the '='
	i := strings.Index(line, "=")

	// Get attribute name
	attributeName := strings.TrimSpace(line[:i])

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
	fullAttribute := strings.TrimSpace(line[i+1:])

	// Discover attribute type based on the value
	var arrayValues []interface{}
	var attributeValue interface{}
	attributeType, err := discoverAttributeType(fullAttribute)
	Check(err)
	switch attributeType {
	case "string":
		attributeValue = fullAttribute[1 : len(fullAttribute)-1]
	case "number":
		if strings.Contains(fullAttribute, ".") || strings.Contains(fullAttribute, ",") {
			attributeValue, _ = strconv.ParseFloat(fullAttribute, 32)
		} else {
			attributeValue, _ = strconv.Atoi(fullAttribute)
		}
	case "boolean":
		attributeValue, err = strconv.ParseBool(fullAttribute)
		Check(err)
	case "array":
		arrayValues, err = parseArrayAttributes(fullAttribute)
		Check(err)
	default:
		return true, Attribute{}, err
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

	for _, line := range data {
		found, newAttr, err := findAttribute(line)
		Check(err)

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

// findBlock looks for a block by searching for the beggining '{' and the closing '}'
func findBlock(data []string, startLine int) (Block, int) {
	var start int
	var end int

	// Declare this as false since blocks can have nested blocks
	blockStarted := false

	// Try to find where the block starts and where it ends
	for i, line := range data {
		// Only start looking when the first byte index is reached
		if i >= startLine {
			// Look for the start of the block
			if strings.Contains(line, "{") && !blockStarted {
				start = i
				blockStarted = true
			}

			// Look for the end of the block
			if strings.Contains(line, "}") && blockStarted {
				end = i
			}
		}
	}

	// If a block was found
	if blockStarted {
		// Get block name
		blockNameRaw := strings.TrimSpace(data[start])
		blockName := strings.TrimSpace(blockNameRaw[:len(blockNameRaw)-1])

		// Get block attributes (if any)
		blockAttributes, err := findAttributes(data[start:end])
		Check(err)

		// Knowing where the block starts and ends, a Block struct can be created
		return Block{
			Name:       blockName,
			RawText:    data[start:end],
			Attributes: blockAttributes,
		}, end
	}

	// No block was found
	return Block{}, 0
}

// findAttributesNoBlock looks for attributes that are outside blocks
func findAttributesNoBlock(data []string) (map[string]Attribute, error) {
	// Declare this to know when line is inside a block
	blockStarted := false

	// Declare an empty map of attributes
	attributes := make(map[string]Attribute)

	// Try to find if there's any block
	for _, line := range data {
		// Look for the start of a block
		if strings.Contains(line, "{") && !blockStarted {
			blockStarted = true
		}

		// If not inside a block, look for attributes
		found, newAttr, err := findAttribute(line)
		Check(err)
		if found && (newAttr.Name != "") {
			attributes[newAttr.Name] = Attribute{
				Name:  newAttr.Name,
				Type:  newAttr.Type,
				Value: newAttr.Value,
				Array: newAttr.Array,
			}
		}

		// Look for the end of a block
		if strings.Contains(line, "}") && blockStarted {
			blockStarted = false
		}
	}

	return attributes, nil
}

// This reads a file as an array of bytes
func readFile(filename string) []string {
	file, err := os.Open(filename)
	Check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var rawData []string

	for scanner.Scan() {
		_, token, err := bufio.ScanLines(scanner.Bytes(), true)
		Check(err)
		rawData = append(rawData, string(token))
	}
	err = scanner.Err()
	Check(err)

	return rawData
}

// ParseNECLFile will read and parse a ".necl" file
func ParseNECLFile(filename string) *File {
	// Read file
	rawText := readFile(filename)

	// Find blocks
	blocks := make(map[string]Block)
	startLine := 0
	for startLine < len(rawText) {
		newBlock, endLine := findBlock(rawText, startLine)
		if newBlock.Name != "" {
			// Add new block to the array of blocks
			blocks[newBlock.Name] = newBlock
		}

		if endLine == 0 {
			startLine += 1
		} else {
			startLine = endLine
		}
	}

	// Find attributes that are not inside blocks
	attributes, err := findAttributesNoBlock(rawText)
	Check(err)

	return &File{
		Attributes: attributes,
		Blocks:     blocks,
		RawText:    rawText,
	}
}
