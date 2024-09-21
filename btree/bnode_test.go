package btree

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBNode(t *testing.T) {
	// Create a sample BNode for testing
	node := make(BNode, 4096)
	node.setHeader(BNODE_NODE, 3)

	t.Run("btype", func(t *testing.T) {
		assert.Equal(t, BNODE_NODE, node.btype())
	})

	t.Run("nkeys", func(t *testing.T) {
		assert.Equal(t, uint16(3), node.nkeys())
	})

	t.Run("setHeader", func(t *testing.T) {
		node.setHeader(BNODE_NODE, 3)
		if node.btype() != BNODE_NODE || node.nkeys() != 3 {
			t.Errorf("setHeader failed, got btype %d and nkeys %d", node.btype(), node.nkeys())
		}
	})

	t.Run("getPtr and setPtr", func(t *testing.T) {
		node.setPtr(0, 12345)
		if node.getPtr(0) != 12345 {
			t.Errorf("Expected ptr 12345, got %d", node.getPtr(0))
		}
	})

	t.Run("getOffset and setOffset", func(t *testing.T) {
		node.setOffset(1, 100)
		if node.getOffset(1) != 100 {
			t.Errorf("Expected offset 100, got %d", node.getOffset(1))
		}
	})

	t.Run("kvPos", func(t *testing.T) {
		pos := node.kvPos(0)
		if pos != 34 {
			t.Errorf("Expected kvPos %d, got %d", 34, pos)
		}
	})

	t.Run("getKey and getVal", func(t *testing.T) {
		// Set up a key-value pair
		key := []byte("testkey")
		val := []byte("testvalue")
		pos := node.kvPos(0) // Get the position for the first key-value pair

		binary.LittleEndian.PutUint16(node[pos:], uint16(len(key)))
		binary.LittleEndian.PutUint16(node[pos+2:], uint16(len(val)))

		copy(node[uint16(38):], key)
		copy(node[uint16(38)+uint16(len(key)):], val)

		// Set offsets based on the positions
		node.setOffset(1, pos) // Offset for the first key (index 1)

		// Debug information
		t.Logf("Node length: %d", len(node))
		t.Logf("pos: %d", pos)
		t.Logf("Offset 1: %d", node.getOffset(1))

		assert.Equal(t, key, node.getKey(0)) // Get key at index 1
		assert.Equal(t, val, node.getVal(0)) // Get value at index 1
	})

	t.Run("nbytes", func(t *testing.T) {
		nbytes := node.nbytes()
		if nbytes <= HEADER {
			t.Errorf("Expected nbytes > %d, got %d", HEADER, nbytes)
		}
	})
}
