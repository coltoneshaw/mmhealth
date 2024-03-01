# Healthcheck tool

Simple Mattermost health check tool. This tool accepts a support packet and generates a markdown file containing the results of the investigation.

## Getting Started

1. Ensure you have `go` installed on your OS. If you do not, follow the guide [here](https://go.dev/doc/install)

2. Download the `mmhealth` tool with `go install`

    ```bash
    go install github.com/coltoneshaw/mmhealth@latest
    ```

3. Ensure you have Docker setup and running on your computer.

## How to use

The easiest way to use this tool is to just point it at a support packet and generate the report, like below.

```bash
mmhealth generate -p ./mattermost_support_packet_2023-09-21-11-55.zip
```

This will output a `healthcheck-report.pdf` into your current directory.

You can modify the report name by passing `--outputName` or `-o`.

## Legend

### Statuses

- `pass` - The Check passed
- `warn` - The check did not pass, but it is just a warning right now. All adoption / proactive checks that fail are a warn.
- `fail` - The check failed and it should be addressed
- `ignore` - The check can be ignored.

### Severity

- `urgent` - Highest priority check and should be addressed immediately.
- `high` - The issue should be scheduled in the next change window.
- `medium` - Can be addressed later, also informational.
- `low` - Can be ignored, but fixing could provide benefits.

### Types

- `Proactive` - A proactive measure to increase the health, reliability, or otherwise inside of Mattermost.
- `Adoption` - Better configuration of Mattermost for optimal usage and adoption.
- `Health` - Environment health checks that should be remediated if failed.

## Contributing

### Making Docker Release

When you need to adjust anything inside of the `./docker/dockerfile`, you'll need to manually adjust the `DOCKER_VERSION` to trigger a rebuild. 

### Making mmhealth Release

Just release by adding a tag on github and publishing, everything else is handled.

### Updating / adding the checks

The checks are provided from the Findings Report, parsed via CSV and updated in the `./mmhealth/files/checks.yaml` file.

To parse a csv:

1. Add a check on the Findings report
2. Download the finding report from Google as a csv
3. Copy to the root dir here and name it `healthcheck.csv`.
4. Run `make build`
5. Run `./bin/mmhealth parsecsv`
6. Check the diff of the new updates
7. Run `make test` to confirm no checks have broke
8. Update or add the checks broken checks
