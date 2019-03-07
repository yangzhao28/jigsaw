package listener

import "strings"

// SegmentFunc  function to split channel id
type SegmentFunc func(path string, start int) (segment string, next int)

func pathSegmenter(path string, start int) (segment string, next int) {
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}
	end := strings.IndexRune(path[start+1:], '.')
	if end == -1 {
		return path[start:], -1
	}
	return path[start : start+end+1], start + end + 1 + 1
}
