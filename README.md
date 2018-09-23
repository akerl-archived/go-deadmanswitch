go-deadmanswitch
=========

[![Build Status](https://img.shields.io/travis/com/akerl/go-deadmanswitch.svg)](https://travis-ci.com/akerl/go-deadmanswitch)
[![GitHub release](https://img.shields.io/github/release/akerl/go-deadmanswitch.svg)](https://github.com/akerl/go-deadmanswitch/releases)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)

Dead man switch via an AWS Lambda. This has two parts:

* It responds to API GW invocations, marking when things "check in"
* It runs on a recurring basis to validate that everything it knows about has "checked in" recently enough. If not, it sends an alert.

## Usage

## Installation

## License

go-deadmanswitch is released under the MIT License. See the bundled LICENSE file for details.
