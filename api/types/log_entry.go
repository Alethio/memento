package types

import "git.aleth.io/alethio/memento/data/storable"

type LogEntry struct {
	TxHash            string                 `json:"txHash"`
	LogIndex          int32                  `json:"logIndex"`
	LogData           storable.ByteArray     `json:"logData"`
	LoggedBy          string                 `json:"loggedBy"`
	HasLogTopics      []string               `json:"hasLogTopics"`
	EventDecoded      map[string]interface{} `json:"eventDecoded"`
	EventDecodedError string                 `json:"eventDecodedError"`
}
