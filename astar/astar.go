package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

/*
地图从左上角开始，水平x 垂直y
  y
x 0,0 1,0 2,0
  0,1 1,1 2,1
  0,2 1,2 2,2
*/

// 节点
type Node struct {
	// 坐标
	X int
	Y int
	// 成本
	F int
	G int
	H int
	// 父节点
	Parent *Node
	// 类型
	Type int
	// 状态
	State int
}

type AStar struct {
	// 启发算法
	Heuristic func(node, end *Node) int
	// 地图大小
	Rows int // y
	Cols int // x
	// 地图节点
	nodes [][]*Node
	start *Node
	end   *Node
	// 开放、关闭列表
	openList  []*Node
	closeList []*Node
	// 相邻节点坐标
	neighborPos [][]int
}

// 移动成本
const (
	COST_STRAIGHT = 10
	COST_DIAGONAL = 14
)

// 节点类型
const (
	NODE_TYPE_NORMAL = iota
	NODE_TYPE_OBSTACLE
)

// 节点状态
const (
	NODE_STATE_CLOSED = iota - 1
	NODE_STATE_NORMAL
	NODE_STATE_OPENED
)

func main() {
	astar := &AStar{
		Rows:      5,
		Cols:      8,
		Heuristic: Diagonal,
	}
	// 5x8地图
	// 0是可移动的网格
	// 1是障碍网格
	mapData := [][]int{
		0: {0, 0, 1, 1, 0, 0, 0, 0},
		1: {0, 0, 0, 0, 1, 0, 0, 0},
		2: {0, 0, 0, 1, 1, 0, 0, 0},
		3: {0, 0, 0, 0, 1, 0, 0, 0},
		4: {0, 0, 0, 0, 0, 0, 0, 0},
	}
	astar.Init(mapData)
	fmt.Println("开始时间", time.Now().UnixNano())
	node := astar.FindPath(
		&Node{X: 0, Y: 0},
		&Node{X: 5, Y: 0},
	)
	fmt.Println("结束时间", time.Now().UnixNano())
	astar.print(node, mapData)
}

func (a *AStar) Init(mapData [][]int) {
	a.nodes = make([][]*Node, a.Cols)
	for i := 0; i < a.Cols; i++ {
		a.nodes[i] = make([]*Node, a.Rows)
	}
	for i := 0; i < len(mapData); i++ {
		for j := 0; j < len(mapData[i]); j++ {
			node := &Node{
				X:    j,
				Y:    i,
				Type: mapData[i][j],
			}
			a.nodes[j][i] = node
		}
	}
	// 如果不允许对角移动，去除对角坐标
	a.neighborPos = [][]int{
		{0, -1},  // 上
		{1, -1},  // 右上
		{1, 0},   // 右
		{1, 1},   // 右下
		{0, 1},   // 下
		{-1, 1},  // 左下
		{-1, 0},  // 左
		{-1, -1}, // 左上
	}
}

func (a *AStar) FindPath(start, end *Node) *Node {
	a.start = a.nodes[start.X][start.Y]
	a.end = a.nodes[end.X][end.Y]
	// 如果起止点是障碍物
	if !a.start.isWalkable() || !a.end.isWalkable() {
		fmt.Println("障碍物不可移动")
		return nil
	}
	// 先把开始节点放进开放列表
	a.openListAppend(start)
	for len(a.openList) > 0 {
		node := a.openListPop()
		// 判断当前节点是否是终点
		if a.isEnd(node) {
			return node
		}
		// 找开放列表的第一个节点的相邻节点
		neighbors := a.findNeighbors(node)
		for _, neighbor := range neighbors {
			// 是否在关闭列表
			if neighbor.isClosed() {
				continue
			}
			// 开始节点移动至当前节点的成本
			// 相邻节点的坐标x,y
			// 开始节点移动至相邻节点的成本
			g, x, y := node.G, neighbor.X, neighbor.Y
			// 判断移动方式是水平（或垂直）、对角，计算成本
			if x == node.X || y == node.Y {
				g += COST_STRAIGHT
			} else {
				g += COST_DIAGONAL
			}
			if !neighbor.isOpened() || g < neighbor.G {
				neighbor.G = g
				neighbor.H = a.Heuristic(neighbor, end)
				neighbor.F = neighbor.G + neighbor.H
				neighbor.Parent = node
				// 优化逻辑，相邻节点是否是终点
				// if a.isEnd(neighbor) {
				// 	return neighbor
				// }
				if !neighbor.isOpened() {
					a.openListAppend(neighbor)
				}
			}
		}
		// 当前节点放进关闭列表
		a.closeListAppend(node)
		// 更新开放列表顺序
		a.openListSort()
	}
	return nil
}

