package common

import "fmt"

var ErrAddressIsRequired = fmt.Errorf("address is required")
var ErrInvalidCoins = fmt.Errorf("coins is invaild")
var ErrInsufficientCoins = fmt.Errorf("coins is insufficient")
