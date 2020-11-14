# Lunar.Music.Web

<div align="center">
    <img src="https://gclstorage.blob.core.windows.net/images/Lunar.Music.Web-banner.png" />
</div>

![Go Build](https://github.com/goh-chunlin/Lunar.Music.Web/workflows/Go%20Build/badge.svg?branch=main)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Donate](https://img.shields.io/badge/$-donate-ff69b4.svg)](https://www.buymeacoffee.com/chunlin)

This is a web app for Lunar Music project. User will be using this web portal to communicate with the Raspberry Pi.

## Key Technologies ##
1. Gin Web Framework
1. Microsoft Azure
   1. Storage
   1. Azure Active Directory
1. RabbitMQ

## Usage (Localhost) ##

1. Clone the project to local;
1. (Optional) Start a local HTTP server such as [http-server](https://www.npmjs.com/package/http-server) in the `static` directory;

   This step is optional because now the static files are served from my Azure Storage. So you can directly use mine. Otherwise, feel free to update all the links pointing to those static in the `app\templates` directory to your own URL.
1. Create a `.env` file in the `app` directory with the following content;
   ```
   AZURE_AD_CALLBACK_URL=
   AZURE_AD_CLIENT_ID=
   AZURE_AD_CLIENT_SECRET=
   RABBITMQ_SERVER_CONNECTION_STRING=
   RABBITMQ_CHANNEL_NAME=
   SECURECOOKIE_HASH_KEY=
   SECURECOOKIE_BLOCK_KEY=
   ```
   
   The RabbitMQ part is optional because it is used only for communicating with my Raspberry Pi.
1. Build the go web project in the `app` directory;
1. Run the output exe.

## Contributing ##
First and foremost, thank you! I appreciate that you want to contribute to this project which is my personal project. Your time is valuable, and your contributions mean a lot to me. You are welcomed to contribute to this project development and make it more awesome every day.

Don't hasitate to contact me, open issue, or even submit a PR if you are intrested to contribute to the project.

Together, we learn better.

## License ##

This library is distributed under the GPL-3.0 License found in the [LICENSE](./LICENSE) file.
