package db

import (
	"goserver/table"
	"testing"
)

func TestInitDefaultDataBase(t *testing.T) {
	type args struct {
		userName  string
		passWord  string
		ip        string
		databases string
		port      int32
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				userName:  "root",
				passWord:  "root",
				ip:        "127.0.0.1",
				port:      3306,
				databases: "sysweb",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitDefaultDataBase(tt.args.userName, tt.args.passWord, tt.args.ip, tt.args.databases, tt.args.port); got != tt.want {
				t.Errorf("InitDefaultDataBase() = %v, want %v", got, tt.want)
			}
			test := table.NewSysTable()
			test.SelectAll()
		})
	}
}
