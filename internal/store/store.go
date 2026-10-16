package store

import "sync"

// Store — потокобезопасное хранилище ключ-значение в памяти.
type Store struct {
	mu sync.RWMutex
	m  map[string][]byte
}

// TODO(комментарий): // New создаёт пустое хранилище.
func New() *Store {
	return &Store{m: make(map[string][]byte)}
}

// TODO(комментарий): // Get возвращает значение и true, если ключ есть.
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

// value сохраняется без копирования. Это безопасно, пока вызывающий код
// больше не меняет этот срез (парсер на каждый аргумент выделяет новый буфер).
// Если когда-нибудь парсер начнёт переиспользовать буферы, тут придётся
// копировать: s.m[key] = append([]byte(nil), value...)
func (s *Store) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

// Проверка и удаление должны идти под одной блокировкой. Если сначала
// вызвать Get, а потом Delete, между ними другая горутина успеет удалить
// ключ, и обе ответят клиенту :1.
func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[key]
	delete(s.m, key)
	return ok
}
