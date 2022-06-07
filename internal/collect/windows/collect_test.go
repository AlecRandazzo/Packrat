// Copyright (c) 2022 Alec Randazzo

package windows

import "testing"

//func TestCollect(t *testing.T) {
//	type args struct {
//		exportList   FileExportList
//		resultWriter ZipResultWriter
//		handler      handler
//	}
//	tests := []struct {
//		name          string
//		args          args
//		wantErr       bool
//		zipTestOutput string
//		wantZipHash   string
//	}{
//		{
//			name: "test1",
//			args: args{
//				exportList: FileExportList{
//					0: {
//						fullPath:      `%SYSTEMDRIVE%:\$MFT`,
//						FullPathRegex: false,
//						FileName:      `$MFT`,
//						FileNameRegex: false,
//					},
//				},
//				resultWriter: ZipResultWriter{},
//				handler: &handler{
//					handle:       nil,
//					volumeLetter: "",
//					vbr:          vbr.VolumeBootRecord{},
//					reader:       nil,
//					lastOffset:   0,
//					filePath:     filepath.FromSlash("../../test/testdata/dummyntfs"),
//				},
//			},
//			wantErr:       false,
//			zipTestOutput: filepath.FromSlash("../../test/testdata/collecttestzip.zip"),
//			wantZipHash:   "29f689d96a790b68df7e84c9e04ef741",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			fileHandle, _ := os.Create(tt.zipTestOutput)
//			zipWriter := zip.NewWriter(fileHandle)
//			tt.args.resultWriter = ZipResultWriter{
//				ZipWriter:  zipWriter,
//				FileHandle: fileHandle,
//			}
//			err := Collect(tt.args.handler, tt.args.exportList, &tt.args.resultWriter)
//			if err != nil {
//				t.Error(err)
//				return
//			}
//			zipWriter.Close()
//			fileHandle.Close()
//
//			// Get file hash
//			file, _ := os.Open(tt.zipTestOutput)
//			hash := md5.New()
//			_, _ = io.Copy(hash, file)
//			hashInBytes := hash.Sum(nil)[:]
//			gotZipHash := hex.EncodeToString(hashInBytes)
//			file.Close()
//			if gotZipHash != tt.wantZipHash {
//				t.Errorf("collect() gotZipHash = %v, want %v", gotZipHash, tt.wantZipHash)
//			}
//
//			// Cleanup
//			_ = os.Remove(tt.zipTestOutput)
//		})
//	}
//}
//
//func Test_getFiles(t *testing.T) {
//	type args struct {
//		handler         *handler
//		resultWriter         ZipResultWriter
//		listOfSearchKeywords searchTermsList
//	}
//	tests := []struct {
//		name        string
//		args        args
//		wantErr     bool
//		testZip     string
//		wantZipHash string
//	}{
//		{
//			name: "test1",
//			args: args{
//				handler: &handler{
//					handle:       nil,
//					volumeLetter: "c",
//					vbr:          vbr.VolumeBootRecord{},
//					reader:       nil,
//					lastOffset:   0,
//					filePath:     filepath.FromSlash("../../test/testdata/dummyntfs"),
//				},
//				resultWriter: ZipResultWriter{},
//				listOfSearchKeywords: searchTermsList{
//					0: searchTerms{
//						fullPathString: `c:\\$mft`,
//						fullPathRegex:  nil,
//						fileNameString: "$mft",
//						fileNameRegex:  nil,
//					},
//					1: searchTerms{
//						fullPathString: `c:\\$mftmirr`,
//						fullPathRegex:  nil,
//						fileNameString: "$mftmirr",
//						fileNameRegex:  nil,
//					},
//				},
//			},
//			testZip:     filepath.FromSlash("../../test/testdata/getFilesTest1.zip"),
//			wantErr:     false,
//			wantZipHash: "a50b885249c709ae97eeba0e2d6ec78d",
//		},
//		{
//			name: "test2",
//			args: args{
//				handler: &handler{
//					handle:       nil,
//					volumeLetter: "c",
//					vbr:          vbr.VolumeBootRecord{},
//					reader:       nil,
//					lastOffset:   0,
//					filePath:     filepath.FromSlash("../../test/testdata/dummyntfs"),
//				},
//				resultWriter: ZipResultWriter{},
//				listOfSearchKeywords: searchTermsList{
//					0: searchTerms{
//						fullPathString: `c:\\$mftmirr`,
//						fullPathRegex:  nil,
//						fileNameString: "$mftmirr",
//						fileNameRegex:  nil,
//					},
//				},
//			},
//			testZip:     filepath.FromSlash("../../test/testdata/getFilesTest2.zip"),
//			wantErr:     false,
//			wantZipHash: "75c57f05d2879cb723dbec6e2e1e8f83",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			fileHandle, _ := os.Create(tt.testZip)
//			zipWriter := zip.NewWriter(fileHandle)
//			tt.args.resultWriter = ZipResultWriter{
//				ZipWriter:  zipWriter,
//				FileHandle: fileHandle,
//			}
//			err := tt.args.handler.GetHandle()
//			if err != nil {
//				t.Errorf("could not load handler file %s: %v", tt.args.handler.filePath, err)
//			}
//			defer tt.args.handler.handle.Close()
//
//			_ = getFiles(tt.args.handler, &tt.args.resultWriter, tt.args.listOfSearchKeywords)
//			zipWriter.Close()
//			fileHandle.Close()
//
//			// Get file hash
//			file, _ := os.Open(tt.testZip)
//			hash := md5.New()
//			_, _ = io.Copy(hash, file)
//			hashInBytes := hash.Sum(nil)[:]
//			gotZipHash := hex.EncodeToString(hashInBytes)
//			file.Close()
//			if gotZipHash != tt.wantZipHash {
//				t.Errorf("getFiles() gotZipHash = %v, want %v", gotZipHash, tt.wantZipHash)
//			}
//
//			// Cleanup
//			//_ = os.Remove(tt.testZip)
//		})
//	}
//}

func Test_normalizeFilePath(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normalize path",
			args: args{
				filepath: `C:\Windows\evil.exe`,
			},
			want: "C__Windows_evil.exe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeFilePath(tt.args.filepath); got != tt.want {
				t.Errorf("normalizeFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
