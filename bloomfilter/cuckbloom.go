package bloomx

import (
	cuckoo "github.com/seiflotfy/cuckoofilter"
)

type CuckBloom struct {
	Bloom *cuckoo.Filter
	size  uint32
}

func NewCuckBloom(_cap uint32) *CuckBloom {
	return &CuckBloom{
		Bloom: cuckoo.NewFilter(uint(_cap)),
		size:  0,
	}
}

func (b *CuckBloom) Insert(bytes []byte) bool {
	ok := b.Bloom.InsertUnique(bytes)
	b.size++
	return ok
}

func (b *CuckBloom) Delete(bytes []byte) bool {
	ok := b.Bloom.Delete(bytes)
	b.size--
	return ok
}

func (b *CuckBloom) Count() uint32 {
	cnt := b.Bloom.Count()
	return uint32(cnt)
}
