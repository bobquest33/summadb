package database

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

func FromLevel(bs []byte) interface{} {
	var val interface{}
	json.Unmarshal(bs, &val)
	return val
}

func ToLevel(val interface{}) []byte {
	var bs []byte
	switch v := val.(type) {
	case []byte:
		bs, _ = json.Marshal(string(v))
	default:
		bs, _ = json.Marshal(v)
	}
	return bs
}

func Random(bytes int) string {
	random := make([]byte, bytes)
	_, err := rand.Read(random)
	if err != nil {
		log.Fatal("Couldn't read random bytes: ", err)
	}
	return fmt.Sprintf("%x", random)
}

func NewRev(oldrev string) string {
	n, _ := strconv.Atoi(strings.Split(oldrev, "-")[0])

	return fmt.Sprintf("%d-%s", (n + 1), Random(5))
}

func GetRev(path string) []byte {
	docs := db.Sub(DOC_STORE)
	oldrev, err := docs.Get([]byte(path+"/_rev"), nil)
	if err != nil {
		oldrev = []byte("0-00000")
	}
	return oldrev
}

func SplitKeys(path string) []string {
	return strings.Split(path, "/")
}

func JoinKeys(keys []string) string {
	return strings.Join(keys, "/")
}

// CleanPath removes _rev, _changes, _deleted and other special things
// from the end of a path.
func CleanPath(path string) string {
	return strings.Split(path, "/_")[0]
}

// GetParent removes the last key from a path, thus returning the parent path
func GetParent(path string) string {
	keys := strings.Split(path, "/")
	parentPath := keys[:len(keys)-2]
}

func EscapeKey(e string) string {
	return url.QueryEscape(e)
}

func UnescapeKey(e string) string {
	v, _ := url.QueryUnescape(e)
	return v
}
