## Setup System ENV Variables from Windows Terminal
`[System.Environment]::SetEnvironmentVariable("DISCORD_CLIENT_ID", "<Your Client ID>", "User")`
`[System.Environment]::SetEnvironmentVariable("DISCORD_CLIENT_SECRET", "<Your Client Secret>", "User")`
`[System.Environment]::SetEnvironmentVariable("DISCORD_CLIENT_AUTH_URL", "<Your Client Auth URL>", "User")`

Note: After, you might need to restart VSCode to have these load.
To Verify they're added.
`echo $env:DISCORD_CLIENT_ID`
`echo $env:DISCORD_CLIENT_SECRET`
`echo $env:DISCORD_CLIENT_AUTH_URL`

To add these variables system wide in lieu of user wide, simply replace "User" with "Machine"

## Discord API Control Panel
https://discord.com/developers/applications