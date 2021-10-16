package blockchain

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

// using fake if changes by IF
// need interface to get same type
func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}
func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}
func (f fakeDB) SaveBlock(hash string, data []byte) {}
func (f fakeDB) SaveChain(data []byte)              {}
func (f fakeDB) EmptyBlocks()                       {}
