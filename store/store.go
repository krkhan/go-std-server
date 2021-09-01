package store

import (
	"crypto/sha512"
	"sync"
	"sync/atomic"
	"time"
)

type Sha512DigestStoreContextKey struct{}

type Sha512DigestStore struct {
	Counter      uint64
	Digests      *map[uint64][sha512.Size]byte
	DigestsLock  *sync.RWMutex
	DelaySeconds time.Duration
}

func (s *Sha512DigestStore) AddDigest(digest [sha512.Size]byte) uint64 {
	newKey := atomic.AddUint64(&(s.Counter), 1)
	go func() {
		time.Sleep(s.DelaySeconds * time.Second)
		s.DigestsLock.Lock()
		defer s.DigestsLock.Unlock()
		digestsMap := s.Digests
		(*digestsMap)[newKey] = digest
	}()
	return newKey
}

func (s *Sha512DigestStore) GetDigest(key uint64) ([sha512.Size]byte, bool) {
	s.DigestsLock.RLock()
	defer s.DigestsLock.RUnlock()
	digestsMap := s.Digests
	v, ok := (*digestsMap)[key]
	return v, ok
}
