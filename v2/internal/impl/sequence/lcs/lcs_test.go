package lcs

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/hbollon/go-edlib/v2/algorithms"
)

func TestLCS(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"AB/empty", args{"AB", ""}, 0},
		{"empty/AB", args{"", "AB"}, 0},
		{"AB/AB", args{"AB", "AB"}, 2},
		{"ABCD/ACBAD", args{"ABCD", "ACBAD"}, 3},
		{"ABCDGH/AEDFHR", args{"ABCDGH", "AEDFHR"}, 3},
		{"AGGTAB/GXTXAYB", args{"AGGTAB", "GXTXAYB"}, 4},
		{"XMJYAUZ/MZJAWXU", args{"XMJYAUZ", "MZJAWXU"}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if got := calc.Distance(tt.args.str1, tt.args.str2); got != tt.want {
				t.Errorf("LCS Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLCSBacktrack(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"AB/empty", args{"AB", ""}, "", true},
		{"empty/AB", args{"", "AB"}, "", true},
		{"AB/AB", args{"AB", "AB"}, "AB", false},
		{"ABCD/ACBAD", args{"ABCD", "ACBAD"}, "ABD", false},
		{"ABCDGH/AEDFHR", args{"ABCDGH", "AEDFHR"}, "ADH", false},
		{"AGGTAB/GXTXAYB", args{"AGGTAB", "GXTXAYB"}, "GTAB", false},
		{"XMJYAUZ/MZJAWXU", args{"XMJYAUZ", "MZJAWXU"}, "MJAU", false},
		{"你好先生/你好夫人", args{"你好先生", "你好夫人"}, "你好", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			got, err := calc.Backtrack(tt.args.str1, tt.args.str2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Backtrack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Backtrack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLCSBacktrackAll(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"AB/empty", args{"AB", ""}, nil, true},
		{"empty/AB", args{"", "AB"}, nil, true},
		{"AB/AB", args{"AB", "AB"}, []string{"AB"}, false},
		{"ABCD/ACBAD", args{"ABCD", "ACBAD"}, []string{"ABD", "ACD"}, false},
		{"ABCDGH/AEDFHR", args{"ABCDGH", "AEDFHR"}, []string{"ADH"}, false},
		{"AGGTAB/GXTXAYB", args{"AGGTAB", "GXTXAYB"}, []string{"GTAB"}, false},
		{"XMJYAUZ/MZJAWXU", args{"XMJYAUZ", "MZJAWXU"}, []string{"MJAU"}, false},
		{"AZBYCWDX/ZAYBWCXD", args{"AZBYCWDX", "ZAYBWCXD"}, []string{"ABCD", "ABCX", "ABWD", "ABWX", "AYCD", "AYCX", "AYWD", "AYWX", "ZBCD", "ZBCX", "ZBWD", "ZBWX", "ZYCD", "ZYCX", "ZYWD", "ZYWX"}, false},
		{"AATCC/ACACG", args{"AATCC", "ACACG"}, []string{"AAC", "ACC"}, false},
		{"您好女士 你好吗？/先生 你好吗？", args{"您好女士 你好吗？", "先生 你好吗？"}, []string{" 你好吗？"}, false},
		{" 是ab是cde22f123g/222222是ab是cd123", args{" 是ab是cde22f123g", "222222是ab是cd123"}, []string{"是ab是cd123"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			got, err := calc.BacktrackAll(tt.args.str1, tt.args.str2)
			if (err != nil) != tt.wantErr {
				t.Errorf("BacktrackAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BacktrackAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLCSDiff(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"AB/empty", args{"AB", ""}, nil, true},
		{"empty/AB", args{"", "AB"}, nil, true},
		{"AB/AB", args{"AB", "AB"}, []string{"AB"}, false},
		{"computer/houseboat", args{"computer", "houseboat"}, []string{" h c o m p u s e b o a t e r", " + -   - -   + + + + +   - -"}, false},
		{"您好女士 你好吗？/先生 你好吗？", args{"您好女士 你好吗？", "先生 你好吗？"}, []string{" 先 生 您 好 女 士   你 好 吗 ？", " + + - - - -          "}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			got, err := calc.Diff(tt.args.str1, tt.args.str2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLCSEditDistance(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"AB/empty", args{"AB", ""}, 2},
		{"empty/AB", args{"", "AB"}, 2},
		{"No characters match", args{"abcd", "effgghh"}, 11},
		{"AB/AB", args{"AB", "AB"}, 0},
		{"CAT/CUT", args{"CAT", "CUT"}, 2},
		{"ACB/AB", args{"ACB", "AB"}, 1},
		{"ABC/ACD", args{"ABC", "ACD"}, 2},
		{"ABCD/ACBAD", args{"ABCD", "ACBAD"}, 3},
		{"ABCDGH/AEDFHR", args{"ABCDGH", "AEDFHR"}, 6},
		{"AGGTAB/GXTXAYB", args{"AGGTAB", "GXTXAYB"}, 5},
		{"XMJYAUZ/MZJAWXU", args{"XMJYAUZ", "MZJAWXU"}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc, err := New()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if got := calc.EditDistance(tt.args.str1, tt.args.str2); got != tt.want {
				t.Errorf("EditDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLCSCaseInsensitive(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Case sensitive (default)
	dist1 := calc.Distance("Hello", "hello")
	if dist1 == 5 {
		t.Errorf("Case-sensitive LCS should not return full length for different cases")
	}

	// Case insensitive
	if err := calc.SetCaseInsensitive(true); err != nil {
		t.Fatalf("SetCaseInsensitive() error = %v", err)
	}
	dist2 := calc.Distance("Hello", "hello")
	if dist2 != 5 {
		t.Errorf("Case-insensitive LCS = %v, want 5", dist2)
	}
}

func TestLCSSimilarity(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name string
		s1   string
		s2   string
		want float64
	}{
		{"Identical strings", "hello", "hello", 1.0},
		{"Empty strings", "", "", 1.0},
		{"Completely different", "abc", "xyz", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := calc.Similarity(tt.s1, tt.s2)
			if sim != tt.want {
				t.Errorf("Similarity() = %v, want %v", sim, tt.want)
			}
		})
	}
}

func TestLCSWithContext(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("Normal execution", func(t *testing.T) {
		ctx := context.Background()
		dist, err := calc.DistanceWithContext(ctx, "ABCD", "ACBAD")
		if err != nil {
			t.Errorf("DistanceWithContext() error = %v", err)
		}
		if dist != 3 {
			t.Errorf("DistanceWithContext() = %v, want 3", dist)
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := calc.DistanceWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("DistanceWithContext() should return error for cancelled context")
		}
	})

	t.Run("Timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(2 * time.Nanosecond)
		_, err := calc.DistanceWithContext(ctx, "hello", "world")
		if err == nil {
			t.Error("DistanceWithContext() should return error for timed out context")
		}
	})
}

func TestLCSProperties(t *testing.T) {
	calc, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if calc.Name() != "LCS" {
		t.Errorf("Name() = %v, want LCS", calc.Name())
	}

	props := calc.Properties()
	if props.Type != algorithms.SequenceType {
		t.Errorf("Properties().Type = %v, want %v", props.Type, algorithms.SequenceType)
	}
	if !props.SupportsUnicode {
		t.Error("Properties().SupportsUnicode should be true")
	}
}

func BenchmarkLCS(b *testing.B) {
	calc, _ := New()
	s1 := "AGGTAB"
	s2 := "GXTXAYB"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Distance(s1, s2)
	}
}

func BenchmarkLCSBacktrack(b *testing.B) {
	calc, _ := New()
	s1 := "AGGTAB"
	s2 := "GXTXAYB"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Backtrack(s1, s2)
	}
}
