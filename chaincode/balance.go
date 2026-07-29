package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) CheckCustomerBalance(ctx contractapi.TransactionContextInterface, customerAddress string) (*BalanceCheckResponse, error) {

	depositQuery := fmt.Sprintf(`{
		"selector": {
			"docType": "D",
			"customerAddress": "%s"
		}
	}`, customerAddress)

	depositsIterator, err := ctx.GetStub().GetQueryResult(depositQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve deposit records: %v", err)
	}
	defer depositsIterator.Close()

	var totalDeposits float64

	for depositsIterator.HasNext() {
		response, err := depositsIterator.Next()
		if err != nil {
			return nil, err
		}

		var deposit FundDepositEntry

		err = json.Unmarshal(response.Value, &deposit)
		if err != nil {
			return nil, fmt.Errorf("failed to decode customer deposit record: %v", err)
		}

		totalDeposits += deposit.Amount
	}

	transfersQuery := fmt.Sprintf(`{
		"selector": {
			"docType": "T",
			"customerAddress": "%s"
		}
	}`, customerAddress)

	transfersIterator, err := ctx.GetStub().GetQueryResult(transfersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve deposit records: %v", err)
	}
	defer transfersIterator.Close()

	var totalTransfers float64

	for transfersIterator.HasNext() {
		response, err := transfersIterator.Next()
		if err != nil {
			return nil, err
		}

		var transfer FundTransferEntry

		err = json.Unmarshal(response.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("failed to decode customer deposit record: %v", err)
		}

		totalTransfers += transfer.Amount
	}

	var balance float64 = totalDeposits - totalTransfers

	response := &BalanceCheckResponse{
		CustomerAddress: customerAddress,
		Balance:         balance,
	}

	return response, nil
}
