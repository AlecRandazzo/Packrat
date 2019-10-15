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

func (dataRunReader DataRunsReader) Read(byteSliceToPopulate []byte) (numberOfBytesRead int, err error) {
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
			log.Debugf("Copying the file '%v' via raw method. The file has %v data runs with a total size of %v. First data run is at the absolute offset of %v",
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
			log.Debugf("Finished copying the file '%v' via raw method.", dataRunReader.fileName)
			return
		}
		log.Debugf("Finished copying via raw method data run number %v for the file '%v'. There are %v data runs left and %v bytes left to read. Moving onto data run number %v.",
			dataRunReader.dataRunTracker,
			dataRunReader.fileName,
			len(dataRunReader.DataRuns)-dataRunReader.dataRunTracker-1,
			dataRunReader.bytesLeftToReadTracker,
			dataRunReader.dataRunTracker+1,
		)

		// Increment our tracker
		dataRunReader.dataRunTracker++

		// Seek to the offset of the next datarun
		_, _ = syscall.Seek(dataRunReader.VolumeHandler.Handle, dataRunReader.DataRuns[dataRunReader.dataRunTracker].AbsoluteOffset, 0)
	}

	return
}
