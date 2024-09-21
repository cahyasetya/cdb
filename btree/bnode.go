package btree

import (
	"encoding/binary"
	"fmt"
	"strings"

	. "github.com/cahyasetya/cdb/util"
)

const (
	BNODE_NODE uint16 = 1
	BNODE_LEAF uint16 = 2
)

// BNode represents a node in the B-tree
type BNode []byte

// btype returns the type of the node (BNODE_NODE or BNODE_LEAF)
//
// Example:
//
//	First two bytes: 01 00 -> Returns: 1 (BNODE_NODE)
//	First two bytes: 02 00 -> Returns: 2 (BNODE_LEAF)
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

// nkeys returns the number of keys in the node
//
// Example:
//
//	First four bytes: 01 00 03 00 -> Returns: 3 (3 keys in the node)
//	First four bytes: 02 00 0A 00 -> Returns: 10 (10 keys in the node)
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// getHeader returns the header of the node
func (node BNode) getHeader() struct {
	typ   uint16
	nkeys uint16
} {
	return struct {
		typ   uint16
		nkeys uint16
	}{
		typ:   node.btype(),
		nkeys: node.nkeys(),
	}
}

// setHeader sets the node type and number of keys in the node header
func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// getPtr retrieves a pointer value at the given index
func (node BNode) getPtr(idx uint16) uint64 {
	Assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

// setPtr sets a pointer value at the given index
func (node BNode) setPtr(idx uint16, val uint64) {
	Assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// getOffset retrieves the offset value for a given index
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := offsetPos(node, idx)
	return binary.LittleEndian.Uint16(node[pos:])
}

// setOffset sets the offset value for a given index
func (node BNode) setOffset(idx uint16, val uint16) {
	Assert(1 <= idx && idx <= node.nkeys())
	binary.LittleEndian.PutUint16(node[offsetPos(node, idx):], val)
}

// kvPos calculates the position of the key-value pair for a given index
func (node BNode) kvPos(idx uint16) uint16 {
	Assert(idx <= node.nkeys())
	offset := node.getOffset(idx)
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + offset
}

// getKey retrieves the key at the given index
func (node BNode) getKey(idx uint16) []byte {
	Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:klen]
}

// getVal retrieves the value at the given index
func (node BNode) getVal(idx uint16) []byte {
	Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	vlen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+klen:][:vlen]
}

// nbytes returns the total number of bytes used in the node
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

func (node BNode) String() string {
	var builder strings.Builder
	for i, b := range node {
		if i > 0 {
			if i%16 == 0 {
				builder.WriteString("\n")
			} else {
				builder.WriteString(" ")
			}
		}
		builder.WriteString(fmt.Sprintf("%3d", b))
	}
	return builder.String()
}
