package mychain

import (
	"testing"
)

func TestTronChain_HexToTronAddress(t *testing.T) {
	tronAddr, err := ImpTron().HexToAddrStr("412c681e6dee9fe1bb764f70efa052d2458aa8f0c8")
	if err != nil {
		t.Errorf("hex to tron addr err:%v", err)
		return
	}

	t.Log("hex to tron addr is :", tronAddr, ImpTron().IsValidAddress(tronAddr))

	hex, err := ImpTron().AddrToHexStr(tronAddr)
	if err != nil {
		t.Errorf("add to hex err:%v", err)
		return
	}

	t.Log("tron add to hex :", hex)
}

func TestTronChain_IsValidAddress(t *testing.T) {
	t.Log(ImpTron().IsValidAddress("BiFvfVHU427V1Bj3TdoafTd7hpECnCjyS1"))

}

func TestTronChain_CreateNewAccount(t *testing.T) {
	r1, r2, err := ImpTron().CreateNewAccount()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("create ok, r1:", r1)
	t.Log("create ok,r2", r2)
}
