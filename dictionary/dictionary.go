package dictionary

import mapset "github.com/deckarep/golang-set"

type Dictionary interface {
	Start()
	Stop()
	Initialize() error
	GetSimilar(key string) mapset.Set
	AddWord(word string)
}
