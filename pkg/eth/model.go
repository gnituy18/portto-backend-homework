package eth

type Hash string

type Block struct {
	BlockNum     int64  `json:"block_num"`
	BlockHash    Hash   `json:"block_hash"`
	BlockTime    int64  `json:"block_time"`
	ParentHash   Hash   `json:"parent_hash"`
	Transactions []Hash `json:"transations"`
}

type Transation struct {
	TxHash Hash   `json:"tx_hash"`
	From   Hash   `json:"from"`
	To     Hash   `json:"to"`
	Nonce  int64  `json:"nonce"`
	Data   Hash   `json:"data"`
	Value  string `json:"value"`
	Logs   []*Log `json:"logs"`
}

type Log struct {
	Index int  `json:"index"`
	Data  Hash `json:"data"`
}
