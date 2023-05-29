package service

import (
	"palo-alto/dictionary"
	"palo-alto/util"
)

type (
	DictionaryHandler interface {
		GetSimilar(word string) []string
	}

	service struct {
		dictionary dictionary.Dictionary
	}
)

func New(dictionary dictionary.Dictionary) DictionaryHandler {
	return &service{dictionary: dictionary}
}

func (s service) GetSimilar(word string) []string {
	return s.dictionary.GetSimilar(util.SortString(word))
}
