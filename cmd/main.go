package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"db-backup-cli/internal/engine"
	"db-backup-cli/internal/notifier"

	"github.com/spf13/cobra"
)

var (
	dbType       string
	host         string
	port         int
	user         string
	password     string
	dbName       string
	outputDir    string
	slackWebhook string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "db-backup",
		Short: "CLI Utility untuk Backup Database",
	}

	var backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Menjalankan proses backup",
		Run:   runBackup,
	}

	backupCmd.Flags().StringVarP(&dbType, "type", "t", "postgres", "Tipe Database (postgres)")
	backupCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host Database")
	backupCmd.Flags().IntVarP(&port, "port", "P", 5432, "Port Database")
	backupCmd.Flags().StringVarP(&user, "user", "u", "postgres", "Username Database")
	backupCmd.Flags().StringVarP(&password, "password", "p", "", "Password Database")
	backupCmd.Flags().StringVarP(&dbName, "db-name", "d", "", "Nama Database Target")
	backupCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "./backups", "Lokasi Penyimpanan")
	backupCmd.Flags().StringVar(&slackWebhook, "slack-webhook", "", "Slack Webhook URL")

	_ = backupCmd.MarkFlagRequired("db-name")

	rootCmd.AddCommand(backupCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runBackup(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	slack := &notifier.SlackNotifier{WebhookURL: slackWebhook}
	cfg := engine.Config{Host: host, Port: port, User: user, Password: password, DBName: dbName}

	var dbEngine engine.DBEngine
	switch dbType {
	case "postgres":
		dbEngine = engine.NewPostgresEngine(cfg)
	default:
		log.Fatalf("DBMS %s belum didukung", dbType)
	}

	log.Println("Menguji koneksi ke database...")
	if err := dbEngine.TestConnection(ctx); err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	_ = os.MkdirAll(outputDir, os.ModePerm)
	fileName := fmt.Sprintf("%s_%d.dump", dbName, time.Now().Unix())
	outputPath := filepath.Join(outputDir, fileName)

	log.Println("Membuat backup...")
	if err := dbEngine.Dump(ctx, outputPath); err != nil {
		_ = slack.Send(fmt.Sprintf("Backup %s GAGAL: %v", dbName, err))
		log.Fatalf("Backup gagal: %v", err)
	}

	successMsg := fmt.Sprintf("Backup %s SUKSES! File: %s", dbName, outputPath)
	log.Println(successMsg)
	_ = slack.Send(successMsg)
}
