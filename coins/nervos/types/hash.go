package types

import "errors"

var (
	HashTypeData  ScriptHashType = "data"
	HashTypeData1 ScriptHashType = "data1"
	HashTypeType  ScriptHashType = "type"

	DepTypeDepGroup DepType = "dep_group"
	DepTypeCode     DepType = "code"
)

func SerializeHashType(hashType ScriptHashType) (string, error) {
	if HashTypeData == hashType {
		return "00", nil
	} else if HashTypeType == hashType {
		return "01", nil
	} else if HashTypeData1 == hashType {
		return "02", nil
	}

	return "", errors.New("Invalid script hash_type: " + string(hashType))
}

func DeserializeHashType(hashType string) (ScriptHashType, error) {
	if "00" == hashType {
		return HashTypeData, nil
	} else if "01" == hashType {
		return HashTypeType, nil
	} else if "02" == hashType {
		return HashTypeData1, nil
	}

	return "", errors.New("Invalid script hash_type: " + hashType)
}
