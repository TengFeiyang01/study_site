package config

import (
	"Training/Study/internal/domain"
	"Training/Study/ioc"
	"context"
	"fmt"
	"log"
	"time"
)

// Hot100Problem Hot100题目结构
type Hot100Problem struct {
	ID         string
	Title      string
	TitleSlug  string
	Difficulty string
}

// Hot100Problems Hot100题目列表
var Hot100Problems = []Hot100Problem{
	{"1", "两数之和", "two-sum", "Easy"},
	{"2", "两数相加", "add-two-numbers", "Medium"},
	{"3", "无重复字符的最长子串", "longest-substring-without-repeating-characters", "Medium"},
	{"4", "寻找两个正序数组的中位数", "median-of-two-sorted-arrays", "Hard"},
	{"5", "最长回文子串", "longest-palindromic-substring", "Medium"},
	{"11", "盛最多水的容器", "container-with-most-water", "Medium"},
	{"15", "三数之和", "3sum", "Medium"},
	{"17", "电话号码的字母组合", "letter-combinations-of-a-phone-number", "Medium"},
	{"19", "删除链表的倒数第 N 个结点", "remove-nth-node-from-end-of-list", "Medium"},
	{"20", "有效的括号", "valid-parentheses", "Easy"},
	{"21", "合并两个有序链表", "merge-two-sorted-lists", "Easy"},
	{"22", "括号生成", "generate-parentheses", "Medium"},
	{"23", "合并K个升序链表", "merge-k-sorted-lists", "Hard"},
	{"25", "K 个一组翻转链表", "reverse-nodes-in-k-group", "Hard"},
	{"31", "下一个排列", "next-permutation", "Medium"},
	{"32", "最长有效括号", "longest-valid-parentheses", "Hard"},
	{"33", "搜索旋转排序数组", "search-in-rotated-sorted-array", "Medium"},
	{"34", "在排序数组中查找元素的第一个和最后一个位置", "find-first-and-last-position-of-element-in-sorted-array", "Medium"},
	{"35", "搜索插入位置", "search-insert-position", "Easy"},
	{"39", "组合总和", "combination-sum", "Medium"},
	{"42", "接雨水", "trapping-rain-water", "Hard"},
	{"46", "全排列", "permutations", "Medium"},
	{"48", "旋转图像", "rotate-image", "Medium"},
	{"49", "字母异位词分组", "group-anagrams", "Medium"},
	{"53", "最大子数组和", "maximum-subarray", "Medium"},
	{"55", "跳跃游戏", "jump-game", "Medium"},
	{"56", "合并区间", "merge-intervals", "Medium"},
	{"62", "不同路径", "unique-paths", "Medium"},
	{"64", "最小路径和", "minimum-path-sum", "Medium"},
	{"70", "爬楼梯", "climbing-stairs", "Easy"},
	{"72", "编辑距离", "edit-distance", "Hard"},
	{"75", "颜色分类", "sort-colors", "Medium"},
	{"76", "最小覆盖子串", "minimum-window-substring", "Hard"},
	{"78", "子集", "subsets", "Medium"},
	{"79", "单词搜索", "word-search", "Medium"},
	{"84", "柱状图中最大的矩形", "largest-rectangle-in-histogram", "Hard"},
	{"85", "最大矩形", "maximal-rectangle", "Hard"},
	{"94", "二叉树的中序遍历", "binary-tree-inorder-traversal", "Easy"},
	{"96", "不同的二叉搜索树", "unique-binary-search-trees", "Medium"},
	{"98", "验证二叉搜索树", "validate-binary-search-tree", "Medium"},
	{"101", "对称二叉树", "symmetric-tree", "Easy"},
	{"102", "二叉树的层序遍历", "binary-tree-level-order-traversal", "Medium"},
	{"104", "二叉树的最大深度", "maximum-depth-of-binary-tree", "Easy"},
	{"105", "从前序与中序遍历序列构造二叉树", "construct-binary-tree-from-preorder-and-inorder-traversal", "Medium"},
	{"114", "二叉树展开为链表", "flatten-binary-tree-to-linked-list", "Medium"},
	{"121", "买卖股票的最佳时机", "best-time-to-buy-and-sell-stock", "Easy"},
	{"124", "二叉树中的最大路径和", "binary-tree-maximum-path-sum", "Hard"},
	{"128", "最长连续序列", "longest-consecutive-sequence", "Medium"},
	{"136", "只出现一次的数字", "single-number", "Easy"},
	{"139", "单词拆分", "word-break", "Medium"},
	{"141", "环形链表", "linked-list-cycle", "Easy"},
	{"142", "环形链表 II", "linked-list-cycle-ii", "Medium"},
	{"146", "LRU 缓存", "lru-cache", "Medium"},
	{"148", "排序链表", "sort-list", "Medium"},
	{"152", "乘积最大子数组", "maximum-product-subarray", "Medium"},
	{"155", "最小栈", "min-stack", "Medium"},
	{"160", "相交链表", "intersection-of-two-linked-lists", "Easy"},
	{"169", "多数元素", "majority-element", "Easy"},
	{"198", "打家劫舍", "house-robber", "Medium"},
	{"200", "岛屿数量", "number-of-islands", "Medium"},
	{"206", "反转链表", "reverse-linked-list", "Easy"},
	{"207", "课程表", "course-schedule", "Medium"},
	{"208", "实现 Trie (前缀树)", "implement-trie-prefix-tree", "Medium"},
	{"215", "数组中的第K个最大元素", "kth-largest-element-in-an-array", "Medium"},
	{"221", "最大正方形", "maximal-square", "Medium"},
	{"226", "翻转二叉树", "invert-binary-tree", "Easy"},
	{"234", "回文链表", "palindrome-linked-list", "Easy"},
	{"236", "二叉树的最近公共祖先", "lowest-common-ancestor-of-a-binary-tree", "Medium"},
	{"238", "除自身以外数组的乘积", "product-of-array-except-self", "Medium"},
	{"239", "滑动窗口最大值", "sliding-window-maximum", "Hard"},
	{"240", "搜索二维矩阵 II", "search-a-2d-matrix-ii", "Medium"},
	{"253", "会议室 II", "meeting-rooms-ii", "Medium"},
	{"279", "完全平方数", "perfect-squares", "Medium"},
	{"283", "移动零", "move-zeroes", "Easy"},
	{"287", "寻找重复数", "find-the-duplicate-number", "Medium"},
	{"297", "二叉树的序列化与反序列化", "serialize-and-deserialize-binary-tree", "Hard"},
	{"300", "最长递增子序列", "longest-increasing-subsequence", "Medium"},
	{"301", "删除无效的括号", "remove-invalid-parentheses", "Hard"},
	{"309", "最佳买卖股票时机含冷冻期", "best-time-to-buy-and-sell-stock-with-cooldown", "Medium"},
	{"312", "戳气球", "burst-balloons", "Hard"},
	{"322", "零钱兑换", "coin-change", "Medium"},
	{"337", "打家劫舍 III", "house-robber-iii", "Medium"},
	{"338", "比特位计数", "counting-bits", "Easy"},
	{"347", "前 K 个高频元素", "top-k-frequent-elements", "Medium"},
	{"394", "字符串解码", "decode-string", "Medium"},
	{"399", "除法求值", "evaluate-division", "Medium"},
	{"406", "根据身高重建队列", "queue-reconstruction-by-height", "Medium"},
	{"416", "分割等和子集", "partition-equal-subset-sum", "Medium"},
	{"437", "路径总和 III", "path-sum-iii", "Medium"},
	{"438", "找到字符串中所有字母异位词", "find-all-anagrams-in-a-string", "Medium"},
	{"448", "找到所有数组中消失的数字", "find-all-numbers-disappeared-in-an-array", "Easy"},
	{"461", "汉明距离", "hamming-distance", "Easy"},
	{"494", "目标和", "target-sum", "Medium"},
	{"538", "把二叉搜索树转换为累加树", "convert-bst-to-greater-tree", "Medium"},
	{"543", "二叉树的直径", "diameter-of-binary-tree", "Easy"},
	{"560", "和为 K 的子数组", "subarray-sum-equals-k", "Medium"},
	{"581", "最短无序连续子数组", "shortest-unsorted-continuous-subarray", "Medium"},
	{"617", "合并二叉树", "merge-two-binary-trees", "Easy"},
	{"621", "任务调度器", "task-scheduler", "Medium"},
	{"647", "回文子串", "palindromic-substrings", "Medium"},
	{"739", "每日温度", "daily-temperatures", "Medium"},
	{"763", "划分字母区间", "partition-labels", "Medium"},
	{"1143", "最长公共子序列", "longest-common-subsequence", "Medium"},
}

