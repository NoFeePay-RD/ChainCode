package repos

import (
	"encoding/json"
	"fmt"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func GetDepositsByCustomer(
	ctx contractapi.TransactionContextInterface,
	customerAddress string,
) ([]FundDepositResponse, error) {

	query := fmt.Sprintf(`{
		"selector": {
			"docType": "D",
			"customerAddress": "%s"
		}
	}`, customerAddress)

	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("[GetDepositsByCustomer] failed to retrieve deposit records: %v", err)
	}
	defer iterator.Close()

	var deposits []FundDepositResponse

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[GetDepositsByCustomer] failed to retrieve deposit records: %v", err)
		}

		var deposit FundDepositEntry
		err = json.Unmarshal(response.Value, &deposit)
		if err != nil {
			return nil, fmt.Errorf("[GetDepositsByCustomer] failed to decode customer deposit record: %v", err)
		}

		var depositResponse FundDepositResponse
		depositResponse.TransactionID = response.Key
		depositResponse.Amount = deposit.Amount
		depositResponse.CustomerAddress = deposit.CustomerAddress
		depositResponse.Timestamp = deposit.Timestamp

		deposits = append(deposits, depositResponse)
	}

	return deposits, nil
}
