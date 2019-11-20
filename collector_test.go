package windowscollector

import "testing"

func TestCollect(t *testing.T) {
	type args struct {
		exportList   ListOfFilesToExport
		resultWriter ResultWriter
		handler      Handler
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Collect(tt.args.handler, tt.args.exportList, tt.args.resultWriter); (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getFiles(t *testing.T) {
	type args struct {
		volumeHandler        *VolumeHandler
		resultWriter         ResultWriter
		listOfSearchKeywords listOfSearchTerms
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getFiles(tt.args.volumeHandler, tt.args.resultWriter, tt.args.listOfSearchKeywords); (err != nil) != tt.wantErr {
				t.Errorf("getFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
