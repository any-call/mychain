package mychain

import (
	"testing"
	"time"
)

func TestBtcChain_GetNowBlock(t *testing.T) {
	blockNum, blockID, err := ImpBtcChain().GetLatestBlockByBlockcypher(time.Second * 10)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("latest block :", blockNum, blockID)
}

func TestBtcChain_GetFullTrans(t *testing.T) {
	list, err := ImpBtcChain().GetFullTxIDByBlockcypher("0000000000000000000239e82a345f6b7f87a816fa66ec888b44e87999015463", time.Second*10, time.Second)
	if err != nil {
		t.Error(err)
		return
	}

	for i, _ := range list {
		t.Logf("%d:%s", i, list[i])
	}
}

func TestBtcChain_GetBatchTrans(t *testing.T) {
	list, err := ImpBtcChain().GetBatchTrans([]string{
		//"97682f372c9a4080a6baa9b564c778390386de615e71df16d205996151fe3f8a",
		//"f8ca0d030836fe643a13c9bcc57ee59994511c926a00ac1457b17cfde8f58f3a",
		//"94e631edb4e553f5fb0645f7785d010d2c0a94ee8eb6f9e4e622a1cdce876504",
		"05d6c47c0b5392b93af37de198e1a2dec4c2490074bae6ed0667e53cf1760165",
		"ad154059ad06bca1f92161d2ceb73553759ec21f748ee1c213e010a362a945c0",
		"5c0fbefe94e9dadf93894924949b8f93bf7fa021b9671cab85d14c5732e44f75",
		//"eec4b294d7a91ba9c9cbbbc8a8dacb23dfddcadc4ebe25f84809f9532e33d814",
		//"bda3a5de79abddcd4584d5caa74a164dc1bfff5b3bdfde9f9a8a4507d216e10c",
		//"6fa04993ba9f2ab573f9256e8803dad47a4d47d19592175cfe08e43a5b792f7c",
		//"c94121fc88fb0b279130e503ae8379e5abb3b6ae643198507e87e2caa1ba235a",
	}, time.Second*10)
	if err != nil {
		t.Error(err)
		return
	}

	for i, _ := range list {
		t.Logf("%s \n from:%s to %s :%.8f %s on %s %d", list[i].TxID, list[i].FromAddress, list[i].ToAddress, list[i].AmountBTC, list[i].Currency, list[i].Time, list[i].TimeStamp())
	}
}
