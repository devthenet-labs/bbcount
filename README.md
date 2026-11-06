# bbcount

Beat Battle's pattern counter: a tiny, standard-library-only Go command-line tool. It reads one Beat Battle pattern
JSON (the file named, else standard input), prints how many steps each track plays and, with `-wav`, renders the
pattern to a mono 16-bit 44100 Hz WAV file. It is the CLI fixture for the agent evaluator of
[maul](https://github.com/devthenet-labs/maul): a project with a `surfaces.cli`, no preview and nothing to serve.

## Usage

```sh
go build -o bbcount .
./bbcount pattern.json
./bbcount -json < pattern.json
./bbcount -wav beat.wav pattern.json
```

```text
kick 4
snare 2
hat 8
total 14
```

`-json` prints `{"kick":4,"snare":2,"hat":8,"total":14}` instead. `-csv` prints a `track,steps` header followed by
one `kick,4` / `snare,2` / `hat,8` / `total,14` line per row instead. With `-wav out.wav` the counts are followed by
`wrote out.wav (2000 ms)`, the rendered length in whole milliseconds (sixteen sixteenth notes at the pattern's bpm).
The WAV is deterministic: the same pattern always gives the same bytes. Swing delays every odd step by that percent of
a step.

The pattern is Beat Battle's shared JSON: `name`, `author`, integer `bpm` 60-200, integer `swing` 0-60 and `tracks`
with exactly `kick`, `snare` and `hat`, 16 booleans each; `id`, `votes` and `createdAt` are accepted and ignored,
any other field is rejected. Exit codes: 0 success, 1 an unreadable or invalid pattern (a message on stderr) or a
write failure, 2 a usage error.

## Develop

Use Go 1.26.6.

```sh
gofmt -l .
go vet ./...
go test -race ./...
```

`.maul/Dockerfile` is the image maul's coding agent works in (and where the evaluator builds the tool): maul's agent
base plus the pinned Go toolchain, offline. `.maul/agent.yaml` declares its immutable tag. The root `Dockerfile` and
the `.github` workflows are the scaffold `maul init app` writes; see `.github/actions/publish/README.md`.
