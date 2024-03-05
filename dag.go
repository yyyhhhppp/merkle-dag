package merkledag

import "hash"


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

