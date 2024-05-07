# Home Assistant Addon: Kasa Reboot

This addon accepts reboot command on stdin, pulls all TPLink devices from `/homeassistant/.storage/core.config_entries` 
on command and reboots all that are not disabled in HA.

Also, see [README](https://github.com/maxosprojects/hass-addons)

## Use in HA

### Install in HA

1. See [README](https://github.com/maxosprojects/hass-addons)
2. Enable `Watchdog`

### Configure in HA

1. Add a Service as an Action in an Automation/Script
2. Type `stdin` to find the `Addon stdin` service
3. In the `Addon` dropdown select `Kasa Reboot` addon
4. Switch to YAML mode, note the addon name (e.g. `c33755f6_kasa_reboot`) and update service call to this 
  (keeping the original addon name):
  ```yaml
  service: hassio.addon_stdin
  data:
    addon: c33755f6_kasa_reboot
    input:
      cmd: reboot
  ```
5. Call the created Automation/Script to verify and set it up to run every day

## Development

See [README](https://github.com/maxosprojects/hass-addons)

### Docker

#### Build image

Image tag is formed after the version in [config.yaml](./config.yaml).

```shell
docker build --tag=maxosprojects/hass-kasa-reboot:1.0.0 .
docker push maxosprojects/hass-kasa-reboot:1.0.0
```

#### Deploy locally

1. See [README](https://github.com/maxosprojects/hass-addons)
2. Click on the addon in the `Local` repo > Build > Start
3. Enable `Watchdog`
