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
}
