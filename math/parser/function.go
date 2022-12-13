package parser

import "math"

// 函数
var CallFunc map[string]func(...Node) float64

// 默认函数
const DEF_FUNC = "__def"

func RegFunc() {
	CallFunc = map[string]func(...Node) float64{
		DEF_FUNC: func(n ...Node) float64 {
			if len(n) == 1 {
				return n[0].Evaluate()
			} else {
				return 0
			}
		},
		"min": func(n ...Node) float64 {
			return math.Min(n[0].Evaluate(), n[1].Evaluate())
		},
		"max": func(n ...Node) float64 {
			return math.Max(n[0].Evaluate(), n[1].Evaluate())
		},
		"floor": func(n ...Node) float64 {
			return math.Floor(n[0].Evaluate())
		},
		"round": func(n ...Node) float64 {
			return math.Round(n[0].Evaluate())
		},
	}
}
