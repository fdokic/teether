package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
)

var dir = "./results/"

type exploit struct {
	balance   float64
	addresses string
	contract  string
}

func main() {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	x := map[string]string{}
	res := make([]exploit, len(files))

	sum := 0.0
	for i, file := range files {
		f, err := ioutil.ReadFile(dir + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(f, &x)
		flo, err := strconv.ParseFloat(x["balance"], 32)
		if err != nil {
			log.Fatal(err)
		}
		res[i] = exploit{addresses: x["addresses"], balance: flo}
		sum += flo
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].balance > res[j].balance
	})
	for _, ex := range res {
		fmt.Println(ex.balance, len(ex.addresses)%42)
	}
	fmt.Printf("total: %v in %v contracts", sum, len(res))
}
