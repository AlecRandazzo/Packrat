package windowscollector

import (
	"github.com/AlecRandazzo/GoFor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/windows"
	"io"
)

type DataRunsReader struct {
	VolumeHandler          VolumeHandler
	DataRuns               mft.DataRuns
	fileName               string
	dataRunTracker         int
	bytesLeftToReadTracker int64
	initialized            bool
}

func (dataRunReader *DataRunsReader) Read(byteSliceToPopulate []byte) (numberOfBytesRead int, err error) {
	// Check if this reader has been initialized, if not, do so.
	if dataRunReader.initialized != true {
		dataRunReader.dataRunTracker = 0
		dataRunReader.bytesLeftToReadTracker = dataRunReader.DataRuns[dataRunReader.dataRunTracker].Length
		_, _ = syscall.Seek(dataRunReader.VolumeHandler.Handle, dataRunReader.DataRuns[dataRunReader.dataRunTracker].AbsoluteOffset, 0)
		dataRunReader.initialized = true

		// These are for debug purposes
		if log.GetLevel() == log.DebugLevel {
			totalSize := int64(0)
			for _, dataRun := range dataRunReader.DataRuns {
				totalSize += dataRun.Length
			}
			log.Debugf("Reading the file '%s' via raw method. The file has %d data runs with a total size of %d. First data run is at the absolute offset of %d",
				dataRunReader.fileName,
				len(dataRunReader.DataRuns),
				totalSize,
				dataRunReader.DataRuns[0].AbsoluteOffset,
			)
		}

	}

	numberOfBytesToRead := len(byteSliceToPopulate)

	// Figure out how many bytes are left to read
	if dataRunReader.bytesLeftToReadTracker-int64(numberOfBytesToRead) == 0 {
		dataRunReader.bytesLeftToReadTracker -= int64(numberOfBytesToRead)
	} else if dataRunReader.bytesLeftToReadTracker-int64(numberOfBytesToRead) < 0 {
		numberOfBytesToRead = int(dataRunReader.bytesLeftToReadTracker)
		dataRunReader.bytesLeftToReadTracker = 0
	} else {
		dataRunReader.bytesLeftToReadTracker -= int64(numberOfBytesToRead)
	}

	// Read from the data run
	buffer := make([]byte, numberOfBytesToRead)
	numberOfBytesRead, _ = syscall.Read(dataRunReader.VolumeHandler.Handle, buffer)
	copy(byteSliceToPopulate, buffer)

	// Check to see if there are any bytes left to read in the current data run
	if dataRunReader.bytesLeftToReadTracker == 0 {
		// Check to see if we have read all the data runs.
		if dataRunReader.dataRunTracker+1 == len(dataRunReader.DataRuns) {
			err = io.EOF
			log.Debugf("Finished reading the file '%s' via raw method.", dataRunReader.fileName)
			return
		}
		log.Debugf("Finished reading via raw method data run number %d for the file '%s'. There are %d data runs left. Moving onto data run number %d.",
			dataRunReader.dataRunTracker,
			dataRunReader.fileName,
			len(dataRunReader.DataRuns)-dataRunReader.dataRunTracker-1,
			dataRunReader.dataRunTracker+1,
		)

		// Increment our tracker
		dataRunReader.dataRunTracker++

		// Get the size of the next datarun
		dataRunReader.bytesLeftToReadTracker = dataRunReader.DataRuns[dataRunReader.dataRunTracker].Length

		// Seek to the offset of the next datarun
		_, _ = syscall.Seek(dataRunReader.VolumeHandler.Handle, dataRunReader.DataRuns[dataRunReader.dataRunTracker].AbsoluteOffset, 0)
	}

	return
}
