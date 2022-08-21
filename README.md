# Boyd Bot

This is the newest installment of the Boyd Bot Trilogy. This repo contains a lot of material from the prior Boyd-Bot and Boyd-2
repos. I created this new repo instead of updating the old ones because I effectively merged the two and wanted to do some cleaning; frankly, it was easier to start anew than to deal with the mess I had made myself. I plan to continue what dev I do on Boyd from here now that I
don't plan to use IRC going forward.

Further improvements and migration to slash commands were contributed by [@Alextopher](https://github.com/Alextopher).

## Installation

```text
git clone https://github.com/bashawhm/boyd_cooper.git
cd boyd_cooper
go build
```

Create the `.env` file and add your Discord bot token.

```text
DISCORD_TOKEN=<YOUR BOT TOKEN>
```

Run the server; this will create `qoutes.txt`

```text
./boyd_copper
```
