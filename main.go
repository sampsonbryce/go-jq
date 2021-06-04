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

	lexedString := filter.Lex(filterString)
	fmt.Println(lexedString)
	tree := filter.Parse(lexedString)
	tree.Print()
	// input, err := readInput()

	// if err != nil {
	// 	log.Fatal("failed to parse input: ", err)
	// }

	// filters, err := filter.CreateFilters(filterString)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Got input %#v\n", input)

	// rootNode, err := filter.CreateJsonNode(&input)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// resultNode, err := processInput(rootNode, filters)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// marshalled, err := resultNode.Marshal()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(marshalled))
}

func processInput(input filter.JsonNode, filters []filter.Filter) (filter.JsonNode, error) {
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
