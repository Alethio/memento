package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/parnurzeal/gorequest"
)

const local = "http://localhost:3001/api"
const remote = "https://blk-api.goerli.ethstats.io/api/v3"

func TestBlock(t *testing.T) {
	start := int64(1350000)
	end := int64(1350100)

	for i := start; i <= end; i++ {
		t.Run(strconv.FormatInt(i, 10), func(tt *testing.T) {
			compareBlockDetails(i, tt)

			if i%300 == 0 {
				if i+300 > end {
					compareBlockRange(i, end, tt)
				} else {
					compareBlockRange(i, i+300, tt)
				}
			}
		})
	}
}

func TestAccountTxs(t *testing.T) {
	accounts := []struct {
		Address, Filters string
	}{
		{"8ced5ad0d8da4ec211c17355ed3dbfec4cf0e5b9", "?includedInBlock=1353543&limit=5"},
		{"454a62a6201d9953b1bbc0303c6ff9599d742ac3", "?includedInBlock=1353595&limit=5"},
		{"a94a6a58dd1fc6a8a3e8e75d123f8e7f30026c00", "?includedInBlock=1353026&limit=5"},
		{"62083c80353df771426d209ef578619ee68d5c7a", "?includedInBlock=1353596&limit=5"},
		{"483b937b035a70526a9e8e60a2dddfcfb5fe3b11", "?includedInBlock=1353589&limit=5"},
	}

	for _, a := range accounts {
		t.Run(a.Address, func(tt *testing.T) {
			compareAccountTxs(a.Address, a.Filters, tt)
		})
	}
}

func compareBlockDetails(number int64, t *testing.T) {
	localData := getBlock(number, local, t)
	remoteData := getBlock(number, remote, t)

	ignored := map[string]struct{}{
		"blockBeneficiaryReward": {},
		"blockTime":              {},
		"numberOfContractMsgs":   {},
		"alethioComment":         {},
		"hasBeneficiaryAlias":    {},
		"blockMixHash":           {},
		"blockNonce":             {},
	}

	for blockProperty, remoteValue := range remoteData {
		if _, ok := ignored[blockProperty]; ok {
			t.Logf("IGNORED: `%s`", blockProperty)
			continue
		}

		if localValue, exists := localData[blockProperty]; exists {
			if blockProperty == "txs" {
				compareBlockTxs(localValue, remoteValue, t)
				continue
			}

			if blockProperty == "includesUncle" {
				continue
			}

			if remoteValue != localValue {
				t.Errorf("ERROR: local property `%s` (=%v) is different on remote (=%v)", blockProperty, localValue, remoteValue)
			}
		} else {
			t.Errorf("ERROR: expected property `%s` in local response but did not find it", blockProperty)
		}
	}
}

func compareBlockTxs(local, remote interface{}, t *testing.T) {
	l := local.([]interface{})
	r := remote.([]interface{})

	ignored := map[string]struct{}{
		"createContractMsgsTriggered": {},
		"totalContractMsgsTriggered":  {},
		"type":                        {},
	}

	for index, remoteTxValue := range r {
		t.Logf("TX: checking tx at index %d", index)
		issues := 0

		remoteTx := remoteTxValue.(map[string]interface{})

		compareTx(remoteTx["txHash"].(string), t)
		compareTxLogEntries(remoteTx["txHash"].(string), t)

		if len(l) > index {
			localTx := l[index].(map[string]interface{})

			for txProperty, remoteValue := range remoteTx {
				if _, ok := ignored[txProperty]; ok {
					t.Logf("IGNORED: `%s`", txProperty)
					continue
				}

				if localValue, exists := localTx[txProperty]; exists {
					if remoteValue != localValue {
						issues++
						t.Errorf("ERROR: local property `txs.%d.%s` (=%v) is different on remote (=%v)", index, txProperty, localValue, remoteValue)
					}
				} else {
					issues++
					t.Errorf("ERROR: expected property `txs.%d.%s` in local response but did not find it", index, txProperty)
				}
			}
		} else {
			issues++
			t.Errorf("ERROR: could not find local tx at index %d", index)
		}

		if issues == 0 {
			t.Logf("TX: all good")
		}
	}
}

