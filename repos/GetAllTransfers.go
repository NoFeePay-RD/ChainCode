package repos

import (
	"encoding/json"
	"fmt"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func GetAllTransfers(
	ctx contractapi.TransactionContextInterface,
) ([]FundTransferResponse, error) {

	query := `{
		"selector": {
			"docType": "T"
		}
	}`

	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("[GetAllTransfers] failed to read customer transfers: %v", err)
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
			return nil, fmt.Errorf("[GetAllTransfers] failed to decode customer transfer record: %v", err)
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
