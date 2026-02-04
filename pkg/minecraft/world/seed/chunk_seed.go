package seed

const A = 341873128712
const B = 132897987541

func ChunkSeed(levelSeed int64, x, z int32) int64 {
	seed := levelSeed
	seed ^= int64(x) * A
	seed ^= int64(z) * B
	seed ^= seed >> 33
	return seed
}
