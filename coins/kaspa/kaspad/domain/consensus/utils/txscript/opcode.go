// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/constants"
	"hash"

	"golang.org/x/crypto/blake2b"
)

// An opcode defines the information related to a txscript opcode. opfunc, if
// present, is the function to call to perform the opcode on the script. The
// current script is passed in as a slice with the first member being the opcode
// itself.
type opcode struct {
	value  byte
	name   string
	length int
	opfunc func(*parsedOpcode, *Engine) error
}

// These constants are the values of the kaspa script opcodes.
const (
	Op0                   = 0x00 // 0
	OpFalse               = 0x00 // 0 - AKA Op0
	OpData1               = 0x01 // 1
	OpData2               = 0x02 // 2
	OpData3               = 0x03 // 3
	OpData4               = 0x04 // 4
	OpData5               = 0x05 // 5
	OpData6               = 0x06 // 6
	OpData7               = 0x07 // 7
	OpData8               = 0x08 // 8
	OpData9               = 0x09 // 9
	OpData10              = 0x0a // 10
	OpData11              = 0x0b // 11
	OpData12              = 0x0c // 12
	OpData13              = 0x0d // 13
	OpData14              = 0x0e // 14
	OpData15              = 0x0f // 15
	OpData16              = 0x10 // 16
	OpData17              = 0x11 // 17
	OpData18              = 0x12 // 18
	OpData19              = 0x13 // 19
	OpData20              = 0x14 // 20
	OpData21              = 0x15 // 21
	OpData22              = 0x16 // 22
	OpData23              = 0x17 // 23
	OpData24              = 0x18 // 24
	OpData25              = 0x19 // 25
	OpData26              = 0x1a // 26
	OpData27              = 0x1b // 27
	OpData28              = 0x1c // 28
	OpData29              = 0x1d // 29
	OpData30              = 0x1e // 30
	OpData31              = 0x1f // 31
	OpData32              = 0x20 // 32
	OpData33              = 0x21 // 33
	OpData34              = 0x22 // 34
	OpData35              = 0x23 // 35
	OpData36              = 0x24 // 36
	OpData37              = 0x25 // 37
	OpData38              = 0x26 // 38
	OpData39              = 0x27 // 39
	OpData40              = 0x28 // 40
	OpData41              = 0x29 // 41
	OpData42              = 0x2a // 42
	OpData43              = 0x2b // 43
	OpData44              = 0x2c // 44
	OpData45              = 0x2d // 45
	OpData46              = 0x2e // 46
	OpData47              = 0x2f // 47
	OpData48              = 0x30 // 48
	OpData49              = 0x31 // 49
	OpData50              = 0x32 // 50
	OpData51              = 0x33 // 51
	OpData52              = 0x34 // 52
	OpData53              = 0x35 // 53
	OpData54              = 0x36 // 54
	OpData55              = 0x37 // 55
	OpData56              = 0x38 // 56
	OpData57              = 0x39 // 57
	OpData58              = 0x3a // 58
	OpData59              = 0x3b // 59
	OpData60              = 0x3c // 60
	OpData61              = 0x3d // 61
	OpData62              = 0x3e // 62
	OpData63              = 0x3f // 63
	OpData64              = 0x40 // 64
	OpData65              = 0x41 // 65
	OpData66              = 0x42 // 66
	OpData67              = 0x43 // 67
	OpData68              = 0x44 // 68
	OpData69              = 0x45 // 69
	OpData70              = 0x46 // 70
	OpData71              = 0x47 // 71
	OpData72              = 0x48 // 72
	OpData73              = 0x49 // 73
	OpData74              = 0x4a // 74
	OpData75              = 0x4b // 75
	OpPushData1           = 0x4c // 76
	OpPushData2           = 0x4d // 77
	OpPushData4           = 0x4e // 78
	Op1Negate             = 0x4f // 79
	OpReserved            = 0x50 // 80
	Op1                   = 0x51 // 81 - AKA OpTrue
	OpTrue                = 0x51 // 81
	Op2                   = 0x52 // 82
	Op3                   = 0x53 // 83
	Op4                   = 0x54 // 84
	Op5                   = 0x55 // 85
	Op6                   = 0x56 // 86
	Op7                   = 0x57 // 87
	Op8                   = 0x58 // 88
	Op9                   = 0x59 // 89
	Op10                  = 0x5a // 90
	Op11                  = 0x5b // 91
	Op12                  = 0x5c // 92
	Op13                  = 0x5d // 93
	Op14                  = 0x5e // 94
	Op15                  = 0x5f // 95
	Op16                  = 0x60 // 96
	OpNop                 = 0x61 // 97
	OpVer                 = 0x62 // 98
	OpIf                  = 0x63 // 99
	OpNotIf               = 0x64 // 100
	OpVerIf               = 0x65 // 101
	OpVerNotIf            = 0x66 // 102
	OpElse                = 0x67 // 103
	OpEndIf               = 0x68 // 104
	OpVerify              = 0x69 // 105
	OpReturn              = 0x6a // 106
	OpToAltStack          = 0x6b // 107
	OpFromAltStack        = 0x6c // 108
	Op2Drop               = 0x6d // 109
	Op2Dup                = 0x6e // 110
	Op3Dup                = 0x6f // 111
	Op2Over               = 0x70 // 112
	Op2Rot                = 0x71 // 113
	Op2Swap               = 0x72 // 114
	OpIfDup               = 0x73 // 115
	OpDepth               = 0x74 // 116
	OpDrop                = 0x75 // 117
	OpDup                 = 0x76 // 118
	OpNip                 = 0x77 // 119
	OpOver                = 0x78 // 120
	OpPick                = 0x79 // 121
	OpRoll                = 0x7a // 122
	OpRot                 = 0x7b // 123
	OpSwap                = 0x7c // 124
	OpTuck                = 0x7d // 125
	OpCat                 = 0x7e // 126
	OpSubStr              = 0x7f // 127
	OpLeft                = 0x80 // 128
	OpRight               = 0x81 // 129
	OpSize                = 0x82 // 130
	OpInvert              = 0x83 // 131
	OpAnd                 = 0x84 // 132
	OpOr                  = 0x85 // 133
	OpXor                 = 0x86 // 134
	OpEqual               = 0x87 // 135
	OpEqualVerify         = 0x88 // 136
	OpReserved1           = 0x89 // 137
	OpReserved2           = 0x8a // 138
	Op1Add                = 0x8b // 139
	Op1Sub                = 0x8c // 140
	Op2Mul                = 0x8d // 141
	Op2Div                = 0x8e // 142
	OpNegate              = 0x8f // 143
	OpAbs                 = 0x90 // 144
	OpNot                 = 0x91 // 145
	Op0NotEqual           = 0x92 // 146
	OpAdd                 = 0x93 // 147
	OpSub                 = 0x94 // 148
	OpMul                 = 0x95 // 149
	OpDiv                 = 0x96 // 150
	OpMod                 = 0x97 // 151
	OpLShift              = 0x98 // 152
	OpRShift              = 0x99 // 153
	OpBoolAnd             = 0x9a // 154
	OpBoolOr              = 0x9b // 155
	OpNumEqual            = 0x9c // 156
	OpNumEqualVerify      = 0x9d // 157
	OpNumNotEqual         = 0x9e // 158
	OpLessThan            = 0x9f // 159
	OpGreaterThan         = 0xa0 // 160
	OpLessThanOrEqual     = 0xa1 // 161
	OpGreaterThanOrEqual  = 0xa2 // 162
	OpMin                 = 0xa3 // 163
	OpMax                 = 0xa4 // 164
	OpWithin              = 0xa5 // 165
	OpUnknown166          = 0xa6 // 166
	OpUnknown167          = 0xa7 // 167
	OpSHA256              = 0xa8 // 168
	OpCheckMultiSigECDSA  = 0xa9 // 169
	OpBlake2b             = 0xaa // 170
	OpCheckSigECDSA       = 0xab // 171
	OpCheckSig            = 0xac // 172
	OpCheckSigVerify      = 0xad // 173
	OpCheckMultiSig       = 0xae // 174
	OpCheckMultiSigVerify = 0xaf // 175
	OpCheckLockTimeVerify = 0xb0 // 176
	OpCheckSequenceVerify = 0xb1 // 177
	OpUnknown178          = 0xb2 // 178
	OpUnknown179          = 0xb3 // 179
	OpUnknown180          = 0xb4 // 180
	OpUnknown181          = 0xb5 // 181
	OpUnknown182          = 0xb6 // 182
	OpUnknown183          = 0xb7 // 183
	OpUnknown184          = 0xb8 // 184
	OpUnknown185          = 0xb9 // 185
	OpUnknown186          = 0xba // 186
	OpUnknown187          = 0xbb // 187
	OpUnknown188          = 0xbc // 188
	OpUnknown189          = 0xbd // 189
	OpUnknown190          = 0xbe // 190
	OpUnknown191          = 0xbf // 191
	OpUnknown192          = 0xc0 // 192
	OpUnknown193          = 0xc1 // 193
	OpUnknown194          = 0xc2 // 194
	OpUnknown195          = 0xc3 // 195
	OpUnknown196          = 0xc4 // 196
	OpUnknown197          = 0xc5 // 197
	OpUnknown198          = 0xc6 // 198
	OpUnknown199          = 0xc7 // 199
	OpUnknown200          = 0xc8 // 200
	OpUnknown201          = 0xc9 // 201
	OpUnknown202          = 0xca // 202
	OpUnknown203          = 0xcb // 203
	OpUnknown204          = 0xcc // 204
	OpUnknown205          = 0xcd // 205
	OpUnknown206          = 0xce // 206
	OpUnknown207          = 0xcf // 207
	OpUnknown208          = 0xd0 // 208
	OpUnknown209          = 0xd1 // 209
	OpUnknown210          = 0xd2 // 210
	OpUnknown211          = 0xd3 // 211
	OpUnknown212          = 0xd4 // 212
	OpUnknown213          = 0xd5 // 213
	OpUnknown214          = 0xd6 // 214
	OpUnknown215          = 0xd7 // 215
	OpUnknown216          = 0xd8 // 216
	OpUnknown217          = 0xd9 // 217
	OpUnknown218          = 0xda // 218
	OpUnknown219          = 0xdb // 219
	OpUnknown220          = 0xdc // 220
	OpUnknown221          = 0xdd // 221
	OpUnknown222          = 0xde // 222
	OpUnknown223          = 0xdf // 223
	OpUnknown224          = 0xe0 // 224
	OpUnknown225          = 0xe1 // 225
	OpUnknown226          = 0xe2 // 226
	OpUnknown227          = 0xe3 // 227
	OpUnknown228          = 0xe4 // 228
	OpUnknown229          = 0xe5 // 229
	OpUnknown230          = 0xe6 // 230
	OpUnknown231          = 0xe7 // 231
	OpUnknown232          = 0xe8 // 232
	OpUnknown233          = 0xe9 // 233
	OpUnknown234          = 0xea // 234
	OpUnknown235          = 0xeb // 235
	OpUnknown236          = 0xec // 236
	OpUnknown237          = 0xed // 237
	OpUnknown238          = 0xee // 238
	OpUnknown239          = 0xef // 239
	OpUnknown240          = 0xf0 // 240
	OpUnknown241          = 0xf1 // 241
	OpUnknown242          = 0xf2 // 242
	OpUnknown243          = 0xf3 // 243
	OpUnknown244          = 0xf4 // 244
	OpUnknown245          = 0xf5 // 245
	OpUnknown246          = 0xf6 // 246
	OpUnknown247          = 0xf7 // 247
	OpUnknown248          = 0xf8 // 248
	OpUnknown249          = 0xf9 // 249
	OpSmallInteger        = 0xfa // 250
	OpPubKeys             = 0xfb // 251
	OpUnknown252          = 0xfc // 252
	OpPubKeyHash          = 0xfd // 253
	OpPubKey              = 0xfe // 254
	OpInvalidOpCode       = 0xff // 255
)

