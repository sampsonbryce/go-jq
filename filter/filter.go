package filter

import "strings"

func IsSlice(v interface{}) ([]interface{}, bool) {
	val, ok := v.([]interface{})

	return val, ok
}

func IsMap(v interface{}) (map[string]interface{}, bool) {
	val, ok := v.(map[string]interface{})

	return val, ok
}

type Filter interface {
	Filter(input *interface{}) (*interface{}, error)
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
