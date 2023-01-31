# Simple usage
Once everything is set up, you only need to type a simple command into a terminal
to download all activities in your Strava account and save them to a JSON file.

```sh
strava-dl-json dl
```

# View downloaded data
```sh
cd $HOME/strava-data/

head current.json
```
```json
[
  {
    "resource_state": 2,
    "athlete": {
      "id": 123,
      "resource_state": 1
    },
    "name": "Sunday long run 🌦",
    "distance": 20052.4,
    "moving_time": 7325
    // ...
  }
  // ...
]
```
Or simply open `current.json` with `Firefox`, which has a pretty good JSON viewer built-in.

# Initial Setup  

The following steps are a **onetime** process.

To proceed, `git` and the `Go` programming language must be installed on your system.

### 1) Clone this repo and create the binary

Run the following commands to clone this repository and create single executable.
```sh
cd $HOME
git clone https://github.com/ad-8/strava-dl-json.git
cd strava-dl-json
go install
```
By default, Go will put the executable in `$HOME/go/bin`. 
Add this folder to your `$PATH` so you can later on simply type `strava-dl-json` in your terminal
(otherwise you would have to type `$HOME/go/bin/strava-dl-json` each time).

### 2) Get Strava Authorization Codes
Strava describes how to do this here: [developers.strava.com](https://developers.strava.com/docs/getting-started/).

The *Getting your data* section of this 
[blog post](https://towardsdatascience.com/using-the-strava-api-and-pandas-to-explore-your-activity-data-d94901d9bfde)
describes the process in other words, which might be helpful.


You will need to obtain the following three codes in order to set up `strava-dl-json` successfully:
- Client ID
- Client Secret
- Refresh Token

Please note that the Refresh Token needed is *not* the Request Token found under [https://www.strava.com/settings/api](https://www.strava.com/settings/api),
but the one which will be generated in process described on [developers.strava.com](https://developers.strava.com/docs/getting-started/).
This is because the Refresh Token must have *read_all* scope.

### 3) Create a download folder where the exported JSON file will be stored
Currently, the default name for the download folder is `strava-data`,
which must be located in your `$HOME` directory, e.g. `/Users/alex/strava-data`.

The name of the download folder can currently only be changed by changing
the value of `appDataDir` in the `cmd/config.go` file.

```sh
cd $HOME
mkdir strava-data
```

### 4) Create a .env file
Create a single file with the name `.env` and add the codes from step 2) to it.
```sh
cd $HOME/strava-data
touch .env
```

Open the `.env` file with your favorite editor and add the following 3 lines to it:
```txt
CLIENT_ID = 123
CLIENT_SECRET = "foo"
REFRESH_TOKEN = "bar"
```
Make sure to replace `123`, `"foo"` and `"bar"` with the actual codes created in step 2).


### 5) Run `strava-dl-json`
```sh
strava-dl-json dl
```

You should see feedback similar to this if everything went well:
```txt
the token expires in 02:57:10 (will be automatically refreshed)

downloaded 1809 activities in 2.532204571s
sucessfully written to file "/Users/alex/strava-data/current.json"
```

Your download folder now contains a `current.json` file, 
with all Strava activities found in your Strava account. 

# JSON to CSV
This can be done in very few lines of code with the Python `Pandas` libray, e.g. see
[this Stackoverflow answer](https://stackoverflow.com/a/37307324).

# TODO / Planned Updates
- set download folder, name of .env file and name of the JSON file via a config file
- prefix the JSON file with the current date and time (so the previous one won't be overwritten)
