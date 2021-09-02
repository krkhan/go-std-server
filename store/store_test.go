package store

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSha512DigestStore(t *testing.T) {
	var delaySeconds int32 = 5
	totalCallers := 100
	digestsMap := make(map[uint64][sha512.Size]byte)
	store := &Sha512DigestStore{
		Counter:      0,
		Digests:      &digestsMap,
		DigestsLock:  &sync.RWMutex{},
		DelaySeconds: time.Duration(delaySeconds),
	}

	var wg sync.WaitGroup
	wg.Add(totalCallers)

	for i := 0; i < totalCallers; i++ {
		go func() {
			defer wg.Done()
			var digest [sha512.Size]byte
			rand.Read(digest[:])
			key := store.AddDigest(digest)
			time.Sleep(time.Duration(2*delaySeconds) * time.Second)
			returnedDigest, ok := store.GetDigest(key)
			if !ok {
				panic(fmt.Sprintf("Did not find key in store: %d", key))
			}
			if !bytes.Equal(digest[:], returnedDigest[:]) {
				panic(fmt.Sprintf("Did not get the right digest from store, expected: '%v' returned: '%v'", digest, returnedDigest))
			}
		}()
	}

	wg.Wait()
}
