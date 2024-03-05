package merkledag

import "hash"


type File interface {
	Node

	Bytes() []byte
}
type KVStore interface {
	Has(key []byte) (bool, error)
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
}
func Add(store KVStore, node Node, h crypto.Hash) []byte {
    // 将Node中的数据保存在KVStore中
    store.Write(node.Bytes())

    // 计算Merkle Root
    hf := h.New()
    hf.Write(store.Read()) // 从KVStore中读取数据进行哈希计算
    hashed := hf.Sum(nil)
    fmt.Printf("Merkle Root: %x\n", hashed)

    return hashed
}

