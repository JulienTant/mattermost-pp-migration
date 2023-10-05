# Mattermost Profile Picture Migration tool

Currently, Mattermost does not allow admin to migrate profile pictures from one server to another. This tool is a workaround to migrate profile pictures from one server to another.

## Requirements

- Two servers (source and destination)
- Must have export users from source and imported them to destination
- The match will be done by username, so make sure nobody change their username during the migration
- Must enable personal access tokens while running this tool (https://docs.mattermost.com/integrations/cloud-personal-access-tokens.html)
- Must disable rate limiting while running this tool (https://docs.mattermost.com/configure/rate-limiting-configuration-settings.html)

## Usage 

```
Usage:
  mattermost-pp-migration [flags]

Flags:
      --dst-access-token string   Destination server access token (required|env DST_ACCESS_TOKEN)
      --dst-server-url string     Destination server URL (required|env DST_SERVER_URL)
      --src-access-token string   Source server access token (required|env SRC_ACCESS_TOKEN)
      --src-server-url string     Source server URL (required|env SRC_SERVER_URL)
```