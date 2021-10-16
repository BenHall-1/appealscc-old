# AppealsCC

Appeals.cc is a website dedicated to being able to submit appeals for twitch, discord, minecraft and more!

It contains authentication for all tyeps of accounts to allow them to submit an appeal for any of the above services, with links to Discord servers, Twitch channels, Minecraft servers and much more!

## How it works

The workflow is as follows:

- You visit <name>.appeals.cc
- Hit login
- Choose how you want to login (Twitch, Discord, Mojang)
- Login with said OAuth2 service
- On discord, it will check whether you are banned on the server\*
- A pre-configured form will then appear
- Once submitted, a list of notifications will be sent out
  - Email to the person who is appealing
  - Email to the team who runs the appeals account
  - Any push notifications configured\*\*
  - Discord webhooks\*\*

*\* Only if you have the AppealsCC bot*
*\*\* Coming Soon*