package necl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicNECLFileParser(t *testing.T) {
	file, err := ParseNECLFile("./test_data/example-1-test-file.necl")
	assert.NoError(t, err)

	// Assert no-block attributes
	assert.EqualValues(t, "example", file.Attributes["name"].Value)
	assert.EqualValues(t, 3.1414999961853027, file.Attributes["pi"].Value)
	assert.EqualValues(t, false, file.Attributes["no"].Value)
	assert.EqualValues(t, "this is a multiline string", file.Attributes["multiline"].Value)

	// Assert comparison operators
	compArray := []interface{}{"this", "array", "can", "compare", "stuff", true}
	assert.EqualValues(t, true, file.Attributes["c1"].Value)
	assert.EqualValues(t, false, file.Attributes["c2"].Value)
	assert.EqualValues(t, false, file.Attributes["c3"].Value)
	assert.EqualValues(t, true, file.Attributes["c4"].Value)
	assert.EqualValues(t, false, file.Attributes["c5"].Value)
	assert.EqualValues(t, true, file.Attributes["c6"].Value)
	assert.EqualValues(t, compArray, file.Attributes["comp_array"].Array)

	// Assert arithmetic operations
	opArray := []interface{}{"this", "array", "can", "calculate", "stuff", 5}
	assert.EqualValues(t, 2, file.Attributes["sum"].Value)
	assert.EqualValues(t, 3, file.Attributes["subtract"].Value)
	assert.EqualValues(t, 25, file.Attributes["multiply"].Value)
	assert.EqualValues(t, 10, file.Attributes["divide"].Value)
	assert.EqualValues(t, 4, file.Attributes["attOp1"].Value)
	assert.EqualValues(t, 8, file.Attributes["attOp2"].Value)
	assert.EqualValues(t, opArray, file.Attributes["op_array"].Array)

	// Assert block attributes
	assert.EqualValues(t, "bar", file.Blocks["block"].Attributes["foo"].Value)
	assert.EqualValues(t, false, file.Blocks["block"].Attributes["cb1"].Value)

	// Assert array values
	longArray := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, true, false}
	multilineArray := []interface{}{"this", "is", "a", "multiline", "array", 1, false}
	assert.EqualValues(t, []interface{}{"test", 1}, file.Attributes["test_array"].Array)
	assert.EqualValues(t, []interface{}{"test", "block", "array", 1234, false}, file.Blocks["block"].Attributes["block_array"].Array)
	assert.EqualValues(t, multilineArray, file.Attributes["multiArray"].Array)
	assert.EqualValues(t, longArray, file.Attributes["long_array"].Array)
	assert.EqualValues(t, "this is a blocked multiline string", file.Blocks["block"].Attributes["block_multiline"].Value)
}

func TestK8sNECLFileParser(t *testing.T) {
	file, err := ParseNECLFile("./test_data/example-2-kubernetes-deployment.necl")
	assert.NoError(t, err)

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
