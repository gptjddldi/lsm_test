package main

import (
	"bufio"
	"fmt"
	Lsm "github.com/gptjddldi/lsm"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runGenerator() error {
	generatorPath := "./lsm-tree/generator1"

	cmd := exec.Command(generatorPath,
		"--puts", "100000",
		//"--gets", "1000",
		//"--ranges", "10",
		"--deletes", "20",
		//"--gets-misses-ratio", "0.3",
		//"--gets-skewness", "0.2",
	)

	outfile, err := os.Create("workload.txt")
	if err != nil {
		return err
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

const dataFolder = "demo-data"

func eraseDataFolder() {
	err := os.RemoveAll(dataFolder)
	if err != nil {
		panic(err)
	}
}
func main() {
	eraseDataFolder()

	// workload generator 실행
	err := runGenerator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running generator: %v\n", err)
		os.Exit(1)
	}

	db, err := Lsm.Open(dataFolder)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("workload.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	start := time.Now()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		switch command {
		case "p":
			if len(parts) != 3 {
				continue
			}
			key := []byte(parts[1])
			value := []byte(parts[2])
			db.Insert(key, value)
		case "g":
			if len(parts) != 2 {
				continue
			}
			key := []byte(parts[1])
			db.Get(key)
		//case "r":
		//	if len(parts) != 3 {
		//		continue
		//	}
		//	start, _ := strconv.Atoi(parts[1])
		//	end, _ := strconv.Atoi(parts[2])
		//	db.Range(start, end)
		case "d":
			if len(parts) != 2 {
				continue
			}
			key := []byte(parts[1])
			db.Delete(key)
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Command execution time: %s\n", elapsed)
}
