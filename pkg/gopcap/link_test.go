package gopcap

import (
	"bytes"
	"testing"
)

func TestEthernetFrameGood(t *testing.T) {
	// Pull some test data.
	data := []byte{
		0x00, 0x16, 0xE3, 0x19, 0x27, 0x15, 0x00, 0x04, 0x76, 0x96, 0x7B, 0xDA, 0x08, 0x00, 0x45, 0x00, 0x00, 0x52, 0x76, 0xED, 0x40, 0x00, 0x40, 0x06, 0x56, 0xCF,
		0xC0, 0xA8, 0x01, 0x02, 0xD4, 0xCC, 0xD6, 0x72, 0x0B, 0x20, 0x1A, 0x0B, 0x4D, 0xC8, 0x4E, 0xED, 0x54, 0xF1, 0x10, 0x72, 0x80, 0x18, 0x1F, 0x4B, 0x6D, 0x2E,
		0x00, 0x00, 0x01, 0x01, 0x08, 0x0A, 0x00, 0xD8, 0xEA, 0x48, 0x82, 0xE4, 0xDA, 0xB0, 0x49, 0x53, 0x4F, 0x4E, 0x20, 0x54, 0x68, 0x75, 0x6E, 0x66, 0x69, 0x73,
		0x63, 0x68, 0x20, 0x53, 0x6D, 0x69, 0x6C, 0x65, 0x79, 0x20, 0x53, 0x6D, 0x69, 0x6C, 0x65, 0x79, 0x47, 0x0A,
	}
	expectedDst := []byte{0x00, 0x16, 0xE3, 0x19, 0x27, 0x15}
	expectedSrc := []byte{0x00, 0x04, 0x76, 0x96, 0x7B, 0xDA}
	frame := new(EthernetFrame)
	err := frame.FromBytes(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if bytes.Compare(frame.MACDestination, expectedDst) != 0 {
		t.Errorf("Unexpected destination MAC: expected %v, got %v", frame.MACDestination, expectedDst)
	}
	if bytes.Compare(frame.MACSource, expectedSrc) != 0 {
		t.Errorf("Unexpected source MAC: expected %v, got %v.", frame.MACSource, expectedSrc)
	}
	if len(frame.VLANTag) != 0 {
		t.Errorf("Incorrectly received VLAN tag: %v", frame.VLANTag)
	}
	if frame.Length != 0 {
		t.Errorf("Incorrectly received length: %v", frame.Length)
	}
	if frame.EtherType != EtherType(2048) {
		t.Errorf("Unexpected EtherType: expected %v, got %v", 2048, frame.EtherType)
	}
}
