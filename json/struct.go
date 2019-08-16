package json

type variableType uint8

const (
	normalValue      variableType = 0
	propertyVariable variableType = 1
)

// mappedJSON - internal mapped JSON struct
type mappedJSON struct {
	format       string
	formatSize   int
	variableMap  map[string]int
	numVariables int
}

// Serializer - the json serializer
type Serializer struct {
	bufferSize int
	mapping    map[string]*mappedJSON
}

// Parameters - a configuration to render a json
type Parameters struct {
	Name       string
	Parameters []interface{}
}

// New - creates a new JSON serializer
func New(bufferSize int) *Serializer {

	return &Serializer{
		bufferSize: bufferSize,
		mapping:    map[string]*mappedJSON{},
	}
}
