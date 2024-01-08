package stringutil

import "testing"

type TestLengthOfLongestSubstringStruct struct {
	input string
	want  int
}

func TestLengthOfLongestSubstring(t *testing.T) {
	data := []TestLengthOfLongestSubstringStruct{
		{
			input: "abcabcbb",
			want:  3,
		},
		{
			input: "bbbbb",
			want:  1,
		},
		{
			input: "pwwkew",
			want:  3,
		},
		{
			input: "12345678",
			want:  8,
		},
		{
			input: "12345678912",
			want:  9,
		},
		{
			input: "dvdf",
			want:  3,
		},
		{
			input: "abcabd",
			want:  4,
		},
		{
			input: "1234567813abc",
			want:  10,
		},
		{
			input: "1234567813",
			want:  8,
		},
		{
			input: "tmmzuxt",
			want:  5,
		},
		{
			input: "ggububgvfk",
			want:  6,
		},
		{
			input: "aabaab!bb",
			want:  3,
		},
		{
			input: "中文重中之重",
			want:  4,
		},
		{
			input: "nfpdmpi",
			want:  5,
		},
	}
	for _, d := range data {
		if got := LengthOfLongestSubstring(d.input); got != d.want {
			t.Errorf("结果错误，lengthOfLongestSubstring(%s) = %d, want %d", d.input, got, d.want)
		}
	}
}

func TestIsUUIDValid(t *testing.T) {
	testUUIDs := []string{
		"123e4567-e89b-12d3-a456",
		"invalid-uuid",
		"123e4567e89b12d3a456426614174000",
	}

	for _, uuid := range testUUIDs {
		if isUUIDValid(uuid) {
			t.Errorf("isUUIDValid(%s) = true, want false", uuid)
		}
	}
}
