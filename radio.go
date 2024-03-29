package rfm69

import (
	"log"
	"time"
)

const (
	debug         = false
	maxPacketSize = 110
	fifoSize      = 66

	// The fifoThreshold value should allow a maximum-sized packet to be
	// written in two bursts, but be large enough to avoid fifo underflow.
	fifoThreshold = 20

	// Approximate time for one byte to be transmitted, based on the data rate.
	byteDuration = time.Millisecond
)

func init() {
	if debug {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
}

// Send transmits the given packet.
func (r *Radio) Send(data []byte) {
	if r.Error() != nil {
		return
	}
	if len(data) > maxPacketSize {
		log.Panicf("attempting to send %d-byte packet", len(data))
	}
	if debug {
		log.Printf("sending %d-byte packet in %s state", len(data), r.State())
	}
	// Terminate packet with zero byte.
	copy(r.txPacket, data)
	r.txPacket[len(data)] = 0
	packet := r.txPacket[:len(data)+1]
	// Prepare for auto-transmit.
	// (Automode from/to sleep mode is not reliable.)
	r.clearFIFO()
	r.setMode(StandbyMode)
	r.hw.WriteRegister(RegAutoModes, EnterConditionFifoNotEmpty|ExitConditionFifoEmpty|IntermediateModeTx)
	r.transmit(packet)
	r.setMode(StandbyMode)
}

func (r *Radio) transmit(data []byte) {
	avail := fifoSize
	for r.Error() == nil {
		if avail > len(data) {
			avail = len(data)
		}
		if debug {
			log.Printf("writing %d bytes to TX FIFO\n", avail)
		}
		r.hw.WriteBurst(RegFifo, data[:avail])
		data = data[avail:]
		if len(data) == 0 {
			break
		}
		// Wait until there is room for at least fifoSize - fifoThreshold bytes in the FIFO.
		// Err on the short side here to avoid TXFIFO underflow.
		time.Sleep(fifoSize / 4 * byteDuration)
		for r.Error() == nil {
			if !r.fifoThresholdExceeded() {
				avail = fifoSize - fifoThreshold
				break
			}
		}
	}
	r.finishTX(avail)
}

func (r *Radio) finishTX(numBytes int) {
	time.Sleep(time.Duration(numBytes) * byteDuration)
	// Wait for automatic return to standby mode when FIFO is empty.
	for r.Error() == nil {
		s := r.mode()
		if s == StandbyMode {
			break
		}
		if debug || s != TransmitterMode {
			log.Printf("waiting for TX to finish in %s state", stateName(s))
		}
		time.Sleep(byteDuration)
	}
}

func (r *Radio) fifoEmpty() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoNotEmpty == 0
}

func (r *Radio) fifoFull() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoFull != 0
}

func (r *Radio) fifoThresholdExceeded() bool {
	return r.hw.ReadRegister(RegIrqFlags2)&FifoLevel != 0
}

func (r *Radio) clearFIFO() {
	r.hw.WriteRegister(RegIrqFlags2, FifoOverrun)
}

// Receive listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
func (r *Radio) Receive(timeout time.Duration) ([]byte, int) {
	if r.Error() != nil {
		return nil, 0
	}
	r.hw.WriteRegister(RegAutoModes, 0)
	r.setMode(ReceiverMode)
	defer r.setMode(SleepMode)
	if debug {
		log.Printf("waiting for interrupt in %s state", r.State())
	}
	r.hw.AwaitInterrupt(timeout)
	rssi := r.ReadRSSI()
	for r.Error() == nil {
		if r.fifoEmpty() {
			if timeout <= 0 {
				break
			}
			time.Sleep(byteDuration)
			timeout -= byteDuration
			continue
		}
		c := r.hw.ReadRegister(RegFifo)
		if r.Error() != nil {
			break
		}
		if c == 0 {
			// End of packet.
			return r.finishRX(rssi)
		}
		r.err = r.receiveBuffer.WriteByte(c)
	}
	return nil, rssi
}

func (r *Radio) finishRX(rssi int) ([]byte, int) {
	r.setMode(StandbyMode)
	size := r.receiveBuffer.Len()
	if size == 0 {
		return nil, rssi
	}
	p := make([]byte, size)
	_, r.err = r.receiveBuffer.Read(p)
	if r.Error() != nil {
		return nil, rssi
	}
	r.receiveBuffer.Reset()
	if debug {
		log.Printf("received %d-byte packet in %s state", size, r.State())
	}
	return p, rssi
}

// SendAndReceive transmits the given packet,
// then listens with the given timeout for an incoming packet.
// It returns the packet and the associated RSSI.
// (This could be further optimized by using an Automode to go directly
// from TX to RX, rather than returning to standby in between.)
func (r *Radio) SendAndReceive(data []byte, timeout time.Duration) ([]byte, int) {
	r.Send(data)
	if r.Error() != nil {
		return nil, 0
	}
	return r.Receive(timeout)
}
