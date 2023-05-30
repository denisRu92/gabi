package dictionary

type Dictionary interface {
	Start()
	Stop()
	Initialize() error
	GetSimilar(key string) []string
	AddWord(word string)
}
