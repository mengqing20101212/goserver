package config

import (
	"goserver/common/logger"
	"testing"
)

func TestInitConfigManger(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "testload",
			args: args{
				path: "G:\\WORK\\me\\goserver\\conf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Init("../logs", "test1.log")
			InitConfigManger(tt.args.path)
		})
	}
}
