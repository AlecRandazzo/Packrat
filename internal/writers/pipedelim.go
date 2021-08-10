// Copyright (c) 2020 Alec Randazzo

package writers

import (
	"io"
)

// PipeDelimitedWriter receiver used with the ResultWriter method that would write the csv results to csv.
type PipeDelimitedWriter struct {
	writer        *io.Writer
	headerWritten bool
}

func NewPipeDelimitedWriter(writer *io.Writer) PipeDelimitedWriter {
	return PipeDelimitedWriter{writer: writer}
}

// ResultWriter writes the results to csv.
func (pipeDelimitedWriter *PipeDelimitedWriter) Write(data interface{}) error {
	//delimiter := "|"
	//csvHeader := []string{
	//	"Record Number",
	//	"Directory",
	//	"System File",
	//	"Hidden",
	//	"Read-only",
	//	"Deleted",
	//	"File Path",
	//	"File Name",
	//	"File Size",
	//	"File Created",
	//	"File Modified",
	//	"File Accessed",
	//	"File Entry Modified",
	//	"FileName Created",
	//	"FileName Modified",
	//	"Filename Accessed",
	//	"Filename Entry Modified",
	//	"\n",
	//}
	//
	//// Write header if not already written
	//if pipeDelimitedWriter.headerWritten != true {
	//	headerSize := len(csvHeader)
	//	for index, header := range csvHeader {
	//		_ = pipeDelimitedWriter.Write([]byte(header))
	//		if index < headerSize-2 {
	//			_ = pipeDelimitedWriter.Write([]byte(delimiter))
	//		}
	//	}
	//	pipeDelimitedWriter.headerWritten = true
	//}
	//
	//csvRow := []string{
	//	fmt.Sprint(file.RecordNumber),                  //Record Number
	//	strconv.FormatBool(file.DirectoryFlag),         //Directory Flag
	//	strconv.FormatBool(file.SystemFlag),            //System file flag
	//	strconv.FormatBool(file.HiddenFlag),            //Hidden flag
	//	strconv.FormatBool(file.ReadOnlyFlag),          //Read only flag
	//	strconv.FormatBool(file.DeletedFlag),           //Deleted Flag
	//	file.FilePath,                                  //File Directory
	//	file.FileName,                                  //File Name
	//	strconv.FormatUint(file.PhysicalFileSize, 10),  // File Size
	//	file.SiCreated.Format("2006-01-02T15:04:05Z"),  //File Created
	//	file.SiModified.Format("2006-01-02T15:04:05Z"), //File Modified
	//	file.SiAccessed.Format("2006-01-02T15:04:05Z"), //File Accessed
	//	file.SiChanged.Format("2006-01-02T15:04:05Z"),  //File entry Modified
	//	file.FnCreated.Format("2006-01-02T15:04:05Z"),  //FileName Created
	//	file.FnModified.Format("2006-01-02T15:04:05Z"), //FileName Modified
	//	file.FnAccessed.Format("2006-01-02T15:04:05Z"), //FileName Accessed
	//	file.FnChanged.Format("2006-01-02T15:04:05Z"),  //FileName Entry Modified
	//	"\n", // Newline
	//}
	//
	//csvRowSize := len(csvRow)
	//for index, item := range csvRow {
	//	_ = pipeDelimitedWriter.Write([]byte(item))
	//	if index < csvRowSize-2 {
	//		_ = pipeDelimitedWriter.Write([]byte(delimiter))
	//	}
	//}

	return nil
}
