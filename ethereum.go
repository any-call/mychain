package mychain

import (
	"encoding/json"
	"fmt"
	"github.com/any-call/gobase/util/mynet"
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
