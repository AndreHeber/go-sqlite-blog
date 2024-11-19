#!/bin/bash

# if go-sqlite-blog exists, remove it
if [ -f go-sqlite-blog ]; then
	rm go-sqlite-blog
fi

# build the application
go build -o go-sqlite-blog ./cmd/Go-SQLite-Blog 

# run the application
./go-sqlite-blog