package mychain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/any-call/gobase/util/mynet"
	"github.com/mr-tron/base58"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

func (self tronChain) HexToAddrStr(hexStr string) (string, error) {
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

func (self tronChain) AddrToHexStr(tronAddr string) (string, error) {
	//base58 decode
	decoded, err := base58.Decode(tronAddr)
	if err != nil {
		return "", err
	}

	// Extract the checksum (last 4 bytes)
	if len(decoded) < 4 {
		return "", errors.New("invalid TRON address")
	}
	checksum := decoded[len(decoded)-4:]

	// Remove the checksum to get the original data
	data := decoded[:len(decoded)-4]

	// Perform double SHA256 hashing to verify integrity
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])

	// Compare computed checksum with extracted checksum
	if !bytesEqual(hash2[:4], checksum) {
		return "", errors.New("checksum verification failed")
	}

	// Convert the data bytes to hexadecimal string
	hexStr := hex.EncodeToString(data)
	return hexStr, nil
}

func (self tronChain) IsValidAddress(address string) bool {
	// TRON address should be exactly 34 characters long
	if len(address) != 34 {
		return false
	}

	// TRON address should start with 'T'
	if !strings.HasPrefix(address, "T") {
		return false
	}

	// Check if all characters are valid Base58 characters
	match, _ := regexp.MatchString("^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$", address)
	if !match {
		return false
	}

	return true
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
