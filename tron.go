package mychain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/any-call/gobase/util/mynet"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type tronChain struct {
	apiKey string
}

func ImpTron(token string) tronChain {
	return tronChain{apiKey: token}
}

func (self tronChain) GetNowBlock(timeout time.Duration) (info *TronBlock, err error) {
	if err = mynet.DoReq("GET", "https://api.trongrid.io/walletsolidity/getnowblock",
		func(r *http.Request) (isTls bool, tm time.Duration, err error) {
			r.Header.Add("Content-Type", "application/json")
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}
			return true, timeout, nil
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

func (self tronChain) GetBlock(num int32, includeTr bool, tm time.Duration) (info *TronBlock, err error) {
	if err = mynet.DoReq("POST", "https://api.trongrid.io/walletsolidity/getblock",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

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

			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

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

			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

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
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

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
	// Check if all characters are valid Base58 characters
	match, _ := regexp.MatchString("^T[1-9A-HJ-NP-Za-km-z]{33}$", address)
	if !match {
		return false
	}

	return true
}

func (self tronChain) CreateNewAccount() (adddress, privateInfo string, err error) {
	// 使用 ECDSA (secp256k1) 曲线生成私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("生成私钥失败: %v", err)
	}

	// 从私钥生成公钥
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// 对公钥进行 SHA256 哈希
	hashSHA256 := sha256.New()
	hashSHA256.Write(pubKey)
	pubHash := hashSHA256.Sum(nil)

	// 对 SHA256 哈希结果进行 RIPEMD160 哈希
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(pubHash)
	pubRipemd160 := ripemd160Hasher.Sum(nil)

	// 添加地址前缀 0x41（波场地址以 41 开头）
	rawAddress := append([]byte{0x41}, pubRipemd160...)

	// 计算地址的校验和：先 SHA256 再取前 4 字节
	checksum := sha256.Sum256(rawAddress)
	checksum = sha256.Sum256(checksum[:])
	address := append(rawAddress, checksum[:4]...)

	// 使用 Base58 编码地址
	encodedAddress := base58.Encode(address)

	// 将私钥转成十六进制
	privKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	return encodedAddress, privKeyHex, nil
}

func (self tronChain) GetAccountBalance(address string) (*AccountInfo, error) {
	url := "https://api.trongrid.io/v1/accounts/" + address

	var accInfo *AccountInfo = nil
	if err := mynet.DoReq("GET", url, func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
		r.Header.Set("Content-Type", "application/json")
		return true, time.Second * 10, nil
	}, func(ret []byte, httpCode int) error {
		if httpCode != http.StatusOK {
			return fmt.Errorf("%d:%s", httpCode, string(ret))
		}

		var retMap map[string]any
		if err := json.Unmarshal(ret, &retMap); err != nil {
			return err
		}

		if b, ok := retMap["success"]; ok {
			if b.(bool) { //成功
				if v, okk := retMap["data"]; okk {
					if tmpSlice, okk1 := v.([]any); okk1 {
						if len(tmpSlice) > 0 {
							jb, err := json.Marshal(tmpSlice[0])
							if err != nil {
								return err
							}

							if err := json.Unmarshal(jb, &accInfo); err != nil {
								return err
							}
							return nil
						} else {
							//说是链上没有数据，
							accInfo = &AccountInfo{}
							return nil
						}
					}
				}
			}
		}

		return fmt.Errorf("不是期望的数据格式")
	}, nil); err != nil {
		return nil, err
	}

	return accInfo, nil
}

func (self tronChain) GetBlockNumber(txID string) (int64, error) {
	url := fmt.Sprintf("https://apilist.tronscanapi.com/api/transaction-info?hash=%s", txID)

	var blockId int64
	if err := mynet.DoReq("GET", url, func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
		r.Header.Set("Content-Type", "application/json")
		return true, time.Second * 10, nil
	}, func(ret []byte, httpCode int) error {
		if httpCode != http.StatusOK {
			return fmt.Errorf("%d:%s", httpCode, string(ret))
		}

		var result struct {
			Block int64 `json:"block"`
		}

		if err := json.Unmarshal(ret, &result); err != nil {
			return err
		}

		if result.Block == 0 {
			return fmt.Errorf("transaction not found")
		}

		blockId = result.Block
		return nil
	}, nil); err != nil {
		return 0, err
	}

	return blockId, nil
}

