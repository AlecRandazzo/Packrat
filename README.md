[![CircleCI](https://circleci.com/gh/AlecRandazzo/Packrat.svg?style=svg)](https://circleci.com/gh/AlecRandazzo/Packrat) [![codecov](https://codecov.io/gh/AlecRandazzo/Packrat/branch/master/graph/badge.svg)](https://codecov.io/gh/AlecRandazzo/Packrat) [![Go Report Card](https://goreportcard.com/badge/github.com/AlecRandazzo/Packrat)](https://goreportcard.com/report/github.com/AlecRandazzo/Packrat) [![GoDoc](https://godoc.org/github.com/AlecRandazzo/Packrat/pkg/gofor?status.png)](https://godoc.org/github.com/AlecRandazzo/Packrat) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/AlecRandazzo/Packrat/issues)

# Packrat
Packrat is a forensic collector geared towards augmenting EDR toolsets. Unfortunately, not all EDR toolsets have the capability of collecting forensically relevant files from endpoints. The GoFor Collector looks to remedy that.

## Usage

```usage: forensic-collector.exe [<flags>]
   
Flags:
  --help                         Show context-sensitive help (also try
                                 --help-long and --help-man).
  --debug                        Enable debug mode.
  --all                          Collect all forensic artifacts.
  --mft                          Collect the system drive MFT.
  --mft-all                      Collect all attached volume MFTs.
  --mft-letters=MFT-LETTERS ...  Collect volume MFTs by volume letter.
  --reg                          Collect all registry hives, both system and
                                 user hives.
  --events                       Collect all event logs.
  --browser                      Collect browser history
  --custom-config=CUSTOM-CONFIG  Custom configuration file that will overwrite
                                 built in config.
  --throttle                     This setting will limit the process to a single
                                 thread. This will reduce the CPU load.
  --output=OUTPUT                Specify the name of the output file. If not
                                 specified, the file name defaults to the host
                                 name and a timestamp.
```

### Examples

Collect all the things: `forensic-collector.exe --all`

Collect just the system drive MFT and export to a custom name zip file: `forensic-collector.exe --mft --output out.zip`

Collect event logs and registry hives: `forensic-collector.exe --events --reg`

Use a custom configuration for collection (see example config in `config/config.yml`): `forensic-collector.exe --custom-config config.yml`
