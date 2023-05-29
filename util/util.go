package util

import "sort"

func SortString(str string) string {
	chars := []rune(str)
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})
	return string(chars)
}
