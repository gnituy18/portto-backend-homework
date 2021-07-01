package eth

type Block struct {
	BlockNum     uint64   `json:"block_num" gorm:"primaryKey,index"`
	BlockHash    string   `json:"block_hash" gorm:"index"`
	BlockTime    uint64   `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"Transactions,omitempty" gorm:"-"`
}

type Transaction struct {
	TxHash string `json:"tx_hash" gorm:"primaryKey"`
	From   string `json:"from"`
	To     string `json:"to"`
	Nonce  uint64 `json:"nonce"`
	Data   string `json:"data"`
	Value  int64  `json:"value"`
	Logs   []*Log `json:"logs" gorm:"embedded"`
}

type Log struct {
	Index uint   `json:"index"`
	Data  string `json:"data"`
}