// Conditional execution constants.
const (
	OpCondFalse = 0
	OpCondTrue  = 1
	OpCondSkip  = 2
)

// opcodeArray holds details about all possible opcodes such as how many bytes
// the opcode and any associated data should take, its human-readable name, and
// the handler function.
var opcodeArray = [256]opcode{
	// Data push opcodes.
	OpFalse:     {OpFalse, "OP_0", 1, opcodeFalse},
	OpData1:     {OpData1, "OP_DATA_1", 2, opcodePushData},
	OpData2:     {OpData2, "OP_DATA_2", 3, opcodePushData},
	OpData3:     {OpData3, "OP_DATA_3", 4, opcodePushData},
	OpData4:     {OpData4, "OP_DATA_4", 5, opcodePushData},
	OpData5:     {OpData5, "OP_DATA_5", 6, opcodePushData},
	OpData6:     {OpData6, "OP_DATA_6", 7, opcodePushData},
	OpData7:     {OpData7, "OP_DATA_7", 8, opcodePushData},
	OpData8:     {OpData8, "OP_DATA_8", 9, opcodePushData},
	OpData9:     {OpData9, "OP_DATA_9", 10, opcodePushData},
	OpData10:    {OpData10, "OP_DATA_10", 11, opcodePushData},
	OpData11:    {OpData11, "OP_DATA_11", 12, opcodePushData},
	OpData12:    {OpData12, "OP_DATA_12", 13, opcodePushData},
	OpData13:    {OpData13, "OP_DATA_13", 14, opcodePushData},
	OpData14:    {OpData14, "OP_DATA_14", 15, opcodePushData},
	OpData15:    {OpData15, "OP_DATA_15", 16, opcodePushData},
	OpData16:    {OpData16, "OP_DATA_16", 17, opcodePushData},
	OpData17:    {OpData17, "OP_DATA_17", 18, opcodePushData},
	OpData18:    {OpData18, "OP_DATA_18", 19, opcodePushData},
	OpData19:    {OpData19, "OP_DATA_19", 20, opcodePushData},
	OpData20:    {OpData20, "OP_DATA_20", 21, opcodePushData},
	OpData21:    {OpData21, "OP_DATA_21", 22, opcodePushData},
	OpData22:    {OpData22, "OP_DATA_22", 23, opcodePushData},
	OpData23:    {OpData23, "OP_DATA_23", 24, opcodePushData},
	OpData24:    {OpData24, "OP_DATA_24", 25, opcodePushData},
	OpData25:    {OpData25, "OP_DATA_25", 26, opcodePushData},
	OpData26:    {OpData26, "OP_DATA_26", 27, opcodePushData},
	OpData27:    {OpData27, "OP_DATA_27", 28, opcodePushData},
	OpData28:    {OpData28, "OP_DATA_28", 29, opcodePushData},
	OpData29:    {OpData29, "OP_DATA_29", 30, opcodePushData},
	OpData30:    {OpData30, "OP_DATA_30", 31, opcodePushData},
	OpData31:    {OpData31, "OP_DATA_31", 32, opcodePushData},
	OpData32:    {OpData32, "OP_DATA_32", 33, opcodePushData},
	OpData33:    {OpData33, "OP_DATA_33", 34, opcodePushData},
	OpData34:    {OpData34, "OP_DATA_34", 35, opcodePushData},
	OpData35:    {OpData35, "OP_DATA_35", 36, opcodePushData},
	OpData36:    {OpData36, "OP_DATA_36", 37, opcodePushData},
	OpData37:    {OpData37, "OP_DATA_37", 38, opcodePushData},
	OpData38:    {OpData38, "OP_DATA_38", 39, opcodePushData},
	OpData39:    {OpData39, "OP_DATA_39", 40, opcodePushData},
	OpData40:    {OpData40, "OP_DATA_40", 41, opcodePushData},
	OpData41:    {OpData41, "OP_DATA_41", 42, opcodePushData},
	OpData42:    {OpData42, "OP_DATA_42", 43, opcodePushData},
	OpData43:    {OpData43, "OP_DATA_43", 44, opcodePushData},
	OpData44:    {OpData44, "OP_DATA_44", 45, opcodePushData},
	OpData45:    {OpData45, "OP_DATA_45", 46, opcodePushData},
	OpData46:    {OpData46, "OP_DATA_46", 47, opcodePushData},
	OpData47:    {OpData47, "OP_DATA_47", 48, opcodePushData},
	OpData48:    {OpData48, "OP_DATA_48", 49, opcodePushData},
	OpData49:    {OpData49, "OP_DATA_49", 50, opcodePushData},
	OpData50:    {OpData50, "OP_DATA_50", 51, opcodePushData},
	OpData51:    {OpData51, "OP_DATA_51", 52, opcodePushData},
	OpData52:    {OpData52, "OP_DATA_52", 53, opcodePushData},
	OpData53:    {OpData53, "OP_DATA_53", 54, opcodePushData},
	OpData54:    {OpData54, "OP_DATA_54", 55, opcodePushData},
	OpData55:    {OpData55, "OP_DATA_55", 56, opcodePushData},
	OpData56:    {OpData56, "OP_DATA_56", 57, opcodePushData},
	OpData57:    {OpData57, "OP_DATA_57", 58, opcodePushData},
	OpData58:    {OpData58, "OP_DATA_58", 59, opcodePushData},
	OpData59:    {OpData59, "OP_DATA_59", 60, opcodePushData},
	OpData60:    {OpData60, "OP_DATA_60", 61, opcodePushData},
	OpData61:    {OpData61, "OP_DATA_61", 62, opcodePushData},
	OpData62:    {OpData62, "OP_DATA_62", 63, opcodePushData},
	OpData63:    {OpData63, "OP_DATA_63", 64, opcodePushData},
	OpData64:    {OpData64, "OP_DATA_64", 65, opcodePushData},
	OpData65:    {OpData65, "OP_DATA_65", 66, opcodePushData},
	OpData66:    {OpData66, "OP_DATA_66", 67, opcodePushData},
	OpData67:    {OpData67, "OP_DATA_67", 68, opcodePushData},
	OpData68:    {OpData68, "OP_DATA_68", 69, opcodePushData},
	OpData69:    {OpData69, "OP_DATA_69", 70, opcodePushData},
	OpData70:    {OpData70, "OP_DATA_70", 71, opcodePushData},
	OpData71:    {OpData71, "OP_DATA_71", 72, opcodePushData},
	OpData72:    {OpData72, "OP_DATA_72", 73, opcodePushData},
	OpData73:    {OpData73, "OP_DATA_73", 74, opcodePushData},
	OpData74:    {OpData74, "OP_DATA_74", 75, opcodePushData},
	OpData75:    {OpData75, "OP_DATA_75", 76, opcodePushData},
	OpPushData1: {OpPushData1, "OP_PUSHDATA1", -1, opcodePushData},
	OpPushData2: {OpPushData2, "OP_PUSHDATA2", -2, opcodePushData},
	OpPushData4: {OpPushData4, "OP_PUSHDATA4", -4, opcodePushData},
	Op1Negate:   {Op1Negate, "OP_1NEGATE", 1, opcode1Negate},
	OpReserved:  {OpReserved, "OP_RESERVED", 1, opcodeReserved},
	OpTrue:      {OpTrue, "OP_1", 1, opcodeN},
	Op2:         {Op2, "OP_2", 1, opcodeN},
	Op3:         {Op3, "OP_3", 1, opcodeN},
	Op4:         {Op4, "OP_4", 1, opcodeN},
	Op5:         {Op5, "OP_5", 1, opcodeN},
	Op6:         {Op6, "OP_6", 1, opcodeN},
	Op7:         {Op7, "OP_7", 1, opcodeN},
	Op8:         {Op8, "OP_8", 1, opcodeN},
	Op9:         {Op9, "OP_9", 1, opcodeN},
	Op10:        {Op10, "OP_10", 1, opcodeN},
	Op11:        {Op11, "OP_11", 1, opcodeN},
	Op12:        {Op12, "OP_12", 1, opcodeN},
	Op13:        {Op13, "OP_13", 1, opcodeN},
	Op14:        {Op14, "OP_14", 1, opcodeN},
	Op15:        {Op15, "OP_15", 1, opcodeN},
	Op16:        {Op16, "OP_16", 1, opcodeN},

	// Control opcodes.
	OpNop:                 {OpNop, "OP_NOP", 1, opcodeNop},
	OpVer:                 {OpVer, "OP_VER", 1, opcodeReserved},
	OpIf:                  {OpIf, "OP_IF", 1, opcodeIf},
	OpNotIf:               {OpNotIf, "OP_NOTIF", 1, opcodeNotIf},
	OpVerIf:               {OpVerIf, "OP_VERIF", 1, opcodeReserved},
	OpVerNotIf:            {OpVerNotIf, "OP_VERNOTIF", 1, opcodeReserved},
	OpElse:                {OpElse, "OP_ELSE", 1, opcodeElse},
	OpEndIf:               {OpEndIf, "OP_ENDIF", 1, opcodeEndif},
	OpVerify:              {OpVerify, "OP_VERIFY", 1, opcodeVerify},
	OpReturn:              {OpReturn, "OP_RETURN", 1, opcodeReturn},
	OpCheckLockTimeVerify: {OpCheckLockTimeVerify, "OP_CHECKLOCKTIMEVERIFY", 1, opcodeCheckLockTimeVerify},
	OpCheckSequenceVerify: {OpCheckSequenceVerify, "OP_CHECKSEQUENCEVERIFY", 1, opcodeCheckSequenceVerify},

	// Stack opcodes.
	OpToAltStack:   {OpToAltStack, "OP_TOALTSTACK", 1, opcodeToAltStack},
	OpFromAltStack: {OpFromAltStack, "OP_FROMALTSTACK", 1, opcodeFromAltStack},
	Op2Drop:        {Op2Drop, "OP_2DROP", 1, opcode2Drop},
	Op2Dup:         {Op2Dup, "OP_2DUP", 1, opcode2Dup},
	Op3Dup:         {Op3Dup, "OP_3DUP", 1, opcode3Dup},
	Op2Over:        {Op2Over, "OP_2OVER", 1, opcode2Over},
	Op2Rot:         {Op2Rot, "OP_2ROT", 1, opcode2Rot},
	Op2Swap:        {Op2Swap, "OP_2SWAP", 1, opcode2Swap},
	OpIfDup:        {OpIfDup, "OP_IFDUP", 1, opcodeIfDup},
	OpDepth:        {OpDepth, "OP_DEPTH", 1, opcodeDepth},
	OpDrop:         {OpDrop, "OP_DROP", 1, opcodeDrop},
	OpDup:          {OpDup, "OP_DUP", 1, opcodeDup},
	OpNip:          {OpNip, "OP_NIP", 1, opcodeNip},
	OpOver:         {OpOver, "OP_OVER", 1, opcodeOver},
	OpPick:         {OpPick, "OP_PICK", 1, opcodePick},
	OpRoll:         {OpRoll, "OP_ROLL", 1, opcodeRoll},
	OpRot:          {OpRot, "OP_ROT", 1, opcodeRot},
	OpSwap:         {OpSwap, "OP_SWAP", 1, opcodeSwap},
	OpTuck:         {OpTuck, "OP_TUCK", 1, opcodeTuck},

	// Splice opcodes.
	OpCat:    {OpCat, "OP_CAT", 1, opcodeDisabled},
	OpSubStr: {OpSubStr, "OP_SUBSTR", 1, opcodeDisabled},
	OpLeft:   {OpLeft, "OP_LEFT", 1, opcodeDisabled},
	OpRight:  {OpRight, "OP_RIGHT", 1, opcodeDisabled},
	OpSize:   {OpSize, "OP_SIZE", 1, opcodeSize},

	// Bitwise logic opcodes.
	OpInvert:      {OpInvert, "OP_INVERT", 1, opcodeDisabled},
	OpAnd:         {OpAnd, "OP_AND", 1, opcodeDisabled},
	OpOr:          {OpOr, "OP_OR", 1, opcodeDisabled},
	OpXor:         {OpXor, "OP_XOR", 1, opcodeDisabled},
	OpEqual:       {OpEqual, "OP_EQUAL", 1, opcodeEqual},
	OpEqualVerify: {OpEqualVerify, "OP_EQUALVERIFY", 1, opcodeEqualVerify},
	OpReserved1:   {OpReserved1, "OP_RESERVED1", 1, opcodeReserved},
	OpReserved2:   {OpReserved2, "OP_RESERVED2", 1, opcodeReserved},

	// Numeric related opcodes.
	Op1Add:               {Op1Add, "OP_1ADD", 1, opcode1Add},
	Op1Sub:               {Op1Sub, "OP_1SUB", 1, opcode1Sub},
	Op2Mul:               {Op2Mul, "OP_2MUL", 1, opcodeDisabled},
	Op2Div:               {Op2Div, "OP_2DIV", 1, opcodeDisabled},
	OpNegate:             {OpNegate, "OP_NEGATE", 1, opcodeNegate},
	OpAbs:                {OpAbs, "OP_ABS", 1, opcodeAbs},
	OpNot:                {OpNot, "OP_NOT", 1, opcodeNot},
	Op0NotEqual:          {Op0NotEqual, "OP_0NOTEQUAL", 1, opcode0NotEqual},
	OpAdd:                {OpAdd, "OP_ADD", 1, opcodeAdd},
	OpSub:                {OpSub, "OP_SUB", 1, opcodeSub},
	OpMul:                {OpMul, "OP_MUL", 1, opcodeDisabled},
	OpDiv:                {OpDiv, "OP_DIV", 1, opcodeDisabled},
	OpMod:                {OpMod, "OP_MOD", 1, opcodeDisabled},
	OpLShift:             {OpLShift, "OP_LSHIFT", 1, opcodeDisabled},
	OpRShift:             {OpRShift, "OP_RSHIFT", 1, opcodeDisabled},
	OpBoolAnd:            {OpBoolAnd, "OP_BOOLAND", 1, opcodeBoolAnd},
	OpBoolOr:             {OpBoolOr, "OP_BOOLOR", 1, opcodeBoolOr},
	OpNumEqual:           {OpNumEqual, "OP_NUMEQUAL", 1, opcodeNumEqual},
	OpNumEqualVerify:     {OpNumEqualVerify, "OP_NUMEQUALVERIFY", 1, opcodeNumEqualVerify},
	OpNumNotEqual:        {OpNumNotEqual, "OP_NUMNOTEQUAL", 1, opcodeNumNotEqual},
	OpLessThan:           {OpLessThan, "OP_LESSTHAN", 1, opcodeLessThan},
	OpGreaterThan:        {OpGreaterThan, "OP_GREATERTHAN", 1, opcodeGreaterThan},
	OpLessThanOrEqual:    {OpLessThanOrEqual, "OP_LESSTHANOREQUAL", 1, opcodeLessThanOrEqual},
	OpGreaterThanOrEqual: {OpGreaterThanOrEqual, "OP_GREATERTHANOREQUAL", 1, opcodeGreaterThanOrEqual},
	OpMin:                {OpMin, "OP_MIN", 1, opcodeMin},
	OpMax:                {OpMax, "OP_MAX", 1, opcodeMax},
	OpWithin:             {OpWithin, "OP_WITHIN", 1, opcodeWithin},

	// Undefined opcodes.
	OpUnknown166: {OpUnknown166, "OP_UNKNOWN166", 1, opcodeInvalid},
	OpUnknown167: {OpUnknown167, "OP_UNKNOWN167", 1, opcodeInvalid},
	OpUnknown178: {OpUnknown188, "OP_UNKNOWN178", 1, opcodeInvalid},
	OpUnknown179: {OpUnknown189, "OP_UNKNOWN179", 1, opcodeInvalid},
	OpUnknown180: {OpUnknown190, "OP_UNKNOWN180", 1, opcodeInvalid},
	OpUnknown181: {OpUnknown191, "OP_UNKNOWN181", 1, opcodeInvalid},
	OpUnknown182: {OpUnknown192, "OP_UNKNOWN182", 1, opcodeInvalid},
	OpUnknown183: {OpUnknown193, "OP_UNKNOWN183", 1, opcodeInvalid},
	OpUnknown184: {OpUnknown194, "OP_UNKNOWN184", 1, opcodeInvalid},
	OpUnknown185: {OpUnknown195, "OP_UNKNOWN185", 1, opcodeInvalid},
	OpUnknown186: {OpUnknown196, "OP_UNKNOWN186", 1, opcodeInvalid},
	OpUnknown187: {OpUnknown197, "OP_UNKNOWN187", 1, opcodeInvalid},
	OpUnknown188: {OpUnknown188, "OP_UNKNOWN188", 1, opcodeInvalid},
	OpUnknown189: {OpUnknown189, "OP_UNKNOWN189", 1, opcodeInvalid},
	OpUnknown190: {OpUnknown190, "OP_UNKNOWN190", 1, opcodeInvalid},
	OpUnknown191: {OpUnknown191, "OP_UNKNOWN191", 1, opcodeInvalid},
	OpUnknown192: {OpUnknown192, "OP_UNKNOWN192", 1, opcodeInvalid},
	OpUnknown193: {OpUnknown193, "OP_UNKNOWN193", 1, opcodeInvalid},
	OpUnknown194: {OpUnknown194, "OP_UNKNOWN194", 1, opcodeInvalid},
	OpUnknown195: {OpUnknown195, "OP_UNKNOWN195", 1, opcodeInvalid},
	OpUnknown196: {OpUnknown196, "OP_UNKNOWN196", 1, opcodeInvalid},
	OpUnknown197: {OpUnknown197, "OP_UNKNOWN197", 1, opcodeInvalid},
	OpUnknown198: {OpUnknown198, "OP_UNKNOWN198", 1, opcodeInvalid},
	OpUnknown199: {OpUnknown199, "OP_UNKNOWN199", 1, opcodeInvalid},
	OpUnknown200: {OpUnknown200, "OP_UNKNOWN200", 1, opcodeInvalid},
	OpUnknown201: {OpUnknown201, "OP_UNKNOWN201", 1, opcodeInvalid},
	OpUnknown202: {OpUnknown202, "OP_UNKNOWN202", 1, opcodeInvalid},
	OpUnknown203: {OpUnknown203, "OP_UNKNOWN203", 1, opcodeInvalid},
	OpUnknown204: {OpUnknown204, "OP_UNKNOWN204", 1, opcodeInvalid},
	OpUnknown205: {OpUnknown205, "OP_UNKNOWN205", 1, opcodeInvalid},
	OpUnknown206: {OpUnknown206, "OP_UNKNOWN206", 1, opcodeInvalid},
	OpUnknown207: {OpUnknown207, "OP_UNKNOWN207", 1, opcodeInvalid},
	OpUnknown208: {OpUnknown208, "OP_UNKNOWN208", 1, opcodeInvalid},
	OpUnknown209: {OpUnknown209, "OP_UNKNOWN209", 1, opcodeInvalid},
	OpUnknown210: {OpUnknown210, "OP_UNKNOWN210", 1, opcodeInvalid},
	OpUnknown211: {OpUnknown211, "OP_UNKNOWN211", 1, opcodeInvalid},
	OpUnknown212: {OpUnknown212, "OP_UNKNOWN212", 1, opcodeInvalid},
	OpUnknown213: {OpUnknown213, "OP_UNKNOWN213", 1, opcodeInvalid},
	OpUnknown214: {OpUnknown214, "OP_UNKNOWN214", 1, opcodeInvalid},
	OpUnknown215: {OpUnknown215, "OP_UNKNOWN215", 1, opcodeInvalid},
	OpUnknown216: {OpUnknown216, "OP_UNKNOWN216", 1, opcodeInvalid},
	OpUnknown217: {OpUnknown217, "OP_UNKNOWN217", 1, opcodeInvalid},
	OpUnknown218: {OpUnknown218, "OP_UNKNOWN218", 1, opcodeInvalid},
	OpUnknown219: {OpUnknown219, "OP_UNKNOWN219", 1, opcodeInvalid},
	OpUnknown220: {OpUnknown220, "OP_UNKNOWN220", 1, opcodeInvalid},
	OpUnknown221: {OpUnknown221, "OP_UNKNOWN221", 1, opcodeInvalid},
	OpUnknown222: {OpUnknown222, "OP_UNKNOWN222", 1, opcodeInvalid},
	OpUnknown223: {OpUnknown223, "OP_UNKNOWN223", 1, opcodeInvalid},
	OpUnknown224: {OpUnknown224, "OP_UNKNOWN224", 1, opcodeInvalid},
	OpUnknown225: {OpUnknown225, "OP_UNKNOWN225", 1, opcodeInvalid},
	OpUnknown226: {OpUnknown226, "OP_UNKNOWN226", 1, opcodeInvalid},
	OpUnknown227: {OpUnknown227, "OP_UNKNOWN227", 1, opcodeInvalid},
	OpUnknown228: {OpUnknown228, "OP_UNKNOWN228", 1, opcodeInvalid},
	OpUnknown229: {OpUnknown229, "OP_UNKNOWN229", 1, opcodeInvalid},
	OpUnknown230: {OpUnknown230, "OP_UNKNOWN230", 1, opcodeInvalid},
	OpUnknown231: {OpUnknown231, "OP_UNKNOWN231", 1, opcodeInvalid},
	OpUnknown232: {OpUnknown232, "OP_UNKNOWN232", 1, opcodeInvalid},
	OpUnknown233: {OpUnknown233, "OP_UNKNOWN233", 1, opcodeInvalid},
	OpUnknown234: {OpUnknown234, "OP_UNKNOWN234", 1, opcodeInvalid},
	OpUnknown235: {OpUnknown235, "OP_UNKNOWN235", 1, opcodeInvalid},
	OpUnknown236: {OpUnknown236, "OP_UNKNOWN236", 1, opcodeInvalid},
	OpUnknown237: {OpUnknown237, "OP_UNKNOWN237", 1, opcodeInvalid},
	OpUnknown238: {OpUnknown238, "OP_UNKNOWN238", 1, opcodeInvalid},
	OpUnknown239: {OpUnknown239, "OP_UNKNOWN239", 1, opcodeInvalid},
	OpUnknown240: {OpUnknown240, "OP_UNKNOWN240", 1, opcodeInvalid},
	OpUnknown241: {OpUnknown241, "OP_UNKNOWN241", 1, opcodeInvalid},
	OpUnknown242: {OpUnknown242, "OP_UNKNOWN242", 1, opcodeInvalid},
	OpUnknown243: {OpUnknown243, "OP_UNKNOWN243", 1, opcodeInvalid},
	OpUnknown244: {OpUnknown244, "OP_UNKNOWN244", 1, opcodeInvalid},
	OpUnknown245: {OpUnknown245, "OP_UNKNOWN245", 1, opcodeInvalid},
	OpUnknown246: {OpUnknown246, "OP_UNKNOWN246", 1, opcodeInvalid},
	OpUnknown247: {OpUnknown247, "OP_UNKNOWN247", 1, opcodeInvalid},
	OpUnknown248: {OpUnknown248, "OP_UNKNOWN248", 1, opcodeInvalid},
	OpUnknown249: {OpUnknown249, "OP_UNKNOWN249", 1, opcodeInvalid},

	OpSmallInteger: {OpSmallInteger, "OP_SMALLINTEGER", 1, opcodeInvalid},
	OpPubKeys:      {OpPubKeys, "OP_PUBKEYS", 1, opcodeInvalid},
	OpUnknown252:   {OpUnknown252, "OP_UNKNOWN252", 1, opcodeInvalid},
	OpPubKeyHash:   {OpPubKeyHash, "OP_PUBKEYHASH", 1, opcodeInvalid},
	OpPubKey:       {OpPubKey, "OP_PUBKEY", 1, opcodeInvalid},

	OpInvalidOpCode: {OpInvalidOpCode, "OP_INVALIDOPCODE", 1, opcodeInvalid},
}

