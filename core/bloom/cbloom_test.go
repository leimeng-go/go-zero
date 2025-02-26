package bloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomBloomFilter(t *testing.T) {
    bf:=NewWithEstimates(10000, 0.01)	

	bf.Add([]byte("apple"))
	bf.Add([]byte("banana"))
	bf.Add([]byte("cherry"))

	assert.True(t, bf.Contains([]byte("apple")))
	assert.True(t, bf.Contains([]byte("banana")))
	assert.True(t, bf.Contains([]byte("cherry")))
	
	bf.Clear()
	assert.False(t, bf.Contains([]byte("apple")))
}
