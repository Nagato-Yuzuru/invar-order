package product

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinateKey_Golden(t *testing.T) {
	tests := []struct {
		name string
		sku  SKUCoordinate
		want string
	}{
		{
			"base case (No Dims)",
			SKUCoordinate{
				ProductCode: "P1",
				Scope:       "US",
			},
			"US:P1",
		},
		{
			"With Golden Dims",
			SKUCoordinate{
				"P1",
				"US",
				Dimensions{
					"color": "red",
				},
			},
			// Golden about xxhash
			"US:P1:c595431927dde5a6",
		},
		{
			name: "Empty Dims",
			sku: SKUCoordinate{
				ProductCode: "P1",
				Scope:       "US",
				Dims:        Dimensions{},
			},
			want: "US:P1",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, tt.sku.CoordinateKey())
			},
		)
	}
}

func FuzzDimensions_SortedKeys(f *testing.F) {
	f.Add(
		[]byte{
			'a',
			'b', 'c', 'd',
		},
	)
	f.Fuzz(
		func(t *testing.T, data []byte) {
			dims := adapterFuzzToDims(data)
			dimPairs := shuffleMapToSlice(dims)
			shuffledDims := make(Dimensions, len(dimPairs))
			for _, p := range dimPairs {
				shuffledDims[p.key] = p.value
			}

			assert.Equal(t, dims.SortedKeys(), shuffledDims.SortedKeys())
		},
	)
}

func FuzzCoordinateKey_Stability(f *testing.F) {
	f.Add(
		"key1", "val1", "key2", "val2",
	)

	f.Fuzz(
		func(t *testing.T, k1, v1, k2, v2 string) {
			dims1 := Dimensions{}
			dims1[k1] = v1
			dims1[k2] = v2

			dims2 := Dimensions{}
			dims2[k2] = v2
			dims2[k1] = v1

			dims3 := Dimensions{}
			dims3[k1] = v1
			dims3[k2] = v2 + "diff"

			sku1 := SKUCoordinate{ProductCode: "P1", Scope: "US", Dims: dims1}
			sku2 := SKUCoordinate{ProductCode: "P1", Scope: "US", Dims: dims2}
			sku3 := SKUCoordinate{ProductCode: "P1", Scope: "US", Dims: dims3}

			assert.Equal(t, sku1.CoordinateKey(), sku2.CoordinateKey())
			assert.NotEqual(t, sku1.CoordinateKey(), sku3.CoordinateKey())
		},
	)
}

func adapterFuzzToDims(fuzzBytes []byte) Dimensions {
	bytesLen := len(fuzzBytes)
	const step = 10
	if bytesLen < step*2 {
		return Dimensions{}
	}

	dims := make(Dimensions, bytesLen/(2*step))
	i := 0
	for i+step*2 < bytesLen {
		dims[string(fuzzBytes[i:i+step])] = string(fuzzBytes[i+step : i+step*2])
		i += step
	}
	return dims
}

type pair[K comparable, V any] struct {
	key   K
	value V
}

func shuffleMapToSlice[M ~map[K]V, K comparable, V any](m M) []pair[K, V] {
	newMap := make([]pair[K, V], 0, len(m))
	for k, v := range m {
		newMap = append(newMap, pair[K, V]{k, v})
	}
	rand.Shuffle(
		len(newMap), func(i, j int) {
			newMap[i], newMap[j] = newMap[j], newMap[i]
		},
	)
	return newMap
}

func BenchmarkCoordinateKey(b *testing.B) {
	sku := SKUCoordinate{
		ProductCode: "PROD-10086",
		Scope:       "Global-Main-Scope",
		Dims: Dimensions{
			"region":   "us-east-1",
			"env":      "production",
			"customer": "enterprise-a",
			"term":     "YEAR",
		},
	}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = sku.CoordinateKey()
	}
}
