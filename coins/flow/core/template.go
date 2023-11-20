package core

import _ "embed"

//go:embed templates/create-account.cdc
var CreateAccountTpl string

//go:embed templates/transfer-flow.cdc
var TransferFlowTpl string
