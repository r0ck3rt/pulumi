// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buildutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersions(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"v0.12.0":                                "0.12.0",
		"v0.12.0+dirty":                          "0.12.0+dirty",
		"v0.12.0-rc.1":                           "0.12.0rc1",
		"v0.12.0-rc.1+dirty":                     "0.12.0rc1+dirty",
		"v0.12.1-dev.1524606809+gf2f1178b":       "0.12.1.dev1524606809",
		"v0.12.1-dev.1524606809+gf2f1178b.dirty": "0.12.1.dev1524606809+dirty",
	}

	for ver, expected := range cases {
		p, err := PyPiVersionFromNpmVersion(ver)
		require.NoError(t, err)
		assert.Equal(t, expected, p, "failed parsing '%s'", ver)
	}
}

func TestPseduoVersion(t *testing.T) {
	t.Parallel()

	pseudoVersion := "v1.29.1-0.20200403140640-efb5e2a48a86"
	assert.True(t, IsPseudoVersion(pseudoVersion))

	tagVersion := "v1.29.0"
	assert.False(t, IsPseudoVersion(tagVersion))

	betaVersion := "v1.29.0-beta.1"
	assert.False(t, IsPseudoVersion(betaVersion))
}
