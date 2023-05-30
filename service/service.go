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
	return s.dictionary.GetSimilar(util.SortString(strings.ToLower(word)))
}
