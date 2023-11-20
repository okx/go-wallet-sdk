package stacks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

func Serialize(stacksTransaction StacksTransaction) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.Write(getBytes(int64(stacksTransaction.Version), 0))
	chainIdBuffer := bytes.NewBuffer(make([]byte, 0, 4))

	chainIdBuffer.Write(getBytesByLength(stacksTransaction.ChainId, 8))
	buf.Write(sliceByteBuffer(chainIdBuffer))
	buf.Write(serializeAuth(&stacksTransaction.Auth))

	buf.Write(getBytes(int64(stacksTransaction.AnchorMode), 0))
	buf.Write(getBytes(int64(stacksTransaction.PostConditionMode), 0))

	buffer2 := serializeLPList(stacksTransaction.PostConditions)
	buf.Write(buffer2)

	buffer3 := SerializePayload(stacksTransaction.Payload)
	buf.Write(buffer3)

	return sliceByteBuffer(buf)
}

func SerializePayload(payload Payload) []byte {
	bufferArray := bytes.NewBuffer([]byte{})
	bufferArray.Write(getBytes(int64(payload.getPayloadType()), 0))
	switch payload.getPayloadType() {
	case 0:
		tokenTransferPayload := payload
		recipient := tokenTransferPayload.getRecipient()
		serializedCV := serializeCV(*recipient)
		bufferArray.Write(serializedCV)
		bufferArray.Write(toArrayLike(tokenTransferPayload.getAmount(), 8))
		bufferArray.Write(serializeStacksMessage(*tokenTransferPayload.getMemo()))
	case 2:
		contractCallPayload := payload
		bufferArray.Write(serializeStacksMessage(*contractCallPayload.getContractAddress()))
		bufferArray.Write(serializeStacksMessage(*contractCallPayload.getContractName()))
		bufferArray.Write(serializeStacksMessage(*contractCallPayload.getFunctionName()))

		length := len(contractCallPayload.getFunctionArgs())
		bufferArray.Write(toArrayLike(big.NewInt(int64(length)), 4))

		functionArgs := make([]byte, 0)
		for _, cv := range contractCallPayload.getFunctionArgs() {
			s := serializeCV(cv)
			functionArgs = append(functionArgs, s...)
		}
		bufferArray.Write(functionArgs)
	default:
		panic("no support")
	}
	return sliceByteBuffer(bufferArray)

}

func serializeStacksMessage(message Message) []byte {
	switch message.GetType() {
	case MEMOSTRING:
		return serializeMemoString(message.(Memo))
	case ADDRESS:
		return serializeAddress(message.(Address))
	case LENGTHPREFIXEDSTRING:
		return serializeLPString(message.(LengthPrefixedString))
	case PAYLOAD:
		return SerializePayload(message.(Payload))
	case POSTCONDITION:
		return serializePostCondition(message.(PostConditionInterface))
	default:
		panic("unknown stack message type")
	}
}

func serializePostCondition(postCondition PostConditionInterface) []byte {
	switch postCondition.getConditionType() {
	case 0:
		postCondition0 := postCondition.(STXPostCondition)
		bufferArray := bytes.NewBuffer(make([]byte, 0, MaxBufferSize))
		bufferArray.Write(fromHexString(fmt.Sprintf("0x0%d", postCondition0.ConditionType)))
		bufferArray.Write(serializePrincipal(postCondition0.Principal))
		bufferArray.Write(fromHexString(fmt.Sprintf("0x0%d", postCondition0.ConditionCode)))
		bufferArray.Write(toArrayLike(postCondition0.amount, 8))
		return sliceByteBuffer(bufferArray)
	case 1:
		postCondition1 := postCondition.(FungiblePostCondition)

		bufferArray := bytes.NewBuffer(make([]byte, 0, MaxBufferSize))
		bufferArray.Write(fromHexString(fmt.Sprintf("0x0%d", postCondition1.ConditionType)))
		bufferArray.Write(serializePrincipal(postCondition1.Principal))

		bufferArray.Write(serializeAssetInfo(postCondition1.assetInfo))

		bufferArray.Write(fromHexString(fmt.Sprintf("0x0%d", postCondition1.ConditionCode)))
		bufferArray.Write(toArrayLike(postCondition1.amount, 8))

		return sliceByteBuffer(bufferArray)
	default:
		panic("ConditionType unknown")
	}

	return nil
}

