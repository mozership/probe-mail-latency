# probe-mail-latency

A [Probe](https://github.com/mozership/probe) action that reads messages from a Maildir, computes the delivery latency of each one from its `Received` headers, and writes the result as a CSV file.

It was the built-in `mail-latency` action of Probe up to v1.21.0. Probe downloads it the first time a workflow uses it, and it needs Probe v1.21.0 or later, which reads the guard and the params its `action.yml` declares.

## Basic Syntax

The action reads the messages in a directory and writes the measurements as a CSV file.

```yaml
name: Delivery latency
actions:
  mail-latency: github.com/mozership/probe-mail-latency@<commit SHA>
jobs:
- name: latency
  steps:
  - name: Measure delivery latency
    uses: mail-latency
    with:
      mail_dir: "/var/mail/probe/new"
      output_dir: "./reports"
    test: res.status == 0
```

The workflow gives the action a name once, under `actions`, and its steps use the name. Probe only takes a full 40-character commit SHA there. The notes of each [release](https://github.com/mozership/probe-mail-latency/releases) start with the action and its commit to copy.

`actions` needs Probe v1.24.0 or later. With an earlier one, write the action in full in each `uses`, as `uses: github.com/mozership/probe-mail-latency@<commit SHA>`.

## Parameters

Both directories are required: one to read from, one to write to.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mail_dir` | String | Yes | Directory holding the messages to measure |
| `output_dir` | String | Yes | Directory the CSV file is written to |

`mail_dir` is read recursively, subdirectories included. A file is measured when it starts with a `Return-Path:` or `Delivered-To:` header, which is how a delivered message usually begins. Other files are left alone. Only the headers are read, up to the first empty line, and a header may be folded over several lines or not. The action fails, and writes no CSV, when a measured message has no `Date` header or no `Received` header.

Each message becomes one row of the CSV, with these columns:

- `Sent Time Offset (sec)` and `Received Time Offset (sec)`: when the message was sent and received, counted from the earliest `Date` in the directory
- `Sent Time`: the `Date` header
- `Last Received Time` and `First Received Time`: the topmost and the bottom `Received` header
- `End-to-End Latency (sec)`: from `Date` to the topmost `Received`
- `Relay Latency (sec)`: from the bottom `Received` to the topmost one
- `Return Path` and `File Path`: the envelope sender and the file name

Times are written in the local time zone of the machine running probe, as `2006-01-02 15:04:05`.

## Response Object

The result points at the file that was written.

| Field | Type | Description |
|-------|------|-------------|
| `res.output_file` | String | Path of the written CSV file, named `mail-latency.<timestamp>.csv`. A run in the same second as an earlier one writes `mail-latency.<timestamp>-2.csv` and so on instead of overwriting it |
| `res.status` | Integer | `0` on success |
| `rt.duration` | String | Time spent measuring, such as `"1.2ms"` |
| `rt.sec` | Float | Time spent measuring in seconds |
| `status` | Integer | `0` on success |

## Guard

The action connects to no host, so its `action.yml` declares `guard: [allow-host]`, and Probe runs it under `--allow-host` without `--allow-action`. It writes its CSV into `output_dir`, so it declares no `read-only`: under `--read-only`, a step that uses it is refused unless `--allow-action` names it.

## Library

The package `github.com/mozership/probe-mail-latency/maillatency` is the action itself, with `GetLatencies` for a program that computes the latencies of a Maildir on its own. It was `github.com/linyows/probe/actions/mail-latency`.

## Testing

`go test ./...` runs against the Maildir in `maillatency/testdata`. `e2e/workflow.yml` runs the action through Probe under `--allow-host`:

```sh
go build -o e2e/action/probe-mail-latency . && ./probe --allow-host localhost --allow-action shell e2e/workflow.yml
```

## Releasing

Pushing a `v*` tag builds the executables with GoReleaser and publishes them on the release. The workflow then has `probe manifest` write an `action.yml` with their URLs and SHA-256 digests, commits it to `main`, and adds the commit to the release notes. The name, the description, `guard` and `params` are those of the `action.yml` at the tag, so that file is where to change them. That commit is the one to pin: Probe reads `action.yml` at the pinned commit and refuses an executable whose digest differs.

## License

MIT
