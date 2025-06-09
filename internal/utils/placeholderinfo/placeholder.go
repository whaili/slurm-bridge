// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package placeholderinfo

import (
	"bytes"
	"encoding/json"

	"k8s.io/utils/ptr"
)

type PlaceholderInfo struct {
	Pods []string `json:"pods"`
}

func (phInfo *PlaceholderInfo) Equal(cmp PlaceholderInfo) bool {
	a, _ := json.Marshal(phInfo)
	b, _ := json.Marshal(cmp)
	return bytes.Equal(a, b)
}

func (podInfo *PlaceholderInfo) ToString() string {
	b, _ := json.Marshal(podInfo)
	return string(b)
}

func ParseIntoPlaceholderInfo(str *string, out *PlaceholderInfo) error {
	data := ptr.Deref(str, "")
	return json.Unmarshal([]byte(data), out)
}
