package filter

import (
	"fmt"
	"log"
	"strings"
)

// Collect
const T_COLL_START = "T_COLL_START"
const T_COLL_END = "T_COLL_END"

const T_PIPE = "T_PIPE"

const T_DOT = "T_DOT"

const T_IDENTIFIER = "T_IDENTIFIER"

const T_KEY = "T_KEY"

const T_NUMBER = "T_NUMBER"
const T_STRING = "T_STRING"

const NAME_RUNES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
const NUMBER_RUNES = "1234567890"

func is_dot(char rune) bool {
	return char == '.'
}

func is_name(char rune) bool {
	return strings.ContainsRune(NAME_RUNES, char)
}

func get_name(runes []rune, i int) (string, int) {
	length := len(runes)
	j := i
	name := ""
	for j < length {
		currentChar := runes[j]
		if !is_name(currentChar) {
			break
		}

		name += string(currentChar)
		j += 1
	}

	return name, j - 1 // go back to last consumed char
}

func is_collect_start(char rune) bool {
	return char == '['
}

func is_collect_end(char rune) bool {
	return char == ']'
}

func is_num(char rune) bool {
	return strings.ContainsRune(NUMBER_RUNES, char)
}

func is_quote(char rune) bool {
	return char == '"'
}

func get_string(runes []rune, i int) (string, int) {
	length := len(runes)
	j := i + 1 // Skip initial quote
	name := ""
	for j < length {
		currentChar := runes[j]
		if is_quote(currentChar) {
			break
		}

		name += string(currentChar)
		j += 1
	}

	return name, j // Dont go back to last char so we consume closing quote
}

func get_num(runes []rune, i int) (string, int) {
	length := len(runes)
	j := i
	num := ""

	for j < length {
		currentChar := runes[j]
		if !is_num(currentChar) {
			break
		}

		num += string(currentChar)
		j += 1
	}

	return num, j - 1 // Go back to last consume char
}

func Lex(text string) string {
	i := 0
	length := len(text)
	runes := []rune(text)
	newText := ""
	for i < length {
		currentChar := runes[i]

		if is_dot(currentChar) {
			newText += T_DOT + " "
		} else if is_name(currentChar) {
			newText += T_KEY + " "

			name, j := get_name(runes, i)
			i = j
			newText += name + " "
		} else if is_collect_start(currentChar) {
			newText += T_COLL_START + " "
		} else if is_collect_end(currentChar) {
			newText += T_COLL_END + " "
		} else if is_num(currentChar) {
			newText += T_NUMBER + " "
			num, j := get_num(runes, i)
			i = j
			newText += num + " "
		} else if is_quote(currentChar) {
			newText += T_STRING + " "
			str, j := get_string(runes, i)
			i = j
			newText += fmt.Sprintf("%q ", str) // Escape quotes in strings
		} else {
			log.Fatal(fmt.Sprintf("Unexpected character '%s' at char %d\n", string(currentChar), i))
		}

		i += 1
	}

	return newText
}
