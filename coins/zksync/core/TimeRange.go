package core

type TimeRange struct {
	ValidFrom  uint64 `json:"validFrom"`
	ValidUntil uint64 `json:"validUntil"`
}

func DefaultTimeRange() *TimeRange {
	return &TimeRange{
		ValidFrom:  0,
		ValidUntil: 4294967295,
	}
}
