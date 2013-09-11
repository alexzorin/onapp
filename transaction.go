package onapp

import (
	"encoding/json"
	"fmt"
)

type Transactions []Transaction

type Transaction struct {
	client     *Client
	Id         int    `json:"id"`
	Status     string `json:"status"`
	Parent     int    `json:"parent_id"`
	User       int    `json:"user_id"`
	ParentType string `json:"parent_type"`
	Action     string `json:"action"`
	CreatedAt  string `json:"created_at"`
	StartedAt  string `json:"started_at"`
	UpdatedAt  string `json:"updated_at"`
	Dependent  int    `json:"dependent_transaction_id"`
}

func (c *Client) GetTransactions() (Transactions, error) {
	return c.getTransactions(0)
}

func (c *Client) getTransactions(vmId int) (Transactions, error) {
	path := "transactions.json"
	if vmId != 0 {
		path = fmt.Sprintf("virtual_machines/%d/transactions.json", vmId)
	}
	data, err, _ := c.getReq(path)
	if err != nil {
		return nil, err
	}
	var out []map[string]Transaction
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	txs := make([]Transaction, len(out))
	for i := range txs {
		txs[i] = out[i]["transaction"]
	}
	return txs, nil
}

func (t *Transaction) IsValid() bool {
	return t.Id > 0
}
