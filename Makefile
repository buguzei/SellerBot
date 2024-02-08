
goose-up:
	goose postgres "host=localhost port=5432 user=buguzei password=123456 dbname=postgres sslmode=disable" up
goose-down:
	goose postgres "host=localhost port=5432 user=buguzei password=123456 dbname=postgres sslmode=disable" down