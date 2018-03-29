package rfm69

import (
	"github.com/ecc1/radio"
)

var (
	// Ensure that *Radio implements the radio.Interface interface.
	_ radio.Interface = (*Radio)(nil)

	// Ensure that *hwFlavor implements the radio.HardwareFlavor interface.
	_ radio.HardwareFlavor = (*hwFlavor)(nil)
)
