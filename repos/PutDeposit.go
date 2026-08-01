package repos

import (
	"encoding/json"
	"fmt"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func PutDeposit(
	ctx contractapi.TransactionContextInterface,
	txID string,
	customerAddress string,
	amount float64,
	timestamp int64,
) error {

	fundDeposit := FundDepositEntry{
		Amount:          amount,
		CustomerAddress: customerAddress,
		DocType:         "D",
		Timestamp:       timestamp,
	}

	fundDepositJSON, err := json.Marshal(fundDeposit)
	if err != nil {
		return fmt.Errorf("[PutDeposit] failed to marshal trasnfer document: %v", err)
	}

	// Store all fields under transaction id
	err = ctx.GetStub().PutState(txID, fundDepositJSON)

	if err != nil {
		return fmt.Errorf("[PutDeposit] failed to put deposit document: %v", err)
	}
	return nil
}
