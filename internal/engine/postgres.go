package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type PostgresEngine struct {
	cfg Config
}

func NewPostgresEngine(cfg Config) *PostgresEngine {
	if cfg.Port == 0 {
		cfg.Port = 5432
	}

	return &PostgresEngine{cfg: cfg}
}

func (p *PostgresEngine) TestConnection(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "pg_isready",
		"-h", p.cfg.Host,
		"-p", strconv.Itoa(p.cfg.Port),
		"-U", p.cfg.User,
		"-d", p.cfg.DBName,
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.cfg.Password))
	return cmd.Run()
}

func (p *PostgresEngine) Dump(ctx context.Context, outputPath string) error {
	cmd := exec.CommandContext(ctx, "pg_dump",
		"-h", p.cfg.Host,
		"-p", strconv.Itoa(p.cfg.Port),
		"-U", p.cfg.User,
		"-F", "c",
		"-b",
		"-f", outputPath,
		p.cfg.DBName,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.cfg.Password))
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *PostgresEngine) Restore(ctx context.Context, inputPath string, tables []string) error {
	args := []string{
		"-h", p.cfg.Host,
		"-p", strconv.Itoa(p.cfg.Port),
		"-U", p.cfg.User,
		"-d", p.cfg.DBName,
		"-v",
	}

	for _, table := range tables {
		args = append(args, "-t", table)
	}
	args = append(args, inputPath)

	cmd := exec.CommandContext(ctx, "pg_restore", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.cfg.Password))
	cmd.Stderr = os.Stderr

	return cmd.Run()

}
