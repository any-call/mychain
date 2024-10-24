package mychain

import (
	"testing"
)

func TestTronChain_HexToTronAddress(t *testing.T) {
	tronAddr, err := ImpTron().HexToAddrStr("412c681e6dee9fe1bb764f70efa052d2458aa8f0c8")
	if err != nil {
		t.Errorf("hex to tron addr err:%v", err)
		return
	}

	t.Log("hex to tron addr is :", tronAddr, ImpTron().IsValidAddress(tronAddr))

	hex, err := ImpTron().AddrToHexStr(tronAddr)
	if err != nil {
		t.Errorf("add to hex err:%v", err)
		return
	}

	t.Log("tron add to hex :", hex)
}

func TestTronChain_IsValidAddress(t *testing.T) {
	t.Log(ImpTron().IsValidAddress("BiFvfVHU427V1Bj3TdoafTd7hpECnCjyS1"))

}

//func TestTronChain_CreateNewAccount(t *testing.T) {
//	// 生成一个使用 secp256k1 椭圆曲线的 ECDSA 私钥
//	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	//// 从私钥派生公钥
//	publicKey := privateKey.PublicKey
//
//	// 使用 TRON 的地址编码方式生成 TRON 地址
//	tronAddress := address.PubkeyToAddress(publicKey)
//	t.Log(tronAddress)
//	// 打印私钥和 TRON 地址
//	// Get the private key in bytes
//	//privateKeyBytes := privateKey.D.Bytes()
//	//
//	//// Convert the private key to hexadecimal format
//	//privateKeyHex := hex.EncodeToString(privateKeyBytes)
//	//t.Logf("Private Key: %s\n", privateKeyHex)
//	//t.Logf("生成的 TRON 地址: %s\n", tronAddress.String())
//}
