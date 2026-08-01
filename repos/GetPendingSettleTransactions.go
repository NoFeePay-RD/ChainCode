package repos

import (
	"encoding/json"
	"fmt"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func GetPendingSettleTransactions(
	ctx contractapi.TransactionContextInterface,
	merchantAddress string,
	startTime int64,
	endTime int64,
) ([]FundTransferResponse, error) {

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "T",
			"merchantAddress": "%s",
			"timestamp": {
				"$gte": %d,
				"$lt": %d
			}
		}
	}`, merchantAddress, startTime, endTime)

	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("[CalculateSettlement] failed to retrieve transfer records: %v", err)
	}
	defer iterator.Close()

	var results []FundTransferResponse

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[CalculateSettlement] failed to read transfer record: %v", err)
		}

		var transferEntry FundTransferEntry
		err = json.Unmarshal(response.Value, &transferEntry)
		if err != nil {
			return nil, fmt.Errorf("[CalculateSettlement] failed to unmarshal transfer record: %v", err)
		}

		var transferResponse FundTransferResponse
		transferResponse.TransactionID = response.Key
		transferResponse.Amount = transferEntry.Amount
		transferResponse.CustomerAddress = transferEntry.CustomerAddress
		transferResponse.MerchantAddress = transferEntry.MerchantAddress
		transferResponse.Timestamp = transferEntry.Timestamp
		results = append(results, transferResponse)
	}

	return results, nil
}
