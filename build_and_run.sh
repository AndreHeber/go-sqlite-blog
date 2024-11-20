#!/bin/bash

cd cmd/Go-SQLite-Blog

# if go-sqlite-blog exists, remove it
if [ -f go-sqlite-blog ]; then
	rm go-sqlite-blog
fi

# build and run the application
go build -o ../go-sqlite-blog . && ./go-sqlite-blog