func serializeAssetInfo(info AssetInfo) []byte {
	bufferArray := bytes.NewBuffer(make([]byte, 0, MaxBufferSize))
	bufferArray.Write(serializeAddress(info.address))
	bufferArray.Write(serializeLPString(info.contractName))
	bufferArray.Write(serializeLPString(info.assetName))

	return sliceByteBuffer(bufferArray)
}

func serializePrincipal(principal PostConditionPrincipalInterface) []byte {
	bufferArray := bytes.NewBuffer(make([]byte, 0, MaxBufferSize))
	prefix := principal.getPrefix()
	bufferArray.Write(fromHexString(fmt.Sprintf("0x0%d", prefix)))
	if prefix == PostConditionPrincipalIDCONTRACT {
		contractPrincipal := principal.(ContractPrincipal)
		bufferArray.Write(serializeAddress(contractPrincipal.Address))
		bufferArray.Write(serializeLPString(contractPrincipal.contractName))
	} else if prefix == PostConditionPrincipalIDSTANDARD {
		standardPrincipal := principal.(PostConditionPrincipal)
		bufferArray.Write(serializeAddress(standardPrincipal.Address))
	}

	return sliceByteBuffer(bufferArray)
}

func serializeStacksMessageLP(message LengthPrefixedString) []byte {
	switch message.Type {
	case 2:
		return serializeLPString(message)
	default:
		panic("unknown stack message type")
	}
}

func serializeLPList(lpList LPList) []byte {
	bufferArray := bytes.NewBuffer(make([]byte, 0, 4096))
	list := lpList.Values
	bufferArray.Write(fromHexString(intToHexString(len(list), &lpList.LengthPrefixBytes)))
	for i := 0; i < len(list); i++ {
		bufferArray.Write(serializeStacksMessage(list[i]))
	}
	return sliceByteBuffer(bufferArray)
}

func serializeLPString(lps LengthPrefixedString) []byte {
	bufferArray := bytes.NewBuffer(make([]byte, 0, MaxBufferSize))
	b := []byte(lps.Content)
	length := len(b)
	lengthBytes := intToBytes(length, lps.LengthPrefixBytes)

	bufferArray.Write(lengthBytes)
	bufferArray.Write(b)

	return sliceByteBuffer(bufferArray)
}

func serializeMemoString(memo Memo) []byte {
	bufferArray := bytes.NewBuffer([]byte{})
	contentBuffer := []byte(memo.Content)
	hex1 := hex.EncodeToString(contentBuffer)
	MemoMaxLengthBytes := 34
	num := hex1
	width := MemoMaxLengthBytes * 2
	str := fmt.Sprintf("%s", num)
	//str = strings.Repeat("0", width-len(str)) + str // Put 0 on the left
	hex1 = str + strings.Repeat("0", width-len(str)) // Put 0 on the right
	byteHex, _ := hex.DecodeString(hex1)
	bufferArray.Write(byteHex)
	return sliceByteBuffer(bufferArray)
}

