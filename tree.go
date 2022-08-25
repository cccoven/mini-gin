package mini_gin

import "strings"

type node struct {
	pattern  string  // 待匹配路由，如 /p/:lang
	part     string  // 路由中的一部分，如 :lang
	children []*node // 子节点
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为 true
}

// matchChild 返回第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

// matchChildren 返回所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

// insert Trie 树插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 只有在最后一级路径时，才在这个节点上存储完整的路径名，如 /p/:lang/doc 中的 doc 节点的 pattern 为 /p/:lang/doc
		n.pattern = pattern
		return
	}

	// 将 parts 中的每一级路径名依次加入到树中
	// 每一级路径都是后面一级路径的父节点
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

// search Trie 树搜索节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 当查找到了带 * 的路径时，表示所有值都可以匹配
		if n.pattern == "" {
			return nil
		}

		return n
	}

	// 从根结点开始，递归地找到 parts 对应地路径上的最后一个节点
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
