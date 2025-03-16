# Home Assistant Addons

- [Grafana Agent](./grafana-agent/DOCS.md)
- [Kasa Reboot](./kasa-reboot/DOCS.md)
- [Backup to S3](./backup-to-s3/DOCS.md)

## Install in HA

These steps are to be followed only once:
1. Settings > Addons > ADD-ON STORE
2. In the top right menu select `Repositories`
3. Add new repository: `https://github.com/maxosprojects/hass-addons`
4. Reload the page. This should show the newly added repository `maxosprojects` with addons
5. Consult corresponding addon docs

## Update in HA

To update the version, once it's available in DockerHub (or when `image` is commented out in `config.yaml`):
1. Settings > Addons > ADD-ON STORE
2. In the top right menu select `Check for updates` and reload the page

See addon docs for details.

## Development

### Docker

#### Docker Hub Access Token

Unless Docker CLI is logged out from Docker Hub, these steps won't be needed again.
Access Token is shown only once - right after it is generated on Docker Hub. It can't be retrieved again.
In case Docker CLI is logged out from Docker Hub, a new Access Token will have to be created.

1. Sign in to [Docker Hub](https://hub.docker.com)
2. Account Settings > Security > New Access Token
3. Create a 'Read/Write' token `rw-token` and copy it to clipboard
4. Log in Docker CLI and use the token as the password: `docker login -u maxosprojects`

#### Build image

Consult corresponding addon docs.

#### Deploy locally

To test an updated version without having to push the changes:
1. Update addon version in addon's `config.yaml`
2. Temporarily comment out `image` in addon's `config.yaml`
3. Copy the addon folder to HA's `/addons` folder
4. Settings > Addons > ADD-ON STORE
5. In the top right menu select `Check for updates` and reload the page
6. The addon should show under `Local add-ons` section
7. Consult corresponding addon docs

Once confirmed to work, don't forget to uncomment `image` in addon's `config.yaml`.
Otherwise, the image will be built in HA upon installing the final version from the repo and the image will
be included in backups.
