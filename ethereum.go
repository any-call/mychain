package mychain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/any-call/gobase/util/mynet"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ethChain struct {
}

func ImpEth() ethChain {
	return ethChain{}
}

func (self ethChain) GetNowBlockNum(tout time.Duration) (info *EthBlockNum, err error) {
	if err := mynet.DoReq("GET",
		"https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			return true, tout, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				return json.Unmarshal(ret, &info)
			}
			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	return
}

func (self ethChain) GetBlockByNum(num int64, tout time.Duration) (info *EthBlock, err error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&boolean=true&tag=0x%X&apikey=AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P", num)
	if err := mynet.DoReq("GET", url,
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			return true, tout, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				return json.Unmarshal(ret, &info)
			}
			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	return
}

func (self ethChain) IsValidAddress(address string) bool {
	if strings.HasPrefix(address, "0x") == false {
		return false
	}

	// Remove "0x" prefix if present
	address = strings.TrimPrefix(address, "0x")

	// Address should be exactly 40 characters long after removing "0x" prefix
	if len(address) != 40 {
		return false
	}

	// Check if all characters are valid hexadecimal characters
	match, _ := regexp.MatchString("^[0-9a-fA-F]+$", address)
	if !match {
		return false
	}

	return true
}

func (self ethChain) GetETHBalance(address string) (float64, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P", address)
	return self.fetchAndConvert(url, 18)
}

func (self ethChain) GetUSDTBalance(address string) (float64, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress=%s&address=%s&tag=latest&apikey=AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P",
		ContractAddrERCUSDT, address)
	return self.fetchAndConvert(url, 6)
}

func (self ethChain) GetUSDCBalance(address string) (float64, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress=%s&address=%s&tag=latest&apikey=AJES32DY7H7V4PVVPD7YYCJJKP84C37G1P",
		ContractAddrERCUSDC, address)
	return self.fetchAndConvert(url, 6)
}

// 通用查询并转换为 float64
func (self ethChain) fetchAndConvert(url string, decimals int) (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	type balanceResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	var r balanceResp
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, err
	}
	if r.Status != "1" {
		return 0, errors.New("API error: " + r.Message)
	}

	// 使用 big.Int 处理大数字，再转为 float64
	b := new(big.Int)
	b.SetString(r.Result, 10)

	denom := new(big.Float).SetFloat64(1)
	for i := 0; i < decimals; i++ {
		denom.Mul(denom, big.NewFloat(10))
	}

	amount := new(big.Float).Quo(new(big.Float).SetInt(b), denom)

	f64, _ := amount.Float64()
	return f64, nil
}