// InsertHot100Problems 插入Hot100题目数据
func InsertHot100Problems(ctx context.Context, app *ioc.Application) {
	log.Printf("正在插入Hot100题目，共%d题", len(Hot100Problems))

	successCount := 0
	skipCount := 0

	for i, hotProblem := range Hot100Problems {
		// 检查是否已存在
		existing, err := app.CodingProblemRepo.FindBySourceId(ctx, hotProblem.ID)
		if err == nil && len(existing) > 0 {
			skipCount++
			continue
		}

		// 创建新题目
		problem := domain.CodingProblem{
			Title:       hotProblem.Title,
			Difficulty:  hotProblem.Difficulty,
			Tags:        []string{"Hot 100", "算法"},
			Source:      "leetcode",
			SourceId:    hotProblem.ID,
			SourceUrl:   fmt.Sprintf("https://leetcode.cn/problems/%s/", hotProblem.TitleSlug),
			StudyStatus: "not_started",
			IsHot100:    true,
			Ctime:       time.Now(),
			Utime:       time.Now(),
		}

		err = app.CodingProblemRepo.Create(ctx, problem)
		if err != nil {
			log.Printf("插入题目失败: %s, 错误: %v", hotProblem.Title, err)
			continue
		}

		successCount++

		// 每10题输出一次进度
		if (i+1)%10 == 0 || i == len(Hot100Problems)-1 {
			log.Printf("Hot100题目插入进度: %d/%d", i+1, len(Hot100Problems))
		}
	}

	log.Printf("✅ Hot100题目插入完成! 成功插入: %d, 跳过(已存在): %d", successCount, skipCount)
}
