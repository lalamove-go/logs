test: 
	GOCACHE=off go test --cover 

test-race: 
	GOCACHE=off go test -race ./...

bench: 
	go test -bench=. -benchmem ./... 

