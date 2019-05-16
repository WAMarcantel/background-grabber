# Background-Grabber
Grabs a new image from [Unsplash] to keep your backgrounds up to date, no-hassle.

## What do I need to get this to work?
* An [Unsplash] developer account
* [Go] 1.12
* A little bit of time

## Installation
1. Clone this repo
2. Run the following commands:
```shell



```


### Usage

* **-accessKey** *string*
    	Access key to your Unsplash Account
* **-backgroundsDirPath** *string*
    	Path to the backgrounds directory on your machine.
* **-collections** *string*
    	A comma-separated list of collection IDs to filter on. The update will only return items in the specified collections.
* **-count** *int*
    	The number of photos to return. (Default: 1; max: 30) (default 1)
* **-featured**
    	Whether or not to limit an update to featured photos from Unsplash.
* **-logLevel** *string*
    	log levelacceptable levels: panic, fatal, error, warn or warning, info, debug, trace (default "info")
* **-orientation** *string*
    	Filter search results by photo orientation. Valid values are landscape, portrait, and squarish.
* **-query** *string*
    	Limit selection to photos matching a search term.
* **-refreshMinutes** *int*
    	The number of minutes this program will wait until it refreshes your backgrounds (default 1440)
* **-username** *string*
    	Limit selection to a single user.
    	

[Unsplash]: https://unsplash.com/
[Go]: https://golang.org/