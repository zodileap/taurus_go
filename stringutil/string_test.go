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

func TestIsUUID(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{
			input: "123e4567-e89b-12d3-a456-426614174000",
			valid: true,
		},
		{
			input: "123e4567-e89b-12d3-a456",
			valid: false,
		},
		{
			input: "invalid-uuid",
			valid: false,
		},
		{
			input: "123e4567e89b12d3a456426614174000",
			valid: true,
		},
	}

	for _, tc := range tests {
		err := IsUUID(tc.input)
		if tc.valid && err != nil {
			t.Errorf("IsUUID(%s) unexpected error: %v", tc.input, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("IsUUID(%s) = nil, want error", tc.input)
		}
	}
}

func TestToUpperFirst(t *testing.T) {
	data := []struct {
		input string
		want  string
	}{
		{
			input: "hello",
			want:  "Hello",
		},
		{
			input: "world",
			want:  "World",
		},
	}

	for _, d := range data {
		if got := ToUpperFirst(d.input, "", 1); got != d.want {
			t.Errorf("结果错误，ToUpperFirst(%s) = %s, want %s", d.input, got, d.want)
		}
	}
}

func TestNumberToLetters(t *testing.T) {
	data := []struct {
		input int
		want  string
	}{
		{
			input: 1,
			want:  "B",
		},
		{
			input: 2,
			want:  "C",
		},
		{
			input: 26,
			want:  "AA",
		},
		{
			input: 27,
			want:  "AB",
		},
	}

	for _, d := range data {
		if got := NumberToLetters(d.input); got != d.want {
			t.Errorf("结果错误，NumberToLetters(%d) = %s, want %s", d.input, got, d.want)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	data := []struct {
		input string
		want  string
	}{
		{
			input: "hello world",
			want:  "helloWorld",
		},
		{
			input: "hello_world",
			want:  "helloWorld",
		},
		{
			input: "hello-world",
			want:  "helloWorld",
		},
		{
			input: "helloWorld",
			want:  "helloWorld",
		},
		{
			input: "HelloWorld",
			want:  "helloWorld",
		},
		{
			input: "byId",
			want:  "byId",
		},
	}

	for _, d := range data {
		if got := ToCamelCase(d.input); got != d.want {
			t.Errorf("结果错误，ToCamelCase(%s) = %s, want %s", d.input, got, d.want)
		}
	}
}
