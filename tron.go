package mychain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/any-call/gobase/frame/myctrl"
	"github.com/any-call/gobase/util/mynet"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/sha3"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	// 1) 使用 secp256k1 生成私钥（正确！！）
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	// 2) 提取未压缩公钥 0x04 + X + Y
	pub := crypto.FromECDSAPub(&privateKey.PublicKey)

	// 3) keccak256(pub[1:]) 取后 20 字节
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pub[1:])
	pubHash := hash.Sum(nil)
	addr20 := pubHash[12:] // 最后20字节

	// 4) Tron 前缀 0x41
	addr21 := append([]byte{0x41}, addr20...)

	// 5) Base58Check 校验
	h1 := sha256.Sum256(addr21)
	h2 := sha256.Sum256(h1[:])
	checksum := h2[:4]
	addressBytes := append(addr21, checksum...)

	address := base58.Encode(addressBytes)
	privHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	return address, privHex, nil
}

func (self tronChain) PrivKeyHexToTronAddress(privHex string) (string, error) {
	// 清理输入
	privHex = strings.TrimSpace(privHex)
	if strings.HasPrefix(privHex, "0x") || strings.HasPrefix(privHex, "0X") {
		privHex = privHex[2:]
	}

	// 私钥应为 32 字节 (64 hex chars)
	if len(privHex) != 64 {
		return "", errors.New("private key hex must be 64 hex characters (32 bytes)")
	}

	// decode hex -> bytes
	privBytes, err := hex.DecodeString(privHex)
	if err != nil {
		return "", fmt.Errorf("invalid hex private key: %w", err)
	}

	// 使用 go-ethereum 来转换为 ecdsa 私钥
	privKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return "", fmt.Errorf("ToECDSA failed: %w", err)
	}

	// 获取未压缩公钥字节 (65 bytes, 0x04 || X(32) || Y(32))
	pubBytes := crypto.FromECDSAPub(&privKey.PublicKey)
	if len(pubBytes) != 65 || pubBytes[0] != 0x04 {
		return "", errors.New("unexpected public key format")
	}

	// keccak256(pubBytes[1:])，取最后 20 字节
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(pubBytes[1:]) // 仅 X||Y
	hash := keccak.Sum(nil)    // 32 bytes
	addr20 := hash[12:]        // last 20 bytes

	// Tron 前缀 0x41 ，拼成 21 字节地址
	addr21 := append([]byte{0x41}, addr20...)

	// Base58Check: 先做两次 sha256 获取 checksum 的前 4 字节
	h1 := sha256.Sum256(addr21)
	h2 := sha256.Sum256(h1[:])
	checksum := h2[:4]

	// 拼接并 base58 编码
	addrWithChecksum := append(addr21, checksum...)
	base58Addr := base58.Encode(addrWithChecksum)

	return base58Addr, nil
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

