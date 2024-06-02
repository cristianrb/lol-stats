run:
	go run main.go -port=8081


mocks:
	mockgen -destination mocks/http_client.go -package mocks lol-stats/cristianrb/api HTTPClient
	mockgen -destination mocks/cache.go -package mocks lol-stats/cristianrb/internal Cache