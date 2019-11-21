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
						ParentDirRecordNumber: 0,
						LogicalFileSize:       0,
						PhysicalFileSize:      0,
						FileNameFlags:         mft.FileNameFlags{},
						FileNamespace:         "",
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
						fullPathString: `c:\$LogFile`,
						fullPathRegex:  nil,
						fileNameString: "$LogFile",
						fileNameRegex:  nil,
					},
				},
			},
			dummyFile:                 `test\testdata\dummyntfs`,
			wantErr:                   false,
			wantDirectoryTree:         mft.DirectoryTree{},
			wantListOfPossibleMatches: possibleMatches{},
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
			gotListOfPossibleMatches, gotDirectoryTree, err := findPossibleMatches(tt.args.volumeHandler, tt.args.listOfSearchKeywords)
			if (err != nil) != tt.wantErr {
				t.Errorf("findPossibleMatches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotListOfPossibleMatches, tt.wantListOfPossibleMatches) {
				t.Errorf("findPossibleMatches() gotListOfPossibleMatches = %v, want %v", gotListOfPossibleMatches, tt.wantListOfPossibleMatches)
			}
			if !reflect.DeepEqual(gotDirectoryTree, tt.wantDirectoryTree) {
				t.Errorf("findPossibleMatches() gotDirectoryTree = %v, want %v", gotDirectoryTree, tt.wantDirectoryTree)
			}
		})
	}
}
