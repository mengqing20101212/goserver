package db

import (
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestDBManger_Execute(t *testing.T) {
	type fields struct {
		dbUrl       string
		db          *sqlx.DB
		connectFlag bool
	}
	type args struct {
		opt TableInterface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &DBManger{
				dbUrl:       tt.fields.dbUrl,
				db:          tt.fields.db,
				connectFlag: tt.fields.connectFlag,
			}
			if got := self.Execute(tt.args.opt); got != tt.want {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBManger_ExecuteSql(t *testing.T) {
	type fields struct {
		dbUrl       string
		db          *sqlx.DB
		connectFlag bool
	}
	type args struct {
		sql string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &DBManger{
				dbUrl:       tt.fields.dbUrl,
				db:          tt.fields.db,
				connectFlag: tt.fields.connectFlag,
			}
			if got := self.ExecuteSql(tt.args.sql); got != tt.want {
				t.Errorf("ExecuteSql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitDataBase(t *testing.T) {
	type args struct {
		manger    *DBManger
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
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitDataBase(tt.args.manger, tt.args.userName, tt.args.passWord, tt.args.ip, tt.args.databases, tt.args.port); got != tt.want {
				t.Errorf("InitDataBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				databases: "mysql",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitDefaultDataBase(tt.args.userName, tt.args.passWord, tt.args.ip, tt.args.databases, tt.args.port); got != tt.want {
				t.Errorf("InitDefaultDataBase() = %v, want %v", got, tt.want)
			}
		})
	}
}
