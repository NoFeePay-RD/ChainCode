package chaincode

import (
	"fmt"
	"nofeepay/repos"
	. "nofeepay/types"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *Contract) CalculateTodaysSettlement(
	ctx contractapi.TransactionContextInterface,
	merchantAddress string,
	date_txt string,
) (*SettlementResponse, error) {

	//convert given data string to unix micro format
	t, err := time.Parse("20060102", date_txt)
	if err != nil {
		return nil, fmt.Errorf("[CalculateTodaysSettlement] failed to convert time formats: %v", err)
	}

	dateMicro := t.UnixMicro()
	nextDateMicro := t.AddDate(0, 0, 1).UnixMicro()
	results, err := repos.GetPendingSettleTransactions(ctx, merchantAddress, dateMicro, nextDateMicro)
	if err != nil {
		return nil, fmt.Errorf("[CalculateTodaysSettlement] failed to read pending settle transactions: %v", err)
	}

	var total float64
	count := len(results)

	for _, transfer := range results {
		total += transfer.Amount
	}

	response := &SettlementResponse{
		MerchantAddress:  merchantAddress,
		TotalAmount:      total,
		TransactionCount: count,
		TransactionDate:  date_txt,
	}

	return response, nil
}
