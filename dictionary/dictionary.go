package dictionary

type Dictionary interface {
	Initialize() error
	GetSimilar(key string) []string
}
