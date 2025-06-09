package string_group

import "unicode"

// GroupType 表示分组类型
type GroupType int

// 分组类型常量
const (
	GroupTypeChinese GroupType = iota
	GroupTypeLetters
	GroupTypeDigits
	GroupTypeOthers
)

// StringGroups 存储分组后的字符串段
type StringGroups struct {
	Chinese []StringSegment // 汉字段
	Letters []StringSegment // 字母段
	Digits  []StringSegment // 数字段
	Others  []StringSegment // 其他非空白字符段
}

// StringSegment 表示原始字符串中的一个子串段
type StringSegment struct {
	Start int // 子串在原始字符串中的起始索引
	End   int // 子串在原始字符串中的结束索引（不包含）
}

// String 返回段的字符串表示
func (s StringSegment) String(original string) string {
	if s.Start >= 0 && s.End <= len(original) && s.Start < s.End {
		return original[s.Start:s.End]
	}
	return ""
}

// GetSegmentsByType 根据分组类型获取对应的分组
func (sg *StringGroups) GetSegmentsByType(groupType GroupType) []StringSegment {
	switch groupType {
	case GroupTypeChinese:
		return sg.Chinese
	case GroupTypeLetters:
		return sg.Letters
	case GroupTypeDigits:
		return sg.Digits
	case GroupTypeOthers:
		return sg.Others
	default:
		return nil
	}
}

// MergeMultiGroups 合并多个指定类型的分组
func (sg *StringGroups) MergeMultiGroups(types ...GroupType) []StringSegment {
	if len(types) == 0 {
		return nil
	}

	// 创建一个新的切片来存储合并后的结果
	var merged []StringSegment

	// 合并所有指定类型的分组
	for _, t := range types {
		group := sg.GetSegmentsByType(t)
		merged = append(merged, group...)
	}

	// 按照起始位置排序
	for i := 0; i < len(merged)-1; i++ {
		for j := i + 1; j < len(merged); j++ {
			if merged[i].Start > merged[j].Start {
				merged[i], merged[j] = merged[j], merged[i]
			}
		}
	}

	return merged
}

// MergeMultiGroupsWithContinuousIntervals 合并多个指定类型的分组，并连接连续的区间
func (sg *StringGroups) MergeMultiGroupsWithContinuousIntervals(maxlength, minLength int, types ...GroupType) []StringSegment {
	// 先合并多个分组
	merged := sg.MergeMultiGroups(types...)

	// 如果合并后的结果为空或只有一个元素，则直接返回
	if len(merged) <= 1 {
		return merged
	}

	// 连接连续的区间
	return connectContinuousIntervals(maxlength, minLength, merged)
}

// connectContinuousIntervals 连接连续的区间
func connectContinuousIntervals(maxlength, minLength int, segments []StringSegment) []StringSegment {
	if len(segments) <= 1 {
		return segments
	}

	// 创建一个新的切片来存储连接后的结果
	result := make([]StringSegment, 0, len(segments))

	// 当前处理的区间
	current := segments[0]

	// 遍历所有区间，连接连续的区间
	for i := 1; i < len(segments); i++ {
		// 如果当前区间的结束位置等于下一个区间的开始位置，则连接它们
		if current.End == segments[i].Start {
			current.End = segments[i].End
		} else {
			// 否则，将当前区间添加到结果中，并开始处理下一个区间
			if length := current.End - current.Start; length >= minLength && (maxlength == 0 || length <= maxlength) {
				result = append(result, current)
			}
			current = segments[i]
		}
	}

	// 添加最后一个处理的区间
	if length := current.End - current.Start; length >= minLength && (maxlength == 0 || length <= maxlength) {
		result = append(result, current)
	}

	return result
}

// 字符类型常量
const (
	typeUnknown = iota // 初始状态或空白字符
	typeChinese        // 汉字
	typeLetters        // 字母
	typeDigits         // 数字
	typeOther          // 其他非空白字符
)

