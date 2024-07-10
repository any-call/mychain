package mychain

import "testing"

func TestEthChain_IsValidAddress(t *testing.T) {
	b := ImpEth().IsValidAddress("0x4d9cf1320060d198579220548877da4885e6497c")
	t.Log("check result :", b)
}
