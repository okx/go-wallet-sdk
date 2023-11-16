// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc
package types

const (
	// base58 prefixes for 4 byte hash magics
	CHAIN_ID_PREFIX = "Net"

	// base58 prefixes for 20 byte hash magics
	ED25519_PUBLIC_KEY_HASH_PREFIX   = "tz1"
	SECP256K1_PUBLIC_KEY_HASH_PREFIX = "tz2"
	P256_PUBLIC_KEY_HASH_PREFIX      = "tz3"
	NOCURVE_PUBLIC_KEY_HASH_PREFIX   = "KT1"  // originated contract identifier
	BLINDED_PUBLIC_KEY_HASH_PREFIX   = "btz1" // blinded tz1

	// base58 prefixes for 16 byte hash magics
	ID_HASH_PREFIX = "id"

	// base58 prefixes for 32 byte hash magics
	BLOCK_HASH_PREFIX               = "B"
	OPERATION_HASH_PREFIX           = "o"
	OPERATION_LIST_HASH_PREFIX      = "Lo"
	OPERATION_LIST_LIST_HASH_PREFIX = "LLo"
	PROTOCOL_HASH_PREFIX            = "P"
	CONTEXT_HASH_PREFIX             = "Co"
	NONCE_HASH_PREFIX               = "nce"
	ED25519_SEED_PREFIX             = "edsk"
	ED25519_PUBLIC_KEY_PREFIX       = "edpk"
	SECP256K1_SECRET_KEY_PREFIX     = "spsk"
	P256_SECRET_KEY_PREFIX          = "p2sk"

	// base58 prefixes for 64 byte hash magics
	ED25519_SECRET_KEY_PREFIX  = "edsk"
	ED25519_SIGNATURE_PREFIX   = "edsig"
	SECP256K1_SIGNATURE_PREFIX = "spsig1"
	P256_SIGNATURE_PREFIX      = "p2sig"
	GENERIC_SIGNATURE_PREFIX   = "sig"

	// base58 prefixes for 33 byte hash magics
	SECP256K1_PUBLIC_KEY_PREFIX = "sppk"
	P256_PUBLIC_KEY_PREFIX      = "p2pk"
	SECP256K1_SCALAR_PREFIX     = "SSp"
	SECP256K1_ELEMENT_PREFIX    = "GSp"

	// base58 prefixes for 54 byte hash magics
	SCRIPT_EXPR_HASH_PREFIX = "expr"

	// base58 prefixes for 56 byte hash magics
	ED25519_ENCRYPTED_SEED_PREFIX         = "edesk"
	SECP256K1_ENCRYPTED_SECRET_KEY_PREFIX = "spesk"
	P256_ENCRYPTED_SECRET_KEY_PREFIX      = "p2esk"

	BLOCK_PAYLOAD_HASH_PREFIX                = "vh"  // "\001\106\242" (* vh(52) *)
	BLOCK_METADATA_HASH_PREFIX               = "bm"  // "\234\249" (* bm(52) *)
	OPERATION_METADATA_HASH_PREFIX           = "r"   // "\005\183" (* r(51) *)
	OPERATION_METADATA_LIST_HASH_PREFIX      = "Lr"  // "\134\039" (* Lr(52) *)
	OPERATION_METADATA_LIST_LIST_HASH_PREFIX = "LLr" // "\029\159\182" (* LLr(53) *)

	// base58 prefixes for 60 byte hash magics
	SECP256K1_ENCRYPTED_SCALAR_PREFIX = "seesk" // "\001\131\036\086\248" (* seesk(93) *)

	// base58 prefixes for Sapling byte hash magics
	SAPLING_SPENDING_KEY_PREFIX = "sask" // "\011\237\020\092" (* sask(241) *) // 169 bytes
	SAPLING_ADDRESS_PREFIX      = "zet1" // "\018\071\040\223" (* zet1(69) *) // 43 bytes

	// base58 prefixes for rollup hash magics
	BLS12_381_PUBLIC_KEY_HASH_PREFIX      = "tz4"
	GENERIC_AGGREGATE_SIGNATURE_PREFIX    = "asig"
	BLS12_381_SIGNATURE_PREFIX            = "BLsig"
	BLS12_381_PUBLIC_KEY_PREFIX           = "BLpk"
	BLS12_381_SECRET_KEY_PREFIX           = "BLsk"
	BLS12_381_ENCRYPTED_SECRET_KEY_PREFIX = "BLesk"
	TORU_ADDRESS_PREFIX                   = "txr1"
	TORU_INBOX_HASH_PREFIX                = "txi"
	TORU_MESSAGE_HASH_PREFIX              = "txm"
	TORU_COMMITMENT_HASH_PREFIX           = "txc"
	TORU_MESSAGE_RESULT_HASH_PREFIX       = "txmr"
	TORU_MESSAGE_RESULT_LIST_HASH_PREFIX  = "txM"
	TORU_WITHDRAW_LIST_HASH_PREFIX        = "txw"
)

