package string_group

import (
	"slices"
	"testing"
	"unicode"
)

var example = "** 用户名User123在2023年10月15日购买了iPhone15手机，价格为¥6999元。 **Xx-()hello"

// 分组测试
func TestGroups(t *testing.T) {
	groups := SplitIntoGroups(example)
	for _, seg := range groups.GetSegmentsByType(GroupTypeChinese) {
		for _, c := range seg.String(example) {
			if !unicode.Is(unicode.Han, c) {
				t.Fatal("汉字子串包含非汉字字符")
			}
		}
	}
}

func TestCommonGroups(t *testing.T) {
	groups := SplitIntoGroups(example)
	for _, seg := range groups.GetSegmentsByType(GroupTypeCommon) {
		if !slices.Contains([]string{"**", "**Xx-()"}, seg.String(example)) {
			t.Fatal("分组错误")
		}
	}
}

// 筛选指定长度的子串
func TestMergeMultiGroupsWithContinuousIntervals(t *testing.T) {
	groups := SplitIntoGroups(example)
	allButChineseWithContinuous := groups.MergeMultiGroupsWithContinuousIntervals(4, 4, GroupTypeLetters, GroupTypeDigits, GroupTypeOthers)
	for _, seg := range allButChineseWithContinuous {
		if seg.End-seg.Start != 4 {
			t.Fatal("筛选指定长度的子串错误")
		}
	}
}

func TestFilterSegmentsByIntervals(t *testing.T) {
	segments := []StringSegment{
		{Start: 0, End: 5},
		{Start: 10, End: 20},
		{Start: 25, End: 30},
	}
	intervals := [][]int{
		{3, 12},  // 与第一个和第二个有交集
		{15, 18}, // 与第二个有交集
		{28, 35}, // 与第三个有交集
	}
	result := FilterSegmentsByIntervals(intervals, segments)
	want := []StringSegment{
		{Start: 3, End: 5},   // 第一个段与第一个区间交集
		{Start: 10, End: 12}, // 第二个段与第一个区间交集
		{Start: 15, End: 18}, // 第二个段与第二个区间交集
		{Start: 28, End: 30}, // 第三个段与第三个区间交集
	}
	if len(result) != len(want) {
		t.Fatalf("期望%d个结果，实际%d", len(want), len(result))
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("第%d个结果错误，期望%+v，实际%+v", i, want[i], result[i])
		}
	}

	// 测试无交集
	intervals2 := [][]int{{40, 50}}
	result2 := FilterSegmentsByIntervals(intervals2, segments)
	if len(result2) != 0 {
		t.Errorf("无交集时应返回空，实际%+v", result2)
	}

	// 测试完全包含
	intervals3 := [][]int{{0, 100}}
	result3 := FilterSegmentsByIntervals(intervals3, segments)
	if len(result3) != len(segments) {
		t.Errorf("完全包含时应返回所有原始段，实际%+v", result3)
	}
	for i := range segments {
		if result3[i] != segments[i] {
			t.Errorf("完全包含时第%d个段错误，期望%+v，实际%+v", i, segments[i], result3[i])
		}
	}
}
