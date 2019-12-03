/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package windowscollector

import (
	mft "github.com/Go-Forensics/MFT-Parser"
	vbr "github.com/Go-Forensics/VBR-Parser"
	"log"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func Test_checkForPossibleMatch(t *testing.T) {
	type args struct {
		listOfSearchKeywords listOfSearchTerms
		fileNameAttributes   mft.FileNameAttributes
	}
	tests := []struct {
		name                  string
		args                  args
		wantResult            bool
		wantFileNameAttribute mft.FileNameAttribute
		wantErr               bool
	}{
		{
			name:    "null keywords",
			wantErr: true,
			args: args{
				listOfSearchKeywords: nil,
				fileNameAttributes: mft.FileNameAttributes{
					0: mft.FileNameAttribute{
						FnCreated:             time.Time{},
						FnModified:            time.Time{},
						FnAccessed:            time.Time{},
						FnChanged:             time.Time{},
						FlagResident:          true,
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "",
						FileNameLength:        16,
						FileName:              "test",
					},
				},
			},
		},
		{
			name:    "null fn attribute",
			wantErr: true,
			args: args{
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\test`,
						fullPathRegex:  nil,
						fileNameString: "test",
						fileNameRegex:  nil,
					},
				},
				fileNameAttributes: nil,
			},
		},
		{
			name:    "file name exact match",
			wantErr: false,
			args: args{
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\test`,
						fullPathRegex:  nil,
						fileNameString: "test",
						fileNameRegex:  nil,
					},
				},
				fileNameAttributes: mft.FileNameAttributes{
					0: mft.FileNameAttribute{
						FnCreated:             time.Time{},
						FnModified:            time.Time{},
						FnAccessed:            time.Time{},
						FnChanged:             time.Time{},
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "WIN32",
						FileName:              "nope",
					},
					1: mft.FileNameAttribute{
						FnCreated:             time.Time{},
						FnModified:            time.Time{},
						FnAccessed:            time.Time{},
						FnChanged:             time.Time{},
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "WIN32",
						FileName:              "test",
					},
				},
			},
			wantResult: true,
			wantFileNameAttribute: mft.FileNameAttribute{
				FnCreated:             time.Time{},
				FnModified:            time.Time{},
				FnAccessed:            time.Time{},
				FnChanged:             time.Time{},
				ParentDirRecordNumber: 0,
				LogicalFileSize:       0,
				PhysicalFileSize:      0,
				FileNameFlags:         mft.FileNameFlags{},
				FileNamespace:         "WIN32",
				FileName:              "test",
			},
		},
		{
			name:    "file name regex match",
			wantErr: false,
			args: args{
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\test`,
						fullPathRegex:  nil,
						fileNameString: "",
						fileNameRegex:  regexp.MustCompile("^test$"),
					},
				},
				fileNameAttributes: mft.FileNameAttributes{
					0: mft.FileNameAttribute{
						FnCreated:             time.Time{},
						FnModified:            time.Time{},
						FnAccessed:            time.Time{},
						FnChanged:             time.Time{},
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "WIN32",
						FileName:              "nope",
					},
					1: mft.FileNameAttribute{
						FnCreated:             time.Time{},
						FnModified:            time.Time{},
						FnAccessed:            time.Time{},
						FnChanged:             time.Time{},
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "WIN32",
						FileName:              "test",
					},
				},
			},
			wantResult: true,
			wantFileNameAttribute: mft.FileNameAttribute{
				FnCreated:             time.Time{},
				FnModified:            time.Time{},
				FnAccessed:            time.Time{},
				FnChanged:             time.Time{},
				ParentDirRecordNumber: 0,
				LogicalFileSize:       0,
				PhysicalFileSize:      0,
				FileNameFlags:         mft.FileNameFlags{},
				FileNamespace:         "WIN32",
				FileName:              "test",
			},
		},
		{
			name:    "file name no match",
			wantErr: false,
			args: args{
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\test`,
						fullPathRegex:  nil,
						fileNameString: "test",
						fileNameRegex:  nil,
					},
				},
				fileNameAttributes: mft.FileNameAttributes{
					0: mft.FileNameAttribute{
						FnCreated:             time.Time{},
						FnModified:            time.Time{},
						FnAccessed:            time.Time{},
						FnChanged:             time.Time{},
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "WIN32",
						FileName:              "nope",
					},
					1: mft.FileNameAttribute{
						FnCreated:             time.Time{},
						FnModified:            time.Time{},
						FnAccessed:            time.Time{},
						FnChanged:             time.Time{},
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "WIN32",
						FileName:              "test2",
					},
				},
			},
			wantResult:            false,
			wantFileNameAttribute: mft.FileNameAttribute{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotFileNameAttribute, err := checkForPossibleMatch(tt.args.listOfSearchKeywords, tt.args.fileNameAttributes)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkForPossibleMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("checkForPossibleMatch() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotFileNameAttribute, tt.wantFileNameAttribute) {
				t.Errorf("checkForPossibleMatch() gotFileNameAttribute = %v, want %v", gotFileNameAttribute, tt.wantFileNameAttribute)
			}
		})
	}
}