func compareBlockRange(start, end int64, t *testing.T) {
	localData := getBlockRange(start, end, local, t)
	remoteData := getBlockRange(start, end, remote, t)

	ignored := map[string]struct{}{
		"totalTxValue": {},
	}

	for index, remoteBlockValue := range remoteData {
		t.Logf("BLOCK-RANGE: checking block at index %d", index)
		issues := 0

		remoteBlock := remoteBlockValue.(map[string]interface{})

		if len(localData) > index {
			localBlock := localData[index].(map[string]interface{})

			for property, remoteValue := range remoteBlock {
				if _, ok := ignored[property]; ok {
					t.Logf("IGNORED: `%s`", property)
					continue
				}

				if localValue, exists := localBlock[property]; exists {
					if remoteValue != localValue {
						issues++
						t.Errorf("ERROR: local property `block-range.%d.%s` (=%v) is different on remote (=%v)", index, property, localValue, remoteValue)
					}
				} else {
					issues++
					t.Errorf("ERROR: expected property `block-range.%d.%s` in local response but did not find it", index, property)
				}
			}
		} else {
			issues++
			t.Errorf("ERROR: could not find local block at index %d in block range", index)
		}

		if issues == 0 {
			t.Logf("BLOCK-RANGE: all good")
		}
	}
}

func compareTx(hash string, t *testing.T) {
	localData := getTx(hash, local, t)
	remoteData := getTx(hash, remote, t)

	issues := 0
	t.Logf("TX: checking %s", hash)

	ignored := map[string]struct{}{
		"msgOutput":                  {},
		"contractMsgsTriggered":      {},
		"firstSeenAt":                {},
		"msgError":                   {},
		"msgErrorString":             {},
		"blockMsgValidationIndex":    {},
		"msgPayloadDecoded":          {},
		"tokenTransfersTriggered":    {},
		"totalContractMsgsTriggered": {},
		"type":                       {},
	}

	for property, remoteValue := range remoteData {
		if _, ok := ignored[property]; ok {
			t.Logf("IGNORED: `%s`", property)
			continue
		}

		if localValue, exists := localData[property]; exists {
			if remoteValue != localValue {
				issues++
				t.Errorf("ERROR: local property `%s` (=%v) is different on remote (=%v)", property, localValue, remoteValue)
			}
		} else {
			issues++
			t.Errorf("ERROR: expected property `%s` in local response but did not find it", property)
		}
	}

	if issues == 0 {
		t.Logf("TX: all good")
	}
}

func compareTxLogEntries(hash string, t *testing.T) {
	localData := getTxLogEntries(hash, local, t)
	remoteData := getTxLogEntries(hash, remote, t)

	ignored := map[string]struct{}{
		"txMsgValidationIndex": {},
		"eventDecoded":         {},
		"eventDecodedError":    {},
	}

	t.Logf("LOG-ENTRIES: checking for tx %s", hash)

	for index, remoteLogValue := range remoteData {
		t.Logf("LOG-ENTRIES: checking log at index %d", index)
		issues := 0

		remoteLog := remoteLogValue.(map[string]interface{})

		if len(localData) > index {
			localLog := localData[index].(map[string]interface{})

			for property, remoteValue := range remoteLog {
				if _, ok := ignored[property]; ok {
					t.Logf("IGNORED: `%s`", property)
					continue
				}

				if localValue, exists := localLog[property]; exists {
					if property == "hasLogTopics" {
						localTopics := localValue.([]interface{})
						for i, topic := range remoteValue.([]interface{}) {
							if len(localTopics) > i {
								if topic != localTopics[i] {
									issues++
									t.Errorf("ERROR: local property `tx.log-entries.%d.%s` (=%v) is different on remote (=%v)", index, property, localValue, remoteValue)
								}
							} else {
								issues++
								t.Errorf("ERROR: log topic at index %d does not exist on local", i)
							}
						}
						continue
					}

					if remoteValue != localValue {
						issues++
						t.Errorf("ERROR: local property `tx.log-entries.%d.%s` (=%v) is different on remote (=%v)", index, property, localValue, remoteValue)
					}
				} else {
					issues++
					t.Errorf("ERROR: expected property `tx.log-entries.%d.%s` in local response but did not find it", index, property)
				}
			}
		} else {
			issues++
			t.Errorf("ERROR: could not find local log at index %d in tx log entries", index)
		}

		if issues == 0 {
			t.Logf("LOG-ENTRIES: all good")
		}
	}
}

