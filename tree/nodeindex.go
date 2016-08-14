package tree

type NodeIndex interface {
	GetNode(name string) (*Node, bool)
}

type nodeIndex struct {
	index map[string]*Node
}

func NewNodeIndex(t *Tree) NodeIndex {

	nodeindex := &nodeIndex{
		index: make(map[string]*Node, 0),
	}

	nodes := t.Nodes()

	for _, n := range nodes {
		// tip
		if len(n.neigh) == 1 {
			nodeindex.index[n.Name()] = n
		}
	}

	return nodeindex
}

// returns the node associated to the name in argument
// it may correspond to a tip node or an internal node
// with a name
func (ni *nodeIndex) GetNode(name string) (*Node, bool) {
	n, ok := ni.index[name]
	return n, ok
}