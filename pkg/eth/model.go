package eth

type Block struct {
	BlockNum     uint64   `json:"block_num"`
	BlockHash    string   `json:"block_hash"`
	BlockTime    uint64   `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"transations,omitempty"`
}

type Transation struct {
	TxHash string `json:"tx_hash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Nonce  uint64 `json:"nonce"`
	Data   string `json:"data"`
	Value  string `json:"value"`
	Logs   []*Log `json:"logs"`
}

type Log struct {
	Index uint   `json:"index"`
	Data  string `json:"data"`
}
