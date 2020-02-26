<p align="center"><img src="https://user-images.githubusercontent.com/4198311/74932238-4912ae80-541c-11ea-92eb-9f9ab40337dd.png" width="50%"></p>


# MrSans

[![Build Status](https://travis-ci.com/skyzh/MrSans.svg?branch=master)](https://travis-ci.com/skyzh/MrSans)

Mr. Sans is part of [the BlueSense project](https://github.com/skyzh/BlueSense). He is the climate reporter of BlueSense.

He does hourly report in Telegram channel [Sans loves Monitoring](https://t.me/thebluesense).

## Configuration

Create `config.toml`. Refer to `config.example.toml` for more information.

Add these environment variables:
```bash
# Set Firebase credentials
GOOGLE_APPLICATION_CREDENTIALS=/opt/bluesense-9e31b-firebase-adminsdk-5vv96-31ac2e637a.json
# If you're using Mr. Sans behind a proxy, set HTTP proxy
http_proxy=http://127.0.0.1:8118
https_proxy=http://127.0.0.1:8118
no_proxy="127.0.0.1,localhost"
```

## Features

### Reporter in Telegram

Mr. Sans reports hourly and daily to the Telegram channel.


<img src="https://user-images.githubusercontent.com/4198311/74337523-6374d880-4ddb-11ea-991f-a984d265e649.png" width="48%"><img src="https://user-images.githubusercontent.com/4198311/74337637-a0d96600-4ddb-11ea-9996-cd95e9175a98.png" width="48%">

### Report Incident

Mr. Sans will report incident in Telegram Channel (refer to config `telegram.log_chat_id`).

### Checkpoint in Firebase

Mr. Sans will periodically checkpoint data from prometheus to firebase
for achieve use.

### Maintain BlueSense services

With the help of Grafana webhook, Mr. Sans will help restart BlueSense
service when there's something wrong with Bluetooth connection. Attach
tag `mrsans-do` in Grafana alert with these values: `restart-systemctl`
or `reboot`.

## Related Projects

[BlueSense](https://github.com/skyzh/BlueSense) is the web frontend.

[BlueMarine](https://github.com/skyzh/BlueMarine) runs on Raspberry Pi. It collects data via Bluetooth from Arduino.

[BlueSensor](https://github.com/skyzh/BlueSensor) runs on Arduino. It collects data from sensors.
