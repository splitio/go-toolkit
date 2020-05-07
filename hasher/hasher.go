package hasher

const (
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
	r1 uint32 = 15
	r2 uint32 = 13
	m  uint32 = 5
	n  uint32 = 0xe6546b64
)

// Hasher interface
type Hasher interface {
	Hash(data []byte) uint32
}

// Murmur32Hasher is a hash function implementing the murmur3 32 bits algorithm
type Murmur32Hasher struct {
	seed uint32
}

// Hash returns the murmur3 (32 bits) hash of a byte slice.
func (h *Murmur32Hasher) Hash(data []byte) uint32 {
	hash := h.seed
	nblocks := len(data) / 4
	for i := 0; i < nblocks; i++ {
		k := uint32(data[i*4+0])<<0 | uint32(data[i*4+1])<<8 | uint32(data[i*4+2])<<16 | uint32(data[i*4+3])<<24
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		hash ^= k
		hash = ((hash<<r2)|(hash>>(32-r2)))*m + n
	}

	l := nblocks * 4
	k1 := uint32(0)
	switch len(data) & 3 {
	case 3:
		k1 ^= uint32(data[l+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(data[l+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(data[l+0])
		k1 *= c1
		k1 = (k1 << r1) | (k1 >> (32 - r1))
		k1 *= c2
		hash ^= k1
	}

	hash ^= uint32(len(data))
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16
	return hash
}

// NewMurmur332Hasher returns a new instance of the Murmur32Hasher
func NewMurmur332Hasher(seed uint32) *Murmur32Hasher {
	return &Murmur32Hasher{seed: seed}
}
