# Flash
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Actions Status](https://github.com/getgauge/flash/workflows/build/badge.svg)](https://github.com/getgauge/flash/actions)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

Execution progress reporter for [Gauge](http://getgauge.io).

## Usage

### Install through Gauge
Run the following command in the Gauge project directory to install and add the plugin to project.
```
gauge install flash
```

* Installing specific version
```
gauge install flash --version 0.0.1
```

#### Offline installation
* Download the plugin from [Releases](https://github.com/getgauge/flash/releases)
```
gauge install flash --file flash-0.0.0-linux.x86_64.zip
```

#### Usage 
- Ensure  that you project's `manifest.json` file contains `flash` in the Plugins list.

Something like :  
```
{
"Language": "java",
"Plugins": [
"html-report",
"flash"
]
}	
```

- Execute specs and open the URL in browser shown in **console output**. 	e.g. http://127.0.0.1:[FLASH_SERVER_PORT]
- FLASH_SERVER_PORT is a random available port, but can be configured using the config below 


### Configuration

* To use a specific port, set `FLASH_SERVER_PORT={port}` as environment variable or in `env/default/flash.properties` file.

## Build from Source

### Requirements
* [Golang](http://golang.org/)

### Compiling

Dependencies
```
go get ./...
```
Compilation
```
go run build/make.go
```

For cross-platform compilation

```
go run build/make.go --all-platforms
```

### Installing
After compilation

```
go run build/make.go --install
```

Installing to a CUSTOM_LOCATION

```
go run build/make.go --install --plugin-prefix CUSTOM_LOCATION
```

### Creating distributable

Note: Run after compiling

```
go run build/make.go --distro
```

For distributable across platforms: Windows and Linux for both x86 and x86_64

```
go run build/make.go --distro --all-platforms
```

## License

![GNU Public License version 3.0](http://www.gnu.org/graphics/gplv3-127x51.png)
`Flash` is released under [GNU Public License version 3.0](http://www.gnu.org/licenses/gpl-3.0.txt)

## Copyright

Copyright 2017 ThoughtWorks, Inc.
