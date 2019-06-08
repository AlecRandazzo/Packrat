[![GoDoc](https://godoc.org/github.com/AlecRandazzo/GoFor?status.png)](https://godoc.org/github.com/AlecRandazzo/GoFor) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwyl/esta/issues)

# GoFor
GoFor (Go Forensics) is a forensic library geared towards augmenting EDR toolsets. Unfortunately, not all EDR toolsets has the capability of collecting forensically relevant files from endpoints.

## Usage

```gofor-collector.exe --zip out.zip```

## Currently Available Features
- GoFor Collector: Windows command line collector that can acquire the files listed below and write them to a zip file.
  - OS Drive $MFT
  - All user NTUSER.DAT and USRCLASS.DAT
  - SYSTEM and SOFTWARE registry hives
  - All Windows event EVTX files

## Future Plans
- Add support to the GoFor collector for uploading to GCP and AWS.
- Finish capabilities of stand alone MFT parser.
- Develop an EVTX parser.
- Develop a registry parser.
- Develop a kubernetes cluster that can be deployed to GCP that will autoparse data uploaded to a storage bucket and insert the results into a BigTable cluster.
- Develop a frontend to review and markdown forensic data stored in BigTable.
