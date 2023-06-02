package service

import (
	"palo-alto/dictionary"
	"palo-alto/util"
	"strings"
)

type (
	DictionaryReader interface {
		GetSimilar(word string) []interface{}
	}

	service struct {
		dictionary dictionary.Dictionary
	}
)

func New(dictionary dictionary.Dictionary) DictionaryReader {
	return &service{dictionary: dictionary}
}

func (s service) GetSimilar(word string) []interface{} {
	lower := strings.ToLower(word)
	permutations := s.dictionary.GetSimilar(util.SortString(lower)).Clone()

	if permutations.Contains(lower) {
		permutations.Remove(lower)
	}

	return permutations.ToSlice()
}

//func (s service) removeWord(word string, permutations []string) []string {
//	for i := 0; i < len(permutations); i++ {
//		if word == permutations[i] {
//			result := make([]string, 0, len(permutations)-1)
//			result = append(result, permutations[:i]...)
//			return append(result, permutations[i+1:]...)
//		}
//	}
//
//	return permutations
//}
