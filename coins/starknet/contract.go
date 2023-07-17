package starknet

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
)

type ABI struct {
	Members []struct {
		Name   string `json:"name"`
		Offset int    `json:"offset"`
		Type   string `json:"type"`
	} `json:"members,omitempty"`
	Inputs []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"inputs"`
	Name    string `json:"name"`
	Outputs []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"outputs"`
	Type            string `json:"type"`
	Size            int    `json:"size,omitempty"`
	StateMutability string `json:"stateMutability,omitempty"`
}

type EntryPointsByType struct {
	Constructor []struct {
		Offset   string `json:"offset"`
		Selector string `json:"selector"`
	} `json:"CONSTRUCTOR"`
	External []struct {
		Offset   string `json:"offset"`
		Selector string `json:"selector"`
	} `json:"EXTERNAL"`
	L1Handler []interface{} `json:"L1_HANDLER"`
}

type ContractDefinition struct {
	ABI               []ABI                  `json:"abi"`
	EntryPointsByType EntryPointsByType      `json:"entry_points_by_type"`
	Program           map[string]interface{} `json:"program"`
}

type rawContractDefinition struct {
	Abi               []ABI             `json:"abi"`
	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`
	Program           string            `json:"program"`
}

func (cd *ContractDefinition) getRawContractDefinition() (*rawContractDefinition, error) {
	compiledContract, err := CompressCompiledContract(cd.Program)
	if err != nil {
		return nil, err
	}

	return &rawContractDefinition{
		Abi:               cd.ABI,
		EntryPointsByType: cd.EntryPointsByType,
		Program:           compiledContract,
	}, nil
}

func CompressCompiledContract(program map[string]interface{}) (string, error) {
	pay, err := json.Marshal(program)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(pay)
	if err != nil {
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
