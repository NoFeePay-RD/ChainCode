package types

type FundDepositEntry struct {
	Amount          float64 `json:"amount"`
	CustomerAddress string  `json:"customerAddress"`
	DocType         string  `json:"docType"`
	Timestamp       int64   `json:"timestamp"`
}

type FundTransferEntry struct {
	Amount          float64 `json:"amount"`
	CustomerAddress string  `json:"customerAddress"`
	DocType         string  `json:"docType"`
	MerchantAddress string  `json:"merchantAddress"`
	Timestamp       int64   `json:"timestamp"`
}
