package main

import "github.com/k1ender/task-master-go/internal/config"

func main() {
	cfg := config.MustInit(".env")
	
	_ = cfg
}
