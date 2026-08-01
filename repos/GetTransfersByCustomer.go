package repos

import (
	"encoding/json"
	"fmt"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func GetTransfersByCustomer(
	ctx contractapi.TransactionContextInterface,
	customerAddress string,
) ([]FundTransferResponse, error) {

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "T",
			"customerAddress": "%s"
		}
	}`, customerAddress)

	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("[GetTransfersByCustomer] failed to retrieve transfer records: %v", err)
	}
	defer iterator.Close()

	var transfers []FundTransferResponse

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var transfer FundTransferEntry

		err = json.Unmarshal(response.Value, &transfer)
		if err != nil {
			return nil, fmt.Errorf("[GetTransfersByCustomer] failed to decode transfer records: %v", err)
		}

		var transferResponse FundTransferResponse
		transferResponse.TransactionID = response.Key
		transferResponse.Amount = transfer.Amount
		transferResponse.CustomerAddress = transfer.CustomerAddress
		transferResponse.MerchantAddress = transfer.MerchantAddress
		transferResponse.Timestamp = transfer.Timestamp

		transfers = append(transfers, transferResponse)
	}

	return transfers, nil
}