func serializeCV(value Message) []byte {
	switch value.GetType() {
	case 0:
		return serializeIntCV(value.(*IntCV))
	case 1:
		return serializeUintCV(value.(*UintCV))
	case 2:
		return serializeBufferCV(value.(*BufferCV))
	case 5:
		return bufferWithTypeID(value.GetType(), serializeAddress(*value.(StandardPrincipalCV).Address))
	case 12:
		return serializeTupleCV(value.(*TupleCV))
	case 3, 4:
		return serializeBoolCV(value.(*BooleanCV))
	case 10:
		return serializeOptionalCV(value.(*SomeCV))
	case 9:
		return serializeNoneCV()
	case 6:
		return serializeContractPrincipalCV(value.(*ContractPrincipalCV))
	case 7, 8:
		return serializeResponseCV(value.(*ResponseCV))
	case 11:
		return serializeListCV(value.(*ListCV))
	case 13:
		return serializeStringCV(value.(*StringCV), "ascii")
	case 14:
		return serializeStringCV(value.(*StringCV), "utf8")
	default:
		panic("Unable to Serialize. Invalid Clarity Value.")
	}
}

func serializeStringCV(value *StringCV, encoding string) []byte {
	bufferArray := bytes.NewBuffer([]byte{})
	var str []byte
	if encoding == "ascii" {
		str = asciiToBytes(value.Data)
	} else if encoding == "utf8" {
		str = []byte(value.Data)
	} else {
		panic("error encoding type")
	}
	length := make([]byte, 4)
	writeUInt32BE(length, uint32(len(str)), 0)
	bufferArray.Write(length)
	bufferArray.Write(str)

	return bufferWithTypeID(value.GetType(), sliceByteBuffer(bufferArray))
}

func serializeListCV(value *ListCV) []byte {
	bufferArray := bytes.NewBuffer([]byte{})
	length := make([]byte, 4)
	writeUInt32BE(length, uint32(len(value.List)), 0)
	bufferArray.Write(length)
	for _, v := range value.List {
		serializedValue := serializeCV(v)
		bufferArray.Write(serializedValue)
	}
	return bufferWithTypeID(value.GetType(), sliceByteBuffer(bufferArray))
}

func serializeBoolCV(value *BooleanCV) []byte {
	return []byte{byte(value.GetType())}
}

func serializeResponseCV(value *ResponseCV) []byte {
	return bufferWithTypeID(value.Type_, serializeCV(value.Value))
}

func serializeNoneCV() []byte {
	return []byte{OptionalNone}
}

func serializeContractPrincipalCV(contract *ContractPrincipalCV) []byte {
	buffer := serializeAddress(contract.Address)
	buffer = append(buffer, serializeLPString(contract.ContractName)...)
	return bufferWithTypeID(contract.Type_, buffer)
}

func serializeOptionalCV(value *SomeCV) []byte {
	return bufferWithTypeID(value.Type_, serializeCV(value.Value))
}

func serializeAddress(address Address) []byte {
	bufferArray := make([]byte, MaxBufferSize)
	bufferArray[0] = byte(address.Version)
	ha, _ := hex.DecodeString(address.Hash160)
	copy(bufferArray[1:], ha)
	return bufferArray[:len(ha)+1]
}

func serializeIntCV(cv *IntCV) []byte {
	b := toArrayLike(cv.Value, 16)
	return bufferWithTypeID(cv.Type_, b)
}

func serializeUintCV(cv *UintCV) []byte {
	b := toArrayLike(cv.Value, 16)
	return bufferWithTypeID(cv.Type_, b)
}

func serializeTupleCV(cv *TupleCV) []byte {
	bufferArray := bytes.NewBuffer([]byte{})
	length := make([]byte, 4)
	writeUInt32BE(length, uint32(len(cv.Data)), 0)
	bufferArray.Write(length)

	// Sort the data according to the lexicographical order of the keys
	pairs := make([]KeyValuePair, 0, len(cv.Data))
	for name, value := range cv.Data {
		pairs = append(pairs, KeyValuePair{name, value})
	}
	sortByKey(pairs)

	for _, pair := range pairs {
		nameWithLength := createLPString(pair.Name, nil, nil)
		bufferArray.Write(serializeLPString(*nameWithLength))

		serializedValue := serializeCV(pair.Value)
		bufferArray.Write(serializedValue)
	}

	s := bufferWithTypeID(cv.GetType(), sliceByteBuffer(bufferArray))
	return s // Assuming bufferWithTypeID function exists
}

