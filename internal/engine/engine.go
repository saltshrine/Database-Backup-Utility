package engine

import "context"

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type DBEngine interface {
	TestConnection(ctx context.Context) error
	Dump(ctx context.Context, filePath string) error
}
