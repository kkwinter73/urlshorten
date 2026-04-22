package main

import (
	"crypto/rand"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("url not found")

const (
	idLength = 6
	idChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type URLStore struct {
	mu   sync.Mutex
	urls map[string]*urlEntry
}

type urlEntry struct {
	URL   string
	Count int
}

type Stats struct {
	URL   string
	Count int
}

func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]*urlEntry),
	}
}

// ランダムな６文字IDを生成する
func generateID() (string, error) {
	b := make([]byte, idLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = idChars[int(b[i])%len(idChars)]
	}

	return string(b), nil
}

// URLを登録して衝突しないIDを返す
func (s *URLStore) Create(url string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 最大10回まで衝突リトライ
	for i := 0; i < 10; i++ {
		id, err := generateID()
		if err != nil {
			return "", err
		}

		if _, exists := s.urls[id]; !exists {
			s.urls[id] = &urlEntry{URL: url, Count: 0}
			return id, nil
		}
	}

	return "", errors.New("failed to generate unique id")
}

// IDから元URLを取り出す
func (s *URLStore) Get(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.urls[id]
	if !ok {
		return "", ErrNotFound
	}
	return entry.URL, nil
}

func (s *URLStore) IncrementCount(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.urls[id]; !ok {
		return ErrNotFound
	}
	s.urls[id].Count++
	return nil
}

func (s *URLStore) GetStats(id string) (*Stats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.urls[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &Stats{URL: entry.URL, Count: entry.Count}, nil
}
