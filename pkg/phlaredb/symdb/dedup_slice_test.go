package symdb

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func generateStrings(b *testing.B, n int) []string {
	b.Helper()

	text, err := os.ReadFile("testdata/electricity.txt")
	if err != nil {
		b.FailNow()
		return nil
	}

	items := make([]string, 0, n)
	lines := strings.Split(string(text), "\n")
	for i := 0; i < n; i++ {
		// Select a random item that's non-empty.
		var item string
		for item == "" {
			item = lines[rand.Intn(len(lines))]
			item = strings.TrimSpace(item)
		}

		singleWord := rand.Int()%2 == 0
		if singleWord {
			words := strings.Fields(item)
			item = words[rand.Intn(len(words))]
		}

		items = append(items, item)
	}

	return items
}

func BenchmarkDedupSlice_append(b *testing.B) {
	var strings deduplicatingSlice[string, string, *stringsHelper]
	strings.init()

	var swissStrings swissSlice[string, string, *stringsHelper]
	swissStrings.init()

	const max = 1 << 20 // 1_048_576
	items := generateStrings(b, max)

	for i := 1 << 1; i <= max; i <<= 1 {
		subItems := items[:i]
		dst := make([]uint32, len(subItems))

		b.Run(fmt.Sprintf("stdmap_%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				strings.append(dst, subItems)
			}
		})

		b.Run(fmt.Sprintf("swiss_%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				swissStrings.append(dst, subItems)
			}
		})
	}
}

func BenchmarkDedupSlice_ingest(b *testing.B) {
	var strings deduplicatingSlice[string, string, *stringsHelper]
	strings.init()

	var swissStrings swissSlice[string, string, *stringsHelper]
	swissStrings.init()

	const max = 1 << 20 // 1_048_576
	items := generateStrings(b, max)

	for i := 1 << 1; i <= max; i <<= 1 {
		subItems := items[:i]
		rewriter := &rewriter{}

		b.Run(fmt.Sprintf("stdmap_%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				strings.ingest(subItems, rewriter)
			}
		})

		b.Run(fmt.Sprintf("swiss_%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				swissStrings.ingest(subItems, rewriter)
			}
		})
	}
}
