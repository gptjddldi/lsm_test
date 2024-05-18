package test_test

import (
	"bufio"
	"fmt"
	"github.com/JyotinderSingh/golsm"
	"github.com/cockroachdb/pebble"
	Lsm "github.com/gptjddldi/lsm"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const benchFolder = "bench-data"

func init() {
	log.SetOutput(io.Discard)
}

func runGenerator() error {
	generatorPath := "../lsm-tree/generator1"

	cmd := exec.Command(generatorPath,
		"--puts", "10000",
		"--gets", "1000",
		//"--ranges", "10",
		//"--deletes", "200",
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
	// 명령어 실행
	if err := cmd.Run(); err != nil {
		return err
	}

	// 파일에 수정 권한 부여
	return os.Chmod("workload.txt", 0666)
}

func eraseBenchFolder() {
	err := os.RemoveAll(benchFolder)
	if err != nil {
		panic(err)
	}
}

func openFile(err error) *os.File {
	file, err := os.Open("workload.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}

	return file
}

func TestMain(m *testing.M) {
	eraseBenchFolder()

	if err := runGenerator(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running generator: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()

	os.Exit(code)
}

func BenchmarkLSMTree(b *testing.B) {
	db, err := Lsm.Open(benchFolder)
	if err != nil {
		log.Fatal(err)
	}

	file := openFile(err)

	scanner := bufio.NewScanner(file)
	b.ResetTimer()
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
			_, err := db.Get(key)
			if err != nil {
				continue
			}
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
}

func BenchmarkPebble(b *testing.B) {
	db, err := pebble.Open(benchFolder, &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}

	file := openFile(err)

	scanner := bufio.NewScanner(file)
	b.ResetTimer()
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		switch command {
		case "p":
			key := []byte(parts[1])
			value := []byte(parts[2])
			db.Set(key, value, pebble.Sync)
		case "g":
			key := []byte(parts[1])
			_, _, err := db.Get(key)
			if err != nil {
				continue
			}
		//case "r":
		//	if len(parts) != 3 {
		//		continue
		//	}
		//	start, _ := strconv.Atoi(parts[1])
		//	end, _ := strconv.Atoi(parts[2])
		//	db.Range(start, end)
		case "d":
			key := []byte(parts[1])
			db.Delete(key, pebble.Sync)
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}
}

func BenchmarkGoLSM(b *testing.B) {
	db, err := golsm.Open("bench-go-lsm", 64_000_000, true)
	if err != nil {
		log.Fatal(err)
	}

	file := openFile(err)

	scanner := bufio.NewScanner(file)
	b.ResetTimer()
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
			key := parts[1]
			value := []byte(parts[2])
			err := db.Put(key, value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error inserting key: %v\n", err)
			}
		case "g":
			if len(parts) != 2 {
				continue
			}
			key := parts[1]
			_, err := db.Get(key)
			if err != nil {
				continue
			}
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
			key := parts[1]
			err := db.Delete(key)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting key: %v\n", err)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}
	os.Remove("bench-go-lsm")
	os.Remove("bench-go-lsm_wal")
}
