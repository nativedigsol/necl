package necl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// findBlock looks for a block by searching for the beggining '{' and the closing '}'
func findBlock(data []string, startLine int, currentAttributes map[string]Attribute) (Block, int, []int, int, []int, error) {
	var start int
	var end int

	// Declare this as false since blocks can have nested blocks
	blockStarted := false

	// Nested blocks
	nestedBlocks := make(map[string]Block)
	nestedBlockStarts := []int{}
	nestedBlockEnds := []int{}

	// Try to find where the block starts and where it ends
	for i, line := range data {
		// Only start looking when the first byte index is reached
		if i >= startLine {
			// Nested block
			if strings.Contains(line, "{") && blockStarted {
				// Check if block is nested from another block
				nested := false
				for _, val := range nestedBlockStarts {
					if i == val {
						nested = true
					}
				}
				if !nested {
					// This is not the prettiest of stuff and it could definitely see an improvement
					// but for now, it works
					newNestedBlock, start, starts, end, ends, err := findBlock(data, i, currentAttributes)
					if err != nil {
						return Block{}, 0, nil, 0, nil, err
					}
					nestedBlockStarts = append(nestedBlockStarts, start)
					nestedBlockStarts = append(nestedBlockStarts, starts...)
					nestedBlockEnds = append(nestedBlockEnds, end)
					nestedBlockEnds = append(nestedBlockEnds, ends...)
					nestedBlocks[newNestedBlock.Name] = newNestedBlock
				}
			}

			// Look for the start of the block
			if strings.Contains(line, "{") && !blockStarted {
				start = i
				blockStarted = true
			}

			// Look for the end of the block
			if strings.Contains(line, "}") && blockStarted {
				nested := false
				for _, val := range nestedBlockEnds {
					if i == val {
						nested = true
					}
				}
				if !nested {
					end = i
					break
				}
			}
		}
	}

	// If a block was found
	if blockStarted {
		// Get block name
		blockNameRaw := strings.TrimSpace(data[start])
		blockName := strings.TrimSpace(blockNameRaw[:len(blockNameRaw)-1])

		// Get block attributes (if any)
		blockAttributes, err := findAttributes(data[start:end], currentAttributes)
		if err != nil {
			return Block{}, 0, nil, 0, nil, err
		}

		// Knowing where the block starts and ends, a Block struct can be created
		return Block{
			Name:       blockName,
			Attributes: blockAttributes,
			Blocks:     nestedBlocks,
		}, start, nestedBlockStarts, end, nestedBlockEnds, nil
	}

	// No block was found
	return Block{}, 0, nil, 0, nil, nil
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

		// Look for the end of a block
		if strings.Contains(line, "}") && blockStarted {
			blockStarted = false
		}

		// Skip looking for attributes if inside a block
		if blockStarted {
			continue
		}

		// If not inside a block, look for attributes
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

// This reads a file as an array of bytes
func readFile(filename string) ([]string, error) {
	trimmedFilename := filename
	if strings.HasPrefix(trimmedFilename, `"`) {
		trimmedFilename = strings.Trim(trimmedFilename, `"`)
	} else if strings.HasPrefix(trimmedFilename, `'`) {
		trimmedFilename = strings.Trim(trimmedFilename, `'`)
	}

	if !strings.HasSuffix(trimmedFilename, ".necl") {
		err := fmt.Errorf("file %s is not in .necl format", trimmedFilename)
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var rawData []string

	for scanner.Scan() {
		_, token, err := bufio.ScanLines(scanner.Bytes(), true)
		if err != nil {
			return nil, err
		}
		rawData = append(rawData, string(token))
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// ParseNECLFile will read and parse a ".necl" file
func ParseNECLFile(filename string) (*File, error) {
	// Read file
	rawText, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	// Remove all comments from the text
	for i, line := range rawText {
		// Break if last line was reached
		// This is needed because this loop should run for the entire length of the text
		// But if comments are being removed, the length of the text is going to be dinamically reduced
		// This is set so no index out of bound errors happen
		if i > len(rawText) {
			break
		}
		// Line comment
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			if i == len(rawText)-1 {
				rawText = rawText[:i-1]
			} else {
				rawText = append(rawText[:i], rawText[i+1:]...)
			}
		} else if strings.Contains(strings.TrimSpace(line), "//") { // Inline comment
			// Find index of where the comments start
			firstSlash := 0
			foundComment := false
			for j, c := range line {
				if string(c) == "/" && firstSlash == 0 && !foundComment {
					firstSlash = j
				}
				if j == firstSlash && string(c) != "/" && !foundComment {
					firstSlash = 0
				}
				if string(c) == "/" && firstSlash != 0 && !foundComment {
					foundComment = true
				}
			}
			rawText[i] = line[:firstSlash-1]
		}
	}

	// Find attributes that are not inside blocks
	attributes, err := findAttributesNoBlock(rawText)
	if err != nil {
		return nil, err
	}

	// Find blocks
	blocks := make(map[string]Block)
	startLine := 0
	for startLine < len(rawText) {
		newBlock, _, _, endLine, _, err := findBlock(rawText, startLine, attributes)
		if err != nil {
			return nil, err
		}

		if newBlock.Name != "" {
			// Add new block to the array of blocks
			blocks[newBlock.Name] = newBlock
		}

		if endLine == 0 {
			startLine += 1
		} else {
			startLine = endLine - 1
		}
	}

	return &File{
		Attributes: attributes,
		Blocks:     blocks,
	}, nil
}
