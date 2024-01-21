# Home Assistant Addon: Grafana Agent

See [README](https://github.com/maxosprojects/hass-addons)

## Use in HA

### Install in HA

1. See [README](https://github.com/maxosprojects/hass-addons)
2. Set the [required configuration options](#required-configuration-options) (and `Save`)

### Required configuration options

1. `ha-integration` Grafana Cloud Token from https://<YOUR GRAFANA CLOUD HOST>/a/grafana-auth-app (and `Save`)
2. From https://noteb5.grafana.net/connections/add-new-connection/linux-node > Run Grafana Agent > Use an existing API token > GCLOUD_HOSTED_LOGS_ID

## Development

See [README](https://github.com/maxosprojects/hass-addons)

### Docker

#### Build image

Image tag is formed after the version of HA where it was confirmed to work and the version of Grafana Agent.

```shell
docker build --tag=maxosprojects/hass-grafana-agent:ha-2024.1.3-ga-v0.39.0 .
docker push maxosprojects/hass-grafana-agent:ha-2024.1.3-ga-v0.39.0
```

#### Deploy locally

1. Click on the addon in the `Local` repo > Build > disable Protection mode
2. Set the [required configuration options](#required-configuration-options) (and `Save`)
3. Info > Start