func compareAccountTxs(address, filter string, t *testing.T) {
	localData := getAccountTxs(address, local, filter, t)
	remoteData := getAccountTxs(address, remote, filter, t)

	ignored := map[string]struct{}{
		"msgErrorString":          {},
		"fee":                     {},
		"type":                    {},
		"blockMsgValidationIndex": {},
	}

	t.Logf("ACCOUNT-TXS: checking for account %s", address)

	for index, remoteTxValue := range remoteData {
		t.Logf("ACCOUNT-TXS: checking tx at index %d", index)
		issues := 0

		remoteTx := remoteTxValue.(map[string]interface{})

		if len(localData) > index {
			localTx := localData[index].(map[string]interface{})

			for property, remoteValue := range remoteTx {
				if _, ok := ignored[property]; ok {
					t.Logf("IGNORED: `%s`", property)
					continue
				}

				if localValue, exists := localTx[property]; exists {
					if remoteValue != localValue {
						issues++
						t.Errorf("ERROR: local property `account.txs.%d.%s` (=%v) is different on remote (=%v)", index, property, localValue, remoteValue)
					}
				} else {
					issues++
					t.Errorf("ERROR: expected property `account.txs.%d.%s` in local response but did not find it", index, property)
				}
			}
		} else {
			issues++
			t.Errorf("ERROR: could not find local tx at index %d in account txs", index)
		}

		if issues == 0 {
			t.Logf("ACCOUNT-TXS: all good")
		}
	}
}

func execQuery(url string, t *testing.T) map[string]interface{} {
	request := gorequest.New()
	resp, body, errs := request.Get(url).End()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	if len(errs) > 0 {
		t.Errorf("did not expect errors, got %v", errs)
	}

	res := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		t.Errorf("got error while unmarshalling response: %s", err)
	}

	return res
}

func getBlock(number int64, source string, t *testing.T) map[string]interface{} {
	res := execQuery(fmt.Sprintf("%s/block/%d", source, number), t)

	return res["data"].(map[string]interface{})
}

func getBlockRange(start, end int64, source string, t *testing.T) []interface{} {
	res := execQuery(fmt.Sprintf("%s/block-range/%d/%d", source, start, end), t)

	return res["data"].([]interface{})
}

func getTx(hash string, source string, t *testing.T) map[string]interface{} {
	res := execQuery(fmt.Sprintf("%s/tx/%s", source, hash), t)

	return res["data"].(map[string]interface{})
}

func getTxLogEntries(hash string, source string, t *testing.T) []interface{} {
	res := execQuery(fmt.Sprintf("%s/tx/%s/log-entries", source, hash), t)

	if res["data"] == nil {
		return []interface{}{}
	}

	return res["data"].([]interface{})
}

func getAccountTxs(address string, source string, filter string, t *testing.T) []interface{} {
	res := execQuery(fmt.Sprintf("%s/account/%s/txs%s", source, address, filter), t)

	return res["data"].([]interface{})
}