// opcodeOnelineRepls defines opcode names which are replaced when doing a
// one-line disassembly. This is done to match the output of the reference
// implementation while not changing the opcode names in the nicer full
// disassembly.
var opcodeOnelineRepls = map[string]string{
	"OP_1NEGATE": "-1",
	"OP_0":       "0",
	"OP_1":       "1",
	"OP_2":       "2",
	"OP_3":       "3",
	"OP_4":       "4",
	"OP_5":       "5",
	"OP_6":       "6",
	"OP_7":       "7",
	"OP_8":       "8",
	"OP_9":       "9",
	"OP_10":      "10",
	"OP_11":      "11",
	"OP_12":      "12",
	"OP_13":      "13",
	"OP_14":      "14",
	"OP_15":      "15",
	"OP_16":      "16",
}

// parsedOpcode represents an opcode that has been parsed and includes any
// potential data associated with it.
type parsedOpcode struct {
	opcode *opcode
	data   []byte
}

// isDisabled returns whether or not the opcode is disabled and thus is always
// bad to see in the instruction stream (even if turned off by a conditional).
func (pop *parsedOpcode) isDisabled() bool {
	switch pop.opcode.value {
	case OpCat:
		return true
	case OpSubStr:
		return true
	case OpLeft:
		return true
	case OpRight:
		return true
	case OpInvert:
		return true
	case OpAnd:
		return true
	case OpOr:
		return true
	case OpXor:
		return true
	case Op2Mul:
		return true
	case Op2Div:
		return true
	case OpMul:
		return true
	case OpDiv:
		return true
	case OpMod:
		return true
	case OpLShift:
		return true
	case OpRShift:
		return true
	default:
		return false
	}
}

