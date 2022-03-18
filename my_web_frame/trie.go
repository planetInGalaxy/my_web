/*
 * @Description:
 1、实现动态路由最常用的数据结构，被称为前缀树(Trie树)：
 它的每一个节点的所有的子节点都拥有相同的前缀。
 url恰好是由/分隔的多段构成的，因此，每一段可以作为前缀树的一个节点。
 我们通过树结构查询，找到匹配的路由（模板），而如果中间某一层的节点都不满足条件，
 那么就说明没有匹配到的路由，查询结束。
 2、实现动态路由功能：
 : 参数匹配，只负责匹配一个分段
 例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。
 * 通配，负责匹配在此之后的所有分段
 例如 /static/*filepath，可以匹配/static/fav.ico，
 也可以匹配/static/js/jQuery.js，这种模式常用于静态服务器， 能够递归地匹配子路径。
 * @Author: Tjg
 * @Date: 2022-03-17 21:00:45
 * @LastEditTime: 2022-03-17 22:44:17
 * @LastEditors: Please set LastEditors
*/
package my_web_frame

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否模糊匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part /*|| child.isWild */ {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 实现优先精确匹配
		if part[0] == '*' {
			n.children = append(n.children, child)
		} else if part[0] == ':' {
			length := len(n.children)
			if length > 0 && n.children[length-1].part[0] == '*' {
				n.children = append(n.children[:length-1], append([]*node{child}, n.children[length-1])...)
			} else {
				n.children = append(n.children, child)
			}
		} else {
			length := len(n.children)
			if length > 1 && n.children[length-2].part[0] == ':' && n.children[length-1].part[0] == '*' {
				n.children = append(n.children[:length-2], append([]*node{child}, n.children[length-2:]...)...)
			} else if length > 0 && (n.children[length-1].part[0] == ':' || n.children[length-1].part[0] == '*') {
				n.children = append(n.children[:length-1], append([]*node{child}, n.children[length-1])...)
			} else {
				n.children = append(n.children, child)
			}
		}
	}
	child.insert(pattern, parts, height+1)
}

// 搜索节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

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
