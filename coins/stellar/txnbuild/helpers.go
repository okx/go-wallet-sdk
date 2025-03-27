package txnbuild

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/stellar/amount"
	"github.com/okx/go-wallet-sdk/coins/stellar/strkey"
	"github.com/okx/go-wallet-sdk/coins/stellar/support/errors"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
)

// validateStellarPublicKey returns an error if a public key is invalid. Otherwise, it returns nil.
// It is a wrapper around the IsValidEd25519PublicKey method of the strkey package.
func validateStellarPublicKey(publicKey string) error {
	if publicKey == "" {
		return errors.New("public key is undefined")
	}

	if !strkey.IsValidEd25519PublicKey(publicKey) {
		return errors.Errorf("%s is not a valid stellar public key", publicKey)
	}
	return nil
}

// validateStellarSignerKey returns an error if a signerkey is invalid. Otherwise, it returns nil.
func validateStellarSignerKey(signerKey string) error {
	if signerKey == "" {
		return errors.New("signer key is undefined")
	}

	var xdrKey xdr.SignerKey
	if err := xdrKey.SetAddress(signerKey); err != nil {
		return errors.Errorf("%s is not a valid stellar signer key", signerKey)
	}
	return nil
}

// validateStellarAsset checks if the asset supplied is a valid stellar Asset. It returns an error if the asset is
// nil, has an invalid asset code or issuer.
func validateStellarAsset(asset BasicAsset) error {
	if asset == nil {
		return errors.New("asset is undefined")
	}

	if asset.IsNative() {
		return nil
	}

	_, err := asset.GetType()
	if err != nil {
		return err
	}

	err = validateStellarPublicKey(asset.GetIssuer())
	if err != nil {
		return errors.Errorf("asset issuer: %s", err.Error())
	}

	return nil
}

// validateAmount checks if the provided value is a valid stellar amount, it returns an error if not.
// This is used to validate price and amount fields in structs.
func validateAmount(n interface{}) error {
	var stellarAmount int64
	// type switch can be extended to handle other types. Currently, the types for number values in the txnbuild
	// package are string or int64.
	switch value := n.(type) {
	case int64:
		stellarAmount = value
	case string:
		v, err := amount.ParseInt64(value)
		if err != nil {
			return err
		}
		stellarAmount = v
	default:
		return errors.Errorf("could not parse expected numeric value %v", n)
	}

	if stellarAmount < 0 {
		return errors.New("amount can not be negative")
	}
	return nil
}

// validateAssetCode checks if the provided asset is valid as an asset code.
// It returns an error if the asset is invalid.
// The asset must be non native (XLM) with a valid asset code.
func validateAssetCode(asset BasicAsset) error {
	// Note: we are not using validateStellarAsset() function for AllowTrust operations because it requires the
	//  following :
	// - asset is non-native
	// - asset code is valid
	// - asset issuer is not required. This is actually ignored by the operation
	if asset == nil {
		return errors.New("asset is undefined")
	}

	if asset.IsNative() {
		return errors.New("native (XLM) asset type is not allowed")
	}

	_, err := asset.GetType()
	if err != nil {
		return err
	}
	return nil
}

// validateChangeTrustAsset checks if the provided asset is valid for use in ChangeTrust operation.
// It returns an error if the asset is invalid.
// The asset must be non native (XLM) with a valid asset code and issuer.
func validateChangeTrustAsset(asset ChangeTrustAsset) error {
	// Note: we are not using validateStellarAsset() function for ChangeTrust operations because it requires the
	//  following :
	// - asset is non-native
	// - asset code is valid
	// - asset issuer is valid
	err := validateAssetCode(asset)
	if err != nil {
		return err
	}

	assetType, err := asset.GetType()
	if err != nil {
		return err
	} else if assetType == AssetTypePoolShare {
		// No issuer for these to validate.
		return nil
	}

	err = validateStellarPublicKey(asset.GetIssuer())
	if err != nil {
		return errors.Errorf("asset issuer: %s", err.Error())
	}

	return nil
}

// ValidationError is a custom error struct that holds validation errors of txnbuild's operation structs.
type ValidationError struct {
	Field   string // Field is the struct field on which the validation error occurred.
	Message string // Message is the validation error message.
}

// Error for ValidationError struct implements the error interface.
func (opError *ValidationError) Error() string {
	return fmt.Sprintf("Field: %s, Error: %s", opError.Field, opError.Message)
}

// NewValidationError creates a ValidationError struct with the provided field and message values.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ParseAssetString parses an asset string in canonical form (SEP-11) into an Asset structure.
// https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0011.md#asset
func ParseAssetString(canonical string) (Asset, error) {
	assets, err := xdr.BuildAssets(canonical)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing asset string")
	}

	if len(assets) != 1 {
		return nil, errors.New("error parsing out a single asset")
	}

	// The above returned a list, so we'll need to grab the first element.
	asset, err := assetFromXDR(assets[0])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing asset string via XDR types")
	}

	return asset, nil
}
