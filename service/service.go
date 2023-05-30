package service

import (
	"palo-alto/dictionary"
	"palo-alto/util"
	"strings"
)

type (
	DictionaryReader interface {
		GetSimilar(word string) []string
	}

	service struct {
		dictionary dictionary.Dictionary
	}
)

func New(dictionary dictionary.Dictionary) DictionaryReader {
	return &service{dictionary: dictionary}
}

func (s service) GetSimilar(word string) []string {
	lower := strings.ToLower(word)
	permutations := s.dictionary.GetSimilar(util.SortString(lower))

	return s.removeWord(lower, permutations)
}

func (s service) removeWord(word string, permutations []string) []string {
	for i := 0; i < len(permutations); i++ {
		if word == permutations[i] {
			return append(permutations[:i], permutations[i+1:]...)
		}
	}

	return permutations
}
