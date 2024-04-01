package aptos_types

import (
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type TypeTag interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}
type TypeTagBool struct{}
type TypeTagU8 struct{}
type TypeTagU16 struct{}
type TypeTagU32 struct{}
type TypeTagU64 struct{}
type TypeTagU128 struct{}
type TypeTagU256 struct{}
type TypeTagAddress struct{}
type TypeTagSigner struct{}
type TypeTagVector struct {
	Value TypeTag
}
type TypeTagStruct struct {
	Value StructTag
}
type StructTag struct {
	Address    AccountAddress
	ModuleName Identifier
	Name       Identifier
	TypeArgs   []TypeTag
}

func (o *StructTag) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := o.Address.Serialize(serializer); err != nil {
		return err
	}
	if err := o.ModuleName.Serialize(serializer); err != nil {
		return err
	}
	if err := o.Name.Serialize(serializer); err != nil {
		return err
	}
	if err := serializeVectorTypeTag(o.TypeArgs, serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *StructTag) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
func NewStructTagFromString(structTag string) (*StructTag, error) {
	parser, err := NewTypeTagParser(structTag, nil)
	if err != nil {
		return nil, err
	}
	typeTag, err := parser.ParseTypeTag()
	if err != nil {
		return nil, err
	}
	typeTagStruct, ok := typeTag.(*TypeTagStruct)
	if !ok {
		return nil, errors.New("parse struct tag err")
	}
	return &StructTag{
		Address:    typeTagStruct.Value.Address,
		ModuleName: typeTagStruct.Value.ModuleName,
		Name:       typeTagStruct.Value.Name,
		TypeArgs:   typeTagStruct.Value.TypeArgs,
	}, nil

}

func (o *TypeTagStruct) IsStringTypeTag() bool {
	return *CORE_CODE_ADDRESS == o.Value.Address && o.Value.ModuleName == "string" && o.Value.Name == "String"
}

func (o *TypeTagStruct) ShortFunctionName() string {
	return fmt.Sprintf("%v::%v::%v", o.Value.Address.ToShortString(), o.Value.ModuleName, o.Value.Name)
}

func (o *TypeTagBool) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(0)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagBool) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
func (o *TypeTagU8) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(1)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagU8) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagU16) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(8)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagU16) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagU32) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(9)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagU32) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagU64) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(2)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagU64) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagU128) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(3)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagU128) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagU256) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(10)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagU256) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagAddress) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(4)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagAddress) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagSigner) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(5)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagSigner) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagVector) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(6)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagVector) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TypeTagStruct) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(7)
	if err != nil {
		return err
	}
	err = o.Value.Serialize(serializer)
	if err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TypeTagStruct) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

/*func (o *StructTag) Serialize(serializer serde.Serializer) error {
	err := o.Address.Serialize(serializer)
	if err != nil {
		return err
	}
	err = o.ModuleName.Serialize(serializer)
	if err != nil {
		return err
	}
	err = o.Name.Serialize(serializer)
	if err != nil {
		return err
	}
	err = serializeVector(o.TypeArgs, serializer)
	if err != nil {
		return err
	}
	return nil
}*/

func serializeVector(v any, serializer serde.Serializer) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return fmt.Errorf("invalid vactor %v", v)
	}
	length := rv.Len()
	if err := serializer.SerializeVariantIndex(uint32(length)); err != nil {
		return err
	}
	for i := 0; i < length; i++ {
		t, ok := rv.Index(i).Interface().(TypeTag)
		if !ok {
			return fmt.Errorf("invalid element in vector %v", v)
		}
		err := t.Serialize(serializer)
		if err != nil {
			return err
		}
	}
	return nil
}

type Token struct {
	TokenType  string
	TokenValue string
}

func isWhiteSpace(c byte) bool {
	b, err := regexp.MatchString("\\s", string(c))
	if err != nil {
		return false
	}
	return b
}

