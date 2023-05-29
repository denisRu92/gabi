package fileDictionary

import (
	"bufio"
	"os"
	"palo-alto/conf"
	"palo-alto/dictionary"
	logger "palo-alto/logging"
	"palo-alto/metric"
	"palo-alto/util"
	"strings"
)

type fileDictionary struct {
	cfg          conf.Config
	m            *metric.Metric
	permutations map[string][]string
}

// New return new FileDictionary instance
func New(cfg conf.Config, m *metric.Metric) dictionary.Dictionary {
	fs := &fileDictionary{
		cfg:          cfg,
		m:            m,
		permutations: make(map[string][]string),
	}

	return fs
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
		currWord := strings.TrimSpace(scanner.Text())
		key := util.SortString(currWord)

		fd.permutations[key] = append(fd.permutations[key], currWord)
	}

	fd.m.AddWordCounter(int64(wordCount))
	logger.Log.Infof("Init %d words to dictionary", wordCount)

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// GetSimilar returns array of permutation if exists else empty array
func (fd *fileDictionary) GetSimilar(key string) []string {
	if val, ok := fd.permutations[key]; ok {
		return val
	}
	return []string{}
}