func serializeBufferCV(value *BufferCV) []byte {
	b := toArrayLike(big.NewInt(int64(len(value.Buffer))), 4)
	bufferArray := bytes.NewBuffer(make([]byte, 0, MaxBufferSize))
	bufferArray.Write(b)
	bufferArray.Write(value.Buffer)

	return bufferWithTypeID(value.Type_, sliceByteBuffer(bufferArray))
}

func bufferWithTypeID(typeId int, buffer []byte) []byte {
	bufferArray := bytes.NewBuffer(make([]byte, 0, MaxBufferSize))
	bufferArray.WriteByte(byte(typeId))
	bufferArray.Write(buffer)

	return sliceByteBuffer(bufferArray)
}

func serializeAuth(tx *StandardAuthorization) []byte {
	bufferArray := new(bytes.Buffer)
	binary.Write(bufferArray, binary.LittleEndian, byte(tx.AuthType))
	switch tx.AuthType {
	case 4:
		if tx.SpendingCondition == nil {
			panic("spendingCondition is null")
		}
		bufferArray.Write(serializeSingleSigSpendingCondition(tx.SpendingCondition))
	case 5:
		if tx.SpendingCondition == nil {
			panic("spendingCondition is null")
		}
		if tx.SponsorSpendingCondition == nil {
			panic("spendingCondition is null")
		}
		bufferArray.Write(serializeSingleSigSpendingCondition(tx.SpendingCondition))
		bufferArray.Write(serializeSingleSigSpendingCondition(tx.SponsorSpendingCondition))
	default:
		panic("Unexpected transaction AuthType while serializing")
	}

	return sliceByteBuffer(bufferArray)
}

func serializeSingleSigSpendingCondition(condition *SingleSigSpendingCondition) []byte {
	bufferArray := new(bytes.Buffer)
	binary.Write(bufferArray, binary.LittleEndian, byte(condition.HashMode))
	byteSigner, err := hex.DecodeString(condition.Signer)
	if err != nil {
		panic(err)
	}
	bufferArray.Write(byteSigner)
	binary.Write(bufferArray, binary.LittleEndian, toArrayLike(&condition.Nonce, 8))
	binary.Write(bufferArray, binary.LittleEndian, toArrayLike(&condition.Fee, 8))
	bufferArray.Write([]byte{byte(condition.KeyEncoding)})
	bufferArray.Write(serializeMessageSignature(&condition.Signature))
	return sliceByteBuffer(bufferArray)
}

func serializeMessageSignature(messageSignature *MessageSignature) []byte {
	bufferArray := new(bytes.Buffer)
	b, err := hex.DecodeString(messageSignature.Data)
	if err != nil {
		panic(err)
	}
	bufferArray.Write(b)
	return sliceByteBuffer(bufferArray)
}

func DeserializePostCondition(postCondition string) PostConditionInterface {
	bytesReader := NewBytesReader(hexToBytes(postCondition))

	postConditionType := bytesReader.ReadUInt8()
	principal := deserializePrincipal(bytesReader)

	switch postConditionType {
	case STX:
		conditionCode := bytesReader.ReadUInt8()
		amount, _ := new(big.Int).SetString(hex.EncodeToString(bytesReader.ReadBytes(8)), 10)
		return STXPostCondition{
			PostCondition: PostCondition{
				StacksMessage: StacksMessage{
					Type: POSTCONDITION,
				},
				ConditionType: STX,
				Principal:     principal,
				ConditionCode: int(conditionCode),
			},
			amount: amount,
		}
	case Fungible:
		assetInfo := deserializeAssetInfo(bytesReader)
		conditionCode := bytesReader.ReadUInt8()
		amount, _ := new(big.Int).SetString(hex.EncodeToString(bytesReader.ReadBytes(8)), 16)
		return FungiblePostCondition{
			PostCondition: PostCondition{
				StacksMessage: StacksMessage{
					Type: POSTCONDITION,
				},
				ConditionType: Fungible,
				Principal:     principal,
				ConditionCode: int(conditionCode),
			},
			assetInfo: *assetInfo,
			amount:    amount,
		}
	}
	return nil
}

