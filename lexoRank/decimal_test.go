package lexoRank

import (
	"fmt"
	"lexorank-go/numeralSystems"
	"testing"
)

func TestHalfDecimal(t *testing.T) {
	type args struct {
		sys numeralSystems.NumeralSystem
	}
	tests := []struct {
		name string
		args args
		want *Decimal
	}{
		{
			name: "1",
			args: args{sys: &numeralSystems.System36{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HalfDecimal(tt.args.sys)
			fmt.Println(got.Format())
		})
	}
}
