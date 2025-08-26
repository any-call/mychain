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
	ret, err := ImpTron("").GetAccAllTrxTransactions("TX4x5hbKZLcF3cY3L86YNpTtLTykBy9HbH", 100, time.Second)
	if err != nil {
		t.Error(err)
		return
	}

	for i, _ := range ret {
		//jb, _ := json.Marshal(ret[i])
		t.Logf("get transactions [%d]:from %s to %s :%v %s ;block:%d %s", i, ret[i].From, ret[i].To, ret[i].ToTrx(), "TRX", ret[i].BlockNumber, ret[i].ToTime().Format("2006-01-02 15:04:05"))
	}

}

func TestTronChain_GetAccountTransactions1(t *testing.T) {
	ret, err := ImpTron("").GetAccAllTrc20Transactions("TMQZvoqhaWFSJYHCJMmbjD13bpxieqgCwP", 100, time.Second)
	if err != nil {
		t.Error(err)
		return
	}

	for i, _ := range ret {
		//jb, _ := json.Marshal(ret[i])
		if ret[i].BlockNumber == 0 {
			var err error
			if ret[i].BlockNumber, err = ImpTron("").GetBlockNumber(ret[i].TxID); err != nil {
				t.Log("get block number err:", err)
			}
		}
		t.Logf("get transactions [%d]:from %s to %s :%v %s ;block:%d %s", i, ret[i].From, ret[i].To, ret[i].ToUsdt(), "usdt", ret[i].BlockNumber, ret[i].ToTime().Format("2006-01-02 15:04:05"))
	}

}

func TestTronChain_FreezeEnergy(t *testing.T) {
	ret, err := ImpTron("").CreateFreezeEnergyTrans("TNoP3HyZkfip2H88QkVyF5P3GSX9ax6gyT", "", 5, true)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ret is :", ret)
}

func TestTronChain_GetTotalNetworkRes(t *testing.T) {
	res, err := ImpTron("").GetTotalNetworkRes()
	if err != nil {
		t.Error(err)
		return
	}

	a, _ := res.CalcTrxForEnergy(10000)
	b, _ := res.CalcTrxForBandwidth(10000)
	t.Log("res is :", a, b)
}