// 查找相邻节点位置
func (a *AStar) findNeighbors(node *Node) []*Node {
	neighbors := make([]*Node, 0)
	for _, v := range a.neighborPos {
		x, y := node.X+v[0], node.Y+v[1]
		// 检测节点是否非法
		if !a.isWalkable(x, y) {
			continue
		}
		neighbors = append(neighbors, a.nodes[x][y])
	}
	return neighbors
}

func (a *AStar) isEnd(node *Node) bool {
	return node.X == a.end.X && node.Y == a.end.Y
}

func (a *AStar) isWalkable(x, y int) bool {
	// 最小越界
	if x < 0 || y < 0 {
		return false
	}
	// 最大越界
	if x > a.Cols-1 || y > a.Rows-1 {
		return false
	}
	// 节点是否可行
	if !a.nodes[x][y].isWalkable() {
		return false
	}
	return true
}

func (node *Node) isWalkable() bool {
	return node.Type != NODE_TYPE_OBSTACLE
}

func (node *Node) isOpened() bool {
	return node.State == NODE_STATE_OPENED
}

func (node *Node) isClosed() bool {
	return node.State == NODE_STATE_CLOSED
}

func (a *AStar) openListAppend(node *Node) {
	node.State = NODE_STATE_OPENED
	a.openList = append(a.openList, node)
}

func (a *AStar) openListPop() *Node {
	s := a.openList
	if len(s) == 0 {
		return nil
	}
	v := s[0]
	s[0] = nil
	s = s[1:]
	a.openList = s
	return v
}

func (a *AStar) openListSort() {
	sort.Slice(a.openList, func(i, j int) bool {
		return a.openList[i].F < a.openList[j].F
	})
}

func (a *AStar) closeListAppend(node *Node) {
	node.State = NODE_STATE_CLOSED
	a.closeList = append(a.closeList, node)
}

func (a *AStar) print(node *Node, mapData [][]int) {
	fmt.Println("导航路径：")
	for node != nil {
		fmt.Printf("x,y: %d,%d cost: f%d h%d g%d \n", node.X, node.Y, node.F, node.H, node.G)
		a.nodes[node.X][node.Y].Type = 9
		node = node.Parent
	}
	fmt.Println("导航图：")
	for i := 0; i < len(mapData); i++ {
		for j := 0; j < len(mapData[i]); j++ {
			if a.nodes[j][i].Type == 9 {
				fmt.Print("* ")
			} else {
				fmt.Print(a.nodes[j][i].Type, " ")
			}
		}
		fmt.Print("\n")
	}
	fmt.Println("准备扫描节点：")
	for _, node := range a.openList {
		fmt.Printf("x,y: %d,%d \n", node.X, node.Y)
	}
	fmt.Println("已扫描节点：")
	for _, node := range a.closeList {
		fmt.Printf("x,y: %d,%d \n", node.X, node.Y)
	}
}

// 曼哈顿
func Manhattan(node, end *Node) int {
	x := abs(node.X - end.X)
	y := abs(node.Y - end.Y)
	return (x + y) * COST_STRAIGHT
}

// 对角线
func Diagonal(node, end *Node) int {
	x := abs(node.X - end.X)
	y := abs(node.Y - end.Y)
	min := min(x, y)
	return min*COST_DIAGONAL + abs(x-y)*COST_STRAIGHT
}

// 欧几里得
func Euclidean(node, end *Node) int {
	x := abs(node.X - end.X)
	y := abs(node.Y - end.Y)
	v := float64(x)*float64(x) + float64(y)*float64(y)
	return int(math.Sqrt(v) * COST_STRAIGHT)
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
