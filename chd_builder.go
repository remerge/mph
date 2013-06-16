package mph

import (
	"errors"
	"fmt"
	"hash/crc64"
	"math/rand"
	"sort"
	"time"
)

type CHDHasher struct {
	r       []uint64
	size    uint64
	buckets uint64
	rand    *rand.Rand
}

var (
	crc64table = crc64.MakeTable(crc64.ECMA)
)

func CDHHash(b []byte) uint64 {
	return crc64.Checksum(b, crc64table)
}

type CHDKeyValue struct {
	key   []byte
	value []byte
}

func (c *CHDKeyValue) Key() []byte {
	return c.key
}

func (c *CHDKeyValue) Value() []byte {
	return c.value
}

type bucket struct {
	index uint64
	kv    []*CHDKeyValue
}

func (b *bucket) String() string {
	a := "bucket{"
	for _, kv := range b.kv {
		a += string(kv.key) + ", "
	}
	return a + "}"
}

// Intermediate data structure storing buckets + outer hash index.
type bucketVector []bucket

func (b bucketVector) Len() int           { return len(b) }
func (b bucketVector) Less(i, j int) bool { return len(b[i].kv) > len(b[j].kv) }
func (b bucketVector) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// Build a new CDH MPH.
type CHDBuilder struct {
	kv []*CHDKeyValue
}

// Create a new builder.
func NewCHDBuilder() *CHDBuilder {
	return &CHDBuilder{}
}

// Add a key and value to the hash table.
func (b *CHDBuilder) Add(key []byte, value []byte) {
	b.kv = append(b.kv, &CHDKeyValue{key, value})
}

// Try to find a hash function that does not cause collisions with table, when
// applied to the keys in the bucket.
func tryHash(hasher *CHDHasher, seen map[uint64]bool, table []*CHDKeyValue, indices []uint64, bucket *bucket, ri, r uint64) bool {
	// Track duplicates within this bucket.
	duplicate := make(map[uint64]bool)
	// Make hashes for each entry in the bucket.
	hashes := make([]uint64, len(bucket.kv))
	for i, kv := range bucket.kv {
		h := hasher.Table(r, kv.key)
		hashes[i] = h
		if seen[h] {
			return false
		}
		if duplicate[h] {
			return false
		}
		duplicate[h] = true
	}

	// Update seen hashes
	for _, h := range hashes {
		seen[h] = true
	}

	// Add the hash index.
	indices[bucket.index] = ri

	// Update the the hash table.
	for i, h := range hashes {
		table[h] = bucket.kv[i]
	}
	return true
}

func (b *CHDBuilder) Build() (*CHD, error) {
	n := uint64(len(b.kv))
	m := n / 2
	// We need a minimum bucket size so we get enough variability in the lower
	// bits after modulo is applied.
	for m < 10 {
		m++
	}

	table := make([]*CHDKeyValue, n)
	hasher := NewCHDHasher(n, m)
	buckets := make(bucketVector, m)
	indices := make([]uint64, m)
	// An extra check to make sure we don't use an invalid index
	for i := range indices {
		indices[i] = ^uint64(0)
	}
	// Have we seen a hash before?
	seen := make(map[uint64]bool)
	// Used to ensure there are no duplicate keys.
	duplicates := make(map[string]bool)

	for _, kv := range b.kv {
		k := string(kv.key)
		if duplicates[k] {
			return nil, errors.New("duplicate key " + k)
		}
		duplicates[k] = true
		oh := hasher.HashIndexFromKey(kv.key)

		buckets[oh].index = oh
		buckets[oh].kv = append(buckets[oh].kv, kv)
	}

	// Order buckets by size (retaining the hash index)
	collisions := 0
	sort.Sort(buckets)
nextBucket:
	for i, bucket := range buckets {
		if len(bucket.kv) == 0 {
			continue
		}

		// Check existing hash functions.
		for ri, r := range hasher.r {
			if tryHash(hasher, seen, table, indices, &bucket, uint64(ri), r) {
				continue nextBucket
			}
		}

		// Keep trying new functions until we get one that does not collide.
		// The number of retries here is very high due to our use of the IEEE
		// checksum function. If we use a better hashing function, less tries
		// are necessary, but retrieval is quite a bit slower.
		for i := 0; i < 10000000; i++ {
			if i > collisions {
				collisions = i
			}
			ri, r := hasher.Generate()
			if tryHash(hasher, seen, table, indices, &bucket, ri, r) {
				hasher.Add(r)
				continue nextBucket
			}
		}

		// Failed to find a hash function with no collisions.
		return nil, errors.New(fmt.Sprintf(
			"failed to find a collision-free hash function after ~10000000 attempts, for bucket %d/%d with %d entries: %s",
			i, len(buckets), len(bucket.kv), &bucket))
	}

	// println("max bucket collisions:", collisions)
	// println("keys:", len(table))
	// println("hash functions:", len(hasher.r))

	return &CHD{
		r:       hasher.r,
		indices: indices,
		table:   table,
	}, nil
}

func NewCHDHasher(size uint64, buckets uint64) *CHDHasher {
	rs := rand.NewSource(time.Now().UnixNano())
	c := &CHDHasher{size: size, buckets: buckets, rand: rand.New(rs)}
	c.Add(c.random())
	return c
}

func (c *CHDHasher) random() uint64 {
	return (uint64(c.rand.Uint32()) << 32) | uint64(c.rand.Uint32())
}

// Hash index from key.
func (h *CHDHasher) HashIndexFromKey(b []byte) uint64 {
	return (CDHHash(b) ^ h.r[0]) % h.buckets
}

// Table hash from random value and key. Generate() returns these random values.
func (h *CHDHasher) Table(r uint64, b []byte) uint64 {
	return (CDHHash(b) ^ h.r[0] ^ r) % h.size
}

func (c *CHDHasher) Generate() (uint64, uint64) {
	return c.Len(), c.random()
}

// Add a random value generated by Generate().
func (c *CHDHasher) Add(r uint64) {
	c.r = append(c.r, r)
}

func (c *CHDHasher) Len() uint64 {
	return uint64(len(c.r))
}

func (h *CHDHasher) String() string {
	return fmt.Sprintf("CHDHasher{size: %d, buckets: %d, r: %v}", h.size, h.buckets, h.r)
}
