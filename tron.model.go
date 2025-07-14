package mychain

import (
	"fmt"
	"math/big"
)

type (
	TronTrResult struct {
		ContractRet string `json:"contractRet"`
	}

	TronContractInfo struct {
		Parameter struct {
			Value struct {
				Data            string `json:"data"` //从中解析交易类型，交易对象，交易金额
				OwnerAddress    string `json:"owner_address"`
				ToAddress       string `json:"to_address"`
				ContractAddress string `json:"contract_address"`
				Amount          int    `json:"amount"`
			} `json:"value"`
			TypeUrl string `json:"type_url"`
		} `json:"parameter"`
		Type string `json:"type"`
	}

	TronTrRawInfo struct {
		Contract      []TronContractInfo `json:"contract"`
		RefBlockBytes string             `json:"ref_block_bytes"`
		RefBlockHash  string             `json:"ref_block_hash"`
		Expiration    int64              `json:"expiration"`
		Timestamp     int64              `json:"timestamp"`
		FeeLimit      int                `json:"fee_limit"`
	}

	TronTransactionInfo struct {
		TxID       string         `json:"txID"`
		Ret        []TronTrResult `json:"ret"`
		Signature  []string       `json:"signature"`
		RawDataHex string         `json:"raw_data_hex"`
		RawData    TronTrRawInfo  `json:"raw_data"`
	}

	TronBlock struct {
		BlockID     string `json:"blockID"` //区块编号，
		BlockHeader struct {
			RawData struct {
				Number         int64  `json:"number"`          //区块编号，也就是链上区块的高度
				TxTrieRoot     string `json:"txTrieRoot"`      //交易根节点的哈希值
				WitnessAddress string `json:"witness_address"` //产生该区块的超级代表的帐户地址
				ParentHash     string `json:"parentHash"`      //上一区块的ID
				Version        int    `json:"version"`         //标识链的版本
				Timestamp      int64  `json:"timestamp"`       //创建块的时间戳
			} `json:"raw_data"`
			WitnessSignature string `json:"witness_signature"` //产生区块的超级代表的签名
		} `json:"block_header"`
		Transactions []TronTransactionInfo `json:"transactions"` //打包进该区块的交易清单
	}

	//代币概览
	TokenOverview struct {
		TotalAssetInTrx float64     `json:"totalAssetInTrx"`
		Data            []TokenData `json:"data"`
		TotalTokenCount int         `json:"totalTokenCount"`
		TotalAssetInUsd float64     `json:"totalAssetInUsd"`
	}

	//代币详情
	TokenData struct {
		TokenId         string  `json:"tokenId"`
		TokenName       string  `json:"tokenName"`
		TokenDecimal    int     `json:"tokenDecimal"`
		TokenAbbr       string  `json:"tokenAbbr"`
		TokenCanShow    int     `json:"tokenCanShow"`
		TokenType       string  `json:"tokenType"`
		TokenLogo       string  `json:"tokenLogo"`
		Vip             bool    `json:"vip"`
		Balance         string  `json:"balance"`
		TokenPriceInTrx int     `json:"tokenPriceInTrx"`
		TokenPriceInUsd float64 `json:"tokenPriceInUsd"`
		AssetInTrx      float64 `json:"assetInTrx"`
		AssetInUsd      float64 `json:"assetInUsd"`
		Percent         int     `json:"percent"`
	}

	//账号详情
	AccountInfo struct {
		Balance int64               `json:"balance"` //单位Sun （1 TRX = 1,000,000 sun）
		TRC20   []map[string]string `json:"trc20"`   //合约地址（key）和余额（value），余额单位为最小计量单位
	}

	// TRC20Tx 表示一个 TRC20 交易记录
	TRC20Tx struct {
		TxID      string `json:"transaction_id"`
		From      string `json:"from"`
		To        string `json:"to"`
		TokenInfo struct {
			Symbol   string `json:"symbol"` //代币符号，常用于判断是否为 USDT
			Name     string `json:"name"`   //代币合约地址，最准确的识别方式
			Address  string `json:"address"`
			Decimals int    `json:"decimals"` //精度，用于 value 转换成实际金额
		} `json:"token_info"`
		Value string `json:"value"` // 注意是字符串，需自行转换
	}

	// TRXTx 代表一条 TRX 主币转账交易
	TRXTx struct {
		TxID      string `json:"txID"`
		Timestamp int64  `json:"block_timestamp"`
		From      string
		To        string
		Amount    int64 // 单位：sun（1 TRX = 1_000_000 sun）
	}
)

func (self *AccountInfo) GetTrxBalance() float64 {
	return float64(self.Balance) / 1e6 // sun 转 TRX
}

// 获取 TRC20 代币余额
func (self *AccountInfo) GetTrc20Balance() float64 {
	if self.TRC20 == nil || len(self.TRC20) == 0 {
		return 0
	}

	const usdtDecimals int = 6
	for _, token := range self.TRC20 {
		if balanceStr, exists := token[ContractAddrTronTrc20]; exists {
			var balance float64
			fmt.Sscanf(balanceStr, "%f", &balance)
			return balance / float64(pow(10, usdtDecimals))
		}
	}

	return 0
}

func (t *TRC20Tx) ToUsdt() float64 {
	if t.TokenInfo.Symbol != "USDT" {
		return 0.0
	}

	valueInt := new(big.Int)
	if _, ok := valueInt.SetString(t.Value, 10); !ok {
		return 0.0
	}

	decimals := t.TokenInfo.Decimals
	if decimals <= 0 {
		decimals = 6 // fallback
	}

	valueFloat := new(big.Float).SetInt(valueInt)
	scale := new(big.Float).SetFloat64(float64(1))
	for i := 0; i < decimals; i++ {
		scale.Mul(scale, big.NewFloat(10))
	}
	usdtValue := new(big.Float).Quo(valueFloat, scale)

	result, _ := usdtValue.Float64()
	return result
}

func (t *TRXTx) ToTrx() float64 {
	return float64(t.Amount) / 1_000_000
}

// 快速计算 10 的次方
func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
