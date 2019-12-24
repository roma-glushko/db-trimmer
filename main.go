package main

import "db-trimmer/internal/poc"

func main() {
	nonblockingPoc := poc.NewNonBlockingPoc("mysql", "root:root@tcp(127.0.0.1:3306)/db-trimmer-sample", 300, 2, 4)
	nonblockingPoc.Execute()
}
