// +build !walrus

package rfm69

// Configuration for Raspberry Pi Zero W with Adafruit Radiofruit RFM69HCW.

const (
	spiDevice    = "/dev/spidev0.1"
	spiSpeed     = 6000000 // Hz
	interruptPin = 22      // GPIO for receive interrupts
	resetPin     = 25      // GPIO for hardware reset
)
