# Healthcheck tool

Simple Mattermost health check tool. This tool accepts a support packet and generates a markdown file containing the results of the investigation.

## Statuses

- 🟢 - The Check passed
- 🟡 - The check did not pass, but it is just a warning right now.
- 🔴 - The check failed and it should be addressed
- `-` - The check can be ignored. 

## Types

- `Proactive` - A proactive measure to increase the health, reliability, or otherwise inside of Mattermost.
- `Adoption` - Better configuration of Mattermost for optimal usage and adoption.
- `Health` - Environment health checks that should be remediated if failed. 

## Example Output

### Configuration Checks

|                NAME                 |   TYPE    | STATUS |                  RESULT                  |                                                                                                                                           DESCRIPTION                                                                                                                                            |
|-------------------------------------|-----------|--------|------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Site URL                            | health    | 🟢     | Site URL is set                          | The siteURL is required by many functions of Mattermost. With it not set some features may not work as expected. [documentation](https://docs.mattermost.com/configure/web-server-configuration-settings.html#site-url)                                                                          |
| Extend Session Length with Activity | adoption  | 🟢     | Session Length is extended with activity | For improved end-user login session lifecycle, consider enableing `ExtendSessionLengthWithActivity` Verify with Enterprise policies if this is compatible. [documenation](https://docs.mattermost.com/configure/session-lengths-configuration-settings.html#extend-session-length-with-activity) |
| ID-Only Notifications               | proactive | 🟡     | Notification contents are set to `full`  | Setting notifications to ID-Only keeps data off Google / Apple servers and in turn your server is more secure. [documentation](https://docs.mattermost.com/configure/site-configuration-settings.html#push-notification-contents)                                                                |
| ElasticSearch Live Indexing         | health    | -      | Elasticsearch is not enabled             | Live index batch size controls how often Elasticsearch is indexed in real time. On highly active servers this needs to be increased to prevent an Elasticsearch crash. [documentation](https://docs.mattermost.com/configure/environment-configuration-settings.html#live-indexing-batch-size)   |
| Link Previews                       | adoption  | 🟡     | Link Previews are not enabled            | Link Previews are a feature that allows for a preview of a link to be displayed in the Mattermost client to improve end user experience. [documentation](https://docs.mattermost.com/configure/site-configuration-settings.html#posts-enablemessagelinkpreviews)                                 |
| IPs in SQL Data Sources             | proactive | 🟢     | Data sources do not contain IP addresses | Using IP addresses in your SQL data sources can cause issues with failovers in the event of a database failure. [documentation](https://docs.mattermost.com/configure/environment-configuration-settings.html#database-datasource)                                                               |

### Mattermost.log Checks

|             NAME             |  TYPE  | STATUS |                          RESULT                            |                                                                                                                                                                                                      DESCRIPTION                                                                                                                                                                                                      |
|------------------------------|--------|--------|--------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| context deadline exceeded    | health | 🟢     | context deadline exceeded not found                          | The context deadline exceeded error is a common error in Mattermost. It is usually caused by a slow database or a slow network connection. [documentation](https://docs.mattermost.com/install/troubleshooting.html#context-deadline-exceeded)                                                                                                                                                                        |
| i/o timeout                  | health | 🔴     | i/o timeout found                                            | Further investigation is needed.  Contact your Technical Account Manager for assistance. A common cause of this error is due to connectivity issues.  The root cause can originate from various factors. Depending on the origin of the error, we recommend verifying accessibility of the resource. In some cases, ingress/egress rules might be causing problems, or issues may arise from the nginx configuration. |
| Error while creating session | health | 🟢     | Error while creating session for user access token not found |                                              

## How to use

1. Clone the repo
2. Run `make build`
3. Run `./healthcheck process -f filename.zip` 
4. See the results. 

## Adding a check

1. Run `./healthcheck add`
2. Follow the interactive prompt
3. Once the check has been added to the `yaml` file you need to build the check code.