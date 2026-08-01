package chaincode

import (
	"fmt"
	"nofeepay/repos"
	. "nofeepay/types"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *Contract) GetAllCustomerTransfers(
	ctx contractapi.TransactionContextInterface,
) ([]FundTransferResponse, error) {

	var response []FundTransferResponse
	response, err := repos.GetAllTransfers(ctx)

	if err != nil {
		return nil, fmt.Errorf("[GetAllCustomerTransfers] failed to read customer transfers: %v", err)
	}
	return response, err
}
