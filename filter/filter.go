package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func IsSlice(v interface{}) ([]interface{}, bool) {
	val, ok := v.([]interface{})

	return val, ok
}

func IsMap(v interface{}) (map[string]interface{}, bool) {
	val, ok := v.(map[string]interface{})

	return val, ok
}

func IsJsonNode(v interface{}) (JsonNode, bool) {
	val, ok := v.(JsonNode)

	return val, ok
}

func IsObjJsonNode(v interface{}) (ObjJsonNode, bool) {
	val, ok := v.(ObjJsonNode)

	return val, ok
}

func IsArrayJsonNode(v interface{}) (ArrayJsonNode, bool) {
	val, ok := v.(ArrayJsonNode)

	return val, ok
}

type JsonNode interface {
	At(key interface{}) (JsonNode, error)
	IsLeaf() bool
	Marshal() ([]byte, error)
}

type ObjJsonNode struct {
	wrapped map[string]interface{}
}

func (o ObjJsonNode) At(key interface{}) (JsonNode, error) {
	var k string
	if val, ok := key.(string); ok {
		k = val
	} else {
		return nil, fmt.Errorf("invalid index '%v'", key)
	}

	if val, ok := o.wrapped[k]; ok {
		node, err := CreateJsonNode(&val)

		if err != nil {
			return nil, err
		}

		return node, nil
	}

	return nil, fmt.Errorf("missing value for key '%v'", key)
}

func (o ObjJsonNode) IsLeaf() bool {
	return false
}

func (o ObjJsonNode) Marshal() ([]byte, error) {
	return json.Marshal(o.wrapped)
}

type ArrayJsonNode struct {
	wrapped []interface{}
}

func (o ArrayJsonNode) At(index interface{}) (JsonNode, error) {
	var i int
	if val, ok := index.(int); ok {
		i = val
	} else {
		return nil, fmt.Errorf("invalid index '%v'", index)
	}

	if len(o.wrapped) > i {
		val := o.wrapped[i]

		node, err := CreateJsonNode(&val)

		if err != nil {
			return nil, err
		}

		return node, nil
	}

	return nil, errors.New(fmt.Sprintf("index out of range '%v'", index))
}

func (o ArrayJsonNode) IsLeaf() bool {
	return false
}

func (o ArrayJsonNode) Marshal() ([]byte, error) {
	return json.Marshal(o.wrapped)
}

type StringJsonNode struct {
	wrapped string
}

func (o StringJsonNode) At(_ interface{}) (JsonNode, error) {
	return nil, errors.New("not implemented")
}

func (o StringJsonNode) IsLeaf() bool {
	return true
}

func (o StringJsonNode) Marshal() ([]byte, error) {
	return json.Marshal(o.wrapped)
}

type IntegerJsonNode struct {
	wrapped int
}

func (o IntegerJsonNode) At(_ interface{}) (JsonNode, error) {
	return nil, errors.New("not implemented")
}

func (o IntegerJsonNode) IsLeaf() bool {
	return true
}

func (o IntegerJsonNode) Marshal() ([]byte, error) {
	return json.Marshal(o.wrapped)
}

func CreateJsonNode(input *interface{}) (JsonNode, error) {
	switch v := (*input).(type) {
	default:
		return nil, fmt.Errorf("unexpected type %T", v)
	case int:
		return IntegerJsonNode{wrapped: v}, nil
	case string:
		return StringJsonNode{wrapped: v}, nil
	case map[string]interface{}:
		return ObjJsonNode{wrapped: v}, nil
	case []interface{}:
		return ArrayJsonNode{wrapped: v}, nil
	}
}

type Filter interface {
	Filter(input JsonNode) (JsonNode, error)
}

func CreateFilters(filterString string) ([]Filter, error) {
	filterSections := strings.Split(filterString, "|")

	filters := []Filter{}
	for _, filterString := range filterSections {
		if IsPathFilter(filterString) {
			filter, err := CreatePathFilter(filterString)

			if err != nil {
				return nil, err
			}

			filters = append(filters, filter)
		}
	}

	return filters, nil
}
