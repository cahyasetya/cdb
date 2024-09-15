package btree

import (
	"bytes"
	"encoding/binary"

	. "github.com/cahyasetya/cdb/util"
)

// offsetPos calculates the position of the offset for a given index
func offsetPos(node BNode, idx uint16) uint16 {
    Assert(1 <= idx && idx <= node.nkeys())
    return HEADER + 8*node.nkeys() + 2*(idx-1)
}

// nodeLookup performs a binary search to find the position of a key in the node
func nodeLookupLE(node BNode, key []byte) uint16 {
    nkeys := node.nkeys()
    found := uint16(0)

    for i := uint16(1); i < nkeys; i++ {
        cmp := bytes.Compare(node.getKey(i), key)
        if cmp <= 0 {
            found = i
        }
        if cmp >= 0 {
            break
        }
    }

    return found
}

// leafInsert inserts a new key-value pair into a leaf node
func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
    new.setHeader(BNODE_LEAF, old.nkeys()+1)
    nodeAppendRange(new, old, 0, 0, idx)
    nodeAppendKV(new, idx, 0, key, val)
    nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// leafUpdate updates an existing key-value pair in a leaf node
func leafUpdate(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys())
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-idx-1)
}

// part of the treeInsert(): KV insertion to an internal node
func nodeInsert(
  tree *Btree, new BNode, node BNode, idx uint16,
  key []byte, val []byte,
) {
	kptr := node.getPtr(idx)
	// recursive insertion to the kid node
  knode := treeInsert(tree, tree.get(kptr), key, val)
  // split the result
  nsplit, split := nodeSplit3(knode)
  // deallocate the kid node
  tree.del(kptr)
  // update the kid links
  nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
}

// nodeAppendKV appends a key-value pair to the node
func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
    // Set pointer
    new.setPtr(idx, ptr)

    // Set key-value pair
    pos := new.kvPos(idx)
    binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
    binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))
    copy(new[pos+4:], key)
    copy(new[pos+4+uint16(len(key)):], val)

    // Update offset
    new.setOffset(idx+1, new.getOffset(idx)+4+uint16((len(key)+len(val))))
}

// nodeAppendRange appends a range of nodes from the old BNode to the new BNode
func nodeAppendRange(new BNode, old BNode, newIdx uint16, oldIdx uint16, n uint16) {
    Assert(newIdx+n <= new.nkeys() && oldIdx+n <= old.nkeys())
    if n == 0 {
        return
    }

    // Copy pointers
    copy(new[HEADER+8*newIdx:], old[HEADER+8*oldIdx:HEADER+8*(oldIdx+n)])

    // Copy offsets
    oldStart := offsetPos(old, oldIdx+1)
    newStart := offsetPos(new, newIdx+1)
    copy(new[newStart:], old[oldStart:oldStart+2*n])

    // Copy key-values
    oldPos := old.kvPos(oldIdx)
    newPos := new.kvPos(newIdx)
    length := old.kvPos(oldIdx+n) - oldPos
    copy(new[newPos:], old[oldPos:oldPos+length])

    // Update offset
    if newIdx > 0 {
        newOffset := new.getOffset(newIdx) + length
        new.setOffset(newIdx+n, newOffset)
    }
}

// nodeReplaceKidN replaces a child node with multiple new child nodes
func nodeReplaceKidN(
    tree *Btree, new BNode, old BNode, idx uint16, kids ...BNode,
) {
    inc := uint16(len(kids))
    new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
    nodeAppendRange(new, old, 0, 0, idx)
    for i, node := range kids {
        nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
    }
    nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}

// nodeSplit2 splits an old BNode into two new BNodes (left and right)
func nodeSplit2(left BNode, right BNode, old BNode) {
    Assert(old.nbytes() <= 2*BTREE_PAGE_SIZE)
    left.setHeader(old.btype(), 0)
    right.setHeader(old.btype(), 0)

    for i := uint16(0); i < old.nkeys(); i++ {
        if left.nbytes()+old.getOffset(i+1)-old.getOffset(i) <= BTREE_PAGE_SIZE {
            nodeAppendKV(left, left.nkeys(), old.getPtr(i), old.getKey(i), old.getVal(i))
        } else {
            nodeAppendKV(right, right.nkeys(), old.getPtr(i), old.getKey(i), old.getVal(i))
        }
    }
}

// nodeSplit3 splits an old BNode into up to three new BNodes
func nodeSplit3(old BNode) (uint16, [3]BNode) {
    if old.nbytes() <= BTREE_PAGE_SIZE {
        old = old[:BTREE_PAGE_SIZE]
        return 1, [3]BNode{old}
    }
    left := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
    right := BNode(make([]byte, BTREE_PAGE_SIZE))
    nodeSplit2(left, right, old)
    if left.nbytes() <= BTREE_PAGE_SIZE {
        left = left[:BTREE_PAGE_SIZE]
        return 2, [3]BNode{left, right}
    }
    leftleft := BNode(make([]byte, BTREE_PAGE_SIZE))
    middle := BNode(make([]byte, BTREE_PAGE_SIZE))
    nodeSplit2(leftleft, middle, left)
    Assert(leftleft.nbytes() <= BTREE_PAGE_SIZE)
    return 3, [3]BNode{leftleft, middle, right}
}

func treeInsert(tree *Btree, node BNode, key []byte, val []byte) BNode {
	new := BNode(make([]byte, 2*BTREE_PAGE_SIZE))

	idx := nodeLookupLE(node, key)

	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			leafUpdate(new, node, idx, key, val)
		} else {
			leafInsert(new, node, idx+1, key, val)
		}
	case BNODE_NODE:
		nodeInsert(tree, new, node, idx, key, val)
	default:
		panic("invalid node type")
	}

	return new
}
