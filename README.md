# musync 
Simple utility to sync your spotify library with a playlist. 

## Motivation 
I use spotify only to get recommendations and normally use [music](music.reaper.im) and it's generally easier to work with public playlists when importing the tracks on it. Decided I just sync everything from the user library to a public playlist. 

## Usage 
You execute the binary while pointing it to a config that follows the [config.template.yaml](/config.template.yaml), the default config file is to be in the folder you are running the binary with the name `musync.yaml`

## Auto Sync 
Since it's a simple cli tool , you can add it to a crontab or any other scheduler that you may use on your system. 