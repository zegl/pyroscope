package symdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/pyroscope/pkg/pprof"
)

func Test_ingest_strings(t *testing.T) {
	pp, err := pprof.OpenFile("testdata/profile.pb.gz")
	require.NoError(t, err)

	wantRewriter := &rewriter{}
	want := deduplicatingSlice[string, string, *stringsHelper]{}
	want.init()
	want.ingest(pp.StringTable, wantRewriter)

	gotRewriter := &rewriter{}
	got := stringsTable{}
	got.init()
	got.ingest(pp.StringTable, gotRewriter)

	require.ElementsMatch(t, want.slice, got.slice)
	require.Equal(t, want.size.Load(), got.size.Load())
	require.ElementsMatch(t, wantRewriter.strings, gotRewriter.strings)

	gotLookup := make(map[string]int64)
	for i, elem := range got.slice {
		gotLookup[elem] = got.ptrs[i]
	}
	require.Equal(t, want.lookup, gotLookup)
}

func Benchmark_ingest_strings(b *testing.B) {
	pp, err := pprof.OpenFile("testdata/profile.pb.gz")
	require.NoError(b, err)

	b.Run("old", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			rewriter := &rewriter{}
			old := deduplicatingSlice[string, string, *stringsHelper]{}
			old.init()
			old.ingest(pp.StringTable, rewriter)
		}
	})

	b.Run("new", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			rewriter := &rewriter{}
			new := stringsTable{}
			new.init()
			new.ingest(pp.StringTable, rewriter)
		}
	})
}
