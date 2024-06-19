package mychain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/any-call/gobase/util/mynet"
	"github.com/mr-tron/base58"
	"io"
	"net/http"
	"strconv"
	"time"
)

type tronChain struct {
}

func ImpTron() tronChain {
	return tronChain{}
}

func (self tronChain) GetNowBlock(timeout time.Duration) (info *TronBlock, err error) {
	if err := mynet.GetJson("https://api.trongrid.io/walletsolidity/getnowblock", nil, timeout, func(ret []byte, httpCode int) error {
		if httpCode == http.StatusOK {
			return json.Unmarshal(ret, &info)
		}
		return fmt.Errorf("http err code:%v", httpCode)
	}, nil); err != nil {
		return nil, err
	}

	return
}

func (self tronChain) GetBlock(num int32, includeTr bool, tm time.Duration) (info *TronBlock, err error) {
	if err = mynet.DoReq("POST", "https://api.trongrid.io/walletsolidity/getblock",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")

			if b, err := json.Marshal(map[string]any{
				"id_or_num": fmt.Sprintf("%d", num),
				"detail":    includeTr,
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, tm, nil
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

func (self tronChain) GetBlockByNum(num int32, tm time.Duration) (info *TronBlock, err error) {
	if err = mynet.DoReq("POST", "https://api.trongrid.io/walletsolidity/getblockbynum",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")

			if b, err := json.Marshal(map[string]int32{
				"num": num,
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, tm, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				//mylog.Debug("ret is :   ", string(ret))
				return json.Unmarshal(ret, &info)
			}
			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	return
}

func (self tronChain) GetBlockByLatestNum(num int32, tm time.Duration) (list []TronBlock, err error) {
	if err = mynet.DoReq("POST", "https://api.trongrid.io/walletsolidity/getblockbylatestnum",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")

			if b, err := json.Marshal(map[string]int32{
				"num": num,
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, tm, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				//mylog.Debug("ret is : ", string(ret))
				var tmp map[string][]TronBlock
				if err = json.Unmarshal(ret, &tmp); err != nil {
					return err
				}
				var ok bool
				if list, ok = tmp["block"]; ok {
					return nil
				}

				return fmt.Errorf("empty data")
			}
			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	return
}

func (self tronChain) GetBlockByLimitNext(startNum int32, endNum int32, tm time.Duration) (list []TronBlock, err error) {
	if err = mynet.DoReq("POST", "https://api.trongrid.io/walletsolidity/getblockbylimitnext",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")

			if b, err := json.Marshal(map[string]int32{
				"startNum": startNum,
				"endNum":   endNum,
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, tm, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				//mylog.Debug("ret is : ", string(ret))
				var tmp map[string][]TronBlock
				if err = json.Unmarshal(ret, &tmp); err != nil {
					return err
				}
				var ok bool
				if list, ok = tmp["block"]; ok {
					return nil
				}

				return fmt.Errorf("empty data")
			}
			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	return
}

func (self tronChain) HexToTronAddress(hexStr string) (string, error) {
	// 解码16进制字符串为字节切片
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	// 计算SHA256哈希并取前4个字节作为校验和
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	// Step 2: Take the first 4 bytes of the second hash, this is the checksum
	checksum := hash2[:4]

	// 将地址和校验和连接起来
	dataWithChecksum := append(data, checksum...)

	// 进行Base58编码
	encoded := base58.Encode(dataWithChecksum)
	return encoded, nil
}
