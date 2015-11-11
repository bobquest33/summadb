package database

import (
	"reflect"
	"strings"

	"github.com/fiatjaf/sublevel"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	DOC_STORE = "doc-store"
	BY_SEQ    = "by-seq"
)

func GetValueAt(path string) ([]byte, error) {
	db := Open().MustSub(DOC_STORE)
	defer db.Close()

	bs, err := db.Get([]byte(path), nil)
	if err != nil {
		return []byte(nil), err
	}

	return bs, nil
}

func GetTreeAt(basepath string) (map[string]interface{}, error) {
	db := Open().MustSub(DOC_STORE)
	defer db.Close()

	baseLength := len(basepath)
	bytebasepath := []byte(basepath)
	baseTree := make(map[string]interface{})

	// fetch the "base" key if exists
	val, err := db.Get(bytebasepath, nil)
	if err == nil {
		baseTree["_val"] = FromLevel(val)
	}

	// iterate only through subkeys by adding a "/" to the prefix
	iter := db.NewIterator(util.BytesPrefix(append(bytebasepath, 0x2f)), nil)
	for iter.Next() {
		key := string(iter.Key())[baseLength:]
		val := iter.Value()

		pathKeys := splitKeys(key)

		lastKey := pathKeys[len(pathKeys)-1]
		if strings.HasPrefix(lastKey, "_") {
			/* skip _deleted and _rev */
			continue
		}

		var tree interface{}
		last := TREE
		for _, subkey := range pathKeys {
			if strings.HasPrefix(subkey, "[") {
				/* it's an array */
				var subarray = make([]interface{}, 0)

				if last == TREE {
					currentTree := tree
					tree[subkey] = subarray
				} else if last == ARRAY {
					currentTree := make([]interface{}, 0)
					currentTree = append(currentTree, subarray)
				}
				tree = currentTree
				last = ARRAY
				continue
			}

			/* it must be an object */
			var subobject map[string]interface{}
			var ok bool
			if subobject, ok = tree[subkey].(map[string]interface{}); !ok {
				/* this subobject doesn't exist in our response object yet, create it */
				subobject = make(map[string]interface{})
				tree[subkey] = subobject
			}
			tree = subobject // descend into that level
			last = TREE
		}
		// no more levels to descend into, apply the value to our response object
		tree["_val"] = FromLevel(val)
	}
	iter.Release()
	err = iter.Error()

	if err != nil {
		return make(map[string]interface{}), err
	}

	return baseTree, nil
}

const (
	TREE  = 'T'
	ARRAY = 'A'
)

func SaveValueAt(path string, bs []byte) error {
	db := Open().MustSub(DOC_STORE)
	defer db.Close()

	txn := make(transaction)
	prepare(txn, SAVE, path, bs)
	return commit(db, txn)
}

func DeleteAt(path string) error {
	db := Open().MustSub(DOC_STORE)
	defer db.Close()

	txn := make(transaction)
	prepare(txn, DELETE, path, nil)
	return commit(db, txn)
}

func SaveTreeAt(path string, tree map[string]interface{}) error {
	db, err := Open().Sub(DOC_STORE)
	if err != nil {
		return err
	}
	defer db.Close()

	txn := make(transaction)
	saveObjectAt(db, txn, path, tree)
	return commit(db, txn)
}

func saveObjectAt(db *sublevel.Sublevel, txn transaction, base string, o map[string]interface{}) error {
	for k, v := range o {
		if v == nil {
			/* setting a value to null should delete it */
			prepare(txn, DELETE, base+"/"+string(k), nil)
			continue
		} else if k[0] == 0x5f {
			/* skip secial values, i. e., those starting with "_", except for "_val" */
			if string(k) == "_val" {
				// for _val we save to the same "base" path we get
				prepare(txn, SAVE, base, ToLevel(v))
				/* also, setting _val to nil shouldn't delete the key,
				   but actually set it to null */
			}
		}

		// proceed if not new
		rv := reflect.ValueOf(v)
		kind := rv.Kind()
		if kind == reflect.Slice {
			/* arrays */
			array := make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				array[i] = rv.Index(i).Interface()
			}
			saveArray(db, txn, base+"/"+k, array)
		} else if kind != reflect.Map {
			/* normal values: strings, booleans, numbers */
			prepare(txn, SAVE, base+"/"+k, ToLevel(v))
		} else {
			/* objects, recurse */
			err := saveObjectAt(db, txn, base+"/"+k, v.(map[string]interface{}))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SaveArrayAt(path string, array []interface{}) error {
	db, err := Open().Sub(DOC_STORE)
	if err != nil {
		return err
	}
	defer db.Close()

	txn := make(transaction)
	saveArray(db, txn, path, array)
	return commit(db, txn)
}

func saveArray(db *sublevel.Sublevel, txn transaction, path string, array []interface{}) {
	/* DELETE this path so we get totally clean subpaths */
	prepare(txn, DELETE, path, nil)

	for i, o := range array {
		kind := reflect.TypeOf(o).Kind()
		arraypath := path + "/[" + intToIndexableString(i) + "]"

		if kind == reflect.Slice {
			saveArray(db, txn, arraypath, o.([]interface{}))
		} else if kind == reflect.Map {
			saveObjectAt(db, txn, arraypath, o.(map[string]interface{}))
		} else {
			txn[arrayPath(arraypath, i)] = op{
				kind: SAVE,
				val:  ToLevel(o),
			}
		}
	}
}
