package mychain

import (
	"strings"
	"testing"
	"time"
)

func TestEthChain_IsValidAddress(t *testing.T) {
	b := ImpEth("").IsValidAddress("0x4d9cf1320060d198579220548877da4885e6497c")
	t.Log("check result :", b)
	t.Log(strings.ToLower("1CE24Ad9908A0964acC91b8EdbD104DD6F9FFAC4"))
	t.Log(strings.ToLower("0017da119Ff092F6c3019F6490385E921067f657"))
}

func TestEthChain_GetNowBlockNum(t *testing.T) {
	info, err := ImpEth("").GetNowBlockNum(time.Second * 30)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("info is :", info.ToNumber())
}

func TestEthChain_GetETHTransactions(t *testing.T) {
	v, err := ImpEth("").
		GetETHTransactions("0x0011f4d21657905d3f90945db12b2516c075d2e0")
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range v {
		t.Log(" block numer:", item.BlockNumber)
	}
}

func TestEthChain_GetETHBalance(t *testing.T) {
	v, err := ImpEth("").GetUSDTTransactions("0x0025ca3839103424f84f462351d7e5d2ff1868de")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("get usdt is :%v", v)
}

func TestEthChain_CreateAccount(t *testing.T) {
	addr, pub, priv, err := ImpEth("").CreateAccount()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("addr is :", addr)
	t.Log("pub is :", pub)
	t.Log("priv is :", priv)
}

func TestEthChain_PublicKeyFromPrivateKey(t *testing.T) {
	pub, err := ImpEth("").PublicKeyFromPrivateKey("5e5d2d274ca920f5b4eb09b428650deb7991149bb06da919878e4b02a5708503")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("pub is :", pub)
}
