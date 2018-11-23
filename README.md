# pa-switch-profile

Tool for quick switching profiles of physical card.

## Usage

```
pa-switch-profile <card> <profile>...
```

- **<card>** - pci name of card or just card index
- **<profile>** - name of profile

You can specify `<card>` as `active` if you want to work with just active card,
and not some specific card.

Example:

```
pa-switch-profile active output:hdmi-stereo output:analog-stereo
```

Given command switches specified card from hdmi sound device to speaker analog
sound device.
