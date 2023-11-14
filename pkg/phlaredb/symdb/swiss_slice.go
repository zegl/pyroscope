package symdb

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/colega/zeropool"
	"github.com/dolthub/swiss"

	schemav1 "github.com/grafana/pyroscope/pkg/phlaredb/schemas/v1"
)

var (
	swissInt64SlicePool  zeropool.Pool[[]int64]
	swissUint32SlicePool zeropool.Pool[[]uint32]
)

type swissIDConversionTable struct {
	table *swiss.Map[int64, int64]
}

func (s *swissIDConversionTable) rewrite(idx *int64) {
	pos := *idx
	var ok bool
	*idx, ok = s.table.Get(pos)
	if !ok {
		panic(fmt.Sprintf("unable to rewrite index %d", pos))
	}
}

func (s *swissIDConversionTable) rewriteUint64(idx *uint64) {
	pos := *idx
	v, ok := s.table.Get(int64(pos))
	if !ok {
		panic(fmt.Sprintf("unable to rewrite index %d", pos))
	}
	*idx = uint64(v)
}

func (s *swissIDConversionTable) rewriteUint32(idx *uint32) {
	pos := *idx
	v, ok := s.table.Get(int64(pos))
	if !ok {
		panic(fmt.Sprintf("unable to rewrite index %d", pos))
	}
	*idx = uint32(v)
}

func (s *swissIDConversionTable) keys(fn func(k int64)) {
	s.table.Iter(func(key int64, _ int64) (stop bool) {
		fn(key)
		return false
	})
}

func (s *swissIDConversionTable) entries(fn func(k int64, v int64)) {
	s.table.Iter(func(key int64, value int64) (stop bool) {
		fn(key, value)
		return false
	})
}

type swissSlice[M schemav1.Models, K comparable, H Helper[M, K]] struct {
	lock   sync.RWMutex
	slice  []M
	size   atomic.Uint64
	lookup *swiss.Map[K, int64]

	helper H
}

func (s *swissSlice[M, K, H]) init() {
	s.lookup = swiss.NewMap[K, int64](0)
}

func (s *swissSlice[M, K, H]) ingest(elems []M, rewriter *rewriter) {
	var (
		rewritingMap = swiss.NewMap[int64, int64](0)
		missing      = swissInt64SlicePool.Get()[:0]
	)

	for pos := range elems {
		_ = s.helper.rewrite(rewriter, elems[pos])
	}

	s.lock.RLock()
	for pos := range elems {
		k := s.helper.key(elems[pos])
		posSlice, ok := s.lookup.Get(k)
		if ok {
			rewritingMap.Put(int64(s.helper.setID(uint64(pos), uint64(posSlice), elems[pos])), posSlice)
		} else {
			missing = append(missing, int64(pos))
		}
	}
	s.lock.RUnlock()

	if len(missing) > 0 {
		s.lock.Lock()
		posSlice := int64(len(s.slice))
		for _, pos := range missing {
			k := s.helper.key(elems[pos])
			posSlice2, ok := s.lookup.Get(k)
			if ok {
				rewritingMap.Put(int64(s.helper.setID(uint64(pos), uint64(posSlice2), elems[pos])), posSlice2)
				continue
			}

			s.slice = append(s.slice, s.helper.clone(elems[pos]))
			s.lookup.Put(k, posSlice)
			rewritingMap.Put(int64(s.helper.setID(uint64(pos), uint64(posSlice), elems[pos])), posSlice)
			posSlice++
			s.size.Add(s.helper.size(elems[pos]))
		}
		s.lock.Unlock()
	}

	swissInt64SlicePool.Put(missing)

	s.helper.addToRewriter(rewriter, &swissIDConversionTable{table: rewritingMap})
}

func (s *swissSlice[M, K, H]) append(dst []uint32, elems []M) {
	missing := swissInt64SlicePool.Get()[:0]
	s.lock.RLock()
	for i, v := range elems {
		k := s.helper.key(v)
		if x, ok := s.lookup.Get(k); ok {
			dst[i] = uint32(x)
		} else {
			missing = append(missing, int64(i))
		}
	}
	s.lock.RUnlock()
	if len(missing) > 0 {
		s.lock.RLock()
		p := uint32(len(s.slice))
		for _, i := range missing {
			e := elems[i]
			k := s.helper.key(e)
			x, ok := s.lookup.Get(k)
			if ok {
				dst[i] = uint32(x)
				continue
			}
			s.size.Add(s.helper.size(e))
			s.slice = append(s.slice, s.helper.clone(e))
			s.lookup.Put(k, int64(p))
			dst[i] = p
			p++
		}
		s.lock.RUnlock()
	}
	swissInt64SlicePool.Put(missing)
}
