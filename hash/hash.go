package hash

import (
	"fmt"
	"hash/crc32"
	"sync"
)

var	crc32q *crc32.Table
var once sync.Once

func GetCrc32Table() *crc32.Table {
	once.Do(func () {
		crc32q = crc32.MakeTable(0xD5828281)
	})

	return crc32q
}
func Crc32Hash(s string) string {
	table := GetCrc32Table()
	hash := crc32.Checksum([]byte(s), table)
	
	return fmt.Sprintf("%x", hash)
}