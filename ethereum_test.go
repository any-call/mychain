package mychain

import (
	"strings"
	"testing"
)

func TestEthChain_IsValidAddress(t *testing.T) {
	b := ImpEth("AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P").IsValidAddress("0x4d9cf1320060d198579220548877da4885e6497c")
	t.Log("check result :", b)
	t.Log(strings.ToLower("1CE24Ad9908A0964acC91b8EdbD104DD6F9FFAC4"))
	t.Log(strings.ToLower("0017da119Ff092F6c3019F6490385E921067f657"))
}

func TestEthChain_GetETHTransactions(t *testing.T) {
	v, err := ImpEth("AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P").
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
	v, err := ImpEth("AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P").GetUSDTTransactions("0x0025ca3839103424f84f462351d7e5d2ff1868de")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("get usdt is :%v", v)
}
