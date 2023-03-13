package necl

import (
	"bufio"
	"os"
	"strings"
)

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
	for i, line := range data {
		// Look for the start of a block
		if strings.Contains(line, "{") && !blockStarted {
			blockStarted = true
		}

		// If not inside a block, look for attributes
		found, newAttr, err := findAttribute(data, i)
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
