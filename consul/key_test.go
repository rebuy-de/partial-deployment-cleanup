package consul

import "testing"

func TestKeys(t *testing.T) {
	cases := []struct {
		expect string
		actual Key
	}{
		{"foo", Key("foo")},
		{"foo/bar", Key("foo").Add("bar")},
		{"foo/bar", Key("/foo/").Add("/bar/")},
		{"foo/bar", Key("bish/bash/bosh/foo/bar").Base("bish/bash/bosh")},
		{"foo/bar", Key("/bish/bash/bosh/foo/bar/").Base("bish/bash/bosh")},
		{"foo/bar", Key("bish/bash/bosh/foo/bar").Base("/bish/bash/bosh/")},
		{"bish", Key(Key("bish/bash/bosh/blub/foo/bar").Get(0))},
		{"blub", Key(Key("bish/bash/bosh/blub/foo/bar").Get(3))},
		{"blub", Key(Key("/bish/bash/bosh/blub/foo/bar/").Get(3))},
		{"blub", Key(Key("///bish/bash/bosh/blub/foo/bar/").Get(3))},
		{"", Key(Key("/foo/bar/").Get(3))},
	}

	for i, tc := range cases {
		if string(tc.actual) != tc.expect {
			t.Errorf("Test case %d failed", i)
			t.Errorf("%#v != %#v", tc.expect, string(tc.actual))
		}
	}
}
