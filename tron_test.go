package mychain

import (
	"testing"
	"time"
)

func TestTronChain_HexToTronAddress(t *testing.T) {
	tronAddr, err := ImpTron("").HexToAddrStr("a614f803b6fd780986a42c78ec9c7f77e6ded13c")
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

func TestTronChain_GetAccountBalanceTRC(t *testing.T) {
	ret, err := ImpTron("").GetAccountBalance("TCziceWb4sTNRgiZvz8ZuhDbGUsDtTJBvC")
	if err != nil {
		t.Error(err)
	}

	t.Logf("get balance trx:%.6f", ret.GetTrxBalance())
	t.Logf("get balance trc-20:%.6f U", ret.GetTrc20Balance())
}

func TestTronChain_GetAccountTransactions(t *testing.T) {
	ret, err := ImpTron("").GetAccountTransactions("TTYmKjuWUr5yJh24d66wg3SCdfFmBrsmgR")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("get transactions:", ret)
}
