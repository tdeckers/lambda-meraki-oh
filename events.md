# Meraki events

## Documentation

https://developer.cisco.com/meraki/webhooks/#!introduction/overview
Samples: https://developer.cisco.com/meraki/webhooks/#webhook-sample-alerts

## Examples

Test

```json
{
    "version": "0.1",
    "sharedSecret": "",
    "sentAt": "2019-11-09T08:41:57.907680Z",
    "organizationId": "456954",
    "organizationName": "Tom Deckers",
    "organizationUrl": "https://n207.meraki.com/o/aa9sQd/manage/organization/overview",
    "networkId": "L_665406844943993490",
    "networkName": "Home",
    "networkUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/wired_status",
    "alertId": "",
    "alertType": "Settings changed",
    "occurredAt": "2019-11-09T08:41:57.906664Z",
    "alertData": {}
}
```

Client disconnects from LAN (~ 1 minute delay):

```json
{
    "version": "0.1",
    "sharedSecret": "verysecret",
    "sentAt": "2019-11-09T08:52:39.258809Z",
    "organizationId": "456954",
    "organizationName": "Tom Deckers",
    "organizationUrl": "https://n207.meraki.com/o/aa9sQd/manage/organization/overview",
    "networkId": "L_665406844943993490",
    "networkName": "Home",
    "networkUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/wired_status",
    "deviceSerial": "Q2MN-5J4X-FLQD",
    "deviceMac": "00:18:0a:3c:83:b0",
    "deviceName": "Gateway",
    "deviceUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/new_wired_status",
    "deviceTags": [],
    "deviceModel": "MX64W",
    "alertId": "679480593797143137",
    "alertType": "Client connectivity changed",
    "occurredAt": "2019-11-09T08:52:17.913000Z",
    "alertData": {
        "mac": "3C:28:6D:29:A7:66",
        "ip": "192.168.2.246",
        "connected": "false",
        "clientName": "Pixel 3 (Tom)",
        "clientUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/usage/list#c=k567cb7"
    }
}
```

Client connects to LAN:
```json
{
    "version": "0.1",
    "sharedSecret": "verysecret",
    "sentAt": "2019-11-09T08:56:18.503438Z",
    "organizationId": "456954",
    "organizationName": "Tom Deckers",
    "organizationUrl": "https://n207.meraki.com/o/aa9sQd/manage/organization/overview",
    "networkId": "L_665406844943993490",
    "networkName": "Home",
    "networkUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/wired_status",
    "deviceSerial": "Q2MN-5J4X-FLQD",
    "deviceMac": "00:18:0a:3c:83:b0",
    "deviceName": "Gateway",
    "deviceUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/new_wired_status",
    "deviceTags": [],
    "deviceModel": "MX64W",
    "alertId": "679480593797143141",
    "alertType": "Client connectivity changed",
    "occurredAt": "2019-11-09T08:55:45.422000Z",
    "alertData": {
        "mac": "3C:28:6D:29:A7:66",
        "ip": "192.168.2.246",
        "connected": "true",
        "clientName": "Pixel 3 (Tom)",
        "clientUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/usage/list#c=k567cb7"
    }
}
```