package necl

type File struct {
	Attributes map[string]Attribute
	Blocks     map[string]Block
	RawText    []string
}

type Block struct {
	Name       string
	Attributes map[string]Attribute
	Blocks     []Block
	RawText    []string
}

type Attribute struct {
	Name  string
	Type  string
	Value interface{}
	Array []interface{}
}
