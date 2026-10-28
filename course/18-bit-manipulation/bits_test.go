package bits

import (
	mbits "math/bits"
	"testing"
)

func TestHelpers(t *testing.T) {
	x := uint(0b10110)
	if !Has(x, 1) || Has(x, 0) {
		t.Error("Has")
	}
	if Set(x, 0) != 0b10111 || Clear(x, 2) != 0b10010 || Toggle(x, 4) != 0b00110 {
		t.Error("Set/Clear/Toggle")
	}
	if DropLowest(x) != 0b10100 || LowestBit(x) != 0b10 {
		t.Error("DropLowest/LowestBit")
	}
	if PopCount(x) != 3 || LowMask(3) != 0b111 || LowMask(0) != 0 {
		t.Error("PopCount/LowMask")
	}
	if LowMask(mbits.UintSize) != ^uint(0) {
		t.Error("LowMask(width) should be all ones")
	}
}
