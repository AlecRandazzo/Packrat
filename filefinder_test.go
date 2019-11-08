package windowscollector

import (
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	"reflect"
	"testing"
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
		// TODO: Add test cases.
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
		wantErr            bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFoundFilesList, err := confirmFoundFiles(tt.args.listOfSearchKeywords, tt.args.listOfPossibleMatches, tt.args.directoryTree)
			if (err != nil) != tt.wantErr {
				t.Errorf("confirmFoundFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
