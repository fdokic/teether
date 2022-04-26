package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var teether_path = "/home/fdokic/teether/"
var contract_table_path = []string{"contracts.csv"}
var result_file_dir = "results/"
var contract_offset = 0

func main() {
	fmt.Println("process start")
	var numCPU = runtime.NumCPU()
	fmt.Printf("cpus: %v \n", numCPU)
	for _, p := range contract_table_path {
		analyzeContracts(teether_path + p)
	}
}

func checkContract(i int, c chan int, str []string, wg *sync.WaitGroup) {
	defer wg.Done()
	code := str[2]
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if err := exec.CommandContext(ctx, "python3", teether_path+"bin/gen_exploit.py", code, "0x1234", "0x1000", "+1000", teether_path+result_file_dir+fmt.Sprint(i)+".txt", str[0], str[1]).Run(); err != nil {
		fmt.Println("timeout: " + fmt.Sprint(i) + err.Error())
		<-c
		return
	}
	/*
		cmd := exec.Command("python", teether_path+"bin/gen_exploit.py", code, "0x1234", "0x1000", "+1000", result_file_dir+fmt.Sprint(i)+".txt")
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("cmd err: " + fmt.Sprint(i) + " contract " + err.Error())
		}
		//fmt.Println(string(out))
		outStr := string(out)
		if strings.Contains(outStr, "eth.sendTransaction") {
			s := str[0] + "\n\n" + str[2] + "\n\n" + str[1] + "\n\n" + outStr
			err := os.WriteFile(result_file_dir+fmt.Sprint(i)+".txt", []byte(s), 0644)
			if err != nil {
				panic(err)
			}
		}*/
	fmt.Printf("checked %v \n", i)
	<-c
}

func analyzeContracts(contract_table_path string) {
	f, err := os.Open(contract_table_path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	c := make(chan int, runtime.NumCPU())
	var wg sync.WaitGroup
	for i := 0; i < contract_offset; i++ {
		csvReader.Read()
	}

	for i := contract_offset; ; i++ {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// do something with read line
		c <- i
		wg.Add(1)
		go checkContract(i, c, rec, &wg)

	}

	wg.Wait()
}
