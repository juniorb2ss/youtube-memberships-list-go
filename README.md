## Installation Steps

First you need have Go 1.15 installed, to build the binarie.

```shell script
go build -o youtube
```   

After build, you can run the command:
```shell script 
./youtube -clientid="XXXXXXX" -secret="XXXXXX" memberships-list <channelId>
```   

## Command Usage
> `clientid` and `secret`: You need first create a new APP in Google Console, in order to have access to Youtube Data v3 API.

> `channelId` to filter the members by channelId - Channels that you have access.                           

After you execute the command, an authorization tab will open in your default browser. This page is from Google and is for authorization, you must login with the account that has access to <channelId>. 

After authorizing, Google will redirect the tab to a `localhost` domain, so that the script can obtain the authentication token, all of this is done automatically.

This token will be saved in a file inside the folder, so that subsequent uses do not need authorization again.

The command will generate a .csv file containing up to 1000 members, each member will export member and subscription information.