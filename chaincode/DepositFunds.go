package chaincode

import (
	"fmt"
	"nofeepay/repos"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *Contract) DepositFunds(
	ctx contractapi.TransactionContextInterface,
	customerAddress string,
	amount float64,
) error {

	if amount <= 0 {
		return fmt.Errorf("[DepositFunds] deposit amount must be greater than zero")
	}

	//Same timestamp is retreived for both endosement and commit. So it keeps the deterministic nature of the transaction
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[DepositFunds] failed to get transaction timestamp: %v", err)
	}

	micro := timestamp.AsTime().UnixMicro()
	txID := ctx.GetStub().GetTxID()
	err = repos.PutDeposit(ctx, txID, customerAddress, amount, micro)

	if err != nil {
		return fmt.Errorf("[DepositFunds] failed to put deposit document: %v", err)
	}
	return nil
}
