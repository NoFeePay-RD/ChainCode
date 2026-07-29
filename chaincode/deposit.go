package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) DepositFunds(ctx contractapi.TransactionContextInterface, customerAddress string, amount float64) error {
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
		Amount:          amount,
		CustomerAddress: customerAddress,
		DocType:         "D",
		Timestamp:       timestamp.AsTime().UnixMicro(),
	}

	fundDepositJSON, err := json.Marshal(fundDeposit)
	if err != nil {
		return err
	}

	// Store all fields under transaction id
	return ctx.GetStub().PutState(txID, fundDepositJSON)
}

// GetAllCustomerFunds retrieves all customer deposits from the ledger
func (s *SmartContract) GetAllCustomerDeposits(ctx contractapi.TransactionContextInterface) ([]FundDepositResponse, error) {
	query := `{
		"selector": {
			"docType": "D"
		}
	}`

	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to read customer deposits: %v", err)
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

		var fundDepositResponse FundDepositResponse
		fundDepositResponse.Amount = fund.Amount
		fundDepositResponse.CustomerAddress = fund.CustomerAddress
		fundDepositResponse.TransactionID = response.Key
		fundDepositResponse.Timestamp = fund.Timestamp
		results = append(results, fundDepositResponse)
	}

	return results, nil
}
