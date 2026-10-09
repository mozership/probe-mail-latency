#!/bin/sh
# Print the action.yml for a release: where Probe downloads each executable
# from, and the SHA-256 digest it must have.
#
#   scripts/action-yml.sh <tag> <checksums.txt>
#
# It declares allow-host, as the action connects to no host, and not
# read-only, as it writes its CSV into output_dir.
set -eu

tag=$1
checksums=$2
repo=${GITHUB_REPOSITORY:-mozership/probe-mail-latency}

cat <<YAML
name: mail-latency
description: Measure mail delivery latency from a Maildir
guard: [allow-host]
params: [mail_dir, output_dir]
runs:
  using: binary
  url: https://github.com/${repo}/releases/download/${tag}/probe-mail-latency_{os}_{arch}
  checksums:
YAML
# A line of checksums.txt is "<digest>  probe-mail-latency_<os>_<arch>".
awk '$2 ~ /^probe-mail-latency_/ { sub(/^probe-mail-latency_/, "", $2); printf "    %s: \"%s\"\n", $2, $1 }' "$checksums" | sort
