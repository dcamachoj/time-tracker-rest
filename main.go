package main

import (
	"dcamachoj/time-tracker-rest/cmd"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cmd.Execute()
}
