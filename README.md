# ADZTBotV2

> A little discord bot with a small footprint and easy to use to share music and video with your friends on your discord server

ADZTBotV2 is the successor of [ADZTBot](https://github.com/DarkOnion0/ADZTBot). It has been rewritten in go to make it
much faster, easier to deploy, easier to maintain...

## 🚀 Key Features

- Share music or video with your friends and vote for the post you like
- Support [Slash Commands](https://support.discord.com/hc/en-us/articles/1500000368501-Slash-Commands-FAQ)

## 📖 Usage

To use the bot you just need to type `/` in the message box in discord and the bot command auto-completion will start,
pretty easy right ?!

![img.png](Pictures/UsageScreenCapture.png)

## 💽 Installation

### 🐹 Go Binary

1. Download the binary from the release page
2. Execute the command with the following flags (this flags can be get running the executable with the `-help` flag)

   ```sh
   -chanm #Discord channel id where the post of the music category will be sent to

   -chanv #Discord channel id where the post of the video category will be sent to

   -db #The mongodb database name

   -token #Bot access token

   -url #The mongodb access url

   ```

3. And that's it

### 🐋 Docker

> **_✨ Comming soon...✨_**
