package main

import (
	"fmt"
	"sort"
	"strconv"
)

// 任务一： 136只出现一次的数字
func singleNumber(nums []int) int {
	// 使用map记录出现的数字。存在就删除
	m := make(map[int]int)
	for _, v := range nums {
		if m[v] == 1 {
			delete(m, v)
		} else {
			m[v] = 1
		}
	}
	for k := range m {
		return k
	}
	return 0
}

// 9. 回文数
func isPalindrome(x int) bool {
	// 转成字符串，两边向中间，同时比较，不相等则false
	var s = strconv.FormatInt(int64(x), 10)
	for k := range s {
		if s[k] != s[len(s)-k-1] {
			return false
		}
	}
	return true
}

// 有效的括号
func isValid(s string) bool {
	var kuoHao = map[string]string{
		"{": "}",
		"[": "]",
		"(": ")",
	}
	var zhan string
	for k, v := range s {
		if len(zhan) == 0 {
			zhan = string(v)
		} else {
			// 栈顶元素映射相等，则切片删除最后一个。否则入栈
			if string(s[k]) == kuoHao[string(zhan[len(zhan)-1])] {
				zhan = zhan[:len(zhan)-1]
			} else {
				zhan = zhan + string(v)
			}
		}

	}
	// 栈中没有元素则全部匹配 true
	if zhan == "" {
		return true
	}
	return false
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	// 以第一个作为基准比较
	s := strs[0]
	// 循环比对，不相等则返回，对应的切片
	for k, v := range s {
		for _, str := range strs {
			if k == len(str) || string(v) != string(str[k]) {
				return s[:k]
			}
		}
	}
	return s
}

// 基本值类型  加一
func plusOne(digits []int) []int {
	// 最后一项开始处理，是9进位1，再处理下一位，非9则结束
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] == 9 {
			digits[i] = 0
		} else {
			digits[i] = digits[i] + 1
			return digits
		}
	}
	// 若程序没有结束，处理进位
	digits[0] = 1
	digits = append(digits, 0)
	return digits
}

// 引用类型：切片
// 26. 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	m := make(map[int]int)
	var idx []int
	// 使用map记录重复项，没有重复的放到数组idx中
	for i := 0; i < len(nums); i++ {
		if _, ok := m[nums[i]]; !ok {
			m[nums[i]] = i
			idx = append(idx, nums[i])
		}
		// nums[i] = -200
	}
	// 数组赋值
	for i := 0; i < len(idx); i++ {
		nums[i] = idx[i]
	}
	return len(idx)
}

// 56. 合并区间：
func merge(intervals [][]int) [][]int {
	if len(intervals) <= 1 {
		return intervals
	}
	// 按左区间排序 Slice自定义排序规则
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	// 按顺序逐项合并
	var res [][]int
	temp := intervals[0]
	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] > temp[1] {
			res = append(res, temp)
			temp = intervals[i]
		} else {
			if intervals[i][1] > temp[1] {
				temp[1] = intervals[i][1]
			}
		}
		//处理最后一项
		if len(intervals)-1 == i {
			res = append(res, temp)
		}
	}

	return res
}

// 两数之和
func twoSum1(nums []int, target int) []int {
	for k, num := range nums {
		for l, r := k+1, len(nums)-1; l <= r; l++ {
			if num+nums[l] == target {
				return []int{k, l}
			}
		}
	}
	return nil
}

func twoSum(nums []int, target int) []int {
	// 使用map，key记录nums的值。value为下标。
	m := make(map[int]int)
	for k, num := range nums {
		dis := target - num
		// map 存在直接返回，否则放入map中
		if idx, ok := m[dis]; ok {
			return []int{idx, k}
		} else {
			m[num] = k
		}
	}
	return nil
}

func main() {

	sing := []int{3, 4, 5, 5, 2, 4, 3, 55, 55}
	fmt.Println("只出现一次：", singleNumber(sing))
	fmt.Println("是否回文数：", isPalindrome(123321))
	fmt.Println("有效括号：", isValid("{{}}{}"))
	ss := []string{"ab", "ab"}
	fmt.Println("最长公共前缀: ", longestCommonPrefix(ss))
	one := []int{0}
	fmt.Println("加1：", plusOne(one))
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Println("删除有序数组中的重复项:数组：", removeDuplicates(nums), nums)
	mer := [][]int{{1, 3}, {2, 6}, {15, 18}, {8, 10}}
	fmt.Println("合并区间：", merge(mer))
	sum := []int{3, 2, 4}
	fmt.Println("两数之和,下标是：", twoSum(sum, 6))
}
