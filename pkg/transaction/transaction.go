package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	BASE_URL = "https://mock-node-wgqbnxruha-as.a.run.app"
)

const (
	STATUS_CONFIRM  = "CONFIRMED"
	STATUS_FAILED   = "FAILED"
	STATUS_PENDEING = "PENDING"
	STATUS_DNE      = "DNE"
	STATUS_UNKNOW   = "UNKNOW"
)

type Transaction struct {
	Symbol            string
	Price             uint64
	Timestamp         uint64
	broadcastResponse BroadcastResponse
	txStatusResponse  TxStatusResponse
}

type TransactionPayload struct {
	Symbol    string `json:"symbol"`
	Price     uint64 `json:"price"`
	Timestamp uint64 `json:"timestamp"`
}

type BroadcastResponse struct {
	TxHash string `json:"tx_hash"`
}

type TxStatusResponse struct {
	TxStatus string `json:"tx_status"`
}

type TransactionStatus struct {
	Success bool
	Status  string
	Message string
}

func (txStatus *TxStatusResponse) mapToGeneralModel() TransactionStatus {
	messageMapping := make(map[string]string)
	messageMapping[STATUS_CONFIRM] = "Transaction has been processed and confirmed"
	messageMapping[STATUS_FAILED] = "Transaction failed to process"
	messageMapping[STATUS_PENDEING] = "Transaction is awaiting processing"
	messageMapping[STATUS_DNE] = "Transaction does not exist"

	messageMapping[STATUS_UNKNOW] = "Unknow transaction"

	if msg, exists := messageMapping[txStatus.TxStatus]; exists {
		return TransactionStatus{
			Success: true,
			Status:  txStatus.TxStatus,
			Message: msg,
		}
	} else {
		return TransactionStatus{
			Success: false,
			Status:  STATUS_UNKNOW,
			Message: messageMapping[STATUS_UNKNOW],
		}
	}
}

func NewTransaction(symbol string, price uint64) Transaction {
	return Transaction{
		Symbol:    symbol,
		Price:     price,
		Timestamp: uint64(time.Now().Unix()),
	}
}

func (transaction *Transaction) createPayload() []byte {
	payload := TransactionPayload{
		Symbol:    transaction.Symbol,
		Price:     transaction.Price,
		Timestamp: transaction.Timestamp,
	}

	dataBytes, _ := json.Marshal(payload)

	return dataBytes
}

func (transaction *Transaction) Broadcast() (*BroadcastResponse, error) {
	res, err := http.Post(BASE_URL+"/broadcast", "application/json", bytes.NewBuffer(transaction.createPayload()))

	if err != nil {
		return nil, errors.New("Can not broadcast please contact admin")
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&transaction.broadcastResponse)

	return &transaction.broadcastResponse, nil
}

func (transaction *Transaction) CheckStatus() (*TransactionStatus, error) {
	res, err := http.Get(BASE_URL + "/check/" + transaction.broadcastResponse.TxHash)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&transaction.txStatusResponse)

	result := transaction.txStatusResponse.mapToGeneralModel()

	return &result, nil
}
