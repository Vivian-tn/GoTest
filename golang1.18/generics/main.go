package main

import (
	"fmt"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"golang.org/x/exp/constraints"
	"sort"
	"strconv"
	"strings"
)

func main() {

	//tt
	a := []int{1, 2, 3}
	a = []int{4, 5, 6}
	fmt.Println(a)
	b := []string{"b", "c", "a"}
	fmt.Println(strings.Join(b, "、"))

	aa := []string{"Sbmuel", "Marc", "Samuel"}
	sort.Strings(aa)
	fmt.Println("=======排序", aa)
	// 切片去重
	names := lo.Uniq([]string{"Samuel", "Marc", "Samuel"})
	fmt.Println("切片去重", names)

	// 过滤掉切片中不符合规则的元素
	even := lo.Filter([]int{1, 2, 3, 4}, func(x int, _ int) bool {
		return x%2 == 0
	})
	fmt.Println("过滤掉切片中不符合规则的元素", even)

	// 类型转换 操作一种类型的 Map/Slice，并将其转化为另一种类型的 Map/Slice：
	list := []int64{1, 2, 3, 4}
	result := lo.Map(list, func(nbr int64, index int) string {
		return strconv.FormatInt(nbr*2, 10)
	})
	fmt.Println("类型转换", result)

	// 并发处理:并发起 goroutine 来处理所传入的数据，并在内部会处理好顺序问题。最终结果会以相同的顺序返回：
	s := lop.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
		return strconv.FormatInt(x*2+1, 10)
	})
	fmt.Println("并发处理", s)

	// 包含
	present := lo.Contains[int]([]int{0, 1, 2, 3, 4, 5}, 6)
	fmt.Println("包含", present)

	// 分组和切割
	groups := lo.GroupBy[int, int]([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
		return i % 3
	})
	fmt.Println("分组和切割", groups)

	// 三元运算
	result1 := lo.Ternary[string](true, "a", "b")
	result2 := lo.Ternary[string](false, "a", "b")
	fmt.Println("三元运算", result1, result2)

	minInt := min(1, 2)
	fmt.Println(minInt)

	minFloat := min(1.0, 2.0)
	fmt.Println(minFloat)

	minStr := min("a", "b")
	fmt.Println(minStr)
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
