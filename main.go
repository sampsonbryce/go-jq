package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sampsonbryce/go-jq/filter"
)

func readInput() (interface{}, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	var f interface{}
	err = json.Unmarshal([]byte(text), &f)

	if err != nil {
		return "", err
	}

	return f, nil
}

func main() {
	fmt.Println(len(os.Args), os.Args)

	if len(os.Args) < 2 {
		log.Fatal("missing FILTER argument")
	}

	filterString := os.Args[1]

	input, err := readInput()

	if err != nil {
		log.Fatal("failed to parse input")
	}

	filters, err := filter.CreateFilters(filterString)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got input %#v\n", input)

	result, err := processInput(&input, filters)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Raw result %#v\n", *result)
	marshalled, err := json.Marshal(*result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(marshalled))
}

func processInput(input *interface{}, filters []filter.Filter) (*interface{}, error) {
	current := input
	for _, filter := range filters {
		result, err := filter.Filter(current)

		if err != nil {
			return nil, err
		}

		current = result
	}

	return current, nil
}
