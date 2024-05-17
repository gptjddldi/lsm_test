# Go Project with LSM Tree Workload Builder

## 준비

- GSL library
    - On macOS: `brew install gsl`
- C++ compiler
- Go

## 빌드 및 실행

1. Clone the repository:
    ```sh
    git clone https://github.com/gptjddldi/lsm_test.git
    cd lsm_test
    ```

2. Build the C++ generator:
    ```sh
    gcc -I/opt/homebrew/opt/gsl/include -L/opt/homebrew/opt/gsl/lib -lgsl -lgslcblas -lm generator/generator.c -o generator1
    ```
- `-I/opt/homebrew/opt/gsl/include` tells the compiler where to find the GSL headers.
- `-L/opt/homebrew/opt/gsl/lib` tells the linker where to find the GSL libraries.
- `-lgsl -lgslcblas -lm` links the GSL and GSL CBLAS libraries as well as the math library.

4. Run the Go program:
    ```sh
    go run main.go
    ```

## 결과

``` shell
+---------- CS 265 ----------+
|        WORKLOAD INFO       |
+----------------------------+
| initial-seed: 13141
| puts-total: 1000000
| gets-total: 0
| get-skewness: 0.0000
| ranges: 0
| range distribution: uniform
| deletes: 20000
| gets-misses-ratio: 0.5000
+----------------------------+
Command execution time: 756.765708ms
```

## 참고

workload.txt 예시
```text
p -1109180 699692587
p 1213834231 -226769626
d 1213834231
p 994957275 2082945813
p -182635081 -1098200189
d -182635081
d 1213834231
```

`p`: put (2 arguments)
`d`: delete (1 argument)
`g`: get (1 argument)
`r`: range (2 arguments)

