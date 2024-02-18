
goose-up:
	goose postgres "host=postgres port=5432 user=postgres password=123456 dbname=postgres sslmode=disable" up
goose-down:
	goose postgres "host=postgres port=5432 user=postgres password=123456 dbname=postgres sslmode=disable" down