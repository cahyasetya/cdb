package btree

import (
	"github.com/cahyasetya/cdb/util"
)

func init() {
	node1Max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VALUE_SIZE
	util.Assert(node1Max <= BTREE_PAGE_SIZE)
}
