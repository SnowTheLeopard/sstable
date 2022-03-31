// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/SnowTheLeopard/sstable/sstable"
// )

// const tempFile = "temp_test_blocks.sst"

// // some sort of demo
// func main() {
// 	testBlockWrite()
// 	testBlockRead()
// 	cleanup()
// }

// func testBlockWrite() {
// 	data := make(map[string]string)
// 	data["zero"] = "lattency"
// 	data["middle"] = "name"
// 	data["at"] = "home"

// 	file, err := os.Create(tempFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	sstable.WriteMap(data, file)
// }

// func testBlockRead() {
// 	rf, _ := os.Open(tempFile)
// 	defer rf.Close()

// 	t, err := sstable.NewTable(rf)
// 	if err != nil {
// 		panic(err)
// 	}

// 	k := "middle"
// 	v, err := t.Search(k)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("%v:%v", k, v)
// }

// func cleanup() {
// 	err := os.Remove(tempFile)
// 	if err != nil {
// 		panic(err)
// 	}
// }
package sstable
