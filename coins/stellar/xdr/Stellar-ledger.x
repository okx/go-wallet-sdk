// Copyright 2015 Stellar Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

%#include "xdr/Stellar-SCP.h"
%#include "xdr/Stellar-transaction.h"

namespace stellar
{

typedef opaque UpgradeType<128>;

enum StellarValueType
{
    STELLAR_VALUE_BASIC = 0,
    STELLAR_VALUE_SIGNED = 1
};

struct LedgerCloseValueSignature
{
    NodeID nodeID;       // which node introduced the value
    Signature signature; // nodeID's signature
};

/* StellarValue is the value used by SCP to reach consensus on a given ledger
 */
struct StellarValue
{
    Hash txSetHash;      // transaction set to apply to previous ledger
    TimePoint closeTime; // network close time

    // upgrades to apply to the previous ledger (usually empty)
    // this is a vector of encoded 'LedgerUpgrade' so that nodes can drop
    // unknown steps during consensus if needed.
    // see notes below on 'LedgerUpgrade' for more detail
    // max size is dictated by number of upgrade types (+ room for future)
    UpgradeType upgrades<6>;

    // reserved for future use
    union switch (StellarValueType v)
    {
    case STELLAR_VALUE_BASIC:
        void;
    case STELLAR_VALUE_SIGNED:
        LedgerCloseValueSignature lcValueSignature;
    }
    ext;
};

const MASK_LEDGER_HEADER_FLAGS = 0x7;

enum LedgerHeaderFlags
{
    DISABLE_LIQUIDITY_POOL_TRADING_FLAG = 0x1,
    DISABLE_LIQUIDITY_POOL_DEPOSIT_FLAG = 0x2,
    DISABLE_LIQUIDITY_POOL_WITHDRAWAL_FLAG = 0x4
};

struct LedgerHeaderExtensionV1
{
    uint32 flags; // LedgerHeaderFlags

    union switch (int v)
    {
    case 0:
        void;
    }
    ext;
};

/* The LedgerHeader is the highest level structure representing the
 * state of a ledger, cryptographically linked to previous ledgers.
 */
struct LedgerHeader
{
    uint32 ledgerVersion;    // the protocol version of the ledger
    Hash previousLedgerHash; // hash of the previous ledger header
    StellarValue scpValue;   // what consensus agreed to
    Hash txSetResultHash;    // the TransactionResultSet that led to this ledger
    Hash bucketListHash;     // hash of the ledger state

    uint32 ledgerSeq; // sequence number of this ledger

    int64 totalCoins; // total number of stroops in existence.
                      // 10,000,000 stroops in 1 XLM

    int64 feePool;       // fees burned since last inflation run
    uint32 inflationSeq; // inflation sequence number

    uint64 idPool; // last used global ID, used for generating objects

    uint32 baseFee;     // base fee per operation in stroops
    uint32 baseReserve; // account base reserve in stroops

    uint32 maxTxSetSize; // maximum size a transaction set can be

    Hash skipList[4]; // hashes of ledgers in the past. allows you to jump back
                      // in time without walking the chain back ledger by ledger
                      // each slot contains the oldest ledger that is mod of
                      // either 50  5000  50000 or 500000 depending on index
                      // skipList[0] mod(50), skipList[1] mod(5000), etc

