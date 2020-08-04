package privite_test

import (
	"fmt"
	"sort"
	"testing"
)

func TestRoundRobin(t *testing.T) {

	//服务器集合
	var S = []Node{Node{"A", 4}, Node{"B", 3}, Node{"C", 4}}

	var max = 0 //最大权重
	{
		for i, _ := range S {
			if max < S[i].weight {
				max = S[i].weight
			}
		}
	}
	var gcd = 0 //最大公约数
	var weight = make([]int, len(S))
	{
		for i, _ := range S {
			if S[i].weight > 0 {
				weight[i] = S[i].weight
			}
		}

		gcd = _gcd(weight)
	}
	var i = 0
	var cw = 0
	for j := 0; j < 20; j++ {

		for {
			i = (i + 1) % len(S)

			if i == 0 {
				cw = cw - gcd
				if cw <= 0 {
					cw = max
					if cw == 0 {
						return
					}
				}
			}

			if S[i].weight > cw {
				break
			}
		}
		fmt.Print(S[i].id)

	}
}

type Node struct {
	id     string
	weight int
}

func _gcd(items []int) int {

	sort.Sort(sort.Reverse(sort.IntSlice(items)))

	var gcd = items[0]
	for i := 0; i < len(items)-1; i++ {
		c := gcd % items[i+1]

		if c == 0 {
			gcd = items[i+1]
		} else {
			gcd = c
		}
	}

	return gcd
}
