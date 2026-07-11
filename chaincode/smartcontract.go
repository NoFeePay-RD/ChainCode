package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type FundDepositEntry struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Timestamp int64   `json:"timestamp"`
}

type FundDepositResponse struct {
	Amount          float64 `json:"amount"`
	CustomerAddress string  `json:"customerAddress"`
	Currency        string  `json:"currency"`
	TransactionID   string  `json:"transactionID"`
	Timestamp       int64   `json:"timestamp"`
}

func (s *SmartContract) DepositFunds(ctx contractapi.TransactionContextInterface, customerAddress string, amount float64, currency string) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be greater than zero")
	}

	//Same timestamp is retreived for both endosement and commit. So it keeps the deterministic nature of the transaction
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	txID := ctx.GetStub().GetTxID()

	fundDeposit := FundDepositEntry{
		Amount:    amount,
		Currency:  currency,
		Timestamp: timestamp.AsTime().UnixMicro(),
	}

	fundDepositJSON, err := json.Marshal(fundDeposit)
	if err != nil {
		return err
	}

	// Create a unique composite key: F~CustomerAddress~TransactionID
	compositeKey, err := ctx.GetStub().CreateCompositeKey("F", []string{customerAddress, txID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, fundDepositJSON)
}

// GetAllCustomerFunds retrieves all customer funds from the ledger
func (s *SmartContract) GetAllCustomerFunds(ctx contractapi.TransactionContextInterface) ([]FundDepositResponse, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("F", []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to read customer funds: %v", err)
	}
	defer iterator.Close()

	var results []FundDepositResponse

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var fund FundDepositEntry

		err = json.Unmarshal(response.Value, &fund)
		if err != nil {
			return nil, fmt.Errorf("failed to decode customer fund record: %v", err)
		}

		// Decode composite key
		_, attributes, err := ctx.GetStub().SplitCompositeKey(response.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %v", err)
		}

		var fundDepositResponse FundDepositResponse
		fundDepositResponse.Amount = fund.Amount
		fundDepositResponse.Currency = fund.Currency
		fundDepositResponse.CustomerAddress = attributes[0]
		fundDepositResponse.TransactionID = attributes[1]
		fundDepositResponse.Timestamp = fund.Timestamp
		results = append(results, fundDepositResponse)
	}

	return results, nil
}

type FundTransferEntry struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Timestamp int64   `json:"timestamp"`
}

type FundTransferResponse struct {
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	CustomerAddress string  `json:"customerAddress"`
	MerchantAddress string  `json:"merchantAddress"`
	Timestamp       int64   `json:"timestamp"`
	TransactionID   string  `json:"transactionID"`
}

func (s *SmartContract) TransferFunds(ctx contractapi.TransactionContextInterface, customerAddress string, amount float64, currency string, merchantAddress string) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be greater than zero")
	}

	//check if the customer has enough funds to transfer
	// Get all customer funds
	customerFunds, err := s.GetAllCustomerFunds(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve customer funds: %v", err)
	}

	var totalFunds float64
	for _, fund := range customerFunds {
		if fund.CustomerAddress == customerAddress && fund.Currency == currency {
			totalFunds += fund.Amount
		}
	}

	if totalFunds < amount {
		return fmt.Errorf("insufficient funds for transfer: available %.2f, required %.2f", totalFunds, amount)
	}

	//Same timestamp is retreived for both endosement and commit. So it keeps the deterministic nature of the transaction
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	txID := ctx.GetStub().GetTxID()

	fundTransfer := FundTransferEntry{
		Amount:    amount,
		Currency:  currency,
		Timestamp: timestamp.AsTime().UnixMicro(),
	}

	fundTransferJSON, err := json.Marshal(fundTransfer)
	if err != nil {
		return err
	}

	// Create a unique composite key: T~CustomerAddress~MerchantAddress~TransactionID
	compositeKey, err := ctx.GetStub().CreateCompositeKey("T", []string{customerAddress, merchantAddress, txID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, fundTransferJSON)
}

func (s *SmartContract) GetAllCustomerTransfers(ctx contractapi.TransactionContextInterface) ([]FundTransferResponse, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("T", []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to read customer transfers: %v", err)
	}
	defer iterator.Close()

	var results []FundTransferResponse

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var transfer FundTransferEntry

		err = json.Unmarshal(response.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("failed to decode customer transfer record: %v", err)
		}

		// Decode composite key
		_, attributes, err := ctx.GetStub().SplitCompositeKey(response.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %v", err)
		}

		var fundTransferResponse FundTransferResponse
		fundTransferResponse.Amount = transfer.Amount
		fundTransferResponse.Currency = transfer.Currency
		fundTransferResponse.CustomerAddress = attributes[0]
		fundTransferResponse.MerchantAddress = attributes[1]
		fundTransferResponse.Timestamp = transfer.Timestamp
		fundTransferResponse.TransactionID = attributes[2]
		results = append(results, fundTransferResponse)
	}

	return results, nil
}