var (
	// 20 byte hash magics
	ED25519_PUBLIC_KEY_HASH_ID   = []byte{0x06, 0xA1, 0x9F}       // "\006\161\159" (* tz1(36) *)
	SECP256K1_PUBLIC_KEY_HASH_ID = []byte{0x06, 0xA1, 0xA1}       // "\006\161\161" (* tz2(36) *)
	P256_PUBLIC_KEY_HASH_ID      = []byte{0x06, 0xA1, 0xA4}       // "\006\161\164" (* tz3(36) *)
	NOCURVE_PUBLIC_KEY_HASH_ID   = []byte{0x02, 0x5A, 0x79}       // "\002\090\121" (* KT1(36) *)
	BLINDED_PUBLIC_KEY_HASH_ID   = []byte{0x01, 0x02, 0x31, 0xDF} // "\002\090\121" (* btz1(37) *)

	// 4 byte hash magics
	CHAIN_ID = []byte{0x57, 0x52, 0x00} // "\087\082\000" (* Net(15) *)

	// 16 byte hash magics
	ID_HASH_ID = []byte{0x99, 0x67} // "\153\103" (* id(30) *) cryptobox_public_key_hash

	// 32 byte hash magics
	BLOCK_HASH_ID               = []byte{0x01, 0x34}       // "\001\052" (* B(51) *)
	OPERATION_HASH_ID           = []byte{0x05, 0x74}       // "\005\116" (* o(51) *)
	OPERATION_LIST_HASH_ID      = []byte{0x85, 0xE9}       // "\133\233" (* Lo(52) *)
	OPERATION_LIST_LIST_HASH_ID = []byte{0x1D, 0x9F, 0x6D} // "\029\159\109" (* LLo(53) *)
	PROTOCOL_HASH_ID            = []byte{0x02, 0xAA}       // "\002\170" (* P(51) *)
	CONTEXT_HASH_ID             = []byte{0x4F, 0xC7}       // "\079\199" (* Co(52) *)
	NONCE_HASH_ID               = []byte{0x45, 0xDC, 0xA9} // "\069\220\169" (* nce(53) *)

	ED25519_SEED_ID         = []byte{0x0D, 0x0F, 0x3A, 0x07} // "\013\015\058\007" (* edsk(54) *)
	ED25519_PUBLIC_KEY_ID   = []byte{0x0D, 0x0F, 0x25, 0xD9} // "\013\015\037\217" (* edpk(54) *)
	SECP256K1_SECRET_KEY_ID = []byte{0x11, 0xA2, 0xE0, 0xC9} // "\017\162\224\201" (* spsk(54) *)
	P256_SECRET_KEY_ID      = []byte{0x10, 0x51, 0xEE, 0xBD} // "\016\081\238\189" (* p2sk(54) *)

	// 33 byte hash magics
	SECP256K1_PUBLIC_KEY_ID = []byte{0x03, 0xFE, 0xE2, 0x56} // "\003\254\226\086" (* sppk(55) *)
	P256_PUBLIC_KEY_ID      = []byte{0x03, 0xB2, 0x8B, 0x7F} // "\003\178\139\127" (* p2pk(55) *)
	SECP256K1_SCALAR_ID     = []byte{0x26, 0xF8, 0x88}       // "\038\248\136" (* SSp(53) *)
	SECP256K1_ELEMENT_ID    = []byte{0x05, 0x5C, 0x00}       // "\005\092\000" (* GSp(54) *)

	// 64 byte hash magics
	ED25519_SECRET_KEY_ID  = []byte{0x2B, 0xF6, 0x4E, 0x07}       // "\043\246\078\007" (* edsk(98) *)
	ED25519_SIGNATURE_ID   = []byte{0x09, 0xF5, 0xCD, 0x86, 0x12} // "\009\245\205\134\018" (* edsig(99) *)
	SECP256K1_SIGNATURE_ID = []byte{0x0D, 0x73, 0x65, 0x13, 0x3F} // "\013\115\101\019\063" (* spsig1(99) *)
	P256_SIGNATURE_ID      = []byte{0x36, 0xF0, 0x2C, 0x34}       // "\054\240\044\052" (* p2sig(98) *)
	GENERIC_SIGNATURE_ID   = []byte{0x04, 0x82, 0x2B}             // "\004\130\043" (* sig(96) *)

	// 54 byte hash magics
	SCRIPT_EXPR_HASH_ID = []byte{0x0D, 0x2C, 0x40, 0x1B} // "\013\044\064\027" (* expr(54) *)

	BLOCK_PAYLOAD_HASH_ID                = []byte{0x01, 0x6A, 0xF2} // "\001\106\242" (* vh(52) *)
	BLOCK_METADATA_HASH_ID               = []byte{0xEA, 0xF9}       // "\234\249" (* bm(52) *)
	OPERATION_METADATA_HASH_ID           = []byte{0x05, 0xB7}       // "\005\183" (* r(51) *)
	OPERATION_METADATA_LIST_HASH_ID      = []byte{0x86, 0x27}       // "\134\039" (* Lr(52) *)
	OPERATION_METADATA_LIST_LIST_HASH_ID = []byte{0x1D, 0x9F, 0xB6} // "\029\159\182" (* LLr(53) *)

	// 60 byte hash magics
	SECP256K1_ENCRYPTED_SCALAR_ID = []byte{0x01, 0x83, 0x24, 0x56, 0xF8} // "\001\131\036\086\248" (* seesk(93) *)

	// Sapling magics
	SAPLING_SPENDING_KEY_ID = []byte{0x0b, 0xED, 0x14, 0x5C} // "\011\237\020\092" (* sask(241) *)
	SAPLING_ADDRESS_ID      = []byte{0x12, 0x47, 0x28, 0xDF} // "\018\071\040\223" (* zet1(69) *)

	// Rollup hash magics
	BLS12_381_PUBLIC_KEY_HASH_ID      = []byte{6, 161, 166}       // "\006\161\166" tz4(36) 20
	GENERIC_AGGREGATE_SIGNATURE_ID    = []byte{2, 075, 234, 101}  // "\002\075\234\101" asig(96) 141
	BLS12_381_SIGNATURE_ID            = []byte{40, 171, 064, 207} // "\040\171\064\207" (* BLsig(96) *) 142
	BLS12_381_PUBLIC_KEY_ID           = []byte{6, 149, 135, 204}  // "\006\149\135\204" (* BLpk(48) *) 76
	BLS12_381_SECRET_KEY_ID           = []byte{3, 150, 192, 40}   // "\003\150\192\040" (* BLsk(32) *) 54
	BLS12_381_ENCRYPTED_SECRET_KEY_ID = []byte{2, 5, 30, 53, 25}  // "\002\005\030\053\025" (* BLesk(58) *) 88
	TORU_ADDRESS_ID                   = []byte{1, 128, 120, 31}   // "\001\128\120\031" txr1(37) 20
	TORU_INBOX_HASH_ID                = []byte{79, 148, 196}      // "\079\148\196" txi(53) 32
	TORU_MESSAGE_HASH_ID              = []byte{79, 149, 30}       // "\079\149\030" txm(53) 32
	TORU_COMMITMENT_HASH_ID           = []byte{79, 148, 17}       // "\079\148\017" txc(53) 32
	TORU_MESSAGE_RESULT_HASH_ID       = []byte{18, 7, 206, 87}    // "\018\007\206\087" txmr(54) 32
	TORU_MESSAGE_RESULT_LIST_HASH_ID  = []byte{79, 146, 82}       // "\079\146\082" txM(53) 32
	TORU_WITHDRAW_LIST_HASH_ID        = []byte{79, 150, 72}       // "\079\150\072" txw(53) 32

	// 56 byte hash magics
	ED25519_ENCRYPTED_SEED_ID         = []byte{0x07, 0x5A, 0x3C, 0xB3, 0x29} // "\007\090\060\179\041" (* edesk(88) *)
	SECP256K1_ENCRYPTED_SECRET_KEY_ID = []byte{0x09, 0xED, 0xF1, 0xAE, 0x96} // "\009\237\241\174\150" (* spesk(88) *)
	P256_ENCRYPTED_SECRET_KEY_ID      = []byte{0x09, 0x30, 0x39, 0x73, 0xAB} // "\009\048\057\115\171" (* p2esk(88) *)
)
