package ntrees

import "strings"

const htmlSpace = "&nbsp;"

// MakeHTMLSpace returns a text containing '&nsbp;' the provided
// number 'count' times.
func MakeHTMLSpace(count int) string {
	if count <= 0 {
		count = 1
	}
	var spaces []string
	for i := 0; i < count; i++ {
		spaces = append(spaces, htmlSpace)
	}
	return strings.Join(spaces, "")
}
