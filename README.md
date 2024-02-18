# Healthcheck tool

Simple Mattermost health check tool. This tool accepts a support packet and generates a markdown file containing the results of the investigation.

## Getting Started

1. Make a directory to run the healthcheck tool inside of:

    ```bash
    mkdir healthcheck
    cd healthcheck
    ```

2. Download the healthcheck tool for your OS

    - mac arm - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/darwin_arm64.tar.gz`
    - mac amd - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/darwin_amd64.tar.gz`
    - windows - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/windows_amd64.zip`
    - linux   - `wget https://github.com/coltoneshaw/mm-healthcheck/releases/download/v0.1.2/linux_amd64.tar.gz`

3. Unpack the tar file

    ```bash
    tar -xzf darwin_arm64.tar.gz | rm darwin_arm64.tar.gz
    ```

4. Initialize the environment.

    ```bash
    ./mmhealthcli init
    ```

    This will generate a docker compose file inside the directory you're in.

5. Move the support packet you want to do a health check on into the repo

    ```bash
    cp <packet location> .
    ```

6. Pull the docker image

    ```bash
    docker pull ghcr.io/coltoneshaw/mm-healthcheck:latest
    ```

7. Generate the PDF from the support packet.

    ```bash
    ./mmhealthcli generate -p ./mattermost_support_packet_2023-09-21-11-55.zip
    ```

8. View the pdf in your directory.

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
