package product

import (
	"fmt"
	"hash/fnv"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/nagato-yuzuru/invar-order/types"
)

// Dimensions must sort then compare
type Dimensions map[string]string

func (d Dimensions) String() string {
	const exprLength = 10
	keys := d.SortedKeys()
	var b strings.Builder
	b.Grow(len(keys) * exprLength)
	for _, key := range keys {
		b.Write([]byte(key))
		b.Write([]byte("="))
		b.Write([]byte(d[key]))
		b.Write([]byte(";"))
	}
	return b.String()
}

func (d Dimensions) SortedKeys() []string {
	keys := maps.Keys(d)
	return slices.Sorted(keys)
}

// SKUCoordinate uniquely define a specific SKU.
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
	const (
		sep     = ":"
		kvSep   = "="
		pairSep = "|"
	)

	// Scope + sep + Code
	baseLen := len(s.Scope) + len(sep) + len(s.ProductCode)

	var b strings.Builder
	b.Grow(baseLen)
	b.WriteString(string(s.Scope))
	b.WriteString(sep)
	b.WriteString(string(s.ProductCode))

	if len(s.Dims) == 0 {
		return b.String()
	}

	b.Grow(baseLen + 1 + 16)

	b.WriteString(string(s.Scope))
	b.WriteString(sep)
	b.WriteString(string(s.ProductCode))
	b.WriteString(sep)

	h := fnv.New64a()
	h.Write([]byte(s.Dims.String()))

	b.WriteString(strconv.FormatUint(h.Sum64(), 16))

	return b.String()
}

type SKU struct {
	ID types.InternalSKUID
	SKUCoordinate
}
