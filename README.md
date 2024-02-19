# Healthcheck tool

Simple Mattermost health check tool. This tool accepts a support packet and generates a markdown file containing the results of the investigation.

## Getting Started

1. Ensure you have `go` installed on your OS. If you do not, follow the guide [here](https://go.dev/doc/install)

2. Download the `mmhealth` tool with `go install`

    ```bash
    go install github.com/coltoneshaw/mmhealth@latest
    ```

3. Run the initialize command. 

    This command downloads the necessary docker container.

    ```bash
    mmhealth init
    ```

## How to use

The easiest way to use this tool is to just point it at a support packet and generate the report, like below.

```bash
mmhealth generate -p ./mattermost_support_packet_2023-09-21-11-55.zip
```

This will output a `healthcheck-report.pdf` into your current directory.

## Statuses

- pass - The Check passed
- warn - The check did not pass, but it is just a warning right now.
- fail - The check failed and it should be addressed
- ignore - The check can be ignored.

## Types

- `Proactive` - A proactive measure to increase the health, reliability, or otherwise inside of Mattermost.
- `Adoption` - Better configuration of Mattermost for optimal usage and adoption.
- `Health` - Environment health checks that should be remediated if failed.



## Adding a check

1. Run `./healthcheck add`
2. Follow the interactive prompt
3. Once the check has been added to the `yaml` file you need to build the check code.
