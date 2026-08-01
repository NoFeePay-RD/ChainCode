package types

type FundDepositResponse struct {
	Amount          float64 `json:"amount"`
	CustomerAddress string  `json:"customerAddress"`
	TransactionID   string  `json:"transactionID"`
	Timestamp       int64   `json:"timestamp"`
}

type FundTransferResponse struct {
	Amount          float64 `json:"amount"`
	CustomerAddress string  `json:"customerAddress"`
	MerchantAddress string  `json:"merchantAddress"`
	Timestamp       int64   `json:"timestamp"`
	TransactionID   string  `json:"transactionID"`
}

type BalanceCheckResponse struct {
	CustomerAddress string  `json:"customerAddress"`
	Balance         float64 `json:"balance"`
}

type SettlementResponse struct {
	MerchantAddress  string  `json:"merchantAddress"`
	TotalAmount      float64 `json:"amount"`
	TransactionCount int     `json:"transactionCount"`
	TransactionDate  string  `json:"transactionDate"`
}
