import Crypto

transaction(publicKeys: [Crypto.KeyListEntry], contracts: {String: String}) {
	prepare(signer: AuthAccount) {
		let account = AuthAccount(payer: signer)

		// add all the keys to the account
		for key in publicKeys {
			account.keys.add(publicKey: key.publicKey, hashAlgorithm: key.hashAlgorithm, weight: key.weight)
		}
		
		// add contracts if provided
		for contract in contracts.keys {
			account.contracts.add(name: contract, code: contracts[contract]!.decodeHex())
		}
	}
}
 