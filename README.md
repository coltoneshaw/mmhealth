# Healthcheck tool

Simple Mattermost health check tool. This tool accepts a support packet and generates a markdown file containing the results of the investigation.

## Statuses

- pass - The Check passed
- warn - The check did not pass, but it is just a warning right now.
- fail - The check failed and it should be addressed
- ignore - The check can be ignored.

## Types

- `Proactive` - A proactive measure to increase the health, reliability, or otherwise inside of Mattermost.
- `Adoption` - Better configuration of Mattermost for optimal usage and adoption.
- `Health` - Environment health checks that should be remediated if failed.

## How to use

1. Clone the repo
2. Run `make build`
3. Run `make buildDockerPdf`
4. Run `./healthcheck process -f filename.zip`. This outputs a `report.md` file that is the raw markdown of the report.
5. Add or make any changes to the report file now, before generating the pdf. 
6. Run `./healthcheck pdf -f report.md` to convert it to a pdf report for publishing.

## Adding a check

1. Run `./healthcheck add`
2. Follow the interactive prompt
3. Once the check has been added to the `yaml` file you need to build the check code.
