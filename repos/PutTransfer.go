package repos

import (
	"encoding/json"
	"fmt"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func PutTransfer(
	ctx contractapi.TransactionContextInterface,
	txID string,
	customerAddress string,
	amount float64,
	merchantAddress string,
	timestamp int64,
) error {

	fundTransfer := FundTransferEntry{
		Amount:          amount,
		CustomerAddress: customerAddress,
		DocType:         "T",
		MerchantAddress: merchantAddress,
		Timestamp:       timestamp,
	}

	fundTransferJSON, err := json.Marshal(fundTransfer)
	if err != nil {
		return fmt.Errorf("[PutTransfer] Failed to marshal transfer record: %v", err)
	}

	err = ctx.GetStub().PutState(txID, fundTransferJSON)

	if err != nil {
		return fmt.Errorf("[PutTransfer] failed to put transfer document: %v", err)
	}
	return nil
}
