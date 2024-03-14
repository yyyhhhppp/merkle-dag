
package merkledag

import (
	"encoding/json"
	"hash"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	switch node.Type() {
	case DIR:
		file := node.(File)
		tmp := StoreFile(store, file, h)
		jsonMS, _ := json.Marshal(tmp)
		hash := CHash(jsonMS, h)
		return hash
	case FILE:
		dir := node.(Dir)
		tmp := StoreDir(store, dir, h)
		jsonMS, _ := json.Marshal(tmp)
		hash := CHash(jsonMS, h)
		return hash
	}
	panic("unknown node")
}

func CHash(data []byte, h hash.Hash) []byte {
	h.Reset()
	hash := h.Sum(data)
	h.Reset()
	return hash
}

func StoreFile(store KVStore, file File, h hash.Hash) *Object {
	data := file.Bytes()
	blob := Object{Data: data, Links: nil}
	jsonMS, _ := json.Marshal(blob)
	hash := CHash(jsonMS, h)
	store.Put(hash, data)
	return &blob
}

func StoreDir(store KVStore, dir Dir, h hash.Hash) *Object {
	it := dir.It()
	treeObject := &Object{}
	for it.Next() {
		n := it.Node()
		switch n.Type() {
		case FILE:
			file := n.(File)
			tmp := StoreFile(store, file, h)
			jsonMS, _ := json.Marshal(tmp)
			hash := CHash(jsonMS, h)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: hash,
				Size: int(file.Size()),
				Name: file.Name(),
			})
			typeName := "link"
			if tmp.Links == nil {
				typeName = "blob"
			}
			treeObject.Data = append(treeObject.Data, []byte(typeName)...)
		case DIR:
			dir := n.(Dir)
			tmp := StoreDir(store, dir, h)
			jsonMS, _ := json.MS(tmp)
			hash := CHash(jsonMS, h)
			treeObject.Links = append(treeObject.Links, Link{
				Hash: hash,
				Size: int(dir.Size()),
				Name: dir.Name(),
			})
			typeName := "tree"
			treeObject.Data = append(treeObject.Data, []byte(typeName)...)
		}
	}
	jsonMS, _ := json.Marshal(treeObject)
	hash := CHash(jsonMS, h)
	store.Put(hash, jsonMS)
	return treeObject
}
