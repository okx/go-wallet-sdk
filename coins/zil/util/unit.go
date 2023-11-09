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
package util

import (
	"math"
)

const (
	ZIL = iota
	LI
	QA
)

func FromQa(qa float64, unit int, is_pack bool) float64 {
	rate := 1.0

	switch unit {
	case ZIL:
		rate = 1000000000000.0
	case LI:
		rate = 1000000.0
	}

	ret := qa / rate

	if is_pack {
		ret = math.Round(ret)
	}

	return ret
}

func ToQa(qa float64, unit int) float64 {
	rate := 1.0

	switch unit {
	case ZIL:
		rate = 1000000000000.0
	case LI:
		rate = 1000000.0
	}

	return qa * rate
}
