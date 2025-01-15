package mychain

import (
	"testing"
	"time"
)

func TestTronChain_HexToTronAddress(t *testing.T) {
	tronAddr, err := ImpTron("").HexToAddrStr("412c681e6dee9fe1bb764f70efa052d2458aa8f0c8")
	if err != nil {
		t.Errorf("hex to tron addr err:%v", err)
		return
	}

	t.Log("hex to tron addr is :", tronAddr, ImpTron("").IsValidAddress(tronAddr))

	hex, err := ImpTron("").AddrToHexStr(tronAddr)
	if err != nil {
		t.Errorf("add to hex err:%v", err)
		return
	}

	t.Log("tron add to hex :", hex)
}

func TestTronChain_GetNowBlock(t *testing.T) {
	info, err := ImpTron("").
		//GetBlock(68755146, true, time.Second*10)
		GetNowBlock(time.Second * 5)
	if err != nil {
		t.Error("get err:", err)
		return
	}

	t.Log("new block :", info.BlockID)
}

func TestTronChain_CreateNewAccount(t *testing.T) {
	r1, r2, err := ImpTron("").CreateNewAccount()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("create ok, r1:", r1)
	t.Log("create ok,r2", r2)
}

func TestTronChain_GetAccountBalance(t *testing.T) {
	ret, err := ImpTron("").GetAccountBalanceTRX("TTYmKjuWUr5yJh24d66wg3SCdfFmBrsmgR")
	if err != nil {
		t.Error(err)
		return
	}

	ret2, err := ImpTron("").GetAccountBalanceTRC("TTYmKjuWUr5yJh24d66wg3SCdfFmBrsmgR")
	if err != nil {
		t.Error(err)
	}

	t.Log("get balance:", ret, ret2)
}

func TestTronChain_GetAccountBalanceTRC(t *testing.T) {
	ret, err := ImpTron("").GetAccountBalanceTRC("TTYmKjuWUr5yJh24d66wg3SCdfFmBrsmgR")
	if err != nil {
		t.Error(err)
	}

	t.Log("get balance:", ret)
}

func TestTronChain_GetAccountTransactions(t *testing.T) {
	ret, err := ImpTron("").GetAccountTransactions("TTYmKjuWUr5yJh24d66wg3SCdfFmBrsmgR")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("get transactions:", ret)
}
