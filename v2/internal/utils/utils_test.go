package utils

import (
	"runtime"
	"testing"
)

func TestMin(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Min between 2/4", args{2, 4}, 2},
		{"Min between 4/2", args{4, 2}, 2},
		{"Min between -25/-42", args{-25, -42}, -42},
		{"Min between 0/0", args{0, 0}, 0},
		{"Min between -1/1", args{-1, 1}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Max between 2/4", args{2, 4}, 4},
		{"Max between 4/2", args{4, 2}, 4},
		{"Max between -25/-42", args{-25, -42}, -25},
		{"Max between 0/0", args{0, 0}, 0},
		{"Max between -1/1", args{-1, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin3(t *testing.T) {
	type args struct {
		a int
		b int
		c int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Min3 between 1/2/3", args{1, 2, 3}, 1},
		{"Min3 between 3/2/1", args{3, 2, 1}, 1},
		{"Min3 between 2/1/3", args{2, 1, 3}, 1},
		{"Min3 between -1/-2/-3", args{-1, -2, -3}, -3},
		{"Min3 between 5/5/5", args{5, 5, 5}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min3(tt.args.a, tt.args.b, tt.args.c); got != tt.want {
				t.Errorf("Min3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin4(t *testing.T) {
	type args struct {
		a int
		b int
		c int
		d int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Min4 between 1/2/3/4", args{1, 2, 3, 4}, 1},
		{"Min4 between 4/3/2/1", args{4, 3, 2, 1}, 1},
		{"Min4 between 3/1/4/2", args{3, 1, 4, 2}, 1},
		{"Min4 between -1/-2/-3/-4", args{-1, -2, -3, -4}, -4},
		{"Min4 between 5/5/5/5", args{5, 5, 5, 5}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min4(tt.args.a, tt.args.b, tt.args.c, tt.args.d); got != tt.want {
				t.Errorf("Min4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
		a []rune
		b []rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Equal between toto/test", args{[]rune("toto"), []rune("test")}, false},
		{"Equal between toto/toto", args{[]rune("toto"), []rune("toto")}, true},
		{"Equal between Toto/toto", args{[]rune("Toto"), []rune("toto")}, false},
		{"Equal between empty/empty", args{[]rune(""), []rune("")}, true},
		{"Equal between empty/non-empty", args{[]rune(""), []rune("test")}, false},
		{"Equal between 🙂😄/🙂😄", args{[]rune("🙂😄"), []rune("🙂😄")}, true},
		{"Equal between 🙂😄/🙂😄🙂😄", args{[]rune("🙂😄"), []rune("🙂😄🙂😄")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []rune
	}{
		{"ToLower HELLO", args{"HELLO"}, []rune("hello")},
		{"ToLower Hello", args{"Hello"}, []rune("hello")},
		{"ToLower hello", args{"hello"}, []rune("hello")},
		{"ToLower empty", args{""}, []rune("")},
		{"ToLower ÀÉÈÊË", args{"ÀÉÈÊË"}, []rune("àéèêë")},
		{"ToLower 123ABC", args{"123ABC"}, []rune("123abc")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLower(tt.args.s)
			if !Equal(got, tt.want) {
				t.Errorf("ToLower() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Normalize single space", args{"hello world"}, "hello world"},
		{"Normalize multiple spaces", args{"hello    world"}, "hello world"},
		{"Normalize tabs and spaces", args{"hello\t\tworld"}, "hello world"},
		{"Normalize leading spaces", args{"  hello world"}, "hello world"},
		{"Normalize trailing spaces", args{"hello world  "}, "hello world"},
		{"Normalize newlines", args{"hello\nworld"}, "hello world"},
		{"Normalize empty", args{""}, ""},
		{"Normalize only spaces", args{"   "}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.args.s); got != tt.want {
				t.Errorf("Normalize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNumCPU(t *testing.T) {
	got := NumCPU()
	actualCPUs := runtime.NumCPU()
	
	if actualCPUs <= 32 {
		if got != actualCPUs {
			t.Errorf("NumCPU() = %v, want %v", got, actualCPUs)
		}
	} else {
		if got != 32 {
			t.Errorf("NumCPU() = %v, want 32 (capped)", got)
		}
	}
}

func TestOptimalWorkers(t *testing.T) {
	type args struct {
		inputSize    int
		minPerWorker int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Small input", args{10, 100}, 1},
		{"Medium input", args{1000, 100}, Min(10, NumCPU())},
		{"Large input exceeds max workers", args{10000, 100}, NumCPU()},
		{"Zero input", args{0, 100}, 1},
		{"Exact boundary", args{100, 100}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OptimalWorkers(tt.args.inputSize, tt.args.minPerWorker)
			if got < 1 {
				t.Errorf("OptimalWorkers() = %v, must be at least 1", got)
			}
			if got > NumCPU() {
				t.Errorf("OptimalWorkers() = %v, must not exceed NumCPU() = %v", got, NumCPU())
			}
			// For predictable tests, check specific values
			if tt.args.inputSize < tt.args.minPerWorker && got != 1 {
				t.Errorf("OptimalWorkers() with small input = %v, want 1", got)
			}
		})
	}
}

func TestClampFloat64(t *testing.T) {
	type args struct {
		v   float64
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"Value in range", args{0.5, 0.0, 1.0}, 0.5},
		{"Value below min", args{-0.5, 0.0, 1.0}, 0.0},
		{"Value above max", args{1.5, 0.0, 1.0}, 1.0},
		{"Value equals min", args{0.0, 0.0, 1.0}, 0.0},
		{"Value equals max", args{1.0, 0.0, 1.0}, 1.0},
		{"Negative range", args{-2.0, -5.0, -1.0}, -2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClampFloat64(tt.args.v, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("ClampFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunesToString(t *testing.T) {
	type args struct {
		runes []rune
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Simple string", args{[]rune("hello")}, "hello"},
		{"Empty runes", args{[]rune("")}, ""},
		{"Nil runes", args{nil}, ""},
		{"Unicode", args{[]rune("こんにちは")}, "こんにちは"},
		{"Emojis", args{[]rune("🙂😄🙂")}, "🙂😄🙂"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunesToString(tt.args.runes); got != tt.want {
				t.Errorf("RunesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToRunes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []rune
	}{
		{"Simple string", args{"hello"}, []rune("hello")},
		{"Empty string", args{""}, nil},
		{"Unicode", args{"こんにちは"}, []rune("こんにちは")},
		{"Emojis", args{"🙂😄🙂"}, []rune("🙂😄🙂")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringToRunes(tt.args.s)
			if tt.want == nil && got != nil {
				t.Errorf("StringToRunes() = %v, want nil", got)
			} else if tt.want != nil && !Equal(got, tt.want) {
				t.Errorf("StringToRunes() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestShingle(t *testing.T) {
	type args struct {
		str string
		k   int
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{"shingle 1", args{"Radiohead", 2}, map[string]int{"Ra": 1, "ad": 2, "di": 1, "ea": 1, "he": 1, "io": 1, "oh": 1}},
		{"shingle 1-1", args{"Radiohead", 3}, map[string]int{"Rad": 1, "adi": 1, "dio": 1, "ead": 1, "hea": 1, "ioh": 1, "ohe": 1}},
		{"shingle 2", args{"I love horror movies", 2}, map[string]int{" h": 1, " l": 1, " m": 1, "I ": 1, "e ": 1, "es": 1, "ho": 1, "ie": 1, "lo": 1, "mo": 1, "or": 2, "ov": 2, "r ": 1, "ro": 1, "rr": 1, "ve": 1, "vi": 1}},
		{"shingle 3", args{"私の名前はジョンです", 2}, map[string]int{"です": 1, "の名": 1, "はジ": 1, "ジョ": 1, "ョン": 1, "ンで": 1, "前は": 1, "名前": 1, "私の": 1}},
		{"shingle 4", args{"🙂😄🙂😄 😄🙂😄", 2}, map[string]int{" 😄": 1, "😄 ": 1, "😄🙂": 2, "🙂😄": 3}},
		{"shingle 5", args{"", 100}, make(map[string]int)},
		{"shingle 6", args{"hello", 0}, make(map[string]int)},
		{"shingle 7", args{"四畳半神話大系", 7}, map[string]int{"四畳半神話大系": 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Shingle(tt.args.str, tt.args.k)
			if len(got) != len(tt.want) {
				t.Errorf("Shingle() = %v, want %v", got, tt.want)
				return
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("Shingle() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestShingleSlice(t *testing.T) {
	type args struct {
		str string
		k   int
	}
	tests := []struct {
		name string
		args args
		want int // number of unique k-grams
	}{
		{"shingle slice 1", args{"Radiohead", 2}, 7},
		{"shingle slice 2", args{"I love horror movies", 2}, 17},
		{"shingle slice 3", args{"私の名前はジョンです", 2}, 9},
		{"shingle slice 4", args{"🙂😄🙂😄 😄🙂😄", 2}, 4},
		{"shingle slice 5", args{"", 100}, 0},
		{"shingle slice 6", args{"hello", 0}, 0},
		{"shingle slice 7", args{"四畳半神話大系", 7}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShingleSlice(tt.args.str, tt.args.k)
			if len(got) != tt.want {
				t.Errorf("ShingleSlice() returned %d unique k-grams, want %d", len(got), tt.want)
			}
		})
	}
}
