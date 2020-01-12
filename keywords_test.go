// Copyright (c) 2020 Alec Randazzo

package windowscollector

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_setupSearchTerms(t *testing.T) {
	type args struct {
		exportList ListOfFilesToExport
	}
	tests := []struct {
		name                     string
		args                     args
		wantListOfSearchKeywords listOfSearchTerms
		wantErr                  bool
	}{
		{
			name: "this should be successful",
			args: args{exportList: ListOfFilesToExport{
				0: FileToExport{
					FullPath:        `C:\windows`,
					IsFullPathRegex: false,
					FileName:        "fake.exe",
					IsFileNameRegex: false,
				},
				1: FileToExport{
					FullPath:        `C:\\windows\\.*`,
					IsFullPathRegex: true,
					FileName:        `.*\.evtx`,
					IsFileNameRegex: true,
				},
			}},
			wantErr: false,
			wantListOfSearchKeywords: listOfSearchTerms{
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
			args: args{exportList: ListOfFilesToExport{
				0: FileToExport{
					FullPath:        "",
					IsFullPathRegex: false,
					FileName:        "blah.exe",
					IsFileNameRegex: false,
				},
			}},
			wantErr:                  true,
			wantListOfSearchKeywords: nil,
		},
		{
			name: "empty filename string",
			args: args{exportList: ListOfFilesToExport{
				0: FileToExport{
					FullPath:        `C:\windows`,
					IsFullPathRegex: false,
					FileName:        "",
					IsFileNameRegex: false,
				},
			}},
			wantErr:                  true,
			wantListOfSearchKeywords: nil,
		},
		{
			name: "trailing slash in file path non regex",
			args: args{exportList: ListOfFilesToExport{
				0: FileToExport{
					FullPath:        `C:\windows\`,
					IsFullPathRegex: false,
					FileName:        "whoa.exe",
					IsFileNameRegex: false,
				},
			}},
			wantErr:                  true,
			wantListOfSearchKeywords: nil,
		},
		{
			name: "trailing slash in file path regex",
			args: args{exportList: ListOfFilesToExport{
				0: FileToExport{
					FullPath:        `C:\windows\`,
					IsFullPathRegex: true,
					FileName:        "whoa.exe",
					IsFileNameRegex: false,
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