    // reserved for future use
    union switch (int v)
    {
    case 0:
        void;
    case 1:
        LedgerHeaderExtensionV1 v1;
    }
    ext;
};

/* Ledger upgrades
note that the `upgrades` field from StellarValue is normalized such that
it only contains one entry per LedgerUpgradeType, and entries are sorted
in ascending order
*/
enum LedgerUpgradeType
{
    LEDGER_UPGRADE_VERSION = 1,
    LEDGER_UPGRADE_BASE_FEE = 2,
    LEDGER_UPGRADE_MAX_TX_SET_SIZE = 3,
    LEDGER_UPGRADE_BASE_RESERVE = 4,
    LEDGER_UPGRADE_FLAGS = 5,
    LEDGER_UPGRADE_CONFIG = 6,
    LEDGER_UPGRADE_MAX_SOROBAN_TX_SET_SIZE = 7
};

struct ConfigUpgradeSetKey {
    Hash contractID;
    Hash contentHash;
};

union LedgerUpgrade switch (LedgerUpgradeType type)
{
case LEDGER_UPGRADE_VERSION:
    uint32 newLedgerVersion; // update ledgerVersion
case LEDGER_UPGRADE_BASE_FEE:
    uint32 newBaseFee; // update baseFee
case LEDGER_UPGRADE_MAX_TX_SET_SIZE:
    uint32 newMaxTxSetSize; // update maxTxSetSize
case LEDGER_UPGRADE_BASE_RESERVE:
    uint32 newBaseReserve; // update baseReserve
case LEDGER_UPGRADE_FLAGS:
    uint32 newFlags; // update flags
case LEDGER_UPGRADE_CONFIG:
    // Update arbitrary `ConfigSetting` entries identified by the key.
    ConfigUpgradeSetKey newConfig;
case LEDGER_UPGRADE_MAX_SOROBAN_TX_SET_SIZE:
    // Update ConfigSettingContractExecutionLanesV0.ledgerMaxTxCount without
    // using `LEDGER_UPGRADE_CONFIG`.
    uint32 newMaxSorobanTxSetSize;
};

struct ConfigUpgradeSet {
    ConfigSettingEntry updatedEntry<>;
};

/* Entries used to define the bucket list */
enum BucketEntryType
{
    METAENTRY =
        -1, // At-and-after protocol 11: bucket metadata, should come first.
    LIVEENTRY = 0, // Before protocol 11: created-or-updated;
                   // At-and-after protocol 11: only updated.
    DEADENTRY = 1,
    INITENTRY = 2 // At-and-after protocol 11: only created.
};

struct BucketMetadata
{
    // Indicates the protocol version used to create / merge this bucket.
    uint32 ledgerVersion;

    // reserved for future use
    union switch (int v)
    {
    case 0:
        void;
    }
    ext;
};

union BucketEntry switch (BucketEntryType type)
{
case LIVEENTRY:
case INITENTRY:
    LedgerEntry liveEntry;

case DEADENTRY:
    LedgerKey deadEntry;
case METAENTRY:
    BucketMetadata metaEntry;
};

enum TxSetComponentType
{
  // txs with effective fee <= bid derived from a base fee (if any).
  // If base fee is not specified, no discount is applied.
  TXSET_COMP_TXS_MAYBE_DISCOUNTED_FEE = 0
};

union TxSetComponent switch (TxSetComponentType type)
{
case TXSET_COMP_TXS_MAYBE_DISCOUNTED_FEE:
  struct
  {
    int64* baseFee;
    TransactionEnvelope txs<>;
  } txsMaybeDiscountedFee;
};

union TransactionPhase switch (int v)
{
case 0:
    TxSetComponent v0Components<>;
};

// Transaction sets are the unit used by SCP to decide on transitions
// between ledgers
struct TransactionSet
{
    Hash previousLedgerHash;
    TransactionEnvelope txs<>;
};

struct TransactionSetV1
{
    Hash previousLedgerHash;
    TransactionPhase phases<>;
};

union GeneralizedTransactionSet switch (int v)
{
// We consider the legacy TransactionSet to be v0.
case 1:
    TransactionSetV1 v1TxSet;
};

struct TransactionResultPair
{
    Hash transactionHash;
    TransactionResult result; // result for the transaction
};

// TransactionResultSet is used to recover results between ledgers
struct TransactionResultSet
{
    TransactionResultPair results<>;
};

// Entries below are used in the historical subsystem

struct TransactionHistoryEntry
{
    uint32 ledgerSeq;
    TransactionSet txSet;

    // when v != 0, txSet must be empty
    union switch (int v)
    {
    case 0:
        void;
    case 1:
        GeneralizedTransactionSet generalizedTxSet;
    }
    ext;
};

struct TransactionHistoryResultEntry
{
    uint32 ledgerSeq;
    TransactionResultSet txResultSet;

    // reserved for future use
    union switch (int v)
    {
    case 0:
        void;
    }
    ext;
};

struct LedgerHeaderHistoryEntry
{
    Hash hash;
    LedgerHeader header;

