# cosmos-proposals-checker

![Latest release](https://img.shields.io/github/v/release/freak12techno/cosmos-proposals-checker)
[![Actions Status](https://github.com/freak12techno/cosmos-proposals-checker/workflows/test/badge.svg)](https://github.com/freak12techno/cosmos-proposals-checker/actions)

cosmos-proposals-checker is a tool that checks all configured chains for new proposals in voting period, then notifies you if one or more of the configured wallets haven't voted on it.

## How can I set it up?

Download the latest release from [the releases page](https://github.com/freak12techno/cosmos-proposals-checker/releases/). After that, you should unzip it and you are ready to go:

```sh
wget <the link from the releases page>
tar <downloaded file>
./cosmos-proposals-checker --config <path to config>
```

Alternatively, install `golang` (>1.18), clone the repo and build it. This will generate a `./main` binary file in the repository folder:
```
git clone https://github.com/freak12techno/cosmos-proposals-checker
cd cosmos-proposals-checker
go build
```

What you probably want to do is to have it running in the background. For that, first of all, we have to copy the file to the system apps folder:

```sh
sudo cp ./cosmos-proposals-checker /usr/bin
```

Then we need to create a systemd service for our app:

```sh
sudo nano /etc/systemd/system/cosmos-proposals-checker.service
```

You can use this template (change the user to whatever user you want this to be executed from. It's advised to create a separate user for that instead of running it from root):

```
[Unit]
Description=Cosmos Proposals Checker
After=network-online.target

[Service]
User=<username>
TimeoutStartSec=0
CPUWeight=95
IOWeight=95
ExecStart=cosmos-proposals-checker --config <config path>
Restart=always
RestartSec=2
LimitNOFILE=800000
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

Then we'll add this service to the autostart and run it:

```sh
sudo systemctl daemon-reload # reload config to reflect changed
sudo systemctl enable cosmos-proposals-checker # put service to autostart
sudo systemctl start cosmos-proposals-checker # start the service
sudo systemctl status cosmos-proposals-checker # validate it's running
```

If you need to, you can also see the logs of the process:

```sh
sudo journalctl -u cosmos-proposals-checker -f --output cat
```

## How does it work?

It queries LCD nodes for the proposals list in voting period, then for each wallet it queries its vote. If you haven't voted, it spawns an alert and sends it to configured notifiers.

## How can I configure it?

All configuration is done via `.toml` config file, which is mandatory. Run the app with `--config <path/to/config.toml>` to specify config. Check out `config.example.toml` to see the params that can be set.

## Notifiers

Currently this program supports the following notifications channels:
1) Telegram

Go to @BotFather in Telegram and create a bot. After that, there are two options:
- you want to send messages to a user. This user should write a message to @getmyid_bot, then copy the `Your user ID` number. Also keep in mind that the bot won't be able to send messages unless you contact it first, so write a message to a bot before proceeding.
- you want to send messages to a channel. Write something to a channel, then forward it to @getmyid_bot and copy the `Forwarded from chat` number. Then add the bot as an admin.

To have fancy commands auto-suggestion, go to @BotFather again, select your bot -> Edit bot -> Edit description and paste the following:
```
proposals - List proposals and wallets' votes on them
proposals_mute - Mutes a proposal
proposals_mutes - List active proposal mutes
help - Displays help
```

Then add a Telegram config to your config file (see `config.example.toml` for reference).

2) PagerDuty

Go to your PaderDuty page, then go to Services. Create a service if you haven't created one already. Select this service, then go to "Integrations" tab, add an integration there. Copy the integration key and add it to the `pagerduty` part in config (see `config.example.toml` for reference). Additionally, override PagerDuty URL if you are using EU version.


## Which networks this is guaranteed to work?

In theory, it should work on a Cosmos-based blockchains that expose a REST server.

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.
