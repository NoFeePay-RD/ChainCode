package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) TransferFunds(ctx contractapi.TransactionContextInterface, customerAddress string, amount float64, merchantAddress string) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be greater than zero")
	}

	//check if the customer has enough funds to transfer
	var response, err = s.CheckCustomerBalance(ctx, customerAddress)
	if err != nil {
		return fmt.Errorf("failed to check customer balance")
	}

	if response.Balance < amount {
		return fmt.Errorf("insufficient funds for transfer: available %.2f, required %.2f", response.Balance, amount)
	}

	//Same timestamp is retreived for both endosement and commit. So it keeps the deterministic nature of the transaction
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	txID := ctx.GetStub().GetTxID()

	fundTransfer := FundTransferEntry{
		Amount:          amount,
		CustomerAddress: customerAddress,
		DocType:         "T",
		MerchantAddress: merchantAddress,
		Timestamp:       timestamp.AsTime().UnixMicro(),
	}

	fundTransferJSON, err := json.Marshal(fundTransfer)
	if err != nil {
		return err
	}

	// Store all fields under transaction id
	return ctx.GetStub().PutState(txID, fundTransferJSON)
}

func (s *SmartContract) GetAllCustomerTransfers(ctx contractapi.TransactionContextInterface) ([]FundTransferResponse, error) {
	query := `{
		"selector": {
			"docType": "T"
		}
	}`

	iterator, err := ctx.GetStub().GetQueryResult(query)
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

		var fundTransferResponse FundTransferResponse
		fundTransferResponse.TransactionID = response.Key
		fundTransferResponse.Amount = transfer.Amount
		fundTransferResponse.CustomerAddress = transfer.CustomerAddress
		fundTransferResponse.MerchantAddress = transfer.MerchantAddress
		fundTransferResponse.Timestamp = transfer.Timestamp
		results = append(results, fundTransferResponse)
	}

	return results, nil
}
