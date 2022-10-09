package main

import (
	"fmt"
	"strings"
)

type Player struct {
	Id    uint64
	Name  string
	X     uint
	Y     uint
	Model string // w、m、wm （Watcher、Marker）
}

type Aoi struct {
	Players      map[uint64]*Player
	PlayersX     map[uint]map[uint64]*Player
	PlayersY     map[uint]map[uint64]*Player
	VisibleRange uint
}

type Callback = func(p1, p2 *Player)

const (
	AOI_WATCHER = "w"
	AOI_MARKER  = "m"
)

const (
	// 地图尺寸
	MAP_ROWS = 20 // y
	MAP_COLS = 20 // x
)

func main() {
	aoi := &Aoi{
		Players:      make(map[uint64]*Player),
		PlayersX:     make(map[uint]map[uint64]*Player),
		PlayersY:     make(map[uint]map[uint64]*Player),
		VisibleRange: 5,
	}

	p1 := &Player{
		Id:    1,
		Name:  "pp",
		X:     0,
		Y:     0,
		Model: "wm",
	}
	p2 := &Player{
		Id:    2,
		Name:  "wl",
		X:     2,
		Y:     20,
		Model: "wm",
	}
	p3 := &Player{
		Id:    3,
		Name:  "sd",
		X:     0,
		Y:     0,
		Model: "wm",
	}
	aoi.Enter(p1, func(p1, p2 *Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Enter(p2, func(p1, p2 *Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Enter(p3, func(p1, p2 *Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Move(p2, 2, 10,
		func(p1, p2 *Player) {
			fmt.Printf("玩家[%s]移动坐标，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *Player) {
			fmt.Printf("玩家[%s]离开视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *Player) {
			fmt.Printf("玩家[%s]进入视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
	)
	aoi.Leave(p3, func(p1, p2 *Player) {
		fmt.Printf("玩家[%s]离开视野，通知玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Move(p2, 20, 10,
		func(p1, p2 *Player) {
			fmt.Printf("玩家[%s]移动坐标，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *Player) {
			fmt.Printf("玩家[%s]离开视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *Player) {
			fmt.Printf("玩家[%s]进入视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
	)
	aoi.Leave(p2, func(p1, p2 *Player) {
		fmt.Printf("玩家[%s]离开地图，通知玩家[%s] \n", p1.Name, p2.Name)
	})
}

func (r *Aoi) Enter(p *Player, f Callback) map[uint64]*Player {
	r.Players[p.Id] = p
	if _, ok := r.PlayersX[p.X]; !ok {
		r.PlayersX[p.X] = make(map[uint64]*Player)
	}
	r.PlayersX[p.X][p.Id] = p
	if _, ok := r.PlayersY[p.Y]; !ok {
		r.PlayersY[p.Y] = make(map[uint64]*Player)
	}
	r.PlayersY[p.Y][p.Id] = p
	fmt.Printf("玩家[%s]进入地图 x%d,y%d \n", p.Name, p.X, p.Y)
	// 如果玩家是被观察者，广播消息给视野内所有观察者
	if r.IsMarker(p) {
		r.Broadcast(p, f)
	}
	// 如果玩家是观察者，广播消息给视野内所有被观察者
	if r.IsWatcher(p) {
		return r.findNeighbors(p, AOI_MARKER)
	}
	// fmt.Printf("内存[aoi] %d", unsafe.Sizeof(r))
	return nil
}

func (r *Aoi) Move(p *Player, x, y uint, move, leave, enter Callback) []*Player {
	fmt.Printf("玩家[%s]移动坐标 x%d,y%d ->  x%d,y%d \n", p.Name, p.X, p.Y, x, y)
	// 获取当前坐标视野内的观察者、被观察者
	bWatchers := r.findNeighbors(p, AOI_WATCHER)
	bMarkers := r.findNeighbors(p, AOI_MARKER)
	// 移动
	p.X, p.Y = x, y
	// 获取移动后坐标视野内的观察者、被观察者
	aWatchers := r.findNeighbors(p, AOI_WATCHER)
	aMarkers := r.findNeighbors(p, AOI_MARKER)
	//
	if r.IsMarker(p) {
		// 离开对方视野
		for id, p1 := range bWatchers {
			if _, ok := aWatchers[id]; ok {
				continue
			}
			leave(p, p1)
		}
		// 进入对方视野
		for id, p1 := range aWatchers {
			if _, ok := bWatchers[id]; ok {
				move(p, p1)
			} else {
				enter(p, p1)
			}
		}
	}
	// 新的视野邻居
	players := []*Player{}
	if r.IsWatcher(p) {
		for id := range aMarkers {
			if p1, ok := bMarkers[id]; !ok {
				players = append(players, p1)
			}
		}
	}
	return players
}

func (r *Aoi) Leave(p *Player, f Callback) {
	delete(r.Players, p.Id)
	delete(r.PlayersX[p.X], p.Id)
	delete(r.PlayersY[p.Y], p.Id)
	// 如果玩家是被观察者，广播消息给视野内所有观察者
	if r.IsMarker(p) {
		r.Broadcast(p, f)
	}
}

func (r *Aoi) Broadcast(p *Player, f Callback) {
	players := r.findNeighbors(p, AOI_MARKER)
	for _, p1 := range players {
		f(p, p1)
	}
}

func (r *Aoi) Get(id uint64) *Player {
	return r.Players[id]
}

func (r *Aoi) IsMarker(p *Player) bool {
	return strings.Contains(p.Model, "m")
}

func (r *Aoi) IsWatcher(p *Player) bool {
	return strings.Contains(p.Model, "w")
}

func (r *Aoi) findNeighbors(p *Player, model string) map[uint64]*Player {
	// 地图边界
	xMin := int64(p.X - r.VisibleRange)
	if xMin < 0 {
		xMin = 0
	}
	xMax := p.X + r.VisibleRange
	if xMax > MAP_ROWS {
		xMax = MAP_ROWS
	}
	// 感兴趣的邻居
	neighbors := map[uint64]*Player{}
	for x := uint(xMin); x < uint(xMax); x++ {
		players, ok := r.PlayersX[x]
		if !ok {
			continue
		}
		for _, p1 := range players {
			if p1.Id == p.Id {
				continue
			}
			// 判断玩家aoi模型
			if model == "w" {
				if !r.IsWatcher(p1) {
					continue
				}
			} else {
				if !r.IsMarker(p1) {
					continue
				}
			}
			// 判断Y轴是否在视野内
			if abs(int(p.Y-p1.Y)) > int(r.VisibleRange) {
				continue
			}
			neighbors[p1.Id] = p1
		}
	}
	return neighbors
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}
