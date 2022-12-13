package main

import "fmt"

type Tire struct {
	Root *Node
	Size int
}

type Node struct {
	Char   string
	Count  int
	IsWord bool
	Next   map[rune]*Node
}

func main() {
	node := &Node{
		Char: "/",
		Next: make(map[rune]*Node, 0),
	}
	tire := &Tire{
		Root: node,
	}
	tire.Append("123")
	tire.Append("12")
	tire.Remove("123")
	fmt.Printf("共有[%d]个字符 \n", tire.Size)
	if tire.Find("12", true) {
		fmt.Println("能查到")
	} else {
		fmt.Println("不能查到")
	}
}

func (r *Tire) Append(s string) {
	if r.Find(s, true) {
		return
	}
	node := r.Root
	for _, c := range s {
		v, ok := node.Next[c]
		if !ok {
			v = &Node{
				Char:   string(c),
				IsWord: false,
				Count:  1,
				Next:   make(map[rune]*Node, 0),
			}
			node.Next[c] = v
			r.Size += 1
		} else {
			v.Count += 1
		}
		node = v
	}
	node.IsWord = true
}

func (r *Tire) Find(s string, isFullMatch bool) bool {
	node := r.Root
	for _, c := range s {
		v, ok := node.Next[c]
		if !ok {
			return false
		}
		node = v
	}
	return (!isFullMatch && !node.IsWord) || node.IsWord
}

func (r *Tire) Remove(s string) bool {
	if !r.Find(s, true) {
		return false
	}
	node := r.Root
	prev := r.Root
	for i, c := range s {
		node = node.Next[c]
		node.Count -= 1
		if node.Count == 0 {
			r.Size -= (len(s) - i)
			break
		}
		prev = node
	}
	if node.Count == 0 {
		delete(prev.Next, rune(node.Char[0]))
	} else {
		node.IsWord = false
	}
	return true
}
