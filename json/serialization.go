package json

import (
	"fmt"
	"reflect"
	"strings"
)

/**
* Has all serialization methods from the JSON serializer.
* @author rnojiri
**/

// SerializeArray - serializes an array of jsons
func (j *Serializer) SerializeArray(items ...ArrayItem) (string, error) {

	numItems := len(items)
	if numItems == 0 {
		return "", nil
	}

	var err error
	var totalSize int
	jsons := make([]string, numItems)

	for i := 0; i < numItems; i++ {
		jsons[i], err = j.Serialize(items[i].Name, items[i].Parameters...)
		if err != nil {
			return "", err
		}
		totalSize += len(jsons[i])
	}

	var b strings.Builder
	b.Grow(totalSize + (numItems - 1) + 2)

	b.WriteString("[")

	for i := 0; i < numItems; i++ {

		b.WriteString(jsons[i])

		if i < numItems-1 {
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

		key, ok := m.variableMap[varName]
		if !ok {
			return "", fmt.Errorf(`variable "%s" does not exist`, varName)
		}

		if kind == reflect.Map {

			strMap, err := j.serializeMap(&value)
			if err != nil {
				return "", err
			}

			params[key] = strMap

		} else if kind == reflect.Array || kind == reflect.Slice {

			strArray, err := j.serializeArray(&value)
			if err != nil {
				return "", err
			}

			params[key] = strArray

		} else {

			if kind == reflect.String {

				str := value.String()

				var b strings.Builder
				b.Grow(len(str) + 2 + (strings.Count(str, `"`) * 2))

				b.WriteByte(doubleQuote)

				for _, c := range []byte(str) {
					if c == doubleQuote {
						b.WriteString(escapedDoubleQuote)
					} else {
						b.WriteByte(c)
					}
				}

				b.WriteByte(doubleQuote)

				params[key] = b.String()

			} else {

				params[key] = genericValue
			}
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
