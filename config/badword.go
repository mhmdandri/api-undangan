package config

import (
	"log"
	"os"
	"strings"
	"sync"
)

var (
	badWordMu sync.RWMutex
	BadWords  map[string]struct{}
)

func LoadBadWords(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Peringatan: Gagal membaca file badword %s: %v\n", path, err)
		return err
	}
	badWordMu.Lock()
	defer badWordMu.Unlock()
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
	badWordMu.RLock()
	defer badWordMu.RUnlock()
	if len(BadWords) == 0 {
		return false
	}
	msg := strings.ToLower(message)
	// Add spaces for word boundary matching
	paddedMsg := " " + msg + " "
	for bad := range BadWords {
		// Match whole words only: " bad " or at start/end with space
		badPadded := " " + bad + " "
		if strings.Contains(paddedMsg, badPadded) || msg == bad || strings.HasPrefix(msg, bad+" ") || strings.HasSuffix(msg, " "+bad) {
			return true
		}
	}
	return false
}