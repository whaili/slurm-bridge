// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package utils

// Get keys from map
func Keys[K comparable, V any](items map[K]V) []K {
	keys := make([]K, len(items))
	i := 0
	for k := range items {
		keys[i] = k
		i++
	}
	return keys
}

// Get values from map
func Values[K comparable, V any](items map[K]V) []V {
	vals := make([]V, len(items))
	i := 0
	for _, v := range items {
		vals[i] = v
		i++
	}
	return vals
}
