// Package product is used to define the descriptive information of a product.
package product

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"

	"github.com/nagato-yuzuru/invar-order/types"
)

// Dimensions must sort then compare
type Dimensions map[string]string

func (d Dimensions) String() string {
	keys := d.SortedKeys()
	var b strings.Builder

	// XXX: key=5, value=5, overhead=2 -> 12 * len
	b.Grow(len(keys) * 12)

	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(d[k])
		b.WriteByte(';')
	}
	return b.String()
}

// SortedKeys return sorted keys of [Dimensions]
func (d Dimensions) SortedKeys() []string {
	keys := make([]string, 0, len(d))
	keys = slices.AppendSeq(keys, maps.Keys(d))
	slices.Sort(keys)
	return keys
}

// SKUCoordinate uniquely define a specific SKU.
// It can be guaranteed that the results of dim
// will be consistent regardless of the insertion order.
type SKUCoordinate struct {
	ProductCode types.ProductCode
	Scope       types.ProjectScope
	Dims        Dimensions
}

func (s SKUCoordinate) String() string {
	return fmt.Sprintf("%s:%s:%s", s.Scope, s.ProductCode, s.Dims)
}

// CoordinateKey generates a canonical cache key.
// Format: "Scope:Code[:HashHex]"
func (s SKUCoordinate) CoordinateKey() string {
	const sep = ':'

	// expect:：Scope + : + Code + : + 16(HashHex)
	baseLen := len(s.Scope) + 1 + len(s.ProductCode) + 1 + 16
	var b strings.Builder
	b.Grow(baseLen)

	b.WriteString(string(s.Scope))
	b.WriteByte(sep)
	b.WriteString(string(s.ProductCode))

	if len(s.Dims) == 0 {
		return b.String()
	}

	b.WriteByte(sep)

	// HACK: New64a return Hash interface that take extra escape
	// instead of xxhash
	// h := fnv.New64a()

	h := xxhash.New()
	keys := s.Dims.SortedKeys()

	for _, k := range keys {
		h.WriteString(k)         //nolint:gosec,errcheck
		h.WriteString("=")       //nolint:gosec,errcheck
		h.WriteString(s.Dims[k]) //nolint:gosec,errcheck
		h.WriteString(";")       //nolint:gosec,errcheck
	}

	var buf [16]byte // On Stack

	// AppendUint need slice
	hexBytes := strconv.AppendUint(buf[:0], h.Sum64(), 16)
	b.Write(hexBytes)

	return b.String()
}

// SKU describes a product, identifying a specific item.
// It does not include variable information such as price.
type SKU struct {
	ID types.InternalSKUID
	SKUCoordinate
}
