package mychain

import (
	"strings"
	"testing"
	"time"
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
	v, err := ImpEth("AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P").
		GetETHBalance("0x000ce2e2c276d93ef3df286b28292c2fb9f121e9")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("get eth is :%v", v)
	time.Sleep(time.Second)
	v, err = ImpEth("AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P").GetUSDTBalance("0xAbEd708f795B57Cb31D1ce793f244d58FF8c30a7")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("get usdt is :%v", v)
	time.Sleep(time.Second)
	v, err = ImpEth("AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P").GetUSDCBalance("0xAbEd708f795B57Cb31D1ce793f244d58FF8c30a7")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("get usdc is :%v", v)
}
