package bloom

import (
	"hash/fnv"
	"math"
)

// BloomFilter represents a Bloom filter data structure
type BloomFilter struct {
	bitset    []bool // The bit array
	size      uint   // Size of the bit array
	hashCount uint   // Number of hash functions
}

// NewWithEstimates creates a new Bloom filter optimized for the expected number of elements
// and desired false positive probability
func NewWithEstimates(n uint, falsePositiveRate float64) *BloomFilter {
	// Calculate optimal size and number of hash functions
	size := calculateOptimalSize(n, falsePositiveRate)
	hashCount := calculateOptimalHashCount(size, n)

	return &BloomFilter{
		bitset:    make([]bool, size),
		size:      size,
		hashCount: hashCount,
	}
}

// New creates a new Bloom filter with specified size and hash count
func NewC(size, hashCount uint) *BloomFilter {
	return &BloomFilter{
		bitset:    make([]bool, size),
		size:      size,
		hashCount: hashCount,
	}
}

// Add adds an item to the Bloom filter
func (bf *BloomFilter) Add(item []byte) {
	for i := uint(0); i < bf.hashCount; i++ {
		position := bf.hash(item, i) % bf.size
		bf.bitset[position] = true
	}
}

// Contains checks if an item might be in the set
// False positives are possible, but false negatives are not
func (bf *BloomFilter) Contains(item []byte) bool {
	for i := uint(0); i < bf.hashCount; i++ {
		position := bf.hash(item, i) % bf.size
		if !bf.bitset[position] {
			return false
		}
	}
	return true
}

// Clear resets the Bloom filter
func (bf *BloomFilter) Clear() {
	bf.bitset = make([]bool, bf.size)
}

// hash generates different hash values for the same item
// using the FNV hash with a seed based on the index
func (bf *BloomFilter) hash(item []byte, index uint) uint {
	h := fnv.New64a()
	// Adding the index to the data to create different hash functions
	data := append(item, byte(index))
	h.Write(data)
	return uint(h.Sum64())
}

// calculateOptimalSize calculates the optimal size based on expected elements and false positive rate
func calculateOptimalSize(n uint, falsePositiveRate float64) uint {
	size := -float64(n) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)
	return uint(math.Ceil(size))
}

// calculateOptimalHashCount calculates the optimal number of hash functions
func calculateOptimalHashCount(size, n uint) uint {
	hashCount := float64(size) / float64(n) * math.Log(2)
	return uint(math.Max(1, math.Round(hashCount)))
}