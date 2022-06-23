package table

import (
	"strings"
)

type Generator interface {
	NewTable(text string) EncodingTable
}

// Binary tree can have only 2 heritage
// Left and right are heirs (childs) nodes they are self typed so struct are recursive
// Pay attention that we use pointer to struct cuz elements may not have value and will be nil
type decodingTree struct {
	Value string
	Zero  *decodingTree
	One   *decodingTree
}

type EncodingTable map[rune]string

func (et EncodingTable) Decode(str string) string {
	dt := et.decodingTree()

	return dt.Decode(str)
}

func (et EncodingTable) decodingTree() decodingTree {
	res := decodingTree{}

	for ch, code := range et {
		res.add(code, ch)
	}

	return res
}

func (dt *decodingTree) Decode(str string) string {
	var buf strings.Builder

	// root of the tree
	currNode := dt

	// "01000(x)1101(b)10101010101"
	for _, ch := range str {
		if currNode.Value != "" {
			buf.WriteString(currNode.Value)
			currNode = dt
		}

		switch ch {
		case '0':
			currNode = currNode.Zero
		case '1':
			currNode = currNode.One
		}
	}

	if currNode.Value != "" {
		buf.WriteString(currNode.Value)
	}

	return buf.String()
}

// reciever should be pointer, cuz we will change it's content
func (dt *decodingTree) add(code string, value rune) {
	currNode := dt

	// 01000(0) <- value
	// we shoud iterate each number and put value in last node
	for _, num := range code {
		switch num {
		case '0':
			if currNode.Zero == nil {
				currNode.Zero = &decodingTree{}
			}

			currNode = currNode.Zero
		case '1':
			if currNode.One == nil {
				currNode.One = &decodingTree{}
			}

			currNode = currNode.One
		}
	}

	currNode.Value = string(value)
}
