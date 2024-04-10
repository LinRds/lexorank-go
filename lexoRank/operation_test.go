package lexoRank

import (
	"lexorank-go/numeralSystems"
	"reflect"
	"testing"
)

func Test_subtract(t *testing.T) {
	type args struct {
		sys   numeralSystems.NumeralSystem
		left  []int
		right []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			// 287-82=205
			name: "BiggerAndSameDigits",
			args: args{sys: &numeralSystems.System10{}, right: []int{2, 8}, left: []int{7, 8, 2}},
			want: []int{5, 0, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := subtract(tt.args.sys, tt.args.left, tt.args.right); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}
