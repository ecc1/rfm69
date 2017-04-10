package rfm69

// Configuration for Raspberry Pi Zero W.

const (
	spiDevice    = "/dev/spidev0.0"
	spiSpeed     = 6000000 // Hz
	interruptPin = 23      // GPIO for receive interrupts
	resetPin     = 24      // GPIO for hardware reset
)
