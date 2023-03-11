package necl

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

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
		if found && (newAttr != Attribute{}) {
			attributes[newAttr.Name] = Attribute{
				Name:  newAttr.Name,
				Type:  newAttr.Type,
				Value: newAttr.Value,
			}
		}

		// Look for the end of a block
		if strings.Contains(line, "}") && blockStarted {
			blockStarted = false
		}
	}

	return attributes, nil
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

	// Discover attribute value
	fullAttribute := strings.TrimSpace(line[i+1:])

	// Discover attribute type based on the value
	var attributeValue interface{}
	attributeType, err := discoverAttributeType(fullAttribute)
	Check(err)
	switch attributeType {
	case "string":
		attributeValue = fullAttribute[1 : len(fullAttribute)-1]
	case "number":
		attributeValue, _ = strconv.ParseFloat(fullAttribute, 32)
	case "boolean":
		attributeValue, err = strconv.ParseBool(fullAttribute)
		Check(err)
	case "array":
		//
		//
		// READ ARRAY
		//
		//
		return true, Attribute{}, err
	default:
		return true, Attribute{}, err
	}

	return true, Attribute{
		Name:  attributeName,
		Type:  attributeType,
		Value: attributeValue,
	}, nil
}

// findAttributesInsideBlock looks for attributes definitions in a block
func findAttributes(data []string) (map[string]Attribute, error) {
	attributes := make(map[string]Attribute)

	for _, line := range data {
		found, newAttr, err := findAttribute(line)
		Check(err)

		if found && (newAttr != Attribute{}) {
			attributes[newAttr.Name] = Attribute{
				Name:  newAttr.Name,
				Type:  newAttr.Type,
				Value: newAttr.Value,
			}
		}
	}

	return attributes, nil
}

// Discover the type of an attribute based on the NECL spec
func discoverAttributeType(value string) (string, error) {
	for _, cRune := range value {
		// Transform rune to string
		c := string(cRune)

		// String
		// One quote
		if strings.HasPrefix(c, `'`) && strings.HasSuffix(c, `'`) {
			return "string", nil
		}

		// Double quote
		if strings.HasPrefix(c, `"`) && strings.HasSuffix(c, `"`) {
			return "string", nil
		}

		// Array
		if strings.HasPrefix(c, `[`) && strings.HasSuffix(c, `]`) {
			return "string", nil
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
	}

	// No valid type was found
	err := errors.New("no valid type was found")
	return "", err
}
