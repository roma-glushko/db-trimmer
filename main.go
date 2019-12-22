package main

import "db-trimmer/internal/poc"

func main() {
	//blockingPoc := poc.NewBlockingPoc("mysql", "root:root@tcp(127.0.0.1:3306)/db-trimmer-sample", 1000)
	//blockingPoc.Execute()

	nonblockingPoc := poc.NewNonBlockingPoc("mysql", "root:root@tcp(127.0.0.1:3306)/db-trimmer-sample", 1000, 4, 2)
	nonblockingPoc.Execute()
}
