package main

import (
	"github.com/jeamon/backend-api/cmd"
)

var (
	GitCommit string
	GitTag    string
)

// @title Demo Rest API
// @version 1.0
// @description Demonstrate simple RestFul API to perform CRUD operation on repository scan data.
// @host localhost:8080
// @BasePath /api/v1

func main() {
	cmd.Execute(GitCommit, GitTag)
}