func isValidAlphabetic(c byte) bool {
	b, err := regexp.MatchString("[_A-Za-z0-9]", string(c))
	if err != nil {
		return false
	}
	return b
}

// Returns Token and Token byte size
func nextToken(tagStr string, pos int) (*Token, int, error) {
	c := tagStr[pos]
	switch {
	case c == ':':
		if tagStr[pos+1] == ':' {
			return &Token{"COLON", "::"}, 2, nil
		} else {
			break
		}
	case c == '<':
		return &Token{"LT", "<"}, 1, nil
	case c == '>':
		return &Token{"GT", ">"}, 1, nil
	case c == ',':
		return &Token{"COMMA", ","}, 1, nil
	case isWhiteSpace(c):
		var i = pos + 1
		for ; i < len(tagStr) && isWhiteSpace(tagStr[i]); i++ {
		}
		return &Token{"SPACE", tagStr[pos:i]}, i - pos, nil
	case isValidAlphabetic(c):
		var i = pos + 1
		for ; i < len(tagStr) && isValidAlphabetic(tagStr[i]); i++ {
		}
		ok, err := isGeneric(tagStr[pos:i])
		if err != nil {
			return nil, 0, err
		}
		if ok {
			return &Token{"GENERIC", tagStr[pos:i]}, i - pos, nil
		}
		return &Token{"IDENT", tagStr[pos:i]}, i - pos, nil
	}
	return nil, 0, errors.New("unrecognized token")
}

func isGeneric(s string) (bool, error) {
	return regexp.MatchString("T\\d+", s)
}
func tokenize(tagStr string) ([]Token, error) {
	pos := 0
	tokens := []Token{}
	for pos < len(tagStr) {
		token, size, err := nextToken(tagStr, pos)
		if err != nil {
			return nil, err
		}
		if token.TokenType != "SPACE" {
			tokens = append(tokens, *token)
		}
		pos += size
	}
	return tokens, nil
}

type TypeTagParser struct {
	Tokens   []Token
	TypeTags []string
}

func NewTypeTagParser(tagStr string, typeTags []string) (*TypeTagParser, error) {
	tokens, err := tokenize(tagStr)
	if err != nil {
		return nil, err
	}
	return &TypeTagParser{Tokens: tokens, TypeTags: typeTags}, nil
}

func (p *TypeTagParser) shift() *Token {
	if len(p.Tokens) == 0 {
		return nil
	}
	t := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	return &t
}

func (p *TypeTagParser) consume(targetToken string) error {
	token := p.shift()
	if token == nil || token.TokenValue != targetToken {
		return ErrInvalidTypeTag
	}
	return nil
}

func (p *TypeTagParser) parseCommaList(endToken string, allowTraillingComma bool) ([]TypeTag, error) {
	if len(p.Tokens) <= 0 {
		return nil, ErrInvalidTypeTag
	}
	res := []TypeTag{}

	for p.Tokens[0].TokenValue != endToken {
		tag, err := p.ParseTypeTag()
		if err != nil {
			return nil, err
		}
		res = append(res, tag)

		if len(p.Tokens) > 0 && p.Tokens[0].TokenValue == endToken {
			break
		}
		err = p.consume(",")
		if err != nil {
			return nil, err
		}
		if len(p.Tokens) > 0 && p.Tokens[0].TokenValue == endToken && allowTraillingComma {
			break
		}

		if len(p.Tokens) <= 0 {
			return nil, ErrInvalidTypeTag
		}
	}

	return res, nil
}

