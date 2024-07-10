package mychain

import "testing"

func TestTronChain_HexToTronAddress(t *testing.T) {
	tronAddr, err := ImpTron().HexToAddrStr("412c681e6dee9fe1bb764f70efa052d2458aa8f0c8")
	if err != nil {
		t.Errorf("hex to tron addr err:%v", err)
		return
	}

	t.Log("hex to tron addr is :", tronAddr)

	hex, err := ImpTron().AddrToHexStr(tronAddr)
	if err != nil {
		t.Errorf("add to hex err:%v", err)
		return
	}

	t.Log("tron add to hex :", hex)
}
