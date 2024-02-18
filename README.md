# Healthcheck tool

Simple Mattermost health check tool. This tool accepts a support packet and generates a markdown file containing the results of the investigation.

## Getting Started

1. Download the healthcheck tool for your OS

    - mac arm - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/darwin_arm64.tar.gz`
    - mac amd - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/darwin_amd64.tar.gz`
    - windows - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/windows_amd64.zip`
    - linux   - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/linux_amd64.tar.gz`


2. Navigate to the directory you cloned.

    ```bash
    cd mm-healthcheck
    ```

3. Move the support packet you want to do a health check on into the repo

    ```bash
    cp <packet location> .
    ```

4. Pull the docker image

    ```bash
    docker pull ghcr.io/coltoneshaw/mm-healthcheck:latest
    ```

5. Run the generate command. Replace `packetname` with the packet you're wanting to run against.

    This will output a `report.md` file within the directory.

    ```bash
    docker compose run mm-healthcheck process -f ./packetname
    ```

6. Generate the PDF.

    If you have a `report.md` in the root directory you do not have to do anything extra, it uses this file by default. 

    ```bash
    docker compose run mm-healthcheck pdf
    ```

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
