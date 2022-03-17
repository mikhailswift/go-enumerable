package enumerable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	orig := []string{"test", "123", "456", "longer string"}
	enum := FromSlice(orig).Filter(func(key int, val string) bool {
		return true
	})

	filtered := make([]string, 0)
	for {
		_, v, ok := enum.Next()
		if !ok {
			break
		}

		filtered = append(filtered, v)
	}

	assert.ElementsMatch(t, orig, filtered)

	expected := []string{"123", "456"}
	enum = FromSlice(orig).Filter(func(key int, val string) bool {
		return len(val) <= 3
	})

	filtered = enum.Values()
	assert.ElementsMatch(t, expected, filtered)
}

func TestMap(t *testing.T) {
	orig := map[string]int{
		"test":  1,
		"test2": 2,
		"test3": 4,
	}

	enum := FromMap(orig)
	mappedEnum := Map(enum, func(k string, v int) (int, string) {
		return v, k
	})

	expected := map[int]string{
		1: "test",
		2: "test2",
		4: "test3",
	}

	actual := make(map[int]string)
	for {
		k, v, ok := mappedEnum.Next()
		if !ok {
			break
		}

		actual[k] = v
	}

	assert.Equal(t, expected, actual)
}

func TestMapValues(t *testing.T) {
	orig := []int{1, 2, 4, 8, 16}
	expected := []int{2, 4, 8, 16, 32}
	actual := MapValues(FromSlice(orig), func(v int) int {
		return 2 * v
	}).Values()
	assert.ElementsMatch(t, expected, actual)

	expected = []int{16, 32}
	actual = MapValues(FromSlice(orig), func(v int) int {
		return 2 * v
	}).Filter(func(_ int, v int) bool {
		return v > 10
	}).Values()
	assert.ElementsMatch(t, expected, actual)
}
