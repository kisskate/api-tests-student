package utils

import "unicode"

func MixCase(value string) string {
	runes := []rune(value)
	upper := true
	for i, char := range runes {
		if !unicode.IsLetter(char) {
			continue
		}
		if upper {
			runes[i] = unicode.ToUpper(char)
		} else {
			runes[i] = unicode.ToLower(char)
		}
		upper = !upper
	}
	return string(runes)
}