func (p *TypeTagParser) ParseTypeTag() (TypeTag, error) {
	if len(p.Tokens) == 0 {
		return nil, ErrInvalidTypeTag
	}
	var err error

	token := p.shift()
	switch token.TokenValue {
	case "u8":
		return &TypeTagU8{}, nil
	case "u16":
		return &TypeTagU16{}, nil
	case "u32":
		return &TypeTagU32{}, nil
	case "u64":
		return &TypeTagU64{}, nil
	case "u128":
		return &TypeTagU128{}, nil
	case "u256":
		return &TypeTagU256{}, nil
	case "bool":
		return &TypeTagBool{}, nil
	case "address":
		return &TypeTagAddress{}, nil
	case "vector":
		err = p.consume("<")
		if err != nil {
			return nil, err
		}
		tag, err := p.ParseTypeTag()
		if err != nil {
			return nil, err
		}
		err = p.consume(">")
		if err != nil {
			return nil, err
		}
		return &TypeTagVector{Value: tag}, nil

	case "string":
		st := StructTag{
			Address:    *CORE_CODE_ADDRESS,
			ModuleName: "string",
			Name:       "string",
			TypeArgs:   []TypeTag{},
		}
		return &TypeTagStruct{
			Value: st,
		}, err
	}

	if token.TokenType == "IDENT" && (strings.HasPrefix(token.TokenValue, "0x") || strings.HasPrefix(token.TokenValue, "0X")) {
		address := token.TokenValue

		err = p.consume("::")
		if err != nil {
			return nil, err
		}
		moduleToken := p.shift()
		if moduleToken == nil || moduleToken.TokenType != "IDENT" {
			return nil, ErrInvalidTypeTag
		}

		err = p.consume("::")
		if err != nil {
			return nil, err
		}
		nameToken := p.shift()
		if nameToken == nil || nameToken.TokenType != "IDENT" {
			return nil, ErrInvalidTypeTag
		}

		// Objects can contain either concrete types e.g. 0x1::object::ObjectCore or generics e.g. T
		// Neither matter as we can't do type checks, so just the address applies and we consume the entire generic.
		// TODO: Support parsing structs that don't come from core code address

		addr, err := FromHex(address)
		if err != nil {
			return nil, err
		}
		if *CORE_CODE_ADDRESS == *addr && moduleToken.TokenValue == "object" && nameToken.TokenValue == "Object" {
			err = p.consumeWholeGeneric()
			if err != nil {
				return nil, err
			}
			return &TypeTagAddress{}, nil
		}

		tyArgs := []TypeTag{}
		// Check if the struct has ty args
		if len(p.Tokens) > 0 && p.Tokens[0].TokenValue == "<" {
			err = p.consume("<")
			if err != nil {
				return nil, err
			}
			tyArgs, err = p.parseCommaList(">", true)
			if err != nil {
				return nil, err
			}
			err = p.consume(">")
			if err != nil {
				return nil, err
			}
		}
		addr, err = FromHex(address)
		if err != nil {
			return nil, err
		}
		st := StructTag{
			Address:    *addr,
			ModuleName: Identifier(moduleToken.TokenValue),
			Name:       Identifier(nameToken.TokenValue),
			TypeArgs:   tyArgs,
		}
		return &TypeTagStruct{
			Value: st,
		}, nil
	}
	if token.TokenType == "GENERIC" {
		if len(p.TypeTags) == 0 {
			return nil, errors.New("can't convert generic type since no typeTags were specified")
		}
		// a generic tokenVal has the format of `T<digit>`, for example `T1`.
		// The digit (i.e 1) indicates the the index of this type in the typeTags array.
		// For a tokenVal == T1, should be parsed as the type in typeTags[1]
		idx, err := strconv.Atoi(token.TokenValue[1:])
		if err != nil {
			return nil, err
		}
		parser, err := NewTypeTagParser(p.TypeTags[idx], nil)
		if err != nil {
			return nil, err
		}
		return parser.ParseTypeTag()
	}
	return nil, ErrInvalidTypeTag
}

func (p *TypeTagParser) consumeWholeGeneric() error {
	err := p.consume("<")
	if err != nil {
		return err
	}
	for p.Tokens[0].TokenValue != ">" {
		if p.Tokens[0].TokenValue == "<" {
			err := p.consumeWholeGeneric()
			if err != nil {
				return err
			}
		} else {
			p.shift()
		}
	}
	err = p.consume(">")
	if err != nil {
		return err
	}
	return nil
}
