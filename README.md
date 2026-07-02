# fakedata

> **This repository is archived.** FakeData has moved to a new home with Club Matto (see below).

## Why

FakeData has found a new home with [Club Matto](https://matto.club) and it's been substantially rebuilt in the process.

The new home is [clubmatto/vetrina](https://github.com/clubmatto/vetrina), a monorepo that includes fakedata/ alongside other tools. There's also a [landing page](https://matto.club/vetrina/fakedata/) with an interactive generator browser.

## What's new

The new FakeData is a ground-up rewrite with a dramatically improved UX:

- Column spec syntax — name:generator[:options] for clean, declarative column definitions (e.g. login:email, count:int:10,20)
- Country-specific generators — append a 2-letter country code to filter data: phone_it, city_de, state_us
- Ordered column pairs — start_date/end_date, created_at/updated_at pairs are automatically ordered
- Non-interactive generator browser — fakedata -g <name> shows details and sample output; fakedata --generators lists everything
- Streaming mode — fakedata --stream generates rows indefinitely

## Migrating

```bash
# Old install
go install github.com/lucapette/fakedata@latest

# New install
go install matto.club/vetrina/fakedata@latest

# Or via Homebrew
brew install clubmatto/tap/fakedata
```

The CLI interface has changed. Key differences:

| Old | New |
|-----|-----|
| `fakedata --separator=,` | `fakedata -s ","` |
| `fakedata --format=ndjson` | `fakedata --format ndjson` |
| `fakedata --limit 5` | `fakedata -n 5` |
| `fakedata -c` | `fakedata --generators` (all generators shown inline) |
| `fakedata column=name:generator` | `fakedata name:generator` |

Check the full README (https://github.com/clubmatto/vetrina/blob/main/fakedata/README.md) for all the details.

Thank you to everyone who contributed, filed issues, and used FakeData over the years!

## License

MIT License Copyright (c) [2022] Luca Pette — now maintained by Club Matto.