// alwaysIllegal returns whether or not the opcode is always illegal when passed
// over by the program counter even if in a non-executed branch (it isn't a
// coincidence that they are conditionals).
func (pop *parsedOpcode) alwaysIllegal() bool {
	switch pop.opcode.value {
	case OpVerIf:
		return true
	case OpVerNotIf:
		return true
	default:
		return false
	}
}

// isConditional returns whether or not the opcode is a conditional opcode which
// changes the conditional execution stack when executed.
func (pop *parsedOpcode) isConditional() bool {
	switch pop.opcode.value {
	case OpIf:
		return true
	case OpNotIf:
		return true
	case OpElse:
		return true
	case OpEndIf:
		return true
	default:
		return false
	}
}

// checkMinimalDataPush returns whether or not the current data push uses the
// smallest possible opcode to represent it. For example, the value 15 could
// be pushed with OP_DATA_1 15 (among other variations); however, OP_15 is a
// single opcode that represents the same value and is only a single byte versus
// two bytes.
func (pop *parsedOpcode) checkMinimalDataPush() error {
	data := pop.data
	dataLen := len(data)
	opcode := pop.opcode.value

	if dataLen == 0 && opcode != Op0 {
		str := fmt.Sprintf("zero length data push is encoded with "+
			"opcode %s instead of OP_0", pop.opcode.name)
		return scriptError(ErrMinimalData, str)
	} else if dataLen == 1 && data[0] >= 1 && data[0] <= 16 {
		if opcode != Op1+data[0]-1 {
			// Should have used OP_1 .. OP_16
			str := fmt.Sprintf("data push of the value %d encoded "+
				"with opcode %s instead of OP_%d", data[0],
				pop.opcode.name, data[0])
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen == 1 && data[0] == 0x81 {
		if opcode != Op1Negate {
			str := fmt.Sprintf("data push of the value -1 encoded "+
				"with opcode %s instead of OP_1NEGATE",
				pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 75 {
		if int(opcode) != dataLen {
			// Should have used a direct push
			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OP_DATA_%d", dataLen,
				pop.opcode.name, dataLen)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 255 {
		if opcode != OpPushData1 {
			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OP_PUSHDATA1",
				dataLen, pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 65535 {
		if opcode != OpPushData2 {
			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OP_PUSHDATA2",
				dataLen, pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	}
	return nil
}

// print returns a human-readable string representation of the opcode for use
// in script disassembly.
func (pop *parsedOpcode) print(oneline bool) string {
	// The reference implementation one-line disassembly replaces opcodes
	// which represent values (e.g. OP_0 through OP_16 and OP_1NEGATE)
	// with the raw value. However, when not doing a one-line dissassembly,
	// we prefer to show the actual opcode names. Thus, only replace the
	// opcodes in question when the oneline flag is set.
	opcodeName := pop.opcode.name
	if oneline {
		if replName, ok := opcodeOnelineRepls[opcodeName]; ok {
			opcodeName = replName
		}

		// Nothing more to do for non-data push opcodes.
		if pop.opcode.length == 1 {
			return opcodeName
		}

		return fmt.Sprintf("%x", pop.data)
	}

	// Nothing more to do for non-data push opcodes.
	if pop.opcode.length == 1 {
		return opcodeName
	}

	// Add length for the OP_PUSHDATA# opcodes.
	retString := opcodeName
	switch pop.opcode.length {
	case -1:
		retString += fmt.Sprintf(" 0x%02x", len(pop.data))
	case -2:
		retString += fmt.Sprintf(" 0x%04x", len(pop.data))
	case -4:
		retString += fmt.Sprintf(" 0x%08x", len(pop.data))
	}

	return fmt.Sprintf("%s 0x%02x", retString, pop.data)
}

// bytes returns any data associated with the opcode encoded as it would be in
// a script. This is used for unparsing scripts from parsed opcodes.
func (pop *parsedOpcode) bytes() ([]byte, error) {
	var retbytes []byte
	if pop.opcode.length > 0 {
		retbytes = make([]byte, 1, pop.opcode.length)
	} else {
		retbytes = make([]byte, 1, 1+len(pop.data)-
			pop.opcode.length)
	}

	retbytes[0] = pop.opcode.value
	if pop.opcode.length == 1 {
		if len(pop.data) != 0 {
			str := fmt.Sprintf("internal consistency error - "+
				"parsed opcode %s has data length %d when %d "+
				"was expected", pop.opcode.name, len(pop.data),
				0)
			return nil, scriptError(ErrInternal, str)
		}
		return retbytes, nil
	}
	nbytes := pop.opcode.length
	if pop.opcode.length < 0 {
		l := len(pop.data)
		// tempting just to hardcode to avoid the complexity here.
		switch pop.opcode.length {
		case -1:
			retbytes = append(retbytes, byte(l))
			nbytes = int(retbytes[1]) + len(retbytes)
		case -2:
			retbytes = append(retbytes, byte(l&0xff),
				byte(l>>8&0xff))
			nbytes = int(binary.LittleEndian.Uint16(retbytes[1:])) +
				len(retbytes)
		case -4:
			retbytes = append(retbytes, byte(l&0xff),
				byte((l>>8)&0xff), byte((l>>16)&0xff),
				byte((l>>24)&0xff))
			nbytes = int(binary.LittleEndian.Uint32(retbytes[1:])) +
				len(retbytes)
		}
	}

	retbytes = append(retbytes, pop.data...)

	if len(retbytes) != nbytes {
		str := fmt.Sprintf("internal consistency error - "+
			"parsed opcode %s has data length %d when %d was "+
			"expected", pop.opcode.name, len(retbytes), nbytes)
		return nil, scriptError(ErrInternal, str)
	}

	return retbytes, nil
}

// *******************************************
// Opcode implementation functions start here.
// *******************************************

// opcodeDisabled is a common handler for disabled opcodes. It returns an
// appropriate error indicating the opcode is disabled. While it would
// ordinarily make more sense to detect if the script contains any disabled
// opcodes before executing in an initial parse step, the consensus rules
// dictate the script doesn't fail until the program counter passes over a
// disabled opcode (even when they appear in a branch that is not executed).
func opcodeDisabled(op *parsedOpcode, vm *Engine) error {
	str := fmt.Sprintf("attempt to execute disabled opcode %s",
		op.opcode.name)
	return scriptError(ErrDisabledOpcode, str)
}

// opcodeReserved is a common handler for all reserved opcodes. It returns an
// appropriate error indicating the opcode is reserved.
func opcodeReserved(op *parsedOpcode, vm *Engine) error {
	str := fmt.Sprintf("attempt to execute reserved opcode %s",
		op.opcode.name)
	return scriptError(ErrReservedOpcode, str)
}

// opcodeInvalid is a common handler for all invalid opcodes. It returns an
// appropriate error indicating the opcode is invalid.
func opcodeInvalid(op *parsedOpcode, vm *Engine) error {
	str := fmt.Sprintf("attempt to execute invalid opcode %s",
		op.opcode.name)
	return scriptError(ErrReservedOpcode, str)
}

// opcodeFalse pushes an empty array to the data stack to represent false. Note
// that 0, when encoded as a number according to the numeric encoding consensus
// rules, is an empty array.
func opcodeFalse(op *parsedOpcode, vm *Engine) error {
	vm.dstack.PushByteArray(nil)
	return nil
}

// opcodePushData is a common handler for the vast majority of opcodes that push
// raw data (bytes) to the data stack.
func opcodePushData(op *parsedOpcode, vm *Engine) error {
	vm.dstack.PushByteArray(op.data)
	return nil
}

// opcode1Negate pushes -1, encoded as a number, to the data stack.
func opcode1Negate(op *parsedOpcode, vm *Engine) error {
	vm.dstack.PushInt(scriptNum(-1))
	return nil
}

// opcodeN is a common handler for the small integer data push opcodes. It
// pushes the numeric value the opcode represents (which will be from 1 to 16)
// onto the data stack.
func opcodeN(op *parsedOpcode, vm *Engine) error {
	// The opcodes are all defined consecutively, so the numeric value is
	// the difference.
	vm.dstack.PushInt(scriptNum((op.opcode.value - (Op1 - 1))))
	return nil
}

// opcodeNop is a common handler for the NOP family of opcodes. As the name
// implies it generally does nothing, however, it will return an error when
// the flag to discourage use of NOPs is set for select opcodes.
func opcodeNop(op *parsedOpcode, vm *Engine) error {
	return nil
}

// popIfBool enforces the "minimal if" policy. In order to
// eliminate an additional source of nuisance malleability, we
// require the following: for OP_IF and OP_NOTIF, the top stack item MUST
// either be an empty byte slice, or [0x01]. Otherwise, the item at the top of
// the stack will be popped and interpreted as a boolean.
func popIfBool(vm *Engine) (bool, error) {
	so, err := vm.dstack.PopByteArray()
	if err != nil {
		return false, err
	}

	if len(so) == 1 && so[0] == 0x01 {
		return true, nil
	}

	if len(so) == 0 {
		return false, nil
	}

	str := fmt.Sprintf("with OP_IF or OP_NOTIF top stack item MUST "+
		"be an empty byte array or 0x01, and is instead: %x",
		so)
	return false, scriptError(ErrMinimalIf, str)
}

// opcodeIf treats the top item on the data stack as a boolean and removes it.
//
// An appropriate entry is added to the conditional stack depending on whether
// the boolean is true and whether this if is on an executing branch in order
// to allow proper execution of further opcodes depending on the conditional
// logic. When the boolean is true, the first branch will be executed (unless
// this opcode is nested in a non-executed branch).
//
// <expression> if [statements] [else [statements]] endif
//
// Note that, unlike for all non-conditional opcodes, this is executed even when
// it is on a non-executing branch so proper nesting is maintained.
//
// Data stack transformation: [... bool] -> [...]
// Conditional stack transformation: [...] -> [... OpCondValue]
func opcodeIf(op *parsedOpcode, vm *Engine) error {
	condVal := OpCondFalse
	if vm.isBranchExecuting() {
		ok, err := popIfBool(vm)

		if err != nil {
			return err
		}

		if ok {
			condVal = OpCondTrue
		}
	} else {
		condVal = OpCondSkip
	}
	vm.condStack = append(vm.condStack, condVal)
	return nil
}

// opcodeNotIf treats the top item on the data stack as a boolean and removes
// it.
//
// An appropriate entry is added to the conditional stack depending on whether
// the boolean is true and whether this if is on an executing branch in order
// to allow proper execution of further opcodes depending on the conditional
// logic. When the boolean is false, the first branch will be executed (unless
// this opcode is nested in a non-executed branch).
//
// <expression> notif [statements] [else [statements]] endif
//
// Note that, unlike for all non-conditional opcodes, this is executed even when
// it is on a non-executing branch so proper nesting is maintained.
//
// Data stack transformation: [... bool] -> [...]
// Conditional stack transformation: [...] -> [... OpCondValue]
func opcodeNotIf(op *parsedOpcode, vm *Engine) error {
	condVal := OpCondFalse
	if vm.isBranchExecuting() {
		ok, err := popIfBool(vm)
		if err != nil {
			return err
		}

		if !ok {
			condVal = OpCondTrue
		}
	} else {
		condVal = OpCondSkip
	}
	vm.condStack = append(vm.condStack, condVal)
	return nil
}

// opcodeElse inverts conditional execution for other half of if/else/endif.
//
// An error is returned if there has not already been a matching OP_IF.
//
// Conditional stack transformation: [... OpCondValue] -> [... !OpCondValue]
func opcodeElse(op *parsedOpcode, vm *Engine) error {
	if len(vm.condStack) == 0 {
		str := fmt.Sprintf("encountered opcode %s with no matching "+
			"opcode to begin conditional execution", op.opcode.name)
		return scriptError(ErrUnbalancedConditional, str)
	}

	conditionalIdx := len(vm.condStack) - 1
	switch vm.condStack[conditionalIdx] {
	case OpCondTrue:
		vm.condStack[conditionalIdx] = OpCondFalse
	case OpCondFalse:
		vm.condStack[conditionalIdx] = OpCondTrue
	case OpCondSkip:
		// Value doesn't change in skip since it indicates this opcode
		// is nested in a non-executed branch.
	}
	return nil
}

// opcodeEndif terminates a conditional block, removing the value from the
// conditional execution stack.
//
// An error is returned if there has not already been a matching OP_IF.
//
// Conditional stack transformation: [... OpCondValue] -> [...]
func opcodeEndif(op *parsedOpcode, vm *Engine) error {
	if len(vm.condStack) == 0 {
		str := fmt.Sprintf("encountered opcode %s with no matching "+
			"opcode to begin conditional execution", op.opcode.name)
		return scriptError(ErrUnbalancedConditional, str)
	}

	vm.condStack = vm.condStack[:len(vm.condStack)-1]
	return nil
}

// abstractVerify examines the top item on the data stack as a boolean value and
// verifies it evaluates to true. An error is returned either when there is no
// item on the stack or when that item evaluates to false. In the latter case
// where the verification fails specifically due to the top item evaluating
// to false, the returned error will use the passed error code.
func abstractVerify(op *parsedOpcode, vm *Engine, c ErrorCode) error {
	verified, err := vm.dstack.PopBool()
	if err != nil {
		return err
	}

	if !verified {
		str := fmt.Sprintf("%s failed", op.opcode.name)
		return scriptError(c, str)
	}
	return nil
}

// opcodeVerify examines the top item on the data stack as a boolean value and
// verifies it evaluates to true. An error is returned if it does not.
func opcodeVerify(op *parsedOpcode, vm *Engine) error {
	return abstractVerify(op, vm, ErrVerify)
}

// opcodeReturn returns an appropriate error since it is always an error to
// return early from a script.
func opcodeReturn(op *parsedOpcode, vm *Engine) error {
	return scriptError(ErrEarlyReturn, "script returned early")
}

func verifyLockTimeWithThreshold(txLockTime, threshold, lockTime uint64) error {
	err := verifyLockTimeThreshold(txLockTime, threshold, lockTime)
	if err != nil {
		return err
	}
	return verifyLockTime(txLockTime, lockTime)
}

// checkLockTimeRequirement is a helper function used to validate locktimes.
func verifyLockTime(txLockTime, lockTime uint64) error {
	if lockTime > txLockTime {
		str := fmt.Sprintf("locktime requirement not satisfied -- "+
			"locktime is greater than the transaction locktime: "+
			"%d > %d", lockTime, txLockTime)
		return scriptError(ErrUnsatisfiedLockTime, str)
	}

	return nil
}

// verifyLockTimeThreshold is a helper function used to verify the lockTimes in both the script and transaction have the same type.
func verifyLockTimeThreshold(txLockTime, threshold, lockTime uint64) error {
	if !((txLockTime < threshold && lockTime < threshold) ||
		(txLockTime >= threshold && lockTime >= threshold)) {
		str := fmt.Sprintf("mismatched locktime types -- tx locktime "+
			"%d, stack locktime %d", txLockTime, lockTime)
		return scriptError(ErrUnsatisfiedLockTime, str)
	}

	return nil
}

// opcodeCheckLockTimeVerify compares the top item on the data stack to the
// LockTime field of the transaction containing the script signature
// validating if the transaction outputs are spendable yet.
func opcodeCheckLockTimeVerify(op *parsedOpcode, vm *Engine) error {
	lockTimeBytes, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	// Make sure lockTimeBytes is exactly 8 bytes.
	// If more - return ErrNumberTooBig
	// If less - pad with 0's
	if len(lockTimeBytes) > 8 {
		str := fmt.Sprintf("lockTime value represented as %x is longer then 8 bytes", lockTimeBytes)
		return scriptError(ErrNumberTooBig, str)
	}
	if len(lockTimeBytes) < 8 {
		paddedLockTimeBytes := make([]byte, 8)
		copy(paddedLockTimeBytes, lockTimeBytes)
		lockTimeBytes = paddedLockTimeBytes
	}
	stackLockTime := binary.LittleEndian.Uint64(lockTimeBytes)
	// The lock time field of a transaction is either a DAA score at
	// which the transaction is finalized or a timestamp depending on if the
	// value is before the constants.LockTimeThreshold. When it is under the
	// threshold it is a DAA score.
	err = verifyLockTimeWithThreshold(vm.tx.LockTime, constants.LockTimeThreshold, stackLockTime)
	if err != nil {
		return err
	}

	// The lock time feature can also be disabled, thereby bypassing
	// OP_CHECKLOCKTIMEVERIFY, if every transaction input has been finalized by
	// setting its sequence to the maximum value (constants.MaxTxInSequenceNum). This
	// condition would result in the transaction being allowed into the blockDAG
	// making the opcode ineffective.
	//
	// This condition is prevented by enforcing that the input being used by
	// the opcode is unlocked (its sequence number is less than the max
	// value). This is sufficient to prove correctness without having to
	// check every input.
	//
	// NOTE: This implies that even if the transaction is not finalized due to
	// another input being unlocked, the opcode execution will still fail when the
	// input being used by the opcode is locked.
	if vm.tx.Inputs[vm.txIdx].Sequence == constants.MaxTxInSequenceNum {
		return scriptError(ErrUnsatisfiedLockTime,
			"transaction input is finalized")
	}

	return nil
}

// opcodeCheckSequenceVerify compares the top item on the data stack to the
// LockTime field of the transaction containing the script signature
// validating if the transaction outputs are spendable yet.
func opcodeCheckSequenceVerify(op *parsedOpcode, vm *Engine) error {
	sequenceBytes, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	// Make sure sequenceBytes is exactly 8 bytes.
	// If more - return ErrNumberTooBig
	// If less - pad with 0's
	if len(sequenceBytes) > 8 {
		str := fmt.Sprintf("sequence value represented as %x is longer then 8 bytes", sequenceBytes)
		return scriptError(ErrNumberTooBig, str)
	}
	if len(sequenceBytes) < 8 {
		paddedSequenceBytes := make([]byte, 8)
		copy(paddedSequenceBytes, sequenceBytes)
		sequenceBytes = paddedSequenceBytes
	}

	// Don't use makeScriptNum here, since sequence is not an actual number, minimal encoding rules don't apply to it,
	// and is more convenient to be represented as an unsigned int.
	stackSequence := binary.LittleEndian.Uint64(sequenceBytes)

	// To provide for future soft-fork extensibility, if the
	// operand has the disabled lock-time flag set,
	// CHECKSEQUENCEVERIFY behaves as a NOP.
	if stackSequence&constants.SequenceLockTimeDisabled != 0 {
		return nil
	}

	// Sequence numbers with their most significant bit set are not
	// consensus constrained. Testing that the transaction's sequence
	// number does not have this bit set prevents using this property
	// to get around a CHECKSEQUENCEVERIFY check.
	txSequence := vm.tx.Inputs[vm.txIdx].Sequence
	if txSequence&constants.SequenceLockTimeDisabled != 0 {
		str := fmt.Sprintf("transaction sequence has sequence "+
			"locktime disabled bit set: 0x%x", txSequence)
		return scriptError(ErrUnsatisfiedLockTime, str)
	}

	// Mask off non-consensus bits before doing comparisons.
	maskedTxSequence := txSequence & constants.SequenceLockTimeMask
	maskedStackSequence := stackSequence & constants.SequenceLockTimeMask
	return verifyLockTime(maskedTxSequence, maskedStackSequence)
}

// opcodeToAltStack removes the top item from the main data stack and pushes it
// onto the alternate data stack.
//
// Main data stack transformation: [... x1 x2 x3] -> [... x1 x2]
// Alt data stack transformation:  [... y1 y2 y3] -> [... y1 y2 y3 x3]
func opcodeToAltStack(op *parsedOpcode, vm *Engine) error {
	so, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	vm.astack.PushByteArray(so)

	return nil
}

// opcodeFromAltStack removes the top item from the alternate data stack and
// pushes it onto the main data stack.
//
// Main data stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 y3]
// Alt data stack transformation:  [... y1 y2 y3] -> [... y1 y2]
func opcodeFromAltStack(op *parsedOpcode, vm *Engine) error {
	so, err := vm.astack.PopByteArray()
	if err != nil {
		return err
	}
	vm.dstack.PushByteArray(so)

	return nil
}

// opcode2Drop removes the top 2 items from the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1]
func opcode2Drop(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.DropN(2)
}

// opcode2Dup duplicates the top 2 items on the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x2 x3]
func opcode2Dup(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.DupN(2)
}

// opcode3Dup duplicates the top 3 items on the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x1 x2 x3]
func opcode3Dup(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.DupN(3)
}

// opcode2Over duplicates the 2 items before the top 2 items on the data stack.
//
// Stack transformation: [... x1 x2 x3 x4] -> [... x1 x2 x3 x4 x1 x2]
func opcode2Over(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.OverN(2)
}

// opcode2Rot rotates the top 6 items on the data stack to the left twice.
//
// Stack transformation: [... x1 x2 x3 x4 x5 x6] -> [... x3 x4 x5 x6 x1 x2]
func opcode2Rot(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.RotN(2)
}

// opcode2Swap swaps the top 2 items on the data stack with the 2 that come
// before them.
//
// Stack transformation: [... x1 x2 x3 x4] -> [... x3 x4 x1 x2]
func opcode2Swap(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.SwapN(2)
}

// opcodeIfDup duplicates the top item of the stack if it is not zero.
//
// Stack transformation (x1==0): [... x1] -> [... x1]
// Stack transformation (x1!=0): [... x1] -> [... x1 x1]
func opcodeIfDup(op *parsedOpcode, vm *Engine) error {
	so, err := vm.dstack.PeekByteArray(0)
	if err != nil {
		return err
	}

	// Push copy of data iff it isn't zero
	if asBool(so) {
		vm.dstack.PushByteArray(so)
	}

	return nil
}

// opcodeDepth pushes the depth of the data stack prior to executing this
// opcode, encoded as a number, onto the data stack.
//
// Stack transformation: [...] -> [... <num of items on the stack>]
// Example with 2 items: [x1 x2] -> [x1 x2 2]
// Example with 3 items: [x1 x2 x3] -> [x1 x2 x3 3]
func opcodeDepth(op *parsedOpcode, vm *Engine) error {
	vm.dstack.PushInt(scriptNum(vm.dstack.Depth()))
	return nil
}

// opcodeDrop removes the top item from the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x2]
func opcodeDrop(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.DropN(1)
}

// opcodeDup duplicates the top item on the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x3]
func opcodeDup(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.DupN(1)
}

// opcodeNip removes the item before the top item on the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x3]
func opcodeNip(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.NipN(1)
}

// opcodeOver duplicates the item before the top item on the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x2]
func opcodeOver(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.OverN(1)
}

// opcodePick treats the top item on the data stack as an integer and duplicates
// the item on the stack that number of items back to the top.
//
// Stack transformation: [xn ... x2 x1 x0 n] -> [xn ... x2 x1 x0 xn]
// Example with n=1: [x2 x1 x0 1] -> [x2 x1 x0 x1]
// Example with n=2: [x2 x1 x0 2] -> [x2 x1 x0 x2]
func opcodePick(op *parsedOpcode, vm *Engine) error {
	val, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	return vm.dstack.PickN(val.Int32())
}

// opcodeRoll treats the top item on the data stack as an integer and moves
// the item on the stack that number of items back to the top.
//
// Stack transformation: [xn ... x2 x1 x0 n] -> [... x2 x1 x0 xn]
// Example with n=1: [x2 x1 x0 1] -> [x2 x0 x1]
// Example with n=2: [x2 x1 x0 2] -> [x1 x0 x2]
func opcodeRoll(op *parsedOpcode, vm *Engine) error {
	val, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	return vm.dstack.RollN(val.Int32())
}

// opcodeRot rotates the top 3 items on the data stack to the left.
//
// Stack transformation: [... x1 x2 x3] -> [... x2 x3 x1]
func opcodeRot(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.RotN(1)
}

// opcodeSwap swaps the top two items on the stack.
//
// Stack transformation: [... x1 x2] -> [... x2 x1]
func opcodeSwap(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.SwapN(1)
}

// opcodeTuck inserts a duplicate of the top item of the data stack before the
// second-to-top item.
//
// Stack transformation: [... x1 x2] -> [... x2 x1 x2]
func opcodeTuck(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.Tuck()
}

// opcodeSize pushes the size of the top item of the data stack onto the data
// stack.
//
// Stack transformation: [... x1] -> [... x1 len(x1)]
func opcodeSize(op *parsedOpcode, vm *Engine) error {
	so, err := vm.dstack.PeekByteArray(0)
	if err != nil {
		return err
	}

	vm.dstack.PushInt(scriptNum(len(so)))
	return nil
}

// opcodeEqual removes the top 2 items of the data stack, compares them as raw
// bytes, and pushes the result, encoded as a boolean, back to the stack.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeEqual(op *parsedOpcode, vm *Engine) error {
	a, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	b, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	vm.dstack.PushBool(bytes.Equal(a, b))
	return nil
}

// opcodeEqualVerify is a combination of opcodeEqual and opcodeVerify.
// Specifically, it removes the top 2 items of the data stack, compares them,
// and pushes the result, encoded as a boolean, back to the stack. Then, it
// examines the top item on the data stack as a boolean value and verifies it
// evaluates to true. An error is returned if it does not.
//
// Stack transformation: [... x1 x2] -> [... bool] -> [...]
func opcodeEqualVerify(op *parsedOpcode, vm *Engine) error {
	err := opcodeEqual(op, vm)
	if err == nil {
		err = abstractVerify(op, vm, ErrEqualVerify)
	}
	return err
}

// opcode1Add treats the top item on the data stack as an integer and replaces
// it with its incremented value (plus 1).
//
// Stack transformation: [... x1 x2] -> [... x1 x2+1]
func opcode1Add(op *parsedOpcode, vm *Engine) error {
	m, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	vm.dstack.PushInt(m + 1)
	return nil
}

// opcode1Sub treats the top item on the data stack as an integer and replaces
// it with its decremented value (minus 1).
//
// Stack transformation: [... x1 x2] -> [... x1 x2-1]
func opcode1Sub(op *parsedOpcode, vm *Engine) error {
	m, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}
	vm.dstack.PushInt(m - 1)

	return nil
}

// opcodeNegate treats the top item on the data stack as an integer and replaces
// it with its negation.
//
// Stack transformation: [... x1 x2] -> [... x1 -x2]
func opcodeNegate(op *parsedOpcode, vm *Engine) error {
	m, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	vm.dstack.PushInt(-m)
	return nil
}

// opcodeAbs treats the top item on the data stack as an integer and replaces it
// it with its absolute value.
//
// Stack transformation: [... x1 x2] -> [... x1 abs(x2)]
func opcodeAbs(op *parsedOpcode, vm *Engine) error {
	m, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if m < 0 {
		m = -m
	}
	vm.dstack.PushInt(m)
	return nil
}

// opcodeNot treats the top item on the data stack as an integer and replaces
// it with its "inverted" value (0 becomes 1, non-zero becomes 0).
//
// NOTE: While it would probably make more sense to treat the top item as a
// boolean, and push the opposite, which is really what the intention of this
// opcode is, it is extremely important that is not done because integers are
// interpreted differently than booleans and the consensus rules for this opcode
// dictate the item is interpreted as an integer.
//
// Stack transformation (x2==0): [... x1 0] -> [... x1 1]
// Stack transformation (x2!=0): [... x1 1] -> [... x1 0]
// Stack transformation (x2!=0): [... x1 17] -> [... x1 0]
func opcodeNot(op *parsedOpcode, vm *Engine) error {
	m, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if m == 0 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcode0NotEqual treats the top item on the data stack as an integer and
// replaces it with either a 0 if it is zero, or a 1 if it is not zero.
//
// Stack transformation (x2==0): [... x1 0] -> [... x1 0]
// Stack transformation (x2!=0): [... x1 1] -> [... x1 1]
// Stack transformation (x2!=0): [... x1 17] -> [... x1 1]
func opcode0NotEqual(op *parsedOpcode, vm *Engine) error {
	m, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if m != 0 {
		m = 1
	}
	vm.dstack.PushInt(m)
	return nil
}

// opcodeAdd treats the top two items on the data stack as integers and replaces
// them with their sum.
//
// Stack transformation: [... x1 x2] -> [... x1+x2]
func opcodeAdd(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	vm.dstack.PushInt(v0 + v1)
	return nil
}

// opcodeSub treats the top two items on the data stack as integers and replaces
// them with the result of subtracting the top entry from the second-to-top
// entry.
//
// Stack transformation: [... x1 x2] -> [... x1-x2]
func opcodeSub(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	vm.dstack.PushInt(v1 - v0)
	return nil
}

// opcodeBoolAnd treats the top two items on the data stack as integers. When
// both of them are not zero, they are replaced with a 1, otherwise a 0.
//
// Stack transformation (x1==0, x2==0): [... 0 0] -> [... 0]
// Stack transformation (x1!=0, x2==0): [... 5 0] -> [... 0]
// Stack transformation (x1==0, x2!=0): [... 0 7] -> [... 0]
// Stack transformation (x1!=0, x2!=0): [... 4 8] -> [... 1]
func opcodeBoolAnd(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v0 != 0 && v1 != 0 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}

	return nil
}

// opcodeBoolOr treats the top two items on the data stack as integers. When
// either of them are not zero, they are replaced with a 1, otherwise a 0.
//
// Stack transformation (x1==0, x2==0): [... 0 0] -> [... 0]
// Stack transformation (x1!=0, x2==0): [... 5 0] -> [... 1]
// Stack transformation (x1==0, x2!=0): [... 0 7] -> [... 1]
// Stack transformation (x1!=0, x2!=0): [... 4 8] -> [... 1]
func opcodeBoolOr(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v0 != 0 || v1 != 0 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}

	return nil
}

// opcodeNumEqual treats the top two items on the data stack as integers. When
// they are equal, they are replaced with a 1, otherwise a 0.
//
// Stack transformation (x1==x2): [... 5 5] -> [... 1]
// Stack transformation (x1!=x2): [... 5 7] -> [... 0]
func opcodeNumEqual(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v0 == v1 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}

	return nil
}

// opcodeNumEqualVerify is a combination of opcodeNumEqual and opcodeVerify.
//
// Specifically, treats the top two items on the data stack as integers. When
// they are equal, they are replaced with a 1, otherwise a 0. Then, it examines
// the top item on the data stack as a boolean value and verifies it evaluates
// to true. An error is returned if it does not.
//
// Stack transformation: [... x1 x2] -> [... bool] -> [...]
func opcodeNumEqualVerify(op *parsedOpcode, vm *Engine) error {
	err := opcodeNumEqual(op, vm)
	if err == nil {
		err = abstractVerify(op, vm, ErrNumEqualVerify)
	}
	return err
}

// opcodeNumNotEqual treats the top two items on the data stack as integers.
// When they are NOT equal, they are replaced with a 1, otherwise a 0.
//
// Stack transformation (x1==x2): [... 5 5] -> [... 0]
// Stack transformation (x1!=x2): [... 5 7] -> [... 1]
func opcodeNumNotEqual(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v0 != v1 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}

	return nil
}

// opcodeLessThan treats the top two items on the data stack as integers. When
// the second-to-top item is less than the top item, they are replaced with a 1,
// otherwise a 0.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeLessThan(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v1 < v0 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}

	return nil
}

// opcodeGreaterThan treats the top two items on the data stack as integers.
// When the second-to-top item is greater than the top item, they are replaced
// with a 1, otherwise a 0.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeGreaterThan(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v1 > v0 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeLessThanOrEqual treats the top two items on the data stack as integers.
// When the second-to-top item is less than or equal to the top item, they are
// replaced with a 1, otherwise a 0.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeLessThanOrEqual(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v1 <= v0 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeGreaterThanOrEqual treats the top two items on the data stack as
// integers. When the second-to-top item is greater than or equal to the top
// item, they are replaced with a 1, otherwise a 0.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeGreaterThanOrEqual(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v1 >= v0 {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}

	return nil
}

// opcodeMin treats the top two items on the data stack as integers and replaces
// them with the minimum of the two.
//
// Stack transformation: [... x1 x2] -> [... min(x1, x2)]
func opcodeMin(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v1 < v0 {
		vm.dstack.PushInt(v1)
	} else {
		vm.dstack.PushInt(v0)
	}

	return nil
}

// opcodeMax treats the top two items on the data stack as integers and replaces
// them with the maximum of the two.
//
// Stack transformation: [... x1 x2] -> [... max(x1, x2)]
func opcodeMax(op *parsedOpcode, vm *Engine) error {
	v0, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	v1, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if v1 > v0 {
		vm.dstack.PushInt(v1)
	} else {
		vm.dstack.PushInt(v0)
	}

	return nil
}

// opcodeWithin treats the top 3 items on the data stack as integers. When the
// value to test is within the specified range (left inclusive), they are
// replaced with a 1, otherwise a 0.
//
// The top item is the max value, the second-top-item is the minimum value, and
// the third-to-top item is the value to test.
//
// Stack transformation: [... x1 min max] -> [... bool]
func opcodeWithin(op *parsedOpcode, vm *Engine) error {
	maxVal, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	minVal, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	x, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}

	if x >= minVal && x < maxVal {
		vm.dstack.PushInt(scriptNum(1))
	} else {
		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// calcHash calculates the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

// opcodeSha256 treats the top item of the data stack as raw bytes and replaces
// it with sha256(data).
//
// Stack transformation: [... x1] -> [... sha256(x1)]
func opcodeSha256(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(buf)
	vm.dstack.PushByteArray(hash[:])
	return nil
}

// opcodeBlake2b treats the top item of the data stack as raw bytes and replaces
// it with blake2b(data).
//
// Stack transformation: [... x1] -> [... blake2b(x1)]
func opcodeBlake2b(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	hash := blake2b.Sum256(buf)
	vm.dstack.PushByteArray(hash[:])
	return nil
}

// OpcodeByName is a map that can be used to lookup an opcode by its
// human-readable name (OP_CHECKMULTISIG, OP_CHECKSIG, etc).
var OpcodeByName = make(map[string]byte)

func init() {
	// Initialize the opcode name to value map using the contents of the
	// opcode array. Also add entries for "OP_FALSE" and "OP_TRUE"
	// since they are aliases for "OP_0" and "OP_1" respectively.
	for _, op := range opcodeArray {
		OpcodeByName[op.name] = op.value
	}
	OpcodeByName["OP_FALSE"] = OpFalse
	OpcodeByName["OP_TRUE"] = OpTrue
}
