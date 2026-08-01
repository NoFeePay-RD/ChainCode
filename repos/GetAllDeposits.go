package repos

import (
	"encoding/json"
	"fmt"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func GetAllDeposits(
	ctx contractapi.TransactionContextInterface,
) ([]FundDepositResponse, error) {
	query := `{
		"selector": {
			"docType": "D"
		}
	}`

	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("[GetAllDeposits] failed to read customer deposits: %v", err)
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
			return nil, fmt.Errorf("[GetAllDeposits] failed to decode customer fund record: %v", err)
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
