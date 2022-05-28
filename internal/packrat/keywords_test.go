// Copyright (c) 2020 Alec Randazzo

package packrat

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_setupSearchTerms(t *testing.T) {
	type args struct {
		exportList FileExportList
	}
	tests := []struct {
		name                     string
		args                     args
		wantListOfSearchKeywords searchTermsList
		wantErr                  bool
	}{
		{
			name: "this should be successful",
			args: args{exportList: FileExportList{
				0: FileExport{
					FullPath:      `C:\windows`,
					FullPathRegex: false,
					FileName:      "fake.exe",
					FileNameRegex: false,
				},
				1: FileExport{
					FullPath:      `C:\\windows\\.*`,
					FullPathRegex: true,
					FileName:      `.*\.evtx`,
					FileNameRegex: true,
				},
			}},
			wantErr: false,
			wantListOfSearchKeywords: searchTermsList{
				0: searchTerms{
					fullPathString: `c:\windows`,
					fullPathRegex:  nil,
					fileNameString: "fake.exe",
					fileNameRegex:  nil,
				},
				1: searchTerms{
					fullPathString: "",
					fullPathRegex:  regexp.MustCompile(`c:\\windows\\.*`),
					fileNameString: "",
					fileNameRegex:  regexp.MustCompile(`.*\.evtx`),
				},
			},
		},
		{
			name: "empty filepath string",
			args: args{exportList: FileExportList{
				0: FileExport{
					FullPath:      "",
					FullPathRegex: false,
					FileName:      "blah.exe",
					FileNameRegex: false,
				},
			}},
			wantErr:                  true,
			wantListOfSearchKeywords: nil,
		},
		{
			name: "empty filename string",
			args: args{exportList: FileExportList{
				0: FileExport{
					FullPath:      `C:\windows`,
					FullPathRegex: false,
					FileName:      "",
					FileNameRegex: false,
				},
			}},
			wantErr:                  true,
			wantListOfSearchKeywords: nil,
		},
		{
			name: "trailing slash in file path non regex",
			args: args{exportList: FileExportList{
				0: FileExport{
					FullPath:      `C:\windows\`,
					FullPathRegex: false,
					FileName:      "whoa.exe",
					FileNameRegex: false,
				},
			}},
			wantErr:                  true,
			wantListOfSearchKeywords: nil,
		},
		{
			name: "trailing slash in file path regex",
			args: args{exportList: FileExportList{
				0: FileExport{
					FullPath:      `C:\windows\`,
					FullPathRegex: true,
					FileName:      "whoa.exe",
					FileNameRegex: false,
				},
			}},
			wantErr:                  true,
			wantListOfSearchKeywords: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotListOfSearchKeywords, err := setupSearchTerms(tt.args.exportList)
			if (err != nil) != tt.wantErr {
				t.Errorf("setupSearchTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotListOfSearchKeywords, tt.wantListOfSearchKeywords) {
				t.Errorf("setupSearchTerms() gotListOfSearchKeywords = %v, want %v", gotListOfSearchKeywords, tt.wantListOfSearchKeywords)
			}
		})
	}
}
