package btree

import "errors"

type Btree struct {
	root uint64

	get func(uint64) []byte
	new func([]byte) uint64
	del func(uint64)
}

var ErrKeyTooLarge = errors.New("key or value too large")

func (tree *Btree) Insert(key []byte, val []byte) error {
	// 1. check the length limit imposed by the node format
	if err := checkLimit(key, val); err != nil {
			return err // the only way for an update to fail
	}
	// 2. create the first node
	if tree.root == 0 {
			root := BNode(make([]byte, BTREE_PAGE_SIZE))
			root.setHeader(BNODE_LEAF, 2)
			// a dummy key, this makes the tree cover the whole key space.
			// thus a lookup can always find a containing node.
			nodeAppendKV(root, 0, 0, nil, nil)
			nodeAppendKV(root, 1, 0, key, val)
			tree.root = tree.new(root)
			return nil
	}
	// 3. insert the key
	node := treeInsert(tree, tree.get(tree.root), key, val)
	// 4. grow the tree if the root is split
	nsplit, split := nodeSplit3(node)
	tree.del(tree.root)
	if nsplit > 1 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeader(BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.new(knode), knode.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(split[0])
	}
	return nil
}

func checkLimit(key []byte, val []byte) error {
	// Calculate the total length of key and value
	totalLength := len(key) + len(val)

	// Check if the total length exceeds the maximum allowed
	if totalLength > BTREE_PAGE_SIZE-100 {
		return ErrKeyTooLarge
	}

	// Check if the key or value individually exceed uint16 max value
	if len(key) > 65535 || len(val) > 65535 {
		return ErrKeyTooLarge
	}

	return nil
}
