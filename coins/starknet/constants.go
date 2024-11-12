package starknet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type JsonType struct {
	AccountAddress string        `json:"accountAddress"`
	TypedData      JsonTypedData `json:"typedData"`
}

type JsonTypedData struct {
	Types       map[string][]Definition `json:"types"`
	PrimaryType string                  `json:"primaryType"`
	Domain      Domain                  `json:"domain"`
	Message     interface{}             `json:"message"`
}

type TypedData struct {
	Types       map[string]TypeDef
	PrimaryType string
	Domain      Domain
	Message     TypedMessage
}

type Domain struct {
	Name    string
	Version string
	ChainId string
}

type TypeDef struct {
	Encoding    *big.Int
	Definitions []Definition
}

type Definition struct {
	Name string
	Type string
}

type TypedMessage interface {
	FmtDefinitionEncoding(string) []*big.Int
}

/*
encoding definition for standard StarkNet Domain messages
*/
func (dm Domain) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	switch field {
	case "name":
		fmtEnc = append(fmtEnc, StrToFelt(dm.Name).Big())
	case "version":
		fmtEnc = append(fmtEnc, StrToFelt(dm.Version).Big())
	case "chainId":
		fmtEnc = append(fmtEnc, StrToFelt(dm.ChainId).Big())
	}
	return fmtEnc
}

/*
'typedData' interface for interacting and signing typed data in accordance with https://github.com/0xs34n/starknet.js/tree/develop/src/utils/typedData
*/
func NewTypedData(types map[string]TypeDef, pType string, dom Domain) (td TypedData, err error) {
	td = TypedData{
		Types:       types,
		PrimaryType: pType,
		Domain:      dom,
	}
	if _, ok := td.Types[pType]; !ok {
		return td, fmt.Errorf("invalid primary type: %s", pType)
	}

	for k, v := range td.Types {
		enc, err := td.GetTypeHash(k)
		if err != nil {
			return td, fmt.Errorf("error encoding type hash: %s %w", enc.String(), err)
		}
		v.Encoding = enc
		td.Types[k] = v
	}
	return td, nil
}

// (ref: https://github.com/0xs34n/starknet.js/blob/767021a203ac0b9cdb282eb6d63b33bfd7614858/src/utils/typedData/index.ts#L166)
func (td TypedData) GetMessageHash(account *big.Int, msg TypedMessage, sc StarkCurve) (hash *big.Int, err error) {
	elements := []*big.Int{UTF8StrToBig("StarkNet Message")}

	domEnc, err := td.GetTypedMessageHash("StarkNetDomain", td.Domain, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash domain: %w", err)
	}
	elements = append(elements, domEnc)
	elements = append(elements, account)

	msgEnc, err := td.GetTypedMessageHash(td.PrimaryType, msg, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash message: %w", err)
	}

	elements = append(elements, msgEnc)
	hash, err = ComputeHashOnElements(elements)
	return hash, err
}

func (td TypedData) GetStructHash(account *big.Int, msg string, sc StarkCurve) (hash *big.Int, err error) {
	elements := []*big.Int{UTF8StrToBig("StarkNet Message")}

	domEnc, err := td.GetTypedMessageHash("StarkNetDomain", td.Domain, sc)
	if err != nil {
		return hash, fmt.Errorf("could not hash domain: %w", err)
	}
	elements = append(elements, domEnc)
	elements = append(elements, account)

	msgEnc, err := td.encodeData(td.PrimaryType, msg)
	if err != nil {
		return hash, fmt.Errorf("could not hash message: %w", err)
	}

	elements = append(elements, msgEnc)
	hash, err = ComputeHashOnElements(elements)
	return hash, err
}

func (td TypedData) encodeData(inType string, msg string) (*big.Int, error) {
	var message map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &message); err != nil {
		return nil, err
	}

	prim, ok := td.Types[inType]
	if !ok {
		return nil, fmt.Errorf("type not found in TypedData")
	}

	var elements []*big.Int
	elements = append(elements, prim.Encoding)

	for _, def := range prim.Definitions {
		v, exists := message[def.Name]
		if !exists {
			continue
		}

		if _, ok = td.Types[def.Type]; ok {
			m, err := json.Marshal(message[def.Name])
			if err != nil {
				return nil, err
			}
			h, err := td.encodeData(def.Type, string(m))
			if err != nil {
				return nil, err
			}
			elements = append(elements, h)
		} else if def.Type == "felt" {
			vBn, err := encodeValue(v.(string))
			if err != nil {
				return nil, err
			}
			elements = append(elements, vBn)
		} else if def.Type == "felt*" {
			interfaces, ok := v.([]interface{})
			if !ok {
				return nil, fmt.Errorf("Expected a slice of interfaces")
			}
			strs := make([]string, len(interfaces))
			for i, iface := range interfaces {
				str, ok := iface.(string)
				if !ok {
					return nil, fmt.Errorf("Expected a string")
				}
				strs[i] = str
			}
			ele, err := EncodeValues(strs)
			if err != nil {
				return nil, err
			}
			elements = append(elements, ele)
		}

	}

	hash, err := ComputeHashOnElements(elements)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func encodeValue(v string) (*big.Int, error) {
	if strings.HasPrefix(v, "0x") {
		return HexToBN(v)
	}
	// Try to parse the string as an integer
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return big.NewInt(i), nil
	}

	return StrToFelt(v).Big(), nil
}

