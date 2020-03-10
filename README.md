icannzplit
==========

vpn provider like ipvanish allows you to route your traffic through one of the many endpoints that they provide all accross the world. 
However using a single endpoint may significantly reduce your internet bandwidth, especially with software like bittorrent.

IcannZplit solves this issue by openning up as many openvpn tunnel as you want and then splitting up the whole addressable ipv4 ranges accross those different tunnel interface.
This is an example of a routing plan set up by icannzplit:
```
{
	"routes": {
		"ipvn-ams-a29": [
			"64.0.0.0/3",
			"224.0.0.0/3"
		],
		"ipvn-atl-a67": [
			"96.0.0.0/3"
		],
		"ipvn-chi-a49": [
			"128.0.0.0/3"
		],
		"ipvn-dal-a36": [
			"0.0.0.0/3",
			"160.0.0.0/3"
		],
		"ipvn-nyc-a15": [
			"32.0.0.0/3",
			"192.0.0.0/3"
		]
	}
}
```

Features:

* vpn provider supported: ipvanish
* configuration file to select how to split the connection
    * static list of endpoint
    * fixed amount of endpoint, choosen randomly
    * fixed amount of endpoints per exit country, choosen randomly 
* runs openvpn (it uses [go-openvpn](https://github.com/Marlinski/go-openvpn))
* creates a routing plan and executes it automatically

Todo:

* add custom openvpn config and more vpn provider

# Building

simply run the following command

```
go build
```

# Usage

run it with --help for the list of options

```
./icannzplit --help
NAME:
   icannzplit - A new cli application

USAGE:
   icannzplit [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value         specify icannzplit home directory (default: "/icannzplit")
   --ipvanish-user value  specify ipvanish username (default: "**IPVANISH_USERNAME**")
   --ipvanish-pass value  specify ipvanish password (default: "**IPVANISH_PASSWORD**")
   --help, -h             show help
   --version, -v          print the version
```

Simple example with ipvanish (as root), if the homedir is empty it will create a configuration automatically.
The below example runs icannzplit with 5 randomly selected exit points:

```
./icannzplit --homedir /etc/icannzplit/ --ipvanish-user <USER> --ipvanish-pass <PASS>
2020/03/10 23:04:37  I C A N N    M U L T I P L E X ! 
2020/03/10 23:04:37  =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- 
23:04:37.753 ▶ DEBU 001 icannzplit:ipvanish >> ipvanish: pulling configurations from: "https://www.ipvanish.com/software/configs/configs.zip"
23:04:37.855 ▶ DEBU 002 icannzplit:ipvanish >> ipvanish: unzipping configuration to "/etc/icannzplit/ipvanish"
23:04:37.931 ▶ DEBU 003 icannzplit:ipvanish >> ipvanish: reading configurations in: "/etc/icannzplit/ipvanish"
23:04:37.938 ▶ DEBU 004 icannzplit:zplit >> using default config
23:04:37.938 ▶ DEBU 005 icannzplit:zplit >> saving settings into "plop/settings.json"
23:04:37.939 ▶ NOTI 006 icannzplit:zplit >> routing plan..
{
	"routes": {
		"ipvn-ams-a29": [
			"64.0.0.0/3",
			"224.0.0.0/3"
		],
		"ipvn-atl-a67": [
			"96.0.0.0/3"
		],
		"ipvn-chi-a49": [
			"128.0.0.0/3"
		],
		"ipvn-dal-a36": [
			"0.0.0.0/3",
			"160.0.0.0/3"
		],
		"ipvn-nyc-a15": [
			"32.0.0.0/3",
			"192.0.0.0/3"
		]
	}
}
...
```

# Configuration

By default if settings.json does not exists it will create it. It looks like this:

```
{
   "ipvanish": {
      "filter-type": 0,
      "filter-random": 5,
      "filter-country": {
         "NL": 3,
         "CH": 2
      },
      "filter-static": {
         "sin-a05",
         "fra-a01"
      },
   }
}
```

There are three different ways of choosing its list of exit points:

* filter-type=0: it will take the number set in `filter-random` and randomly select that many exit points from all the configuration available.
* filter-type=1: for each country listed in `filter-country`, it will randomly select as many exist points as it is configured.
* filter-type=2: will use exactly the configuration listes in `filter-static`
