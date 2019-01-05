package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/thetatoken/ukulele/blockchain"
	"github.com/thetatoken/ukulele/core"
	"github.com/thetatoken/ukulele/rlp"
	"github.com/thetatoken/ukulele/store/database/backend"
	"github.com/thetatoken/ukulele/store/trie"
)

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: inspect_data -config=<path_to_config_home> -key=<key>")
}

func main() {
	configPathPtr := flag.String("config", "", "path to ukuele config home")
	keyPtr := flag.String("key", "", "db key")
	flag.Parse()
	configPath := *configPathPtr
	key := *keyPtr

	mainDBPath := path.Join(configPath, "db", "main")
	refDBPath := path.Join(configPath, "db", "ref")
	db, err := backend.NewLDBDatabase(mainDBPath, refDBPath, 256, 0)
	handleError(err)
	// db, _ := backend.NewAerospikeDatabase()
	// db, _ := backend.NewMongoDatabase()

	k := str2hex2bytes(key)
	value, err := db.Get(k)
	handleError(err)

	node, err := trie.DecodeNode(k, value, 0)
	if err == nil {
		fmt.Printf("%v\n", node)
	} else {
		if strings.HasPrefix(err.Error(), "invalid number of list elements") {
			block := core.ExtendedBlock{}
			err = rlp.DecodeBytes(value, &block)
			if err == nil {
				fmt.Printf("%v\n", block)
			} else {
				blockByHeightIndexEntry := blockchain.BlockByHeightIndexEntry{}
				err = rlp.DecodeBytes(value, &blockByHeightIndexEntry)
				if err == nil {
					fmt.Printf("%v\n", blockByHeightIndexEntry)
				} else {
					handleError(err)
				}
			}
		} else {
			handleError(err)
		}
	}

	os.Exit(0)
}

func str2hex2bytes(str string) []byte {
	var bytes []byte
	if strings.HasPrefix(str, "0x") {
		str = strings.TrimPrefix(str, "0x")
	}
	for {
		if len(str) <= 0 {
			break
		}
		s := str[:2]
		i, _ := strconv.ParseUint(s, 16, 64)
		bytes = append(bytes, byte(i))
		str = str[2:]
	}
	return bytes
}
