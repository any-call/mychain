package mychain

import (
	"strings"
	"testing"
)

func TestEthChain_IsValidAddress(t *testing.T) {
	b := ImpEth().IsValidAddress("0x4d9cf1320060d198579220548877da4885e6497c")
	t.Log("check result :", b)
	t.Log(strings.ToLower("1CE24Ad9908A0964acC91b8EdbD104DD6F9FFAC4"))
	t.Log(strings.ToLower("0017da119Ff092F6c3019F6490385E921067f657"))
}
