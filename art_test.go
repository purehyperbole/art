package art

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testvalue struct {
	Key    string
	Value  string
	Prefix string
}

func TestARTInsertLookup(t *testing.T) {
	cases := []struct {
		Name     string
		Existing []testvalue
		Lookups  []testvalue
	}{
		{
			"simple",
			[]testvalue{{"test", "1234", "est"}},
			[]testvalue{{"test", "1234", "est"}},
		},
		{
			"normal",
			[]testvalue{{"too", "1234", "oo"}, {"bad", "5678", "ad"}, {"you'll", "9101112", "ou'll"}, {"never", "13141516", "ever"}, {"be", "17181920", "e"}, {"rad", "21222324", "ad"}},
			[]testvalue{{"too", "1234", "oo"}, {"bad", "5678", "ad"}, {"you'll", "9101112", "ou'll"}, {"never", "13141516", "ever"}, {"be", "17181920", "e"}, {"rad", "21222324", "ad"}},
		},
		{
			"derivative",
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "est"}},
			[]testvalue{{"test1234", "bacon", "234"}},
		},
		{
			"split",
			[]testvalue{{"test1234", "bacon", "234"}, {"test", "1234", "est"}},
			[]testvalue{{"test1234", "bacon", "234"}, {"test1234", "bacon", "234"}},
		},
		{
			"split-single-shared-character",
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "est"}, {"test1000", "egg", "est"}},
			[]testvalue{{"test", "1234", "est"}, {"test1234", "bacon", "34"}, {"test1000", "egg", "00"}},
		},
		{
			"complex",
			[]testvalue{{"test", "1234", "st"}, {"test1234", "bacon", "234"}, {"tomato", "egg", "ato"}, {"tamale", "hash browns", "male"}, {"todo", "beans", ""}, {"todos", "mushrooms", "s"}, {"abalienate", "toast", ""}, {"abalienated", "onions", ""}, {"abalienating", "sausage", "ng"}},
			[]testvalue{{"test", "1234", "st"}, {"test1234", "bacon", "234"}, {"tomato", "egg", "ato"}, {"tamale", "hash browns", "male"}, {"todo", "beans", "o"}, {"todos", "mushrooms", ""}, {"abalienate", "toast", ""}, {"abalienated", "onions", ""}, {"abalienating", "sausage", "ng"}},
		},
		{
			"single-character",
			[]testvalue{{"todo", "toast", "odo"}, {"todos", "bacon", ""}},
			[]testvalue{{"todo", "toast", "odo"}, {"todos", "bacon", ""}},
		},
		{
			"mixed",
			[]testvalue{{"unsophisticatedness", "0", "-"}, {"unsophisticate", "1", "-"}, {"unsophisticatedly", "2", "-"}, {"unsophisticated", "3", "-"}, {"unsophistication", "4", "-"}, {"unsophistic", "5", "-"}, {"unsophistically", "6", "-"}, {"unsophistical", "7", "-"}},
			[]testvalue{{"unsophisticatedness", "0", "-"}, {"unsophisticate", "1", "-"}, {"unsophisticatedly", "2", "-"}, {"unsophisticated", "3", "-"}, {"unsophistication", "4", "-"}, {"unsophistic", "5", "-"}, {"unsophistically", "6", "-"}, {"unsophistical", "7", "-"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			r := New()

			for _, kv := range tc.Existing {
				r.Insert([]byte(kv.Key), kv.Value)
			}

			for _, kv := range tc.Lookups {
				value := r.Lookup([]byte(kv.Key))
				require.NotNil(t, value)
				assert.Equal(t, kv.Value, value)
			}
		})
	}
}

func TestARTInsert(t *testing.T) {
	data, err := ioutil.ReadFile("words.txt")
	if err != nil {
		t.SkipNow()
	}

	r := New()

	for _, w := range bytes.Split(data, []byte("\n")) {
		if len(w) < 1 {
			continue
		}
		r.Insert(w, w)
	}

}
