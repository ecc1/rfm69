![rfm69hcw module](images/rfm69.png)

The `rfm69` package provides a Go interface to an SPI-attached
[RFM69HCW module.](http://www.hoperf.com/rf_transceiver/modules/RFM69HCW.html)

An RFM69HCW module on a convenient breakout board
is [available here](https://www.adafruit.com/products/3070)
[or here.](https://www.sparkfun.com/products/12775)

The current version supports only OOK modulation (on-off keying)
and a proprietary packet format (variable-length, null-terminated).
Patches to support more general use are welcome.

**Note that an antenna must be attached before using the module.**

An edge-mounted female SMA connector can be attached conveniently.

![rmf69 antenna](images/rfm69_sma.png)
