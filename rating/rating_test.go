package rating_test

import (
	"fmt"
	"testing"

	"github.com/mashiike/golicko/rating"
)

func TestRatingCompare(t *testing.T) {

	cases := []struct {
		isDiff     bool
		isWeeker   bool
		isStronger bool
		leftVal    float64
		leftDev    float64
		rightVal   float64
		rightDev   float64
		winProb    float64
	}{
		{
			isDiff:     false,
			isWeeker:   false,
			isStronger: false,
			leftVal:    1500.0,
			leftDev:    350.0,
			rightVal:   1600.0,
			rightDev:   350.0,
			winProb:    0.4233,
		},
		{
			isDiff:     false,
			isWeeker:   false,
			isStronger: false,
			leftVal:    1500.0,
			leftDev:    50.0,
			rightVal:   1600.0,
			rightDev:   50.0,
			winProb:    0.3631,
		},
		{
			isDiff:     true,
			isWeeker:   true,
			isStronger: false,
			leftVal:    1500.0,
			leftDev:    50.0,
			rightVal:   1700.0,
			rightDev:   50.0,
			winProb:    0.2453,
		},
		{
			isDiff:     true,
			isWeeker:   false,
			isStronger: true,
			leftVal:    1580.0,
			leftDev:    42.0,
			rightVal:   1420.0,
			rightDev:   42.0,
			winProb:    0.7119,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			left := rating.New(c.leftVal, c.leftDev, 0.06)
			right := rating.New(c.rightVal, c.rightDev, 0.06)

			t.Logf("left:%v", left)
			t.Logf("right:%v", right)
			got := left.IsDifferent(right)
			if got != c.isDiff {
				t.Errorf("unexpected IsDifferent result: %v", got)
			}

			got = left.IsWeeker(right)
			if got != c.isWeeker {
				t.Errorf("unexpected IsWeeker result: %v", got)
			}

			got = left.IsStronger(right)
			if got != c.isStronger {
				t.Errorf("unexpected IsStronger result: %v", got)
			}

			gotWin := left.WinProb(right)
			if gotWin != c.winProb {
				t.Errorf("unexpected WinProb result: %v", gotWin)
			}
		})
	}
}