// 创建质押能量的交易
func (self tronChain) CreateFreezeEnergyTrans(owner, receiver string, trxAmount int64, isBandwidth bool) (info *TronTransaction, err error) {
	if receiver == "" {
		receiver = owner
	}

	hexOwner, err := self.AddrToHexStr(owner)
	if err != nil {
		return nil, err
	}

	hexReceiver, err := self.AddrToHexStr(receiver)
	if err != nil {
		return nil, err
	}

	if err = mynet.DoReq("POST", "https://api.trongrid.io/wallet/freezebalancev2",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

			if b, err := json.Marshal(map[string]any{
				"owner_address":  hexOwner,
				"frozen_balance": trxAmount * 1_000_000, // TRX -> Sun
				"resource": myctrl.ObjFun(func() string {
					if isBandwidth {
						return "BANDWIDTH"
					}
					return "ENERGY"
				}),
				"receiver_address": hexReceiver,
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, time.Second * 15, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				return json.Unmarshal(ret, &info)
			}

			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	if info.Error != "" {
		return nil, fmt.Errorf("tron error:%s", info.Error)
	}

	return info, nil
}

// 创建质押能量的交易
func (self tronChain) CreateUnFreezeEnergyTrans(owner, receiver string, trxAmount int64, isBandwidth bool) (info *TronTransaction, err error) {
	if receiver == "" {
		receiver = owner
	}

	hexOwner, err := self.AddrToHexStr(owner)
	if err != nil {
		return nil, err
	}

	hexReceiver, err := self.AddrToHexStr(receiver)
	if err != nil {
		return nil, err
	}

	if err = mynet.DoReq("POST", "https://api.trongrid.io/wallet/unfreezebalancev2",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

			if b, err := json.Marshal(map[string]any{
				"owner_address":    hexOwner,
				"unfreeze_balance": trxAmount * 1_000_000, // TRX -> Sun
				"resource": myctrl.ObjFun(func() string {
					if isBandwidth {
						return "BANDWIDTH"
					}
					return "ENERGY"
				}),
				"receiver_address": hexReceiver,
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, time.Second * 15, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				return json.Unmarshal(ret, &info)
			}

			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	if info.Error != "" {
		return nil, fmt.Errorf("tron error:%s", info.Error)
	}

	return info, nil
}

// tron 原生币 trx 交易
func (self tronChain) CreateTrxTrans(from, to string, amount int64) (info *TronTransaction, err error) {
	hexFrom, err := self.AddrToHexStr(from)
	if err != nil {
		return nil, err
	}

	hexTo, err := self.AddrToHexStr(to)
	if err != nil {
		return nil, err
	}

	if err = mynet.DoReq("POST", "https://api.trongrid.io/wallet/createtransaction",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

			if b, err := json.Marshal(map[string]any{
				"owner_address": hexFrom,
				"to_address":    hexTo,
				"amount":        amount * 1_000_000, // TRX -> Sun
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, time.Second * 15, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				return json.Unmarshal(ret, &info)
			}

			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	if info.Error != "" {
		return nil, fmt.Errorf("tron error:%s", info.Error)
	}

	return info, nil
}

// SignTronTransaction 对 Tron 的 raw_data_hex 进行签名
func (self tronChain) SignTrans(tx *TronTransaction, privateKeyHex string) error {
	// 1. 解码 raw_data_hex
	rawData, err := hex.DecodeString(tx.RawDataHex)
	if err != nil {
		return fmt.Errorf("decode raw_data_hex failed: %w", err)
	}

	// 2. keccak256 哈希
	hash := crypto.Keccak256(rawData)

	// 3. 加载私钥
	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	// 4. 签名
	signature, err := crypto.Sign(hash, privKey)
	if err != nil {
		return fmt.Errorf("sign failed: %w", err)
	}

	// 5. 返回 hex 字符串
	tx.Signature = []string{hex.EncodeToString(signature)}
	return nil
}

// 广播交易
func (self tronChain) BroadcastTrans(tx *TronTransaction) (info string, err error) {
	if err = mynet.DoReq("POST", "https://api.trongrid.io/wallet/broadcasttransaction",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

			if b, err := json.Marshal(tx); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, time.Second * 15, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				info = string(ret)
				return nil
			}

			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return "", err
	}

	return info, nil
}

// 发送质押能交易
func (self tronChain) SendFreezeEnergyTrans(owner, receiver string, trxAmount int64, isBandwidth bool, privatedKey string) (string, error) {
	//1：创建交易
	tx, err := self.CreateFreezeEnergyTrans(owner, receiver, trxAmount, isBandwidth)
	if err != nil {
		return "", err
	}

	//2:签名 交易
	err = self.SignTrans(tx, privatedKey)
	if err != nil {
		return "", err
	}

	//3:广播交易
	return self.BroadcastTrans(tx)
}

// 取全网资源参数
func (self tronChain) GetTotalNetworkRes() (res *NetworkRes, err error) {
	if err = mynet.DoReq("POST", "https://api.trongrid.io/wallet/getaccountresource",
		func(r *http.Request) (isTls bool, timeout time.Duration, err error) {
			r.Header.Add("accept", "application/json")
			r.Header.Add("Content-Type", "application/json")
			if self.apiKey != "" {
				r.Header.Set("TRON-PRO-API-KEY", self.apiKey)
			}

			if b, err := json.Marshal(map[string]any{
				"address": "TNoP3HyZkfip2H88QkVyF5P3GSX9ax6gyT",
				"visible": true,
			}); err != nil {
				return false, 0, err
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				r.Header.Add("Content-Length", strconv.Itoa(len(b)))
			}

			return true, time.Second * 15, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode == http.StatusOK {
				if err = json.Unmarshal(ret, &res); err != nil {
					return err
				}

				if res.TotalEnergyLimit <= 0 || res.TotalNetLimit <= 0 {
					return fmt.Errorf("invalid network resource data")
				}

				return nil
			}

			return fmt.Errorf("http err code:%v", httpCode)
		}, nil); err != nil {
		return nil, err
	}

	return res, nil
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
