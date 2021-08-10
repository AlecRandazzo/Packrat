// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

// ResultWriter interface for result writers to allow for output format extensibility.
type ResultWriter interface {
	ResultWriter(streamer io.Writer, outputChannel *chan UsefulMftFields, waitGroup *sync.WaitGroup)
}

// CsvResultWriter receiver used with the ResultWriter method that would write the csv results to csv.
type CsvResultWriter struct{}

// ResultWriter writes the results to csv.
func (csvResultWriter *CsvResultWriter) ResultWriter(streamer io.Writer, outputChannel *chan UsefulMftFields, waitGroup *sync.WaitGroup) {
	delimiter := "|"
	csvHeader := []string{
		"Record Number",
		"Directory",
		"System File",
		"Hidden",
		"Read-only",
		"Deleted",
		"File Path",
		"File Name",
		"File Size",
		"File Created",
		"File Modified",
		"File Accessed",
		"File Entry Modified",
		"FileName Created",
		"FileName Modified",
		"Filename Accessed",
		"Filename Entry Modified",
		"\n",
	}

	// Write CSV header
	headerSize := len(csvHeader)
	for index, header := range csvHeader {
		_, _ = streamer.Write([]byte(header))
		if index < headerSize-2 {
			_, _ = streamer.Write([]byte(delimiter))
		}
	}

	openChannel := true
	for {
		var file UsefulMftFields
		file, openChannel = <-*outputChannel
		if openChannel == false {
			break
		}
		csvRow := []string{
			fmt.Sprint(file.RecordNumber),                  //Record Number
			strconv.FormatBool(file.DirectoryFlag),         //Directory Flag
			strconv.FormatBool(file.SystemFlag),            //System file flag
			strconv.FormatBool(file.HiddenFlag),            //Hidden flag
			strconv.FormatBool(file.ReadOnlyFlag),          //Read only flag
			strconv.FormatBool(file.DeletedFlag),           //Deleted Flag
			file.FilePath,                                  //File Directory
			file.FileName,                                  //File Name
			strconv.FormatUint(file.PhysicalFileSize, 10),  // File Size
			file.SiCreated.Format("2006-01-02T15:04:05Z"),  //File Created
			file.SiModified.Format("2006-01-02T15:04:05Z"), //File Modified
			file.SiAccessed.Format("2006-01-02T15:04:05Z"), //File Accessed
			file.SiChanged.Format("2006-01-02T15:04:05Z"),  //File entry Modified
			file.FnCreated.Format("2006-01-02T15:04:05Z"),  //FileName Created
			file.FnModified.Format("2006-01-02T15:04:05Z"), //FileName Modified
			file.FnAccessed.Format("2006-01-02T15:04:05Z"), //FileName Accessed
			file.FnChanged.Format("2006-01-02T15:04:05Z"),  //FileName Entry Modified
			"\n", // Newline
		}

		csvRowSize := len(csvRow)
		for index, item := range csvRow {
			_, _ = streamer.Write([]byte(item))
			if index < csvRowSize-2 {
				_, _ = streamer.Write([]byte(delimiter))
			}
		}

	}
	waitGroup.Done()
	return
}
