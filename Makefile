spotify-seeder-run:
	@go run cmd/main.go spotify seeder run --db "./db/decibel.db" --dir "./data/Spotify Extended Streaming History" --verbose

spotify-stats-top-artists:
	@go run cmd/main.go spotify stats top-artists --db "./db/decibel.db" --verbose

spotify-stats-top-tracks:
	@go run cmd/main.go spotify stats top-tracks --db "./db/decibel.db" --verbose

spotify-stats-top-albums:
	@go run cmd/main.go spotify stats top-albums --db "./db/decibel.db" --verbose

spotify-stats-most-skipped-tracks:
	@go run cmd/main.go spotify stats most-skipped-tracks --db "./db/decibel.db" --verbose
