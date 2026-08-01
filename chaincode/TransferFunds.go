package chaincode

import (
	"fmt"
	"nofeepay/repos"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func (s *Contract) TransferFunds(
	ctx contractapi.TransactionContextInterface,
	customerAddress string,
	amount float64,
	merchantAddress string,
) error {

	if amount <= 0 {
		return fmt.Errorf("[TransferFunds] transfer amount must be greater than zero")
	}

	//check if the customer has enough funds to transfer
	var response, err = s.CheckCustomerBalance(ctx, customerAddress)
	if err != nil {
		return fmt.Errorf("[TransferFunds] failed to check customer balance")
	}

	if response.Balance < amount {
		return fmt.Errorf("[TransferFunds] insufficient funds for transfer: available %.2f, required %.2f", response.Balance, amount)
	}

	//Same timestamp is retreived for both endosement and commit. So it keeps the deterministic nature of the transaction
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[TransferFunds] failed to get transaction timestamp: %v", err)
	}

	micro := timestamp.AsTime().UnixMicro()
	txID := ctx.GetStub().GetTxID()
	err = repos.PutTransfer(ctx, txID, customerAddress, amount, merchantAddress, micro)

	if err != nil {
		return fmt.Errorf("[TransferFunds] failed to put transfer document: %v", err)
	}
	return nil
}
