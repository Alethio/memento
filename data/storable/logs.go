package storable

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/alethio/web3-go/types"
)

type LogEntriesGroup struct {
	RawBlock    types.Block
	RawReceipts []types.Receipt

	blockNumber int64

	logEntries []*LogEntry
}

type LogEntry struct {
	TxHash          string
	LogIndex        int32
	LogData         ByteArray
	LoggedBy        string
	Topic0          string
	Topic1          string
	Topic2          string
	Topic3          string
	IncludedInBlock int64
}

func NewStorableLogEntries(block types.Block, receipts []types.Receipt) *LogEntriesGroup {
	return &LogEntriesGroup{RawBlock: block, RawReceipts: receipts}
}

func (leg *LogEntriesGroup) ToDB(tx *sql.Tx) error {
	if len(leg.RawReceipts) == 0 {
		return nil
	}

	log.Trace("storing log entries")
	start := time.Now()
	defer func() {
		log.WithField("duration", time.Since(start)).WithField("count", len(leg.logEntries)).Debug("done storing log entries")
	}()

	err := leg.enhance()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("log_entries", "tx_hash", "log_index", "log_data", "logged_by", "topic_0", "topic_1", "topic_2", "topic_3", "included_in_block"))
	if err != nil {
		return err
	}

	for _, log := range leg.logEntries {
		_, err = stmt.Exec(log.TxHash, log.LogIndex, log.LogData, log.LoggedBy, log.Topic0, log.Topic1, log.Topic2, log.Topic3, log.IncludedInBlock)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

// enhance processes the raw receipts data and generates a combined list of LogEntry entities for all the
// transactions included in the block
func (leg *LogEntriesGroup) enhance() error {
	number, err := strconv.ParseInt(leg.RawBlock.Number, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	leg.blockNumber = number

	for _, receipt := range leg.RawReceipts {
		for index, log := range receipt.Logs {
			le, err := leg.buildStorableLogEntry(log, receipt.TransactionHash, int32(index))
			if err != nil {
				return err
			}

			leg.logEntries = append(leg.logEntries, le)
		}
	}

	return nil
}

func (leg *LogEntriesGroup) buildStorableLogEntry(log types.Log, txHash string, index int32) (*LogEntry, error) {
	l := &LogEntry{}
	l.IncludedInBlock = leg.blockNumber
	l.TxHash = Trim0x(txHash)
	l.LogIndex = index // = transaction log index

	l.LogData = ByteArray(Trim0x(log.Data))
	l.LoggedBy = Trim0x(log.Address)

	if len(log.Topics) > 0 {
		l.Topic0 = Trim0x(log.Topics[0])
	}

	if len(log.Topics) > 1 {
		l.Topic1 = Trim0x(log.Topics[1])
	}

	if len(log.Topics) > 2 {
		l.Topic2 = Trim0x(log.Topics[2])
	}

	if len(log.Topics) > 3 {
		l.Topic3 = Trim0x(log.Topics[3])
	}

	return l, nil
}
