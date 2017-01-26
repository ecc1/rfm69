package rfm69

import (
	"testing"
	"unsafe"
)

func TestRFConfiguration(t *testing.T) {
	have := int(unsafe.Sizeof(RFConfiguration{}))
	want := RegTemp2 - RegOpMode + 1
	if have != want {
		t.Errorf("Sizeof(RFConfiguration) == %d, want %d", have, want)
	}
}
