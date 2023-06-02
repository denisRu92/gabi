package file_dictionary

import (
	"bufio"
	mapset "github.com/deckarep/golang-set"
	"os"
	"palo-alto/config"
	"palo-alto/dictionary"
	logger "palo-alto/logging"
	"palo-alto/metric"
	"palo-alto/util"
	"strings"
)

type fileDictionary struct {
	cfg          config.Config
	m            metric.Metric
	permutations map[string]mapset.Set

	permutationsCh chan permutationsReq
	addWordCh      chan string
	stopCh         chan struct{}
}

type permutationsReq struct {
	key    string
	respCh chan mapset.Set
}

// New return new FileDictionary instance
func New(cfg config.Config, m metric.Metric) dictionary.Dictionary {
	fs := &fileDictionary{
		cfg:          cfg,
		m:            m,
		permutations: make(map[string]mapset.Set),

		permutationsCh: make(chan permutationsReq),
		addWordCh:      make(chan string),
		stopCh:         make(chan struct{}),
	}

	return fs
}

func (fd *fileDictionary) Start() {
	for {
		select {
		case word := <-fd.addWordCh:
			fd.addWord(word)
		case req := <-fd.permutationsCh:
			req.respCh <- fd.getSimilar(req.key)
		case <-fd.stopCh:
			return
		}
	}
}

func (fd *fileDictionary) Stop() {
	close(fd.stopCh)
}

// Initialize init FileDictionary
func (fd *fileDictionary) Initialize() error {
	// Open the file
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + fd.cfg.WordsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file contents
	wordCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordCount++
		fd.AddWord(scanner.Text())
	}

	logger.Log.Infof("Init %d words to dictionary", wordCount)

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (fd *fileDictionary) GetSimilar(key string) mapset.Set {
	respCh := make(chan mapset.Set)
	fd.permutationsCh <- permutationsReq{
		respCh: respCh,
		key:    key,
	}

	return <-respCh
}

// getSimilar returns array of permutation if exists else empty array
func (fd *fileDictionary) getSimilar(key string) mapset.Set {
	if val, ok := fd.permutations[key]; ok {
		return val
	}
	return mapset.NewSet()
}

func (fd *fileDictionary) AddWord(word string) {
	fd.addWordCh <- word
}

// addWord adds a new word to the permutations dictionary
func (fd *fileDictionary) addWord(word string) {
	currWord := strings.TrimSpace(word)
	key := util.SortString(currWord)

	if val, ok := fd.permutations[key]; ok {
		val.Add(currWord)
	} else {
		permutations := mapset.NewSet()
		permutations.Add(word)
		fd.permutations[key] = permutations
	}

	fd.m.IncWordCounter()
}
