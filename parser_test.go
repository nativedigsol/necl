package necl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNECLFileParser(t *testing.T) {
	file := ParseNECLFile("./test_data/example-1-simple-file.necl")

	// Assert no-block attributes
	assert.EqualValues(t, file.Attributes["name"].Value, "example")
	assert.EqualValues(t, file.Attributes["pi"].Value, 3.1414999961853027)
	assert.EqualValues(t, file.Attributes["no"].Value, false)

	// Assert block attributes
	assert.EqualValues(t, file.Blocks["block"].Attributes["foo"].Value, "bar")

	// Assert array values
	longArray := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, true, false}
	assert.EqualValues(t, []interface{}{"test", 1}, file.Attributes["test_array"].Array)
	assert.EqualValues(t, []interface{}{"test", "block", "array", 1234, false}, file.Blocks["block"].Attributes["block_array"].Array)
	assert.EqualValues(t, longArray, file.Attributes["long_array"].Array)
}
