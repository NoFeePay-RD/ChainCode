package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *SmartContract) CalculateTodaysSettlement(ctx contractapi.TransactionContextInterface, merchantAddress string, date_txt string) (*SettlementResponse, error) {

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("T", []string{date_txt, merchantAddress})
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
