## Slack iBot ##


### What is slackibot? ###

Slackibot is a 24x7 support bot for slack. It can handle opsgenie/jira/e-mail alerts.

Current features:
- Trigger alerts on special words.
- Set delay periods between events.
- Auto-responses depending on time frames.
- Answer to reporters via direct messages and ask to create jira incidents on behalf.

### How to build  ###

#### Windows/Linux/Mac ####
Go to script folder and execute
./build.sh

#### Build dependencies ####
You will need Docker version 17.05 or higher since "multi-stage builds" feature is used.
Further read: https://docs.docker.com/develop/develop-images/multistage-build/