// SplitIntoGroups 将字符串分为汉字、字母、数字和其他字符四组
func SplitIntoGroups(s string) StringGroups {
	// 预分配切片
	result := StringGroups{
		Chinese: make([]StringSegment, 0, len(s)/8+1), // 假设约1/8的字符是汉字
		Letters: make([]StringSegment, 0, len(s)/8+1), // 假设约1/8的字符是字母
		Digits:  make([]StringSegment, 0, len(s)/8+1), // 假设约1/8的字符是数字
		Others:  make([]StringSegment, 0, len(s)/8+1), // 假设约1/8的字符是其他非空白字符
	}

	start := 0                 // 当前段的起始位置
	currentType := typeUnknown // 当前正在处理的字符类型
	hasSegment := false        // 是否已经开始一个段

	// 遍历字符串中的每个字符
	for i, r := range s {
		// 判断字符类型
		var charType int
		if unicode.Is(unicode.Han, r) {
			charType = typeChinese
		} else if unicode.IsLetter(r) {
			charType = typeLetters
		} else if unicode.IsDigit(r) {
			charType = typeDigits
		} else {
			charType = typeOther
		}

		// 如果这是一个新段或字符类型发生变化
		if !hasSegment {
			// 开始新段
			start = i
			hasSegment = true
			currentType = charType
		} else if charType != currentType {
			// 字符类型变化，结束当前段并开始新段
			seg := StringSegment{Start: start, End: i}
			switch currentType {
			case typeChinese:
				result.Chinese = append(result.Chinese, seg)
			case typeLetters:
				result.Letters = append(result.Letters, seg)
			case typeDigits:
				result.Digits = append(result.Digits, seg)
			case typeOther:
				result.Others = append(result.Others, seg)
			}

			// 开始新段
			start = i
			currentType = charType
		}
	}

	// 处理最后一个分组
	if hasSegment {
		seg := StringSegment{Start: start, End: len(s)}
		switch currentType {
		case typeChinese:
			result.Chinese = append(result.Chinese, seg)
		case typeLetters:
			result.Letters = append(result.Letters, seg)
		case typeDigits:
			result.Digits = append(result.Digits, seg)
		case typeOther:
			result.Others = append(result.Others, seg)
		}
	}

	return result
}

// FilterSegmentsByIntervals 过滤并裁剪[]StringSegment，只保留与给定区间有交集的部分
// intervals: 区间数组，每个元素为[start, end)，与StringSegment定义一致
// segments: 待处理的字符串段
// 返回：只保留与区间有交集的部分，并裁剪为交集区间
func FilterSegmentsByIntervals(intervals [][]int, segments []StringSegment) []StringSegment {
	if len(intervals) == 0 || len(segments) == 0 {
		return nil
	}

	// 先对intervals按start排序，方便遍历
	type interval struct{ start, end int }
	ivls := make([]interval, len(intervals))
	for i, iv := range intervals {
		if len(iv) != 2 || iv[0] >= iv[1] {
			continue // 跳过非法区间
		}
		ivls[i] = interval{iv[0], iv[1]}
	}
	// 排序
	for i := 0; i < len(ivls)-1; i++ {
		for j := i + 1; j < len(ivls); j++ {
			if ivls[i].start > ivls[j].start {
				ivls[i], ivls[j] = ivls[j], ivls[i]
			}
		}
	}

	result := make([]StringSegment, 0)
	ivlIdx := 0
	for _, seg := range segments {
		// 跳过无效段
		if seg.Start >= seg.End {
			continue
		}
		// 区间和段都已排序，双指针遍历
		for ivlIdx < len(ivls) && ivls[ivlIdx].end <= seg.Start {
			ivlIdx++ // 当前区间在段左侧，无交集
		}
		tmpIdx := ivlIdx
		for tmpIdx < len(ivls) && ivls[tmpIdx].start < seg.End {
			// 有交集
			overlapStart := max(seg.Start, ivls[tmpIdx].start)
			overlapEnd := min(seg.End, ivls[tmpIdx].end)
			if overlapStart < overlapEnd {
				result = append(result, StringSegment{Start: overlapStart, End: overlapEnd})
			}
			tmpIdx++
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
