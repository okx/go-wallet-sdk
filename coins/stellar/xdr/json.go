package xdr

import (
	"encoding/json"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/stellar/support/errors"
	"regexp"
	"strconv"
	"time"
)

// iso8601Time is a timestamp which supports parsing dates which have a year outside the 0000..9999 range
type iso8601Time struct {
	time.Time
}

// reISO8601 is the regular expression used to parse date strings in the
// ISO 8601 extended format, with or without an expanded year representation.
var reISO8601 = regexp.MustCompile(`^([-+]?\d{4,})-(\d{2})-(\d{2})`)

// MarshalJSON serializes the timestamp to a string
func (t iso8601Time) MarshalJSON() ([]byte, error) {
	ts := t.Format(time.RFC3339)
	if t.Year() > 9999 {
		ts = "+" + ts
	}

	return json.Marshal(ts)
}

// UnmarshalJSON parses a JSON string into a iso8601Time instance.
func (t *iso8601Time) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == nil {
		return nil
	}

	text := *s
	m := reISO8601.FindStringSubmatch(text)

	if len(m) != 4 {
		return fmt.Errorf("UnmarshalJSON: cannot parse %s", text)
	}
	// No need to check for errors since the regexp guarantees the matches
	// are valid integers
	year, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	day, _ := strconv.Atoi(m[3])

	ts, err := time.Parse(time.RFC3339, "2006-01-02"+text[len(m[0]):])
	if err != nil {
		return errors.Wrap(err, "Could not extract time")
	}

	t.Time = time.Date(year, time.Month(month), day, ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond(), ts.Location())
	return nil
}

func newiso8601Time(epoch int64) *iso8601Time {
	return &iso8601Time{time.Unix(epoch, 0).UTC()}
}

type claimPredicateJSON struct {
	And            *[]claimPredicateJSON `json:"and,omitempty"`
	Or             *[]claimPredicateJSON `json:"or,omitempty"`
	Not            *claimPredicateJSON   `json:"not,omitempty"`
	Unconditional  bool                  `json:"unconditional,omitempty"`
	AbsBefore      *iso8601Time          `json:"abs_before,omitempty"`
	AbsBeforeEpoch *int64                `json:"abs_before_epoch,string,omitempty"`
	RelBefore      *int64                `json:"rel_before,string,omitempty"`
}
