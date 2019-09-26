package api

import (
	"git.aleth.io/alethio/memento/api/types"
	"git.aleth.io/alethio/memento/data/storable"
)

func (a *API) getBlockTxs(number int64) ([]types.Tx, error) {
	var txs = make([]types.Tx, 0)

	rows, err := a.db.Query(`select tx_index, tx_hash, value, "from", "to", msg_gas_limit, tx_gas_used, tx_gas_price from txs where included_in_block = $1 order by tx_index`, number)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for rows.Next() {
		var (
			txIndex  int32
			txHash   string
			value    string
			from     storable.ByteArray
			to       storable.ByteArray
			gasLimit string
			gasUsed  string
			gasPrice string
		)

		err = rows.Scan(&txIndex, &txHash, &value, &from, &to, &gasLimit, &gasUsed, &gasPrice)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		txs = append(txs, types.Tx{
			TxIndex:     &txIndex,
			TxHash:      &txHash,
			Value:       &value,
			From:        &from,
			To:          &to,
			MsgGasLimit: &gasLimit,
			TxGasUsed:   &gasUsed,
			TxGasPrice:  &gasPrice,
		})
	}

	return txs, nil
}
