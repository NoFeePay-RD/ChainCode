package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) CalculateTodaysSettlement(ctx contractapi.TransactionContextInterface, merchantAddress string, date_txt string) (*SettlementResponse, error) {

	//convert given data string to unix micro format
	t, err := time.Parse("20060102", date_txt)
	if err != nil {
		panic(err)
	}

	dateMicro := t.UnixMicro()
	nextDateMicro := t.AddDate(0, 0, 1).UnixMicro()
	query := fmt.Sprintf(`{
		"selector": {
			"docType": "T",
			"merchantAddress": "%s",
			"timestamp": {
				"$gte": %d,
				"$lt": %d
			}
		}
	}`, merchantAddress, dateMicro, nextDateMicro)

	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transfer records: %v", err)
	}
	defer iterator.Close()

	var totalAmount float64
	var transactionCount int

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to read transfer record: %v", err)
		}

		var fund FundTransferEntry
		err = json.Unmarshal(response.Value, &fund)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal transfer record: %v", err)
		}

		totalAmount += fund.Amount
		transactionCount++
	}

	settlement := &SettlementResponse{
		MerchantAddress:  merchantAddress,
		TotalAmount:      totalAmount,
		TransactionCount: transactionCount,
		TransactionDate:  date_txt,
	}

	return settlement, nil
}
