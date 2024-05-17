package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// C++ 프로그램 실행
func runGenerator() error {
	generatorPath := "./lsm-tree/generator1"

	// 실행 권한 부여
	if err := os.Chmod(generatorPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	cmd := exec.Command(generatorPath,
		"--puts", "100000",
		"--gets", "1000",
		"--ranges", "10",
		"--deletes", "20",
		"--gets-misses-ratio", "0.3",
		"--gets-skewness", "0.2")

	// 생성된 파일로 출력 리디렉션
	outfile, err := os.Create("workload.txt")
	if err != nil {
		return err
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr

	// 명령어 실행
	return cmd.Run()
}

// 생성된 파일 읽기
func readWorkload(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	// C++ 프로그램 실행
	err := runGenerator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running generator: %v\n", err)
		os.Exit(1)
	}

	// 생성된 파일 읽기
	err = readWorkload("workload.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading workload: %v\n", err)
		os.Exit(1)
	}
}
