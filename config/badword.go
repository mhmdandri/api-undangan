package config

import (
	"os"
	"strings"
)

var BadWords map[string]struct{}

func LoadBadWords(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	BadWords = make(map[string]struct{})
	for _, line := range strings.Split(string(data), "\n") {
		word := strings.TrimSpace(strings.ToLower(line))
		if word != "" {
			BadWords[word] = struct{}{}
		}
	}
	return nil
}

func ContainsBadword(message string) bool {
	msg := strings.ToLower(message)
	for bad := range BadWords {
		if strings.Contains(msg, bad){
			return true
		}
	}
	return false
}