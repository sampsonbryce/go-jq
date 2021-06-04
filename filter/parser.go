package filter

import (
	"fmt"
	"log"
	"strings"

	"github.com/sampsonbryce/go-jq/util"
)

const TOKEN_RUNES = "ABCDEFGHIJKLMNOPQRSTUVWXYZ_"

type Node struct {
	name     string
	value    string
	parent   *Node
	children []*Node
}

func (n *Node) Print(level int) {
	fmt.Printf("%s%s %s\n", strings.Repeat("\t", level), n.name, n.value)

	// fmt.Printf("children %d %v\n", len(n.children), n.children)

	for _, child := range n.children {
		child.Print(level + 1)
	}
}

type Tree struct {
	root        *Node
	currentNode *Node
}

func (t *Tree) AddNode(name string, value string) *Node {
	node := Node{name: name, value: value}
	// fmt.Printf("Current node %p\n", t.currentNode)
	t.currentNode.children = append(t.currentNode.children, &node)

	// fmt.Printf("Node added %s %s %p\n", name, value, &node)
	return &node
}

func (t *Tree) AddNodeAndRecurse(name string, value string) *Node {
	node := t.AddNode(name, value)

	// fmt.Println("Recursing")
	t.currentNode = node

	return node
}

func (t *Tree) EndNode() {
	t.currentNode = t.currentNode.parent
}

func (t *Tree) Print() {
	t.root.Print(0)
}

func createTree() Tree {
	root := Node{name: "ROOT"}
	tree := Tree{root: &root, currentNode: &root}

	return tree
}

func isTokenRune(char rune) bool {
	return strings.ContainsRune(TOKEN_RUNES, char)
}

func consumeToken(runes []rune, i int) (string, int) {
	length := len(runes)
	j := i
	token := ""
	for j < length {
		currentChar := runes[j]
		if !isTokenRune(currentChar) {
			break
		}

		token += string(currentChar)
		j += 1
	}

	return token, j - 1 // go back to last consumed char
}

func consumeString(runes []rune, i int) (string, int) {
	length := len(runes)
	j := i + 1 // +1 to skip leading quote
	token := ""
	for j < length {
		currentChar := runes[j]
		if currentChar == '"' {
			break
		}

		token += string(currentChar)
		j += 1
	}

	return token, j // Dont go back to last char so we consume closing quote
}

func consumeValue(runes []rune, i int) (string, int) {
	length := len(runes)
	j := i + 1

	token := ""
	for j < length {
		currentChar := runes[j]
		if currentChar == ' ' {
			break
		}

		token += string(currentChar)
		j += 1
	}

	return token, j - 1
}

func isStringToken(token string) bool {
	stringTokens := []string{T_STRING}

	return util.Contains(stringTokens, token)
}

func isValueToken(token string) bool {
	stringTokens := []string{T_IDENTIFIER, T_KEY, T_NUMBER}

	return util.Contains(stringTokens, token)
}

func isRecurseToken(token string) bool {
	return strings.HasSuffix(token, "_START")
}

func getRecurseTokenName(token string) string {
	return strings.TrimSuffix(token, "_START")
}

func isEndToken(token string) bool {
	return strings.HasSuffix(token, "_END")
}

func Parse(text string) Tree {
	tree := createTree()
	runes := []rune(text)
	length := len(runes)
	i := 0

	for i < length {
		currentChar := runes[i]
		// fmt.Printf("At char %d\n", i)

		if currentChar == 'T' {
			token, j := consumeToken(runes, i)
			// fmt.Printf("Got token %s\n", token)
			i = j // +2 to skip extra whitespace
			// fmt.Printf("New i %d\n", i)
			value := ""
			recurse := isRecurseToken(token)
			end := isEndToken(token)

			if isStringToken(token) {
				stringValue, j := consumeString(runes, i+2)
				i = j
				value = stringValue
			} else if isValueToken(token) {
				rawValue, j := consumeValue(runes, i+2)
				i = j
				value = rawValue
			}

			if end {
				tree.EndNode()
			} else if recurse {
				tree.AddNodeAndRecurse(getRecurseTokenName(token), value)
			} else {
				tree.AddNode(token, value)
			}
		} else if currentChar == ' ' {
			// pass
		} else {
			log.Fatal(fmt.Printf("Unexpected char '%s' at char %d\n", string(currentChar), i))
		}

		i += 1
	}

	return tree
}
