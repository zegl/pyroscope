package symdb

import (
	"sync"
	"sync/atomic"

	"golang.org/x/exp/slices"
)

type stringsTable struct {
	// original stream:
	//  b a c
	//  0 1 2
	//
	// ingested:
	//  [ a b c ] (slice)
	//  [ 1 0 2 ] (ptrs)

	slice []string
	ptrs  []int64

	size atomic.Uint64

	lock sync.RWMutex
}

func (s *stringsTable) init() {
	s.slice = make([]string, 0)
	s.ptrs = make([]int64, 0)
}

func (s *stringsTable) ingest(elems []string, rw *rewriter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, elem := range elems {
		idx, ok := slices.BinarySearch(s.slice, elem)
		if ok {
			continue
		}

		// TODO(bryan) this could probably be optimized
		s.slice = append(s.slice[:idx], append([]string{elem}, s.slice[idx:]...)...)
		s.ptrs = append(s.ptrs[:idx], append([]int64{int64(i)}, s.ptrs[idx:]...)...)
		s.size.Add(uint64(len(elem)))

		// pos, ok := s.lookup[elem]
		// if ok {
		// 	rewritingMap[int64(i)] = pos
		// 	continue
		// }

		// s.slice = append(s.slice, elem)
		// s.lookup[elem] = int64(len(s.slice) - 1)
		// rewritingMap[int64(i)] = int64(len(s.slice) - 1)
		// s.size.Add(uint64(len(elem)))
	}

	rw.strings = make(stringConversionTable, len(s.ptrs)+1)
	copy(rw.strings, s.ptrs)

	// s.addToRewriter(rw, rewritingMap)
}

func (s *stringsTable) append(dst []uint32, elems []string) {
	s.lock.Lock()
	for i, elem := range elems {
		idx, ok := slices.BinarySearch(s.slice, elem)
		if ok {
			dst[i] = uint32(s.ptrs[idx])
			continue
		}

		// TODO(bryan) this could probably be optimized
		s.slice = append(s.slice[:idx], append([]string{elem}, s.slice[idx:]...)...)
		s.ptrs = append(s.ptrs[:idx], append([]int64{int64(i)}, s.ptrs[idx:]...)...)
		s.size.Add(uint64(len(elem)))

		// x, ok := s.lookup[elem]
		// if ok {
		// 	dst[i] = uint32(x)
		// 	continue
		// }

		// s.size.Add(uint64(len(elem)))
		// s.slice = append(s.slice, elem)
		// s.lookup[elem] = int64(len(s.slice) - 1)
	}
	s.lock.Unlock()
}

func (s *stringsTable) addToRewriter(r *rewriter, m idConversionTable) {
	var maxID int64
	for id := range m {
		if id > maxID {
			maxID = id
		}
	}
	r.strings = make(stringConversionTable, maxID+1)

	for x, y := range m {
		r.strings[x] = y
	}
}

func (s *stringsTable) Size() uint64 {
	return s.size.Load()
}

func (s *stringsTable) sliceHeaderCopy() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.slice
}
