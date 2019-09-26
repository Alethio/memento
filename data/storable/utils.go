package storable

import (
	"fmt"
	"math/big"
	"strings"
)

// HexStrToBigIntStr transforms a hex sting like "0xff" to a big int string like "15". Arbitrary length values are possible.
func HexStrToBigIntStr(hexString string) (string, error) {
	value, err := HexStrToBigInt(hexString)
	return value.String(), err
}

// HexStrToBigInt transforms a hex sting like "0xff" to a big.Int. Arbitrary length values are possible.
func HexStrToBigInt(hexString string) (*big.Int, error) {
	value := new(big.Int)
	_, ok := value.SetString(Trim0x(hexString), 16)
	if !ok {
		return value, fmt.Errorf("could not transform hex string to big int: %s", hexString)
	}

	return value, nil
}

// Trim0x removes the "0x" prefix of hexes if it exists
func Trim0x(str string) string {
	return strings.TrimPrefix(str, "0x")
}
