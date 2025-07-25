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

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestMarshalMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Value        Map
		ExpectedYAML string
		ExpectedJSON string
	}{
		{
			Value: Map{
				MustMakeKey("my", "testKey"):        NewValue("testValue"),
				MustMakeKey("my", "anotherTestKey"): NewValue("anotherTestValue"),
			},
			ExpectedYAML: "my:anotherTestKey: anotherTestValue\nmy:testKey: testValue\n",
			ExpectedJSON: `{"my:anotherTestKey":"anotherTestValue","my:testKey":"testValue"}`,
		},
		{
			Value: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":12321123131}`),
			},
			ExpectedYAML: "my:parent:\n  nested: 12321123131\n",
			ExpectedJSON: `{"my:parent":{"nested":12321123131}}`,
		},
		{
			Value: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":[12321123131]}`),
			},
			ExpectedYAML: "my:parent:\n  nested:\n  - 12321123131\n",
			ExpectedJSON: `{"my:parent":{"nested":[12321123131]}}`,
		},
		{
			Value: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":{"foo":12321123131}}`),
			},
			ExpectedYAML: "my:parent:\n  nested:\n    foo: 12321123131\n",
			ExpectedJSON: `{"my:parent":{"nested":{"foo":12321123131}}}`,
		},
		{
			Value: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":4.2}`),
			},
			ExpectedYAML: "my:parent:\n  nested: 4.2\n",
			ExpectedJSON: `{"my:parent":{"nested":4.2}}`,
		},
		{
			Value: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":[4.2]}`),
			},
			ExpectedYAML: "my:parent:\n  nested:\n  - 4.2\n",
			ExpectedJSON: `{"my:parent":{"nested":[4.2]}}`,
		},
		{
			Value: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":{"foo":4.2}}`),
			},
			ExpectedYAML: "my:parent:\n  nested:\n    foo: 4.2\n",
			ExpectedJSON: `{"my:parent":{"nested":{"foo":4.2}}}`,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.ExpectedYAML, func(t *testing.T) {
			t.Parallel()

			yamlBytes, err := yaml.Marshal(test.Value)
			require.NoError(t, err)
			assert.Equal(t, test.ExpectedYAML, string(yamlBytes))
			newYAMLMap, err := roundtripMapYAML(test.Value)
			require.NoError(t, err)
			assert.Equal(t, test.Value, newYAMLMap)

			jsonBytes, err := json.Marshal(test.Value)
			require.NoError(t, err)
			assert.Equal(t, test.ExpectedJSON, string(jsonBytes))
			newJSONMap, err := roundtripMapJSON(test.Value)
			require.NoError(t, err)
			assert.Equal(t, test.Value, newJSONMap)
		})
	}
}

func TestMarshalling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Value    map[string]interface{}
		Expected Map
	}{
		{
			Value: map[string]interface{}{
				"my:anotherTestKey": "anotherTestValue",
				"my:testKey":        "testValue",
			},
			Expected: Map{
				MustMakeKey("my", "testKey"):        NewValue("testValue"),
				MustMakeKey("my", "anotherTestKey"): NewValue("anotherTestValue"),
			},
		},
		{
			Value: map[string]interface{}{
				"my:secureTestKey": map[string]interface{}{
					"secure": "securevalue",
				},
			},
			Expected: Map{
				MustMakeKey("my", "secureTestKey"): NewSecureValue("securevalue"),
			},
		},
		{
			Value: map[string]interface{}{
				"my:arrayKey": []string{"a", "b", "c"},
			},
			Expected: Map{
				MustMakeKey("my", "arrayKey"): NewObjectValue(`["a","b","c"]`),
			},
		},
		{
			Value: map[string]interface{}{
				"my:mapKey": map[string]interface{}{
					"a": "b",
					"c": "d",
				},
			},
			Expected: Map{
				MustMakeKey("my", "mapKey"): NewObjectValue(`{"a":"b","c":"d"}`),
			},
		},
		{
			Value: map[string]interface{}{
				"my:servers": []interface{}{
					map[string]interface{}{"port": 80, "host": "example"},
				},
			},
			Expected: Map{
				MustMakeKey("my", "servers"): NewObjectValue(`[{"host":"example","port":80}]`),
			},
		},
		{
			Value: map[string]interface{}{
				"my:mapKey": map[string][]int{
					"nums": {1, 2, 3},
				},
			},
			Expected: Map{
				MustMakeKey("my", "mapKey"): NewObjectValue(`{"nums":[1,2,3]}`),
			},
		},
		{
			Value: map[string]interface{}{
				"my:mapKey": map[string]interface{}{
					"a": map[string]interface{}{"secure": "securevalue"},
					"c": "d",
				},
			},
			Expected: Map{
				MustMakeKey("my", "mapKey"): NewSecureObjectValue(`{"a":{"secure":"securevalue"},"c":"d"}`),
			},
		},
		{
			Value: map[string]interface{}{
				"my:servers": []interface{}{
					map[string]interface{}{
						"port": 80,
						"host": "example",
						"token": map[string]interface{}{
							"secure": "securevalue",
						},
					},
				},
			},
			Expected: Map{
				Key{
					namespace: "my",
					name:      "servers",
				}: NewSecureObjectValue(`[{"host":"example","port":80,"token":{"secure":"securevalue"}}]`),
			},
		},
		{
			Value: map[string]interface{}{
				"my:mapKey": map[string]interface{}{
					"a": map[string]interface{}{"secure": "foo", "bar": "blah"},
					"c": "d",
				},
			},
			Expected: Map{
				MustMakeKey("my", "mapKey"): NewObjectValue(`{"a":{"bar":"blah","secure":"foo"},"c":"d"}`),
			},
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		yamlBytes, err := yaml.Marshal(test.Value)
		require.NoError(t, err)
		t.Run(fmt.Sprintf("YAML: %s", yamlBytes), func(t *testing.T) {
			t.Parallel()

			var m Map
			err := yaml.Unmarshal(yamlBytes, &m)

			require.NoError(t, err)
			assert.Equal(t, test.Expected, m)

			newM, err := roundtripMapYAML(m)
			require.NoError(t, err)
			assert.Equal(t, m, newM)
		})

		jsonBytes, err := json.Marshal(test.Value)
		require.NoError(t, err)
		t.Run(fmt.Sprintf("JSON: %s", jsonBytes), func(t *testing.T) {
			t.Parallel()

			var m Map
			err := json.Unmarshal(jsonBytes, &m)

			require.NoError(t, err)
			assert.Equal(t, test.Expected, m)

			newM, err := roundtripMapJSON(m)
			require.NoError(t, err)
			assert.Equal(t, m, newM)
		})
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Config     Map
		Expected   map[Key]string
		SecureKeys []Key
	}{
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
			Expected: map[Key]string{
				MustMakeKey("my", "testKey"): "testValue",
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewSecureValue("securevalue"),
			},
			Expected: map[Key]string{
				MustMakeKey("my", "testKey"): "[secret]",
			},
			SecureKeys: []Key{MustMakeKey("my", "testKey")},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`{"inner":"value"}`),
			},
			Expected: map[Key]string{
				MustMakeKey("my", "testKey"): `{"inner":"value"}`,
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewSecureObjectValue(`{"inner":{"secure":"securevalue"}}`),
			},
			Expected: map[Key]string{
				MustMakeKey("my", "testKey"): `{"inner":"[secret]"}`,
			},
			SecureKeys: []Key{MustMakeKey("my", "testKey")},
		},
		{
			Config: Map{
				//nolint:lll
				MustMakeKey("my", "testKey"): NewSecureObjectValue(`[{"inner":{"secure":"securevalue"}},{"secure":"securevalue2"}]`),
			},
			Expected: map[Key]string{
				MustMakeKey("my", "testKey"): `[{"inner":"[secret]"},"[secret]"]`,
			},
			SecureKeys: []Key{MustMakeKey("my", "testKey")},
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			decrypter := NewBlindingDecrypter()
			actual, err := test.Config.Decrypt(decrypter)
			require.NoError(t, err)
			assert.Equal(t, test.Expected, actual)

			assert.Equal(t, len(test.SecureKeys) != 0, test.Config.HasSecureValue())
			assert.ElementsMatch(t, test.SecureKeys, test.Config.SecureKeys())
		})
	}
}

func TestGetSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Key            string
		Path           bool
		Config         Map
		Expected       Value
		ExpectNotFound bool
	}{
		{
			Key: "my:testKey",
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
			Expected: NewValue("testValue"),
		},
		{
			Key: "my:testKey",
			Config: Map{
				MustMakeKey("my", "testKey"): NewSecureValue("secureValue"),
			},
			Expected: NewSecureValue("secureValue"),
		},
		{
			Key: "my:test.Key",
			Config: Map{
				MustMakeKey("my", "test.Key"): NewValue("testValue"),
			},
			Expected: NewValue("testValue"),
		},
		{
			Key:  "my:0",
			Path: true,
			Config: Map{
				MustMakeKey("my", "0"): NewValue("testValue"),
			},
			Expected: NewValue("testValue"),
		},
		{
			Key:  `my:["testKey"]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
			Expected: NewValue("testValue"),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
			Expected: NewValue("value"),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":{"secure":"securevalue"}}`),
			},
			Expected: NewSecureValue("securevalue"),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":true}`),
			},
			Expected: NewValue("true"),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":false}`),
			},
			Expected: NewValue("false"),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":100}`),
			},
			Expected: NewValue("100"),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":-2}`),
			},
			Expected: NewValue("-2"),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":{"nested":"foo"}}`),
			},
			Expected: NewObjectValue(`{"nested":"foo"}`),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":{"nested":{"secure":"securevalue"}}}`),
			},
			Expected: NewSecureObjectValue(`{"nested":{"secure":"securevalue"}}`),
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":{"nested":{"a":"b","secure":"val"}}}`),
			},
			Expected: NewObjectValue(`{"nested":{"a":"b","secure":"val"}}`),
		},
		{
			Key:  `my:testKey`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`["a"]`),
			},
			Expected: NewObjectValue(`["a"]`),
		},
		{
			Key:  `my:names[0]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: NewValue("a"),
		},
		{
			Key:  `my:names[1]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: NewValue("b"),
		},
		{
			Key:  `my:names[2]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: NewValue("c"),
		},
		{
			Key:  `my:names[3]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			ExpectNotFound: true,
		},
		{
			Key:  `my:outer.inner.nested`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"hi"}`),
			},
			ExpectNotFound: true,
		},
		{
			Key:  `my:parent.nested`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":12321123131}`),
			},
			Expected: NewValue("12321123131"),
		},
		{
			Key:  `my:parent.nested`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "parent"): NewObjectValue(`{"nested":4.2}`),
			},
			Expected: NewValue("4.2"),
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			key, err := ParseKey(test.Key)
			require.NoError(t, err)

			v, ok, err := test.Config.Get(key, test.Path)
			require.NoError(t, err)
			if test.ExpectNotFound {
				assert.False(t, ok)
				assert.Equal(t, Value{}, v)
			} else {
				assert.True(t, ok)
				assert.Equal(t, test.Expected, v)
			}
		})
	}
}

func TestGetFail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Key           string
		ExpectedError string
	}{
		{
			Key:           `my:["foo`,
			ExpectedError: "invalid config key path: missing closing quote in property name",
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(test.Key, func(t *testing.T) {
			t.Parallel()

			config := make(Map)

			key, err := ParseKey(test.Key)
			require.NoError(t, err)

			_, found, err := config.Get(key, true /*path*/)
			assert.False(t, found)
			assert.EqualError(t, err, test.ExpectedError)
		})
	}
}

func TestRemoveSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Key      string
		Path     bool
		Config   Map
		Expected Map
	}{
		{
			Key:      "my:testKey",
			Config:   Map{},
			Expected: Map{},
		},
		{
			Key: "my:testKey",
			Config: Map{
				MustMakeKey("my", "anotherTestKey"): NewValue("testValue"),
			},
			Expected: Map{
				MustMakeKey("my", "anotherTestKey"): NewValue("testValue"),
			},
		},
		{
			Key: "my:testKey",
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
			Expected: Map{},
		},
		{
			Key: "my:anotherTestKey",
			Config: Map{
				MustMakeKey("my", "testKey"):        NewValue("testValue"),
				MustMakeKey("my", "anotherTestKey"): NewValue("anotherTestValue"),
			},
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
		},
		{
			Key: "my:testKey",
			Config: Map{
				MustMakeKey("my", "testKey"): NewSecureValue("secureValue"),
			},
			Expected: Map{},
		},
		{
			Key:  `my:outer`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
			Expected: Map{},
		},
		{
			Key:  `my:outer.inner`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{}`),
			},
		},
		{
			Key:  `my:names[0]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["b","c"]`),
			},
		},
		{
			Key:  `my:names[1]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","c"]`),
			},
		},
		{
			Key:  `my:names[2]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b"]`),
			},
		},
		{
			Key:  `my:names[3]`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "names"): NewObjectValue(`["a","b","c"]`),
			},
		},
		{
			Key:  `my:outer.inner.nested`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":{"nested": "value"}}`),
			},
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":{}}`),
			},
		},
		{
			Key:  `my:outer[0].nested`,
			Path: true,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`[{"nested": "value"}]`),
			},
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`[{}]`),
			},
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			key, err := ParseKey(test.Key)
			require.NoError(t, err)
			err = test.Config.Remove(key, test.Path)
			require.NoError(t, err)
			assert.Equal(t, test.Expected, test.Config)
		})
	}
}

func TestRemoveFail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Key           string
		Config        Map
		ExpectedError string
	}{
		{
			Key:           `my:["foo`,
			Config:        Map{},
			ExpectedError: "invalid config key path: missing closing quote in property name",
		},
		{
			Key: `my:foo.bar`,
			Config: Map{
				MustMakeKey("my", "foo"): NewObjectValue(`{"bar":"baz","secure":"myvalue"}`),
			},
			ExpectedError: "bar.bar: maps with the single key \"secure\" are reserved",
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			key, err := ParseKey(test.Key)
			require.NoError(t, err)

			err = test.Config.Remove(key, true /*path*/)
			assert.EqualError(t, err, test.ExpectedError)
		})
	}
}

func TestSetSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Key      string
		Value    Value
		Path     bool
		Config   Map
		Expected Map
	}{
		{
			Key:   "my:testKey",
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
		},
		{
			Key:   "my:anotherTestKey",
			Value: NewValue("anotherTestValue"),
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
			Expected: Map{
				MustMakeKey("my", "testKey"):        NewValue("testValue"),
				MustMakeKey("my", "anotherTestKey"): NewValue("anotherTestValue"),
			},
		},
		{
			Key:   "my:0",
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "0"): NewValue("testValue"),
			},
		},
		{
			Key:   "my:true",
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "true"): NewValue("testValue"),
			},
		},
		{
			Key:   "my:test.Key",
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "test.Key"): NewValue("testValue"),
			},
		},
		{
			Key:   "my:testKey",
			Path:  true,
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
		},
		{
			Key:   "my:0",
			Path:  true,
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "0"): NewValue("testValue"),
			},
		},
		{
			Key:   "my:true",
			Path:  true,
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "true"): NewValue("testValue"),
			},
		},
		{
			Key:   `my:["0"]`,
			Path:  true,
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "0"): NewValue("testValue"),
			},
		},
		{
			Key:   `my:["true"]`,
			Path:  true,
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "true"): NewValue("testValue"),
			},
		},
		{
			Key:   `my:["test.Key"]`,
			Path:  true,
			Value: NewValue("testValue"),
			Expected: Map{
				MustMakeKey("my", "test.Key"): NewValue("testValue"),
			},
		},
		{
			Key:   `my:nested["test.Key"]`,
			Path:  true,
			Value: NewValue("value"),
			Expected: Map{
				MustMakeKey("my", "nested"): NewObjectValue(`{"test.Key":"value"}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("value"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "outer"): NewValue("value"),
			},
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "outer"): NewSecureValue("securevalue"),
			},
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("true"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":true}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("false"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":false}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("10"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":10}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("0"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":0}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("-1"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":-1}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("00"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"00"}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("01"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"01"}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("0123456"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"0123456"}`),
			},
		},
		{
			Key:   `my:array[0]`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "array"): NewValue("value"),
			},
			Expected: Map{
				MustMakeKey("my", "array"): NewObjectValue(`["value"]`),
			},
		},
		{
			Key:   `my:array[0]`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "array"): NewSecureValue("value"),
			},
			Expected: Map{
				MustMakeKey("my", "array"): NewObjectValue(`["value"]`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"existing":"existingValue"}`),
			},
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"existing":"existingValue","inner":"value"}`),
			},
		},
		{
			Key:   `my:outer.inner`,
			Path:  true,
			Value: NewSecureValue("securevalue"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewSecureObjectValue(`{"inner":{"secure":"securevalue"}}`),
			},
		},
		{
			Key:   `my:outer.inner.nested`,
			Path:  true,
			Value: NewValue("value"),
			Expected: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":{"nested":"value"}}`),
			},
		},
		{
			Key:   `my:name[0]`,
			Path:  true,
			Value: NewValue("value"),
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["value"]`),
			},
		},
		{
			Key:   `my:name[0][0]`,
			Path:  true,
			Value: NewValue("value"),
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`[["value"]]`),
			},
		},
		{
			Key:   `my:name[0]`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["value","b","c"]`),
			},
		},
		{
			Key:   `my:name[1]`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","value","c"]`),
			},
		},
		{
			Key:   `my:name[2]`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","value"]`),
			},
		},
		{
			Key:   `my:name[3]`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c","value"]`),
			},
		},
		{
			Key:   `my:name[3][0]`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c",["value"]]`),
			},
		},
		{
			Key:   `my:name[3][0]nested`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c",[{"nested":"value"}]]`),
			},
		},
		{
			Key:   `my:name[3].foo.bar`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c",{"foo":{"bar":"value"}}]`),
			},
		},
		{
			Key:   `my:servers[0].name`,
			Path:  true,
			Value: NewValue("foo"),
			Expected: Map{
				MustMakeKey("my", "servers"): NewObjectValue(`[{"name":"foo"}]`),
			},
		},
		{
			Key:   `my:servers[0].host`,
			Path:  true,
			Value: NewValue("example"),
			Config: Map{
				MustMakeKey("my", "servers"): NewObjectValue(`[{"name":"foo"}]`),
			},
			Expected: Map{
				MustMakeKey("my", "servers"): NewObjectValue(`[{"host":"example","name":"foo"}]`),
			},
		},
		{
			Key:   `my:name[0]`,
			Path:  true,
			Value: NewSecureValue("securevalue"),
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewSecureObjectValue(`[{"secure":"securevalue"},"b","c"]`),
			},
		},
		{
			Key:   `my:testKey`,
			Value: NewValue("false"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("false"),
			},
		},
		{
			Key:   `my:testKey`,
			Value: NewValue("true"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("true"),
			},
		},
		{
			Key:   `my:testKey`,
			Value: NewValue("10"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("10"),
			},
		},
		{
			Key:   `my:testKey`,
			Value: NewValue("-1"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("-1"),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("false"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`[false]`),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("true"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`[true]`),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("10"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`[10]`),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("0"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`[0]`),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("-1"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`[-1]`),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("00"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`["00"]`),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("01"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`["01"]`),
			},
		},
		{
			Key:   `my:testKey[0]`,
			Path:  true,
			Value: NewValue("0123456"),
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`["0123456"]`),
			},
		},
		{
			Key:   `my:key.secure`,
			Path:  true,
			Value: NewValue("value"),
			Config: Map{
				MustMakeKey("my", "key"): NewObjectValue(`{"bar":"baz"}`),
			},
			Expected: Map{
				MustMakeKey("my", "key"): NewObjectValue(`{"bar":"baz","secure":"value"}`),
			},
		},
		{
			Key:   `my:special.object`,
			Path:  true,
			Value: NewObjectValue(`{"foo":"bar","fizz":"buzz"}`),
			Config: Map{
				MustMakeKey("my", "special"): NewObjectValue(`{"thing1":1,"thing2":2}`),
			},
			Expected: Map{
				MustMakeKey("my", "special"): NewObjectValue(`{"object":{"fizz":"buzz","foo":"bar"},"thing1":1,"thing2":2}`),
			},
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			if test.Config == nil {
				test.Config = make(Map)
			}

			key, err := ParseKey(test.Key)
			require.NoError(t, err)

			err = test.Config.Set(key, test.Value, test.Path)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, test.Config)
		})
	}
}

func TestSetFail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Key           string
		Config        Map
		ExpectedError string
	}{
		// Syntax errors.
		{
			Key:           "my:root[",
			ExpectedError: "invalid config key path: missing closing bracket in array index",
		},
		{
			Key:           `my:root["nested]`,
			ExpectedError: "invalid config key path: missing closing quote in property name",
		},
		{
			Key:           "my:root.array[abc]",
			ExpectedError: "invalid config key path: invalid array index: strconv.ParseInt: parsing \"abc\": invalid syntax",
		},

		// First path component must be a string.
		{
			Key:           `my:[""]`,
			ExpectedError: "config key is empty",
		},
		{
			Key:           "my:[0]",
			ExpectedError: "first path segement of config key must be a string",
		},

		// Index out of range.
		{
			Key:           `my:name[-1]`,
			ExpectedError: "name[-1]: array index out of range",
		},
		{
			Key: `my:name[4]`,
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`["a","b","c"]`),
			},
			ExpectedError: "name[4]: array index out of range",
		},

		// A "secure" key that is a map with a single string value is reserved by the system.
		{
			Key:           `my:key.secure`,
			ExpectedError: "maps with the single key \"secure\" are reserved",
		},
		{
			Key:           `my:super.nested.map.secure`,
			ExpectedError: "maps with the single key \"secure\" are reserved",
		},

		// Type mismatches.
		{
			Key: `my:outer.inner`,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue("[1,2,3]"),
			},
			ExpectedError: "outer.inner: key for an array must be an int",
		},
		{
			Key: `my:array[0]`,
			Config: Map{
				MustMakeKey("my", "array"): NewObjectValue(`{"inner":"value"}`),
			},
			ExpectedError: "array[0]: key for a map must be a string",
		},
		{
			Key: `my:outer.inner.nested`,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
			ExpectedError: "outer.inner: expected a map",
		},
		{
			Key: `my:outer.inner[0]`,
			Config: Map{
				MustMakeKey("my", "outer"): NewObjectValue(`{"inner":"value"}`),
			},
			ExpectedError: "outer.inner: expected an array",
		},

		// Strict path parsing
		{
			Key:           `my:root.[1]"`,
			ExpectedError: "invalid config key path: expected property name after '.'",
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			if test.Config == nil {
				test.Config = make(Map)
			}

			key, err := ParseKey(test.Key)
			require.NoError(t, err)

			err = test.Config.Set(key, NewValue("value"), true /*path*/)
			assert.EqualError(t, err, test.ExpectedError)
		})
	}
}

func TestCopyMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Config   Map
		Expected Map
	}{
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
			Expected: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewSecureValue("stackAsecurevalue"),
			},
			Expected: Map{
				MustMakeKey("my", "testKey"): NewSecureValue("stackBsecurevalue"),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`{"inner":"value"}`),
			},
			Expected: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`{"inner":"value"}`),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewSecureObjectValue(`{"inner":{"secure":"stackAsecurevalue"}}`),
			},
			Expected: Map{
				MustMakeKey("my", "testKey"): NewSecureObjectValue(`{"inner":{"secure":"stackBsecurevalue"}}`),
			},
		},
		{
			Config: Map{
				//nolint:lll
				MustMakeKey("my", "testKey"): NewSecureObjectValue(`[{"inner":{"secure":"stackAsecurevalue"}},{"secure":"stackAsecurevalue2"}]`),
			},
			Expected: Map{
				//nolint:lll
				MustMakeKey("my", "testKey"): NewSecureObjectValue(`[{"inner":{"secure":"stackBsecurevalue"}},{"secure":"stackBsecurevalue2"}]`),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "test.Key"): NewValue("testValue"),
			},
			Expected: Map{
				MustMakeKey("my", "test.Key"): NewValue("testValue"),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`[["value"]]`),
			},
			Expected: Map{
				MustMakeKey("my", "name"): NewObjectValue(`[["value"]]`),
			},
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			newConfig, err := test.Config.Copy(newPrefixCrypter("stackA"), newPrefixCrypter("stackB"))
			require.NoError(t, err)

			assert.Equal(t, test.Expected, newConfig)
		})
	}
}

func TestPropertyMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Config   Map
		Expected resource.PropertyMap
	}{
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("testValue"),
			},
			Expected: resource.PropertyMap{
				"my:testKey": resource.NewStringProperty("testValue"),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("1"),
			},
			Expected: resource.PropertyMap{
				"my:testKey": resource.NewNumberProperty(1.0),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("18446744073709551615"),
			},
			Expected: resource.PropertyMap{
				"my:testKey": resource.NewNumberProperty(1.8446744073709552e+19),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewValue("true"),
			},
			Expected: resource.PropertyMap{
				"my:testKey": resource.NewBoolProperty(true),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewSecureValue("stackAsecurevalue"),
			},
			Expected: resource.PropertyMap{
				"my:testKey": resource.MakeSecret(resource.NewStringProperty("stackAsecurevalue")),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "testKey"): NewObjectValue(`{"inner":"value"}`),
			},
			Expected: resource.PropertyMap{
				"my:testKey": resource.NewObjectProperty(resource.PropertyMap{
					"inner": resource.NewStringProperty("value"),
				}),
			},
		},
		{
			Config: Map{
				//nolint:lll
				MustMakeKey("my", "testKey"): NewSecureObjectValue(`[{"inner":{"secure":"stackAsecurevalue"}},{"secure":"stackAsecurevalue2"}]`),
			},
			Expected: resource.PropertyMap{
				//nolint:lll
				"my:testKey": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewObjectProperty(resource.PropertyMap{
						"inner": resource.MakeSecret(resource.NewStringProperty("stackAsecurevalue")),
					}),
					resource.MakeSecret(resource.NewStringProperty("stackAsecurevalue2")),
				}),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "test.Key"): NewValue("testValue"),
			},
			Expected: resource.PropertyMap{
				"my:test.Key": resource.NewStringProperty("testValue"),
			},
		},
		{
			Config: Map{
				MustMakeKey("my", "name"): NewObjectValue(`[["value"]]`),
			},
			Expected: resource.PropertyMap{
				"my:name": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewArrayProperty([]resource.PropertyValue{
						resource.NewStringProperty("value"),
					}),
				}),
			},
		},
	}

	//nolint:paralleltest // false positive because range var isn't used directly in t.Run(name) arg
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			t.Parallel()

			decrypter := nopCrypter{}
			propMap, err := test.Config.AsDecryptedPropertyMap(context.Background(), decrypter)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, propMap)
		})
	}
}

func roundtripMapYAML(m Map) (Map, error) {
	return roundtripMap(m, yaml.Marshal, yaml.Unmarshal)
}

func roundtripMapJSON(m Map) (Map, error) {
	return roundtripMap(m, json.Marshal, json.Unmarshal)
}

func roundtripMap(m Map, marshal func(v interface{}) ([]byte, error),
	unmarshal func([]byte, interface{}) error,
) (Map, error) {
	b, err := marshal(m)
	if err != nil {
		return nil, err
	}

	var newM Map
	err = unmarshal(b, &newM)
	return newM, err
}
