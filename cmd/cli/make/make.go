package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func run(name string, args ...string) {
	fmt.Printf("Running command: %s %s\n", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %s\n", err)
		os.Exit(1)
	}
}

func runInDir(dir string, name string, args ...string) {
	fmt.Printf("Running command in directory %s: %s %s\n", dir, name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %s\n", err)
		os.Exit(1)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run make.go [run_server|build|test]")
		return
	}

	switch os.Args[1] {
	case "run_server":
		run("go", "run", "./cmd/server/main.go")
	default:
		fmt.Println("Unknown command. Available commands: run_server, build, test")
	}

}

func GetEnv(key ...string) string {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}
	var envs string
	for _, k := range key {
		env := os.Getenv(k)
		if env == "" {
			fmt.Printf("Environment variable %s is not set\n", k)
			os.Exit(1)
		}

		if k == "GOOSE_DNS" {
			env = fmt.Sprintf(`"%s"`, env)
		}
		envs += env + " "
	}
	fmt.Println(envs)
	return envs
}
