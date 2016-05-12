package consul

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDateFormats(t *testing.T) {
	cases := []string{
		"2002-10-02T15:00:00.05Z",
		"2002-10-02T15:00:00Z",
		"2013-02-04T22:44:30.652Z",
	}

	for i, tc := range cases {
		input := fmt.Sprintf(`{
			"updated_at": "%s",
			"created_at": "%s"
		}`, tc, tc)

		branch := new(Branch)
		err := json.Unmarshal([]byte(input), branch)
		if err != nil {
			t.Errorf("Unmarshaling of case %d failed: %s", i, err.Error())
			continue
		}

	}
}
