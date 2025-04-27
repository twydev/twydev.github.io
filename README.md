# Theme Source

From [Minimal Mistakes Jekyll theme](https://github.com/mmistakes/minimal-mistakes).

# Sync Journal Entries

1. set up `scripts/.env`
2. from `scripts/` execute the command `go run .`
3. existing journal entries will be synced to existing posts
4. new journal entries will be added as new posts

assumptions
- the journal `/knowledge/books` directory will sync with `_posts/notes` directory
- journal entries with prefix `datetime-pub-xxx` will be published
