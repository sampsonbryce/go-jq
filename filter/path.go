package filter

import (
	"fmt"
	"strconv"
	"strings"
)

type PathAccessorError struct {
	Accessor string
}

func (e PathAccessorError) Error() string {
	return fmt.Sprintf("accessor: unable to access json using accessor '%v'", e.Accessor)
}

type PathFilterError struct {
	Path string
}

func (e PathFilterError) Error() string {
	return fmt.Sprintf("path filter: unable to access json using path '%v'", e.Path)
}

type PathError struct {
	Path string
	Err  error
}

func (e PathError) Error() string {
	return fmt.Sprintf("path: invalid path '%v'", e.Path)
}

func (e PathError) Unwrap() error {
	return e.Err
}

type PathAccessor interface {
	Access(input JsonNode) (JsonNode, error)
}

type ObjectPathAccessor struct {
	raw string
	key string
}

func (a ObjectPathAccessor) Access(input JsonNode) (JsonNode, error) {
	if input.IsLeaf() {
		return nil, PathAccessorError{Accessor: a.raw}
	}

	val, err := input.At(a.key)

	if err != nil {
		return nil, err
	}

	return val, nil
}

type ArrayPathAccessor struct {
	raw   string
	index int
}

func (a *ArrayPathAccessor) Access(input JsonNode) (JsonNode, error) {
	if input.IsLeaf() {
		return nil, PathAccessorError{Accessor: a.raw}
	}

	val, err := input.At(a.index)

	if err != nil {
		return nil, err
	}

	return val, nil
}

type ArrayIteratorAccessor struct {
	raw string
}

func (a *ArrayIteratorAccessor) Access(input JsonNode) (JsonNode, error) {
	if input.IsLeaf() {
		return nil, PathAccessorError{Accessor: a.raw}
	}

	if val, ok := input.(ArrayJsonNode); ok {
		return val, nil
	}

	return nil, PathAccessorError{Accessor: a.raw}
}

func createArrayAccessor(accessorString string) (PathAccessor, error) {
	indexString := strings.Trim(accessorString, "[]")

	if len(indexString) == 0 {
		return &ArrayIteratorAccessor{raw: accessorString}, nil
	}

	i, err := strconv.Atoi(indexString)

	if err != nil {
		return nil, PathAccessorError{Accessor: accessorString}
	}

	return &ArrayPathAccessor{index: i, raw: accessorString}, nil
}

func createPathAccessor(accessorString string) (PathAccessor, error) {
	if strings.HasPrefix(accessorString, "[") && strings.HasSuffix(accessorString, "]") {
		accessor, err := createArrayAccessor(accessorString)

		if err != nil {
			return nil, err
		}

		return accessor, nil
	}

	return ObjectPathAccessor{raw: accessorString, key: accessorString}, nil
}

type PathFilter struct {
	accessors []PathAccessor
}

func (f *PathFilter) Filter(input JsonNode) (JsonNode, error) {
	current := input
	for _, accessor := range f.accessors {
		result, err := accessor.Access(current)

		if err != nil {
			return nil, err
		}

		current = result
	}

	return current, nil
}

func CreatePathFilter(filterString string) (*PathFilter, error) {
	// Get path parts. First part will be '' from split so remove it
	rawAccessors := strings.Split(filterString, ".")[1:]
	accessors := []PathAccessor{}
	for _, rawAccessor := range rawAccessors {
		accessor, err := createPathAccessor(rawAccessor)

		if err != nil {
			return nil, PathError{Path: filterString, Err: err}
		}

		accessors = append(accessors, accessor)
	}

	return &PathFilter{accessors: accessors}, nil
}

func IsPathFilter(filterString string) bool {
	if strings.HasPrefix(filterString, ".") {
		return true
	}

	return false
}
