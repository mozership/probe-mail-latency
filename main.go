// Command probe-mail-latency serves a Probe action that measures how long
// mail took to be delivered, from the messages in a Maildir. Probe starts it
// for a step that uses github.com/mozership/probe-mail-latency@<commit>.
package main

import maillatency "github.com/mozership/probe-mail-latency/maillatency"

// version is set at build time.
var version = "dev"

func main() {
	maillatency.Serve()
}