func deserializeAssetInfo(bytesReader *BytesReader) *AssetInfo {
	return &AssetInfo{
		type_:        ASSETINFO,
		address:      *deserializeAddress(bytesReader),
		contractName: *deserializeLPString(bytesReader),
		assetName:    *deserializeLPString(bytesReader),
	}
}

func deserializePrincipal(bytesReader *BytesReader) PostConditionPrincipalInterface {
	prefix := bytesReader.ReadUInt8()
	address := deserializeAddress(bytesReader)

	p := PostConditionPrincipal{
		Type:    PRINCIPAL,
		Prefix:  int(prefix),
		Address: *address,
	}

	if prefix == Standard {
		return p
	}

	contractName := deserializeLPString(bytesReader)
	return ContractPrincipal{
		PostConditionPrincipal: p,
		contractName:           *contractName,
	}
}

func DeserializeCV(serializedClarityValue string) ClarityValue {
	if strings.HasPrefix(serializedClarityValue, "0x") {
		serializedClarityValue = serializedClarityValue[:2]
	}

	bytesReader := NewBytesReader(hexToBytes(serializedClarityValue))
	type_ := bytesReader.ReadUInt8()
	switch type_ {
	case Int:
		readBytes := bytesReader.ReadBytes(16)
		value := new(big.Int).SetBytes(readBytes)
		return *NewIntCV(value)
	case Uint:
		readBytes := bytesReader.ReadBytes(16)
		value := new(big.Int).SetBytes(readBytes)
		return *NewUintCV(value)
	case Buffer:
		bufferLength := bytesReader.readUInt32BE()
		return BufferCV{Buffer, bytesReader.ReadBytes(int(bufferLength))}
	case PrincipalStandard:
		address := deserializeAddress(bytesReader)
		return StandardPrincipalCV{
			Type_:   PrincipalStandard,
			Address: address,
		}
	case PrincipalContract:
		cAddress := deserializeAddress(bytesReader)
		contractName := deserializeLPString(bytesReader)
		return ContractPrincipalCV{
			PrincipalContract,
			*cAddress,
			*contractName,
		}
	case OptionalNone:
		return NoneCV{OptionalNone}
	case OptionalSome:
		return SomeCV{OptionalSome, DeserializeCV(serializedClarityValue[2:])}
	}

	return nil
}

func deserializeAddress(bytesReader *BytesReader) *Address {
	version := hexToInt(hex.EncodeToString(bytesReader.ReadBytes(1)))
	data := hex.EncodeToString(bytesReader.ReadBytes(20))

	return &Address{0, uint64(version), data}
}

func deserializeLPString(bytesReader *BytesReader) *LengthPrefixedString {
	length := hexToInt(hex.EncodeToString(bytesReader.ReadBytes(1)))
	content := string(bytesReader.ReadBytes(length))
	return &LengthPrefixedString{
		Content:           content,
		LengthPrefixBytes: 1,
		MaxLengthBytes:    128,
		Type:              LENGTHPREFIXEDSTRING,
	}
}

func DeserializeJson(args []interface{}) []ClarityValue {
	var res []ClarityValue
	for _, arg := range args {
		argType := getType(arg)
		jsonBytes, err := json.Marshal(arg)
		if err != nil {
			panic(err)
		}
		res = append(res, parseFunctionArgs(argType, jsonBytes))
	}
	return res
}