// 查询账户 Trc20交易记录
func (self tronChain) GetAccAllTrc20Transactions(address string, limit int, freqTimeout time.Duration) ([]TRC20Tx, error) {
	if limit <= 0 || limit > 200 {
		limit = 200
	}

	var allTxs []TRC20Tx
	var fingerprint string

	// trc20Response 表示 API 的响应结构
	type trc20Response struct {
		Data []TRC20Tx `json:"data"`
		Meta struct {
			Fingerprint string `json:"fingerprint"`
		} `json:"meta"`
	}

	for {
		url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=%d", address, limit)
		if fingerprint != "" {
			url += "&fingerprint=" + fingerprint
		}

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("请求失败: %w", err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("状态码错误: %d\n%s", resp.StatusCode, string(body))
		}

		var res trc20Response
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, fmt.Errorf("JSON 解析失败: %w", err)
		}

		if len(res.Data) == 0 {
			break // 没有更多交易
		}

		allTxs = append(allTxs, res.Data...)

		if res.Meta.Fingerprint == "" {
			break
		}
		fingerprint = res.Meta.Fingerprint

		// 防止请求过快被限流
		time.Sleep(freqTimeout)
	}

	return allTxs, nil
}

// 查询账户 Trx交易记录
func (self tronChain) GetAccAllTrxTransactions(address string, limit int, freqTimeout time.Duration) ([]TRXTx, error) {
	if limit <= 0 || limit > 200 {
		limit = 200
	}

	var allTxs []TRXTx
	var fingerprint string

	// rawTransaction 是原始交易结构
	type rawTransaction struct {
		TxID        string `json:"txID"`
		Timestamp   int64  `json:"block_timestamp"`
		BlockNumber int64  `json:"blockNumber"`
		RawData     struct {
			Contract []struct {
				Type      string `json:"type"`
				Parameter struct {
					Value struct {
						OwnerAddress string `json:"owner_address"`
						ToAddress    string `json:"to_address"`
						Amount       int64  `json:"amount"`
					} `json:"value"`
				} `json:"parameter"`
			} `json:"contract"`
		} `json:"raw_data"`
	}

	// trxResponse 用于接收分页返回
	type trxResponse struct {
		Data []rawTransaction `json:"data"`
		Meta struct {
			Fingerprint string `json:"fingerprint"`
		} `json:"meta"`
	}

	for {
		url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions?limit=%d", address, limit)
		if fingerprint != "" {
			url += "&fingerprint=" + fingerprint
		}
		url += "&only_confirmed=true" //只取确认的交易记录

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("请求失败: %w", err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("状态码错误: %d\n%s", resp.StatusCode, string(body))
		}

		var res trxResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, fmt.Errorf("JSON解析失败: %w", err)
		}

		for _, raw := range res.Data {
			for _, c := range raw.RawData.Contract {
				if c.Type == "TransferContract" {
					from, _ := self.HexToAddrStr(c.Parameter.Value.OwnerAddress) //decodeBase58Address(c.Parameter.Value.OwnerAddress)
					to, _ := self.HexToAddrStr(c.Parameter.Value.ToAddress)      //decodeBase58Address(c.Parameter.Value.ToAddress)
					tx := TRXTx{
						TxID:        raw.TxID,
						Timestamp:   raw.Timestamp,
						BlockNumber: raw.BlockNumber,
						From:        from,
						To:          to,
						Amount:      c.Parameter.Value.Amount,
					}
					allTxs = append(allTxs, tx)
				}
			}
		}

		if res.Meta.Fingerprint == "" {
			break
		}
		fingerprint = res.Meta.Fingerprint
		time.Sleep(freqTimeout)
	}

	return allTxs, nil
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
