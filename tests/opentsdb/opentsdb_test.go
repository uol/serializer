package opentsdb

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	serializer "github.com/uol/serializer/opentsdb"
)

// genRandom - generates a random number
func genRandom(min, max int) int {

	return (rand.Intn(max-min+1) + min)
}

// createSerializer - creates a new serializer
func createSerializer() *serializer.Serializer {

	rand.Seed(time.Now().UnixNano())

	return serializer.New(100)
}

// serialize - try to serialize the named type
func serialize(t *testing.T, s *serializer.Serializer, item serializer.ArrayItem) string {

	result, err := s.Serialize(item.Metric, item.Timestamp, item.Value, item.Tags...)
	if !assert.NoError(t, err, "error serializing line") {
		panic(err)
	}

	return result
}

// serializeArray - try to serialize the named type
func serializeArray(t *testing.T, s *serializer.Serializer, items []serializer.ArrayItem) string {

	result, err := s.SerializeArray(items...)
	if !assert.NoError(t, err, "error serializing array") {
		panic(err)
	}

	return result
}

// TestSingleLineStringTags - tests a single line, string tags
func TestSingleLineStringTags(t *testing.T) {

	s := createSerializer()

	line := serializer.ArrayItem{
		Metric:    "single",
		Timestamp: time.Now().Unix(),
		Value:     float64(genRandom(1, 100)),
		Tags: []interface{}{
			"host", "localhost",
			"ttl", "1",
		},
	}

	result := serialize(t, s, line)
	expected := fmt.Sprintf("put %s %d %.0f %s=%s %s=%s\n", line.Metric, line.Timestamp, line.Value, line.Tags[0], line.Tags[1], line.Tags[2], line.Tags[3])

	assert.Equal(t, expected, result, "expected same string")
}

// TestSingleLineMixedTypeTags - tests a single line, mixed type tags
func TestSingleLineMixedTypeTags(t *testing.T) {

	s := createSerializer()

	line := serializer.ArrayItem{
		Metric:    "single",
		Timestamp: time.Now().Unix(),
		Value:     float64(genRandom(10, 100)) + 0.5,
		Tags: []interface{}{
			"host", "localhost",
			"ttl", 1,
			"float", float64(genRandom(300, 1000)) + 0.1,
			"boolean", genRandom(0, 10) >= 5,
		},
	}

	result := serialize(t, s, line)
	expected := fmt.Sprintf("put %s %d %.1f %s=%s %s=%d %s=%.1f %s=%t\n", line.Metric, line.Timestamp, line.Value, line.Tags[0], line.Tags[1], line.Tags[2], line.Tags[3], line.Tags[4], line.Tags[5], line.Tags[6], line.Tags[7])

	assert.Equal(t, expected, result, "expected same string")
}

// TestMultiLineStringTags - tests a multi line, string tags
func TestMultiLineStringTags(t *testing.T) {

	s := createSerializer()

	const size = 21
	format := ""
	lines := make([]serializer.ArrayItem, size)
	args := []interface{}{}

	for i := 0; i < size; i++ {

		lines[i] = serializer.ArrayItem{
			Metric:    "multi" + strconv.Itoa(i),
			Timestamp: time.Now().Unix(),
			Value:     float64(i),
			Tags: []interface{}{
				"host", "host" + strconv.Itoa(i),
				"index", strconv.Itoa(i),
			},
		}

		args = append(args, lines[i].Metric, lines[i].Timestamp, lines[i].Value, lines[i].Tags[0], lines[i].Tags[1], lines[i].Tags[2], lines[i].Tags[3])

		format += "put %s %d %.0f %s=%s %s=%s\n"
	}

	result := serializeArray(t, s, lines)
	expected := fmt.Sprintf(format, args...)

	assert.Equal(t, expected, result, "expected same string")
}

// TestMultiLineMixedTypeTags - tests a multi line, mixed type tags
func TestMultiLineMixedTypeTags(t *testing.T) {

	s := createSerializer()

	const size = 2
	format := ""
	lines := make([]serializer.ArrayItem, size)
	args := []interface{}{}

	for i := 0; i < size; i++ {

		lines[i] = serializer.ArrayItem{
			Metric:    "multi" + strconv.Itoa(i),
			Timestamp: time.Now().Unix(),
			Value:     float64(i),
			Tags: []interface{}{
				"host", "host" + strconv.Itoa(i),
				"index", i,
				"float", float64(genRandom(300, 1000)) + 0.1,
				"boolean", genRandom(0, 10) >= 5,
			},
		}

		args = append(args, lines[i].Metric, lines[i].Timestamp, lines[i].Value, lines[i].Tags[0], lines[i].Tags[1], lines[i].Tags[2], lines[i].Tags[3], lines[i].Tags[4], lines[i].Tags[5], lines[i].Tags[6], lines[i].Tags[7])

		format += "put %s %d %.0f %s=%s %s=%d %s=%.1f %s=%t\n"
	}

	result := serializeArray(t, s, lines)
	expected := fmt.Sprintf(format, args...)

	assert.Equal(t, expected, result, "expected same string")
}
