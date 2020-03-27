package main

import (
	"fmt"
	"log"

	"github.com/ecc1/rfm69"
)

func main() {
	r := rfm69.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	r.Reset()
	dumpRegs(r)

	fmt.Printf("\nTesting individual writes\n")
	hw := r.Hardware()
	fmt.Printf("source: %02X %02X %02X\n", 0x44, 0x55, 0x66)
	hw.WriteRegister(rfm69.RegSyncValue1, 0x44)
	hw.WriteRegister(rfm69.RegSyncValue2, 0x55)
	hw.WriteRegister(rfm69.RegSyncValue3, 0x66)
	readRegs(r)

	r.Reset()
	fmt.Printf("\nTesting burst writes\n")
	data := []byte{0x77, 0x88, 0x99}
	fmt.Printf("source: %02X %02X %02X\n", data[0], data[1], data[2])
	hw.WriteBurst(rfm69.RegSyncValue1, data)
	readRegs(r)
}

func dumpRegs(r *rfm69.Radio) {
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("\nConfiguration registers:\n")
	regs := r.ReadConfiguration().Bytes()
	resetValue := rfm69.ResetRFConfiguration.Bytes()
	for i, v := range regs {
		fmt.Printf("%02X  %02X  %08b", rfm69.RegOpMode+i, v, v)
		r := resetValue[i]
		if v == r {
			fmt.Printf("\n")
		} else {
			fmt.Printf("  **** SHOULD BE %02X  %08b\n", r, r)
		}
	}
	testRegs := []byte{
		rfm69.RegTest,
		rfm69.RegTestLna,
		rfm69.RegTestPa1,
		rfm69.RegTestPa2,
		rfm69.RegTestDagc,
		rfm69.RegTestAfc,
	}
	for _, reg := range testRegs {
		v := r.Hardware().ReadRegister(reg)
		fmt.Printf("%02X  %02X  %08b\n", reg, v, v)
	}
}

func readRegs(r *rfm69.Radio) {
	hw := r.Hardware()
	x := hw.ReadRegister(rfm69.RegSyncValue1)
	y := hw.ReadRegister(rfm69.RegSyncValue2)
	z := hw.ReadRegister(rfm69.RegSyncValue3)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf("single: %02X %02X %02X\n", x, y, z)
	v := hw.ReadBurst(rfm69.RegSyncValue1, 3)
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	fmt.Printf(" burst: %02X %02X %02X\n", v[0], v[1], v[2])
}
