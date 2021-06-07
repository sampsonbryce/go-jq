package filter

import (
	"encoding/json"
	"fmt"
	"log"
)

const RESPONSE_TYPE_MAP = "MAP"
const RESPONSE_TYPE_ARRAY = "ARRAY"
const RESPONSE_TYPE_VALUE = "VALUE"

type Response struct {
	valueMap   map[string]interface{}
	valueArray []interface{}
	value      *interface{}
	valueType  string
}

func (r *Response) Marshal() ([]byte, error) {
	if r.valueType == RESPONSE_TYPE_MAP {
		return json.Marshal(r.valueMap)
	} else if r.valueType == RESPONSE_TYPE_ARRAY {
		return json.Marshal(r.valueArray)
	} else if r.valueType == RESPONSE_TYPE_VALUE {
		return json.Marshal(r.value)
	}

	log.Fatal(fmt.Sprintf("Invalid response valueType %s in marshal\n", r.valueType))

	return []byte{}, nil // should never be reached
}

func ExecTree(input *interface{}, tree *Tree) Response {
	if len(tree.root.children) == 0 {
		log.Fatal("Attempt to execute empty tree")
	}

	currentResponse := Exec(input, tree.root.children[0])

	for i := 1; i < len(tree.root.children); i++ {
		currentResponse = execResponse(currentResponse, tree.root.children[i])
		fmt.Printf("Current response %v\n", currentResponse)
	}

	return currentResponse
}

func Exec(input *interface{}, node *Node) Response {
	if val, ok := IsSlice(*input); ok {
		return execArray(val, node)
	} else if val, ok := IsMap(*input); ok {
		return execMap(val, node)
	} else {
		return execValue(input, node)
	}
}

func execResponse(response Response, node *Node) Response {
	fmt.Printf("Execing response %v %v\n", response, *node)
	if response.valueType == RESPONSE_TYPE_MAP {
		return execMap(response.valueMap, node)
	} else if response.valueType == RESPONSE_TYPE_ARRAY {
		return execArray(response.valueArray, node)
	} else if response.valueType == RESPONSE_TYPE_VALUE {
		return execValue(response.value, node)
	}

	log.Fatal(fmt.Sprintf("Invalid response valueType %s in execResponse\n", response.valueType))

	return Response{} // Should never be reached
}

func execValue(input *interface{}, node *Node) Response {
	log.Fatal("Cannot process a value")
	return createResponse(input)
}

func execArray(input []interface{}, node *Node) Response {
	fmt.Printf("Exec array %v %v\n", input, *node)
	return Response{}
}

func execMap(input map[string]interface{}, node *Node) Response {
	fmt.Printf("Exec map %v %v\n", input, *node)
	if node.name == ACCESS_NAME {
		response := Response{valueMap: input, valueType: RESPONSE_TYPE_MAP}
		for _, childNode := range node.children {
			response = execResponse(response, childNode)
		}

		return response
	} else if node.name == T_KEY {
		if val, ok := input[node.value]; ok {
			fmt.Printf("Got val from object %v\n", val)
			return createResponse(&val)
		}
		log.Fatal(fmt.Sprintf("Missing key '%s' in object %v\n", node.value, input))
	} else if node.name == COLLECT_NAME {
		if len(node.children) == 1 && node.children[0].name == T_STRING {
			key := node.children[0].value

			if val, ok := input[key]; ok {
				return createResponse(&val)
			}

			log.Fatal(fmt.Sprintf("Missing key '%s' in object %v\n", key, input))
		}

		log.Fatal("Invalid value in collect for accessing a map")
	}

	log.Fatal(fmt.Sprintf("Unable to parse map input using node %v\n", node))

	return Response{}
}

func createResponse(v *interface{}) Response {
	if val, ok := IsSlice(*v); ok {
		return Response{valueArray: val, valueType: RESPONSE_TYPE_ARRAY}
	} else if val, ok := IsMap(*v); ok {
		return Response{valueMap: val, valueType: RESPONSE_TYPE_MAP}
	} else {
		return Response{value: v, valueType: RESPONSE_TYPE_VALUE}
	}
}
