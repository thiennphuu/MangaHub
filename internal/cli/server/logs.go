package server

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// logsCmd tails a real log file if present (~/.mangahub/logs/server.log).
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View server logs",
	Long:  `View server logs from ~/.mangahub/logs/server.log with optional filtering and following.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		follow, _ := cmd.Flags().GetBool("follow")
		level, _ := cmd.Flags().GetString("level")

		logPath, err := defaultLogPath()
		if err != nil {
			return err
		}

		file, err := ensureLogFile(logPath)
		if err != nil {
			return err
		}
		defer file.Close()

		if follow {
			fmt.Printf("Following server logs from %s (Ctrl+C to exit)...\n\n", logPath)
			return followLogs(file, level)
		}

		fmt.Printf("Recent server logs from %s:\n\n", logPath)
		return printRecentLogs(file, level, 100)
	},
}

func init() {
	logsCmd.Flags().BoolP("follow", "f", false, "Follow logs in real-time")
	logsCmd.Flags().StringP("level", "l", "", "Filter by log level (debug, info, warn, error)")
}

func defaultLogPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".mangahub", "logs", "server.log"), nil
}

// ensureLogFile guarantees the log directory exists and returns an open file.
// If the log file does not yet exist, it is created empty and a friendly
// message is shown so the command does not fail with ENOENT.
func ensureLogFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	// Open in read-only, creating if missing.
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", path, err)
	}

	// If the file is empty, let the user know instead of silently showing nothing.
	stat, err := file.Stat()
	if err == nil && stat.Size() == 0 {
		fmt.Printf("Log file is empty at %s. Start a server to generate logs.\n\n", path)
	}
	return file, nil
}

func printRecentLogs(file *os.File, level string, maxLines int) error {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if filterLevel(line, level) {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	for _, line := range lines {
		fmt.Println(line)
	}
	return nil
}

func followLogs(file *os.File, level string) error {
	// Seek to end so we only show new entries.
	if _, err := file.Seek(0, os.SEEK_END); err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// Wait for new data to be written
			time.Sleep(500 * time.Millisecond)
			continue
		}
		line = strings.TrimRight(line, "\r\n")
		if filterLevel(line, level) {
			fmt.Println(line)
		}
	}
}

func filterLevel(line, level string) bool {
	if level == "" {
		return true
	}
	level = strings.ToUpper(level)
	tag := fmt.Sprintf("[%s]", level)
	return strings.Contains(line, tag)
}
