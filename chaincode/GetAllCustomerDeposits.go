package chaincode

import (
	"fmt"
	"nofeepay/repos"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// GetAllCustomerFunds retrieves all customer deposits from the ledger
func (s *Contract) GetAllCustomerDeposits(
	ctx contractapi.TransactionContextInterface,
) ([]FundDepositResponse, error) {

	var response []FundDepositResponse
	response, err := repos.GetAllDeposits(ctx)

	if err != nil {
		return nil, fmt.Errorf("[GetAllCustomerDeposits] failed to read customer deposits: %v", err)
	}
	return response, nil
}