func Test_confirmFoundFiles(t *testing.T) {
	type args struct {
		listOfSearchKeywords  listOfSearchTerms
		listOfPossibleMatches possibleMatches
		directoryTree         mft.DirectoryTree
	}
	tests := []struct {
		name               string
		args               args
		wantFoundFilesList foundFiles
	}{
		{
			name: "matches and no matches",
			wantFoundFilesList: foundFiles{
				0: foundFile{
					dataRuns: nil,
					fullPath: `c:\exactmatch`,
					fileSize: 0,
				},
				1: foundFile{
					dataRuns: nil,
					fullPath: `c:\regexmatch`,
					fileSize: 0,
				},
			},
			args: args{
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\exactmatch`,
						fullPathRegex:  nil,
						fileNameString: "exactmatch",
						fileNameRegex:  nil,
					},
					1: searchTerms{
						fullPathString: "",
						fullPathRegex:  regexp.MustCompile(`^c:\\regexmatch$`),
						fileNameString: "",
						fileNameRegex:  regexp.MustCompile("^regexmatch$"),
					},
					2: searchTerms{
						fullPathString: `c:\nomatch`,
						fullPathRegex:  nil,
						fileNameString: "nomatch",
						fileNameRegex:  nil,
					},
				},
				listOfPossibleMatches: possibleMatches{
					0: possibleMatch{
						fileNameAttribute: mft.FileNameAttribute{
							FnCreated:             time.Time{},
							FnModified:            time.Time{},
							FnAccessed:            time.Time{},
							FnChanged:             time.Time{},
							ParentDirRecordNumber: 5,
							LogicalFileSize:       0,
							PhysicalFileSize:      0,
							FileNameFlags:         mft.FileNameFlags{},
							FileNamespace:         "WIN32",
							FileName:              "exactmatch",
						},
						dataRuns: nil,
					},
					1: possibleMatch{
						fileNameAttribute: mft.FileNameAttribute{
							FnCreated:             time.Time{},
							FnModified:            time.Time{},
							FnAccessed:            time.Time{},
							FnChanged:             time.Time{},
							ParentDirRecordNumber: 5,
							LogicalFileSize:       0,
							PhysicalFileSize:      0,
							FileNameFlags:         mft.FileNameFlags{},
							FileNamespace:         "WIN32",
							FileName:              "regexmatch",
						},
						dataRuns: nil,
					},
					2: possibleMatch{
						fileNameAttribute: mft.FileNameAttribute{
							FnCreated:             time.Time{},
							FnModified:            time.Time{},
							FnAccessed:            time.Time{},
							FnChanged:             time.Time{},
							ParentDirRecordNumber: 7,
							LogicalFileSize:       0,
							PhysicalFileSize:      0,
							FileNameFlags:         mft.FileNameFlags{},
							FileNamespace:         "WIN32",
							FileName:              "exactmatch", // this wont be confirmed since parent dir record num is 7 not 5
						},
						dataRuns: nil,
					},
					3: possibleMatch{
						fileNameAttribute: mft.FileNameAttribute{
							FnCreated:             time.Time{},
							FnModified:            time.Time{},
							FnAccessed:            time.Time{},
							FnChanged:             time.Time{},
							ParentDirRecordNumber: 6,
							LogicalFileSize:       0,
							PhysicalFileSize:      0,
							FileNameFlags:         mft.FileNameFlags{},
							FileNamespace:         "WIN32",
							FileName:              "exactmatch", // this wont be confirmed since parent dir record num is 6 not 5
						},
						dataRuns: nil,
					},
				},
				directoryTree: mft.DirectoryTree{
					5: `c:`,
					6: `d:`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFoundFilesList := confirmFoundFiles(tt.args.listOfSearchKeywords, tt.args.listOfPossibleMatches, tt.args.directoryTree)
			if !reflect.DeepEqual(gotFoundFilesList, tt.wantFoundFilesList) {
				t.Errorf("confirmFoundFiles() gotFoundFilesList = %v, want %v", gotFoundFilesList, tt.wantFoundFilesList)
			}
		})
	}
}

func Test_findPossibleMatches(t *testing.T) {
	type args struct {
		volumeHandler        *VolumeHandler
		listOfSearchKeywords listOfSearchTerms
	}
	tests := []struct {
		name                      string
		args                      args
		wantListOfPossibleMatches possibleMatches
		wantDirectoryTree         mft.DirectoryTree
		wantErr                   bool
		dummyFile                 string
	}{
		{
			name: "find possible matches",
			args: args{
				volumeHandler: &VolumeHandler{},
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\$mftmirr`,
						fullPathRegex:  nil,
						fileNameString: "$mftmirr",
						fileNameRegex:  nil,
					},
					1: searchTerms{
						fullPathString: `c:\software`,
						fullPathRegex:  nil,
						fileNameString: "software",
						fileNameRegex:  nil,
					},
				},
			},
			dummyFile: `test\testdata\dummyntfs`,
			wantErr:   false,
			wantDirectoryTree: mft.DirectoryTree{
				5:  `c:\`,
				11: `c:\$Extend`,
			},
			wantListOfPossibleMatches: possibleMatches{
				0: possibleMatch{
					fileNameAttribute: mft.FileNameAttribute{
						FnCreated:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FnModified:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FnAccessed:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FnChanged:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FlagResident:            true,
						ParentDirRecordNumber:   5,
						ParentDirSequenceNumber: 5,
						LogicalFileSize:         4096,
						PhysicalFileSize:        4096,
						FileNameFlags: mft.FileNameFlags{
							ReadOnly:          false,
							Hidden:            true,
							System:            true,
							Archive:           false,
							Device:            false,
							Normal:            false,
							Temporary:         false,
							Sparse:            false,
							Reparse:           false,
							Compressed:        false,
							Offline:           false,
							NotContentIndexed: false,
							Encrypted:         false,
							Directory:         false,
							IndexView:         false,
						},
						AttributeSize:  112,
						FileNameLength: 16,
						FileNamespace:  "WIN32 & DOS",
						FileName:       "$MFTMirr",
					},
					dataRuns: mft.DataRuns{
						0: mft.DataRun{
							AbsoluteOffset: 8192,
							Length:         4096,
						},
					},
				},
				1: possibleMatch{
					fileNameAttribute: mft.FileNameAttribute{
						FnCreated:    time.Date(2019, 8, 21, 6, 43, 46, 194743600, time.UTC),
						FnModified:   time.Date(2019, 8, 21, 6, 43, 46, 194743600, time.UTC),
						FnAccessed:   time.Date(2019, 8, 21, 6, 43, 46, 194743600, time.UTC),
						FnChanged:    time.Date(2019, 8, 21, 6, 43, 46, 194743600, time.UTC),
						FlagResident: true,
						NameLength: mft.NameLength{
							FlagNamed: false,
							NamedSize: 0,
						},
						AttributeSize:           112,
						ParentDirRecordNumber:   506651,
						ParentDirSequenceNumber: 27,
						LogicalFileSize:         0,
						PhysicalFileSize:        0,
						FileNameFlags: mft.FileNameFlags{
							ReadOnly:          false,
							Hidden:            false,
							System:            false,
							Archive:           true,
							Device:            false,
							Normal:            false,
							Temporary:         false,
							Sparse:            false,
							Reparse:           false,
							Compressed:        false,
							Offline:           false,
							NotContentIndexed: false,
							Encrypted:         false,
							Directory:         false,
							IndexView:         false,
						},
						FileNameLength: 16,
						FileNamespace:  "POSIX",
						FileName:       "SOFTWARE",
					},
					dataRuns: mft.DataRuns{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := dummyHandler{
				Handle:               nil,
				VolumeLetter:         "",
				Vbr:                  vbr.VolumeBootRecord{},
				mftReader:            nil,
				lastReadVolumeOffset: 0,
				filePath:             tt.dummyFile,
			}

			var err error
			*tt.args.volumeHandler, err = GetVolumeHandler("c", handle)
			if err != nil {
				log.Panic(err)
			}
			defer tt.args.volumeHandler.Handle.Close()

			mftRecord0, _ := parseMFTRecord0(tt.args.volumeHandler)
			_, _ = tt.args.volumeHandler.Handle.Seek(tt.args.volumeHandler.Vbr.MftByteOffset, 0)

			foundFile := foundFile{
				dataRuns: mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns,
				fullPath: "$mft",
			}
			tt.args.volumeHandler.mftReader = rawFileReader(tt.args.volumeHandler, foundFile)

			gotListOfPossibleMatches, gotDirectoryTree, err := findPossibleMatches(tt.args.volumeHandler, tt.args.listOfSearchKeywords)
			if (err != nil) != tt.wantErr {
				t.Errorf("findPossibleMatches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotListOfPossibleMatches, tt.wantListOfPossibleMatches) {
				t.Errorf("findPossibleMatches() gotListOfPossibleMatches = %+v, want %+v", gotListOfPossibleMatches, tt.wantListOfPossibleMatches)
			}
			if !reflect.DeepEqual(gotDirectoryTree, tt.wantDirectoryTree) {
				t.Errorf("findPossibleMatches() gotDirectoryTree = %v, want %v", gotDirectoryTree, tt.wantDirectoryTree)
			}
		})
	}
}
