package mychain

import (
	"encoding/json"
	"fmt"
	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/mynet"
	"net/http"
	"strings"
	"time"
)

type (
	BtcTxInfo struct {
		TxID        string  // 交易ID
		Time        string  // 交易时间
		FromAddress string  // 发送地址
		ToAddress   string  // 接收地址
		AmountBTC   float64 // 金额（单位 BTC）
		Currency    string  // 币种（固定 BTC）
	}
)

func (self BtcTxInfo) TimeStamp() int64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", self.Time, time.UTC)
	if err != nil {
		return 0
	}

	return t.UnixMilli()
}

type btcChain struct {
}

func ImpBtcChain() btcChain {
	return btcChain{}
}

func (self btcChain) GetNowBlock(timeout time.Duration) (info int64, err error) {
	type TmpBlock struct {
		Data []struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}

	if err = mynet.DoReq("GET", "https://api.blockchair.com/bitcoin/blocks?limit=1",
		func(r *http.Request) (isTls bool, tm time.Duration, err error) {
			r.Header.Add("Content-Type", "application/json")
			return true, timeout, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode != http.StatusOK {
				return fmt.Errorf("http err code:%v[%s]", httpCode, string(ret))
			}

			var tmp TmpBlock
			if err := json.Unmarshal(ret, &tmp); err != nil {
				return err
			}
			if tmp.Data != nil || len(tmp.Data) > 0 {
				info = tmp.Data[0].ID
			}
			return nil
		}, nil); err != nil {
		return 0, err
	}

	return
}

func (self btcChain) GetFullTxID(blockId int64, timeout time.Duration, sleepOnPage time.Duration) (list []string, err error) {
	type TransactionData struct {
		Data []struct {
			Hash       string `json:"hash"`
			IsCoinBase bool   `json:"is_coinbase"`
		} `json:"data"`
	}
	offseet := 0
	list = make([]string, 0, 8000)

	for {
		//mylog.Info("offset is :", offseet)
		var txData TransactionData
		err = mynet.DoReq("GET", fmt.Sprintf("https://api.blockchair.com/bitcoin/transactions?q=block_id(%d)&limit=100&offset=%d", blockId, offseet),
			func(r *http.Request) (isTls bool, tm time.Duration, err error) {
				r.Header.Add("Content-Type", "application/json")
				return true, timeout, nil
			}, func(ret []byte, httpCode int) error {
				if httpCode != http.StatusOK {
					return fmt.Errorf("http err code:%v[%s]", httpCode, string(ret))
				}

				if err := json.Unmarshal(ret, &txData); err != nil {
					return err
				}
				return nil
			}, nil)
		if err != nil {
			break
		}

		if txData.Data == nil || len(txData.Data) == 0 {
			break
		}
		for i, _ := range txData.Data {
			if txData.Data[i].IsCoinBase {
				mylog.Debug("coin base is :", txData.Data[i].Hash)
			} else {
				list = append(list, txData.Data[i].Hash)
			}
		}

		if len(txData.Data) < 100 { //说明是最后一页了
			break
		}

		offseet += 100
		time.Sleep(sleepOnPage) //加点延迟防止被限流
	}

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (self btcChain) GetTrans(txID string, timeout time.Duration) (info []BtcTxInfo, err error) {
	var result struct {
		Data map[string]struct {
			Transaction struct {
				Time       string `json:"time"`
				IsCoinBase bool   `json:"is_coinbase"`
			} `json:"transaction"`
			Inputs []struct {
				Recipient string `json:"recipient"`
				Value     int64  `json:"value"`
			} `json:"inputs"`
			Outputs []struct {
				Recipient string `json:"recipient"`
				Value     int64  `json:"value"`
			} `json:"outputs"`
		} `json:"data"`
	}

	if err = mynet.DoReq("GET", fmt.Sprintf("https://api.blockchair.com/bitcoin/dashboards/transaction/%s", txID),
		func(r *http.Request) (isTls bool, tm time.Duration, err error) {
			r.Header.Add("Content-Type", "application/json")
			return true, timeout, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode != http.StatusOK {
				return fmt.Errorf("http err code:%v[%s]", httpCode, string(ret))
			}

			if err := json.Unmarshal(ret, &result); err != nil {
				return err
			}

			txData, exists := result.Data[txID]
			if !exists {
				return fmt.Errorf("transaction %s not found in response", txID)
			}

			fromAddress := "Multiple"
			if len(txData.Inputs) == 1 {
				fromAddress = txData.Inputs[0].Recipient
			}

			for _, out := range txData.Outputs {
				// 忽略无效输出
				if out.Recipient == "" || out.Value == 0 {
					continue
				}

				info = append(info, BtcTxInfo{
					TxID:        txID,
					Time:        txData.Transaction.Time,
					FromAddress: fromAddress,
					ToAddress:   out.Recipient,
					AmountBTC:   float64(out.Value) / 1e8,
					Currency:    "BTC",
				})
			}

			return nil
		}, nil); err != nil {
		return nil, err
	}

	return
}

func (self btcChain) GetBatchTrans(txids []string, timeout time.Duration) (info []BtcTxInfo, err error) {
	type blockCypherTx struct {
		Hash      string `json:"hash"`
		Confirmed string `json:"confirmed"`
		Inputs    []struct {
			Addresses []string `json:"addresses"`
		} `json:"inputs"`
		Outputs []struct {
			Addresses []string `json:"addresses"`
			Value     int64    `json:"value"` // 单位为 Satoshi
		} `json:"outputs"`
	}

	url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/main/txs/%s", strings.Join(txids, ";"))
	//mylog.Info("url is :", url)
	if err = mynet.DoReq("GET", url,
		func(r *http.Request) (isTls bool, tm time.Duration, err error) {
			r.Header.Add("Content-Type", "application/json")
			return true, timeout, nil
		}, func(ret []byte, httpCode int) error {
			if httpCode != http.StatusOK {
				return fmt.Errorf("http err code:%v[%s]", httpCode, string(ret))
			}

			var txList []blockCypherTx
			// 当只有一笔交易时，返回的是 object 而不是 array
			if strings.HasPrefix(string(ret), "{") {
				var singleTx blockCypherTx
				if err := json.Unmarshal(ret, &singleTx); err != nil {
					return fmt.Errorf("unmarshal single tx failed: %v", err)
				}
				txList = append(txList, singleTx)
			} else {
				if err := json.Unmarshal(ret, &txList); err != nil {
					return fmt.Errorf("unmarshal tx list failed: %v", err)
				}
			}

			for _, tx := range txList {
				from := "Multiple"
				if len(tx.Inputs) == 1 && len(tx.Inputs[0].Addresses) == 1 {
					from = tx.Inputs[0].Addresses[0]
				}
				to := ""
				amount := float64(0)
				for _, out := range tx.Outputs {
					if len(out.Addresses) > 0 {
						to = out.Addresses[0]
						amount = float64(out.Value) / 1e8 // Satoshi 转 BTC
						break
					}
				}
				info = append(info, BtcTxInfo{
					TxID:        tx.Hash,
					Time:        tx.Confirmed,
					FromAddress: from,
					ToAddress:   to,
					AmountBTC:   amount,
					Currency:    "BTC",
				})
			}
			return nil
		}, nil); err != nil {
		return nil, err
	}

	return
}