func parseFunctionArgs(argType int, jsonBytes []byte) ClarityValue {
	switch argType {
	case Int:
		var cv IntCV
		err := json.Unmarshal(jsonBytes, &cv)
		if err != nil {
			panic(err)
		}
		return &cv
	case Uint:
		var cv UintCV
		err := json.Unmarshal(jsonBytes, &cv)
		if err != nil {
			panic(err)
		}
		return &cv
	case Buffer:
		var cv BufferCV
		err := json.Unmarshal(jsonBytes, &cv)
		if err != nil {
			panic(err)
		}
		return &cv
	case BoolTrue, BoolFalse:
		var cv BooleanCV
		err := json.Unmarshal(jsonBytes, &cv)
		if err != nil {
			panic(err)
		}
		return &cv
	case PrincipalStandard:
		var cv StandardPrincipalCV
		err := json.Unmarshal(jsonBytes, &cv)
		if err != nil {
			panic(err)
		}
		return cv
	case PrincipalContract:
		var cv ContractPrincipalCV
		err := json.Unmarshal(jsonBytes, &cv)
		if err != nil {
			panic(err)
		}
		return &cv
	case ResponseOk, ResponseErr:
		var cv ResponseCV
		cv.Type_ = argType
		newType, newJsonBytes := getTypeWithJson(jsonBytes)
		cv.Value = parseFunctionArgs(newType, newJsonBytes)
		return &cv
	case OptionalNone:
		return &NoneCV{OptionalNone}
	case OptionalSome:
		var cv SomeCV
		cv.Type_ = OptionalSome
		newType, newJsonBytes := getTypeWithJson(jsonBytes)
		cv.Value = parseFunctionArgs(newType, newJsonBytes)
		return &cv
	case List:
		var cv ListCV
		cv.Type_ = List
		var lists []ClarityValue
		jsonData := make(map[string]interface{})
		err := json.Unmarshal(jsonBytes, &jsonData)
		if err != nil {
			panic(err)
		}
		fmt.Println(jsonData)
		value, ok := jsonData["list"].([]interface{})
		if !ok {
			panic("Invalid argument")
		}
		for _, v := range value {
			newJsonBytes, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			newType := getType(v)
			l := parseFunctionArgs(newType, newJsonBytes)
			lists = append(lists, l)
		}
		cv.List = lists
		return &cv
	case Tuple:
		var cv TupleCV
		cv.Type = Tuple
		jsonData := make(map[string]interface{})
		err := json.Unmarshal(jsonBytes, &jsonData)
		if err != nil {
			panic(err)
		}
		value, ok := jsonData["data"].(map[string]interface{})
		if !ok {
			panic("Invalid argument")
		}
		dataMap := make(map[string]ClarityValue)
		for k, v := range value {
			newJsonBytes, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			newType := getType(v)
			tuple := parseFunctionArgs(newType, newJsonBytes)
			dataMap[k] = tuple
		}
		cv.Data = dataMap
		return &cv
	case IntASCII, IntUTF8:
		var cv StringCV
		cv.Type_ = argType
		err := json.Unmarshal(jsonBytes, &cv)
		if err != nil {
			panic(err)
		}
		return &cv
	}

	return nil
}

func getTypeWithJson(jsonBytes []byte) (int, []byte) {
	jsonData := make(map[string]interface{})
	err := json.Unmarshal(jsonBytes, &jsonData)
	if err != nil {
		panic(err)
	}
	v, ok := jsonData["value"].(map[string]interface{})
	if !ok {
		panic("Invalid argument")
	}
	newType, ok := v["type"].(float64)
	if !ok {
		panic("Invalid argument")
	}
	newJsonBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return int(newType), newJsonBytes
}

func getType(arg interface{}) int {
	argMap, ok := arg.(map[string]interface{})
	if !ok {
		panic("Invalid argument")
	}
	argType, ok := argMap["type"].(float64)
	if !ok {
		panic("Invalid argument")
	}
	return int(argType)
}

func getFunctionArgs(jsonData string) []interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		panic(err)
	}
	args, ok := data["functionArgs"].([]interface{})
	if !ok {
		panic("Invalid functionArgs")
	}
	return args
}
