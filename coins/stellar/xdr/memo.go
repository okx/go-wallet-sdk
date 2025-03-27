/*
 * Copyright 2016 Stellar Development Foundation and contributors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file includes portions of third-party code from [https://github.com/stellar/go].
 * The original code is licensed under the Apache License 2.0.
 */

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