func EncodeValues(values []string) (*big.Int, error) {
	var ele []*big.Int
	for _, v := range values {
		vBn, err := encodeValue(v)
		if err != nil {
			return nil, err
		}
		ele = append(ele, vBn)
	}
	res, err := ComputeHashOnElements(ele)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (td TypedData) GetTypedMessageHash(inType string, msg TypedMessage, sc StarkCurve) (hash *big.Int, err error) {
	prim := td.Types[inType]
	elements := []*big.Int{prim.Encoding}

	for _, def := range prim.Definitions {
		if def.Type == "felt" {
			fmtDefinitions := msg.FmtDefinitionEncoding(def.Name)
			elements = append(elements, fmtDefinitions...)
			continue
		}

		innerElements := []*big.Int{}
		encType := td.Types[def.Type]
		innerElements = append(innerElements, encType.Encoding)
		fmtDefinitions := msg.FmtDefinitionEncoding(def.Name)
		innerElements = append(innerElements, fmtDefinitions...)
		innerElements = append(innerElements, big.NewInt(int64(len(innerElements))))

		innerHash, err := sc.HashElements(innerElements)
		if err != nil {
			return hash, fmt.Errorf("error hashing internal elements: %v %w", innerElements, err)
		}
		elements = append(elements, innerHash)
	}

	hash, err = ComputeHashOnElements(elements)
	return hash, err
}

func (td TypedData) GetTypeHash(inType string) (ret *big.Int, err error) {
	enc, err := td.EncodeType(inType)
	if err != nil {
		return ret, err
	}
	sel := GetSelectorFromName(enc)
	return sel, nil
}

func (td TypedData) EncodeType(inType string) (enc string, err error) {
	var typeDefs TypeDef
	var ok bool
	if typeDefs, ok = td.Types[inType]; !ok {
		return enc, fmt.Errorf("can't parse type %s from types %v", inType, td.Types)
	}
	var buf bytes.Buffer
	customTypes := make(map[string]TypeDef)
	buf.WriteString(inType)
	buf.WriteString("(")
	for i, def := range typeDefs.Definitions {
		if def.Type != "felt" && def.Type != "felt*" {
			var customTypeDef TypeDef
			if customTypeDef, ok = td.Types[def.Type]; !ok {
				return enc, fmt.Errorf("can't parse type %s from types %v", def.Type, td.Types)
			}
			customTypes[def.Type] = customTypeDef
		}
		buf.WriteString(fmt.Sprintf("%s:%s", def.Name, def.Type))
		if i != (len(typeDefs.Definitions) - 1) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")

	for customTypeName, customType := range customTypes {
		buf.WriteString(fmt.Sprintf("%s(", customTypeName))
		for i, def := range customType.Definitions {
			buf.WriteString(fmt.Sprintf("%s:%s", def.Name, def.Type))
			if i != (len(customType.Definitions) - 1) {
				buf.WriteString(",")
			}
		}
		buf.WriteString(")")
	}
	return buf.String(), nil
}

func GetMessageHashWithJson(jsonStr string) (string, error) {
	var data JsonType
	json.Unmarshal([]byte(jsonStr), &data)

	exampleTypes := make(map[string]TypeDef)
	for k, v := range data.TypedData.Types {
		exampleTypes[k] = TypeDef{Definitions: v}
	}
	ttd, err := NewTypedData(exampleTypes, data.TypedData.PrimaryType, data.TypedData.Domain)
	if err != nil {
		return "", err
	}
	messageData, err := json.Marshal(data.TypedData.Message)
	if err != nil {
		return "", err
	}
	address, err := HexToBN(data.AccountAddress)
	if err != nil {
		return "", err
	}
	hash, err := ttd.GetStructHash(address, string(messageData), SC())
	if err != nil {
		return "", err
	}

	return BigToHex(hash), nil
}
