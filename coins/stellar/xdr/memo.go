package xdr

func MemoText(text string) Memo {
	return Memo{Type: MemoTypeMemoText, Text: &text}
}

func MemoID(id uint64) Memo {
	idObj := Uint64(id)
	return Memo{Type: MemoTypeMemoId, Id: &idObj}
}

func MemoHash(hash Hash) Memo {
	return Memo{Type: MemoTypeMemoHash, Hash: &hash}
}

func MemoRetHash(hash Hash) Memo {
	return Memo{Type: MemoTypeMemoReturn, RetHash: &hash}
}
