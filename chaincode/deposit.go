package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

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
