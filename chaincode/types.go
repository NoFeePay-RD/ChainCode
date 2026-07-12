package chaincode

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type FundDepositEntry struct {
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
}

type FundDepositResponse struct {
	Amount          float64 `json:"amount"`
	CustomerAddress string  `json:"customerAddress"`
	TransactionID   string  `json:"transactionID"`
	Timestamp       int64   `json:"timestamp"`
}

type FundTransferEntry struct {
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
}

type FundTransferResponse struct {
	Amount          float64 `json:"amount"`
	CustomerAddress string  `json:"customerAddress"`
	MerchantAddress string  `json:"merchantAddress"`
	Timestamp       int64   `json:"timestamp"`
	TransactionID   string  `json:"transactionID"`
}
