package chaincode

import (
	"fmt"
	"nofeepay/repos"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *Contract) CheckCustomerBalance(
	ctx contractapi.TransactionContextInterface,
	customerAddress string,
) (*BalanceCheckResponse, error) {

	depositsResponses, err := repos.GetDepositsByCustomer(ctx, customerAddress)
	if err != nil {
		return nil, fmt.Errorf("[CheckCustomerBalance] failed to read deposits for the customer: %v", err)
	}

	var totalDeposits float64
	for _, deposit := range depositsResponses {
		totalDeposits += deposit.Amount
	}

	transferResponses, err := repos.GetTransfersByCustomer(ctx, customerAddress)
	if err != nil {
		return nil, fmt.Errorf("[CheckCustomerBalance] failed to calculate the total transfers: %v", err)
	}

	var totalTransfers float64
	for _, transfer := range transferResponses {
		totalTransfers += transfer.Amount
	}

	balance := totalDeposits - totalTransfers

	response := &BalanceCheckResponse{
		CustomerAddress: customerAddress,
		Balance:         balance,
	}

	return response, nil
}
