/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package zil

import (
	"regexp"
	"strconv"
)

func IsBech32(addr string) bool {
	match, _ := regexp.MatchString("^zil1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}", addr)
	return match
}
func IsPublicKey(public_key string) bool {
	match, _ := regexp.MatchString("^(0x)?[[:xdigit:]]{66}$", public_key)
	return match
}

func IsPrivateKey(private_key string) bool {
	match, _ := regexp.MatchString("^(0x)?[[:xdigit:]]{64}$", private_key)
	return match
}

func IsAddress(address string) bool {
	match, _ := regexp.MatchString("^(0x)?[[:xdigit:]]{40}$", address)
	return match
}

func IsSignature(signature string) bool {
	match, _ := regexp.MatchString("^(0x)?[[:xdigit:]]{128}$", signature)
	return match
}

func IsByteString(str string, len int) bool {
	pattern := "^(0x)?[0-9a-fA-F]{" + strconv.FormatInt(int64(len), 10) + "}"
	match, _ := regexp.MatchString(pattern, str)
	return match
}