    // reserved for future use
    union switch (int v)
    {
    case 0:
        void;
    }
    ext;
};

// historical SCP messages

struct LedgerSCPMessages
{
    uint32 ledgerSeq;
    SCPEnvelope messages<>;
};

// note: ledgerMessages may refer to any quorumSets encountered
// in the file so far, not just the one from this entry
struct SCPHistoryEntryV0
{
    SCPQuorumSet quorumSets<>; // additional quorum sets used by ledgerMessages
    LedgerSCPMessages ledgerMessages;
};

// SCP history file is an array of these
union SCPHistoryEntry switch (int v)
{
case 0:
    SCPHistoryEntryV0 v0;
};

// represents the meta in the transaction table history

// STATE is emitted every time a ledger entry is modified/deleted
// and the entry was not already modified in the current ledger

enum LedgerEntryChangeType
{
    LEDGER_ENTRY_CREATED = 0, // entry was added to the ledger
    LEDGER_ENTRY_UPDATED = 1, // entry was modified in the ledger
    LEDGER_ENTRY_REMOVED = 2, // entry was removed from the ledger
    LEDGER_ENTRY_STATE = 3    // value of the entry
};

union LedgerEntryChange switch (LedgerEntryChangeType type)
{
case LEDGER_ENTRY_CREATED:
    LedgerEntry created;
case LEDGER_ENTRY_UPDATED:
    LedgerEntry updated;
case LEDGER_ENTRY_REMOVED:
    LedgerKey removed;
case LEDGER_ENTRY_STATE:
    LedgerEntry state;
};

typedef LedgerEntryChange LedgerEntryChanges<>;

struct OperationMeta
{
    LedgerEntryChanges changes;
};

struct TransactionMetaV1
{
    LedgerEntryChanges txChanges; // tx level changes if any
    OperationMeta operations<>;   // meta for each operation
};

struct TransactionMetaV2
{
    LedgerEntryChanges txChangesBefore; // tx level changes before operations
                                        // are applied if any
    OperationMeta operations<>;         // meta for each operation
    LedgerEntryChanges txChangesAfter;  // tx level changes after operations are
                                        // applied if any
};

enum ContractEventType
{
    SYSTEM = 0,
    CONTRACT = 1,
    DIAGNOSTIC = 2
};

struct ContractEvent
{
    // We can use this to add more fields, or because it
    // is first, to change ContractEvent into a union.
    ExtensionPoint ext;

    Hash* contractID;
    ContractEventType type;

    union switch (int v)
    {
    case 0:
        struct
        {
            SCVal topics<>;
            SCVal data;
        } v0;
    }
    body;
};

struct DiagnosticEvent
{
    bool inSuccessfulContractCall;
    ContractEvent event;
};

typedef DiagnosticEvent DiagnosticEvents<>;

struct SorobanTransactionMetaExtV1
{
    ExtensionPoint ext;

    // The following are the components of the overall Soroban resource fee
    // charged for the transaction.
    // The following relation holds:
    // `resourceFeeCharged = totalNonRefundableResourceFeeCharged + totalRefundableResourceFeeCharged`
    // where `resourceFeeCharged` is the overall fee charged for the 
    // transaction. Also, `resourceFeeCharged` <= `sorobanData.resourceFee` 
    // i.e.we never charge more than the declared resource fee.
    // The inclusion fee for charged the Soroban transaction can be found using 
    // the following equation:
    // `result.feeCharged = resourceFeeCharged + inclusionFeeCharged`.

    // Total amount (in stroops) that has been charged for non-refundable
    // Soroban resources.
    // Non-refundable resources are charged based on the usage declared in
    // the transaction envelope (such as `instructions`, `readBytes` etc.) and 
    // is charged regardless of the success of the transaction.
    int64 totalNonRefundableResourceFeeCharged;
    // Total amount (in stroops) that has been charged for refundable
    // Soroban resource fees.
    // Currently this comprises the rent fee (`rentFeeCharged`) and the
    // fee for the events and return value.
    // Refundable resources are charged based on the actual resources usage.
    // Since currently refundable resources are only used for the successful
    // transactions, this will be `0` for failed transactions.
    int64 totalRefundableResourceFeeCharged;
    // Amount (in stroops) that has been charged for rent.
    // This is a part of `totalNonRefundableResourceFeeCharged`.
    int64 rentFeeCharged;
};

union SorobanTransactionMetaExt switch (int v)
{
case 0:
    void;
case 1:
    SorobanTransactionMetaExtV1 v1;
};

struct SorobanTransactionMeta 
{
    SorobanTransactionMetaExt ext;

