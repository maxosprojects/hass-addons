package main

import "testing"

func Test_displayData(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "display data",
		},
	}
	baseDir = "./test-data/"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configs, err := getData(getFilename(configsFile), &ConfigEntries{})
			noErr(err)
			tplinkEntries := filterTplinkDevices(configs)
			displayData(tplinkEntries)
		})
	}
}

func Test_runCommand(t *testing.T) {
	type args struct {
		cmdString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "reboot",
			args: args{
				cmdString: "{\"cmd\":\"reboot\"}",
			},
			wantErr: false,
		},
	}
	baseDir = "./test-data/"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCommand(tt.args.cmdString); (err != nil) != tt.wantErr {
				t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
