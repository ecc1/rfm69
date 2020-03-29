// +build !walrus

package rfm69

// Configuration for Raspberry Pi Zero W with Adafruit RFM69HCW bonnet:
// https://www.adafruit.com/product/4072
// https://learn.adafruit.com/adafruit-radio-bonnets/pinouts

const (
	spiDevice    = "/dev/spidev0.1"
	spiSpeed     = 6000000 // Hz
	interruptPin = 22      // GPIO for receive interrupts (DIO0)
	resetPin     = 25      // GPIO for hardware reset
)