    ContractEvent events<>;             // custom events populated by the
                                        // contracts themselves.
    SCVal returnValue;                  // return value of the host fn invocation

    // Diagnostics events that are not hashed.
    // This will contain all contract and diagnostic events. Even ones
    // that were emitted in a failed contract call.
    DiagnosticEvent diagnosticEvents<>;
};

struct TransactionMetaV3
{
    ExtensionPoint ext;

    LedgerEntryChanges txChangesBefore;  // tx level changes before operations
                                         // are applied if any
    OperationMeta operations<>;          // meta for each operation
    LedgerEntryChanges txChangesAfter;   // tx level changes after operations are
                                         // applied if any
    SorobanTransactionMeta* sorobanMeta; // Soroban-specific meta (only for 
                                         // Soroban transactions).
};

// This is in Stellar-ledger.x to due to a circular dependency 
struct InvokeHostFunctionSuccessPreImage
{
    SCVal returnValue;
    ContractEvent events<>;
};

// this is the meta produced when applying transactions
// it does not include pre-apply updates such as fees
union TransactionMeta switch (int v)
{
case 0:
    OperationMeta operations<>;
case 1:
    TransactionMetaV1 v1;
case 2:
    TransactionMetaV2 v2;
case 3:
    TransactionMetaV3 v3;
};

// This struct groups together changes on a per transaction basis
// note however that fees and transaction application are done in separate
// phases
struct TransactionResultMeta
{
    TransactionResultPair result;
    LedgerEntryChanges feeProcessing;
    TransactionMeta txApplyProcessing;
};

// this represents a single upgrade that was performed as part of a ledger
// upgrade
struct UpgradeEntryMeta
{
    LedgerUpgrade upgrade;
    LedgerEntryChanges changes;
};

struct LedgerCloseMetaV0
{
    LedgerHeaderHistoryEntry ledgerHeader;
    // NB: txSet is sorted in "Hash order"
    TransactionSet txSet;

    // NB: transactions are sorted in apply order here
    // fees for all transactions are processed first
    // followed by applying transactions
    TransactionResultMeta txProcessing<>;

    // upgrades are applied last
    UpgradeEntryMeta upgradesProcessing<>;

    // other misc information attached to the ledger close
    SCPHistoryEntry scpInfo<>;
};

struct LedgerCloseMetaExtV1
{
    ExtensionPoint ext;
    int64 sorobanFeeWrite1KB;
};

union LedgerCloseMetaExt switch (int v)
{
case 0:
    void;
case 1:
    LedgerCloseMetaExtV1 v1;
};

struct LedgerCloseMetaV1
{
    LedgerCloseMetaExt ext;

    LedgerHeaderHistoryEntry ledgerHeader;

    GeneralizedTransactionSet txSet;

    // NB: transactions are sorted in apply order here
    // fees for all transactions are processed first
    // followed by applying transactions
    TransactionResultMeta txProcessing<>;

    // upgrades are applied last
    UpgradeEntryMeta upgradesProcessing<>;

    // other misc information attached to the ledger close
    SCPHistoryEntry scpInfo<>;

    // Size in bytes of BucketList, to support downstream
    // systems calculating storage fees correctly.
    uint64 totalByteSizeOfBucketList;

    // Temp keys that are being evicted at this ledger.
    LedgerKey evictedTemporaryLedgerKeys<>;

    // Archived restorable ledger entries that are being
    // evicted at this ledger.
    LedgerEntry evictedPersistentLedgerEntries<>;
};

union LedgerCloseMeta switch (int v)
{
case 0:
    LedgerCloseMetaV0 v0;
case 1:
    LedgerCloseMetaV1 v1;
};
}
