# Mattermost Profile Picture Migration tool

Currently, Mattermost does not allow administrators to migrate profile pictures from one server to another. This tool is a workaround to migrate profile pictures from one server to another.

## Requirements

- Two servers (source and destination)
- Must have exported users from the source server and import them to the destination server
- The match will be done by username, so make sure nobody changes their username during the migration
- Must be a sysadmin on the two servers
- Must enable personal access tokens while running this tool (https://docs.mattermost.com/integrations/cloud-personal-access-tokens.html) - and create one on each server.
- Must disable rate limiting while running this tool (https://docs.mattermost.com/configure/rate-limiting-configuration-settings.html)

## Usage 

You can either use environment variables or flags to pass the arguments. The command will match users on both servers and update the users on the destination server with the profile picture of the source server.

```
Usage:
  mattermost-pp-migration [flags]

Flags:
      --dst-access-token string   Destination server access token (required|env DST_ACCESS_TOKEN)
      --dst-server-url string     Destination server URL (required|env DST_SERVER_URL)
      --src-access-token string   Source server access token (required|env SRC_ACCESS_TOKEN)
      --src-server-url string     Source server URL (required|env SRC_SERVER_URL)
```
