// Copyright 2015 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSensitiveRawMessageImplements(t *testing.T) {
	var _ fmt.Stringer = SensitiveRawMessage{}
	var _ fmt.GoStringer = SensitiveRawMessage{}
	var _ json.Marshaler = SensitiveRawMessage{}
	var _ json.Unmarshaler = &SensitiveRawMessage{}
}

func TestSensitiveRawMessage(t *testing.T) {
	sensitive := NewSensitiveRawMessage(json.RawMessage("secret"))

	for i, str := range []string{
		sensitive.String(),
		sensitive.GoString(),
		fmt.Sprintf("%v", sensitive),
		fmt.Sprintf("%#v", sensitive),
		fmt.Sprint(sensitive),
	} {
		if str != "[redacted]" {
			t.Errorf("#%v: expected redacted, got %s", i, str)
		}
	}
}
