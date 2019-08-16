package json

import (
	"fmt"
	"reflect"
	"strings"
)

// SerializeArray - serializes an array of jsons
func (j *Serializer) SerializeArray(parameters ...Parameters) (string, error) {

	numParameters := len(parameters)
	if numParameters == 0 {
		return "", nil
	}

	var err error
	var totalSize int
	jsons := make([]string, numParameters)

	for i := 0; i < numParameters; i++ {
		jsons[i], err = j.Serialize(parameters[i].Name, parameters[i].Parameters...)
		if err != nil {
			return "", err
		}
		totalSize += len(jsons[i])
	}

	var b strings.Builder
	b.Grow(totalSize + (numParameters - 1) + 2)

	b.WriteString("[")

	for i := 0; i < numParameters; i++ {

		b.WriteString(jsons[i])

		if i < numParameters-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String(), nil
}

// Serialize - serializes a mapped JSON
func (j *Serializer) Serialize(name string, parameters ...interface{}) (string, error) {

	m, ok := j.mapping[name]
	if !ok {
		return "", fmt.Errorf("no json mapping with name \"%s\"", name)
	}

	if m.numVariables != len(parameters)/2 {
		return "", fmt.Errorf("wrong number of variables")
	}

	params := make([]interface{}, m.numVariables)
	for i := 0; i < len(parameters); i += 2 {

		varName, ok := parameters[i].(string)
		if !ok {
			return "", fmt.Errorf("error casting variable index %d to string", i)
		}

		genericValue := parameters[i+1]
		value := reflect.ValueOf(genericValue)
		kind := value.Kind()

		if kind == reflect.Map {

			strMap, err := j.serializeMap(&value)
			if err != nil {
				return "", err
			}

			params[m.variableMap[varName]] = strMap

		} else if kind == reflect.Array || kind == reflect.Slice {

			strArray, err := j.serializeArray(&value)
			if err != nil {
				return "", err
			}

			params[m.variableMap[varName]] = strArray

		} else {

			params[m.variableMap[varName]] = genericValue
		}
	}

	return fmt.Sprintf(m.format, params...), nil
}

// serializeMap - serializes a map to JSON format
func (j *Serializer) serializeMap(value *reflect.Value) (string, error) {

	it := value.MapRange()

	hasNext := it.Next()

	var b strings.Builder

	for hasNext {

		key := it.Key().String()
		val := it.Value()

		strVal, err := j.getValueFromField(nil, &val)
		if err != nil {
			return "", err
		}

		j.writeStringValue(key, &b)
		b.WriteString(":")
		b.WriteString(strVal)

		hasNext = it.Next()
		if hasNext {
			b.WriteString(",")
		}
	}

	return b.String(), nil
}

// serializeArray - serializes an array to JSON format
func (j *Serializer) serializeArray(value *reflect.Value) (string, error) {

	arraySize := value.Len()

	var b strings.Builder

	for i := 0; i < arraySize; i++ {

		val := value.Index(i)

		strVal, err := j.getValueFromField(nil, &val)
		if err != nil {
			return "", err
		}

		b.WriteString(strVal)

		if i < arraySize-1 {
			b.WriteString(",")
		}
	}

	return b.String(), nil
}
