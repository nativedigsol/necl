package necl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicNECLFileParser(t *testing.T) {
	file := ParseNECLFile("./test_data/example-1-simple-file.necl")

	// Assert no-block attributes
	assert.EqualValues(t, "example", file.Attributes["name"].Value)
	assert.EqualValues(t, 3.1414999961853027, file.Attributes["pi"].Value)
	assert.EqualValues(t, false, file.Attributes["no"].Value)
	assert.EqualValues(t, "this is a multiline string", file.Attributes["multiline"].Value)

	// Assert block attributes
	assert.EqualValues(t, file.Blocks["block"].Attributes["foo"].Value, "bar")

	// Assert array values
	longArray := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, true, false}
	assert.EqualValues(t, []interface{}{"test", 1}, file.Attributes["test_array"].Array)
	assert.EqualValues(t, []interface{}{"test", "block", "array", 1234, false}, file.Blocks["block"].Attributes["block_array"].Array)
	assert.EqualValues(t, longArray, file.Attributes["long_array"].Array)
	assert.EqualValues(t, "this is a blocked multiline string", file.Blocks["block"].Attributes["block_multiline"].Value)
}

func TestK8sNECLFileParser(t *testing.T) {
	file := ParseNECLFile("./test_data/example-2-kubernetes-deployment.necl")

	// Global attributes
	assert.EqualValues(t, "apps/v1", file.Attributes["apiVersion"].Value)
	assert.EqualValues(t, "Deployment", file.Attributes["kind"].Value)

	// Metadata block
	assert.EqualValues(t, "nginx-deployment", file.Blocks["metadata"].Attributes["name"].Value)
	assert.EqualValues(t, "nginx", file.Blocks["metadata"].Blocks["labels"].Attributes["app"].Value)

	// Spec block
	assert.EqualValues(t, 3, file.Blocks["spec"].Attributes["replicas"].Value)
	assert.EqualValues(t, "nginx", file.Blocks["spec"].Blocks["selector"].Blocks["matchLabels"].Attributes["app"].Value)
	assert.EqualValues(t, "nginx", file.Blocks["spec"].Blocks["template"].Blocks["metadata"].Blocks["labels"].Attributes["app"].Value)
	assert.EqualValues(t, "nginx:1.14.2", file.Blocks["spec"].Blocks["template"].Blocks["spec"].Blocks["containers"].Blocks["nginx"].Attributes["image"].Value)
	assert.EqualValues(t, 80, file.Blocks["spec"].Blocks["template"].Blocks["spec"].Blocks["containers"].Blocks["nginx"].Blocks["ports"].Attributes["containerPort"].Value)
}
