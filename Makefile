spotify-seeder-run:
	@go run cmd/main.go spotify seeder run --db "./db/decibel.db" --dir "./data/Spotify Extended Streaming History" --verbose

spotify-stats-artists:
	@go run cmd/main.go spotify stats artists --db "./db/decibel.db" --verbose

spotify-stats-tracks:
	@go run cmd/main.go spotify stats tracks --db "./db/decibel.db" --verbose
