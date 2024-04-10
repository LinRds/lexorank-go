package lexoRank

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestNewRank(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    *Rank
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				str: "0|i00007:",
			},
			want:    &Rank{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRank(tt.args.str)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRank() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_between(t *testing.T) {
	type args struct {
		left  *Decimal
		right *Decimal
	}
	tests := []struct {
		name string
		args args
		want *Decimal
	}{
		{
			name: "TestValueBetweenMaxAndMin",
			args: args{
				left:  MinDecimal(),
				right: MaxDecimal(),
			},
		},
	}
	for _, tt := range tests {
		fmt.Println(tt.args.left.Format())
		fmt.Println(tt.args.right.Format())
		t.Run(tt.name, func(t *testing.T) {
			got := between(tt.args.left, tt.args.right)
			fmt.Println(got)
		})
	}
}

func TestMinRank(t *testing.T) {
	tests := []struct {
		name string
		want *Rank
	}{
		{
			name: "1",
			want: nil,
		},
	}
	for _, tt := range tests {
		got := MinRank()
		fmt.Println(tt.name, ": ", got.Value)
	}
}

func TestMaxRank(t *testing.T) {
	type args struct {
		bucket *Bucket
	}
	tests := []struct {
		name string
		args args
		want *Rank
	}{
		{
			name: "maxRank",
			args: args{bucket: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxRank(tt.args.bucket)
			fmt.Println(got.Format())
		})
	}
}

func TestRank_Between(t *testing.T) {
	tests := []struct {
		name  string
		cur   *Rank
		other *Rank
		want  *Rank
	}{
		{
			name:  "1",
			cur:   MinRank(),
			other: MaxRank(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cur.Between(tt.other)
			next := tt.cur.GenNext()
			fmt.Println(tt.cur.Format(), " ", tt.other.Format())
			fmt.Println("mid is: ", got.Format())
			fmt.Println("next is: ", next.Format())
		})
	}
}

func TestRank_Between2(t *testing.T) {
	tests := []struct {
		name  string
		cur   *Rank
		other *Rank
		want  string
	}{
		{
			name:  "1",
			cur:   NewRank("0|i0000f:"),
			other: NewRank("0|i0000v:"),
			want:  "0|i0000n:",
		},
		{
			name:  "1",
			cur:   NewRank("0|hzzzzb:"),
			other: NewRank("0|hzzzzb:i"),
			want:  "0|hzzzzb:9",
		},
		{
			name:  "1",
			cur:   NewRank("0|aaaaay:aaaac"),
			other: NewRank("0|aaaaaz:"),
			want:  "0|aaaaay:n5556",
		},
	}
	for _, tt := range tests {
		got := tt.cur.Between(tt.other)
		str := []string{tt.cur.Format(), tt.other.Format(), got.Format()}
		sort.Strings(str)
		fmt.Println(str)
	}

}
