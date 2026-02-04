// I got this code from https://github.com/dgryski/go-xoroshiro/blob/master/xoro.go, and i made some changes
package xoroshiro

type State [2]uint64

const (
	SPLIT_MIX_GAMMA uint64 = 0x9E3779B97F4A7C15
	SPLIT_MIX_M1    uint64 = 0xBF58476D1CE4E5B9
	SPLIT_MIX_M2    uint64 = 0x94D049BB133111EB
)
const (
	CHUNK_X_MULTI uint64 = 341873128712
	CHUNK_Z_MULTI uint64 = 132897987541
)

// New returns a new RNG
func New(seed *uint64) State {
	var s State
	s[0] = splitMix64(seed)
	s[1] = splitMix64(seed)

	return s
}

func (s *State) Next() uint64 {
	s0 := s[0]
	s1 := s[1]

	result := rotl(s0+s1, 17) + s0

	s1 ^= s0
	s[0] = rotl(s0, 49) ^ s1 ^ (s1 << 21)
	s[1] = rotl(s1, 28)

	return result
}

func Fork(parent, salt uint64) State {
	new := parent + salt

	return New(&new)
}

func At(parent uint64, x, z int64) State {

	parent ^= uint64(int64(x)) * CHUNK_X_MULTI
	parent ^= uint64(int64(z)) * CHUNK_Z_MULTI

	return New(&parent)
}

func splitMix64(levelSeed *uint64) uint64 {
	*levelSeed += SPLIT_MIX_GAMMA

	z := *levelSeed
	z = (z ^ (z >> 30)) * SPLIT_MIX_M1
	z = (z ^ (z >> 27)) * SPLIT_MIX_M2

	return z ^ (z >> 31)
}

func rotl(x uint64, k uint) uint64 {
	return (x << k) | (x >> (64 - k))
}
