package mychain

import (
	"testing"
	"time"
)

func TestBtcChain_GetNowBlock(t *testing.T) {
	info, err := ImpBtcChain().GetNowBlock(time.Second * 10)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("latest block :", info)
}

func TestBtcChain_GetFullTrans(t *testing.T) {
	list, err := ImpBtcChain().GetFullTxID(898363, time.Second*10, time.Second)
	if err != nil {
		t.Error(err)
		return
	}

	for i, _ := range list {
		t.Logf("%d:%s", i, list[i])
	}
}

func TestBtcChain_GetTrans(t *testing.T) {
	list, err := ImpBtcChain().GetTrans("2db804356c23d012038ab1d7f4e9b1836f2f352ed1a285f6d7c0ebbd0aae9cf4", time.Second*10)
	if err != nil {
		t.Error(err)
		return
	}

	for i, _ := range list {
		t.Logf("from:%s to %s :%.8f %s", list[i].FromAddress, list[i].ToAddress, list[i].AmountBTC, list[i].Currency)
	}
}
