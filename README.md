[![GoDoc](https://godoc.org/github.com/AlecRandazzo/GoFor/pkg/gofor?status.png)](https://godoc.org/github.com/AlecRandazzo/GoFor/pkg/gofor) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/AlecRandazzo/GoFor-Windows-Collector/issues)

# GoFor Collector
GoFor (Go Forensics) Collector is geared towards augmenting EDR toolsets. Unfortunately, not all EDR toolsets has the capability of collecting forensically relevant files from endpoints. The GoFor Collector looks to remedy that.

## Usage

### GoFor Collector

To collect all forensic files:
```gofor-collector.exe /z whatever.zip /g a```

To collect just event logs:
```gofor-collector.exe /z whatever.zip /g e```

To collect $MFT and registry hives: ```gofor-collector.exe /z whatever.zip /g mr```

For `/g` concatenate the abbreviation characters together for what you want. The order doesn't matter. Valid values are `a` for all (defaults to this if you don't use `/g`), `m` for $MFT, `r` for system registries, `u` for user registries, `e` for event logs.

## Currently Available Features
- GoFor Collector: Windows command line collector that can acquire the files listed below and write them to a zip file.
  - OS Drive $MFT
  - All user NTUSER.DAT and USRCLASS.DAT
  - SYSTEM and SOFTWARE registry hives
  - All Windows event EVTX files

## Future Plans
- Add support to the GoFor collector for uploading to GCP and AWS.
