package main

import (
	"context"
	"fmt"
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"sync"
    "syscall"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// BuildResult represents the build progress result
type BuildResult struct {
	Message   string `json:"message"`
	Percent   int    `json:"percent"`
	Status    string `json:"status"`
	PreStatus int    `json:"prestatus"`
}

// App struct
type App struct {
	ctx context.Context
	cancelBuild context.CancelFunc // Add this field to hold the cancel function
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectFileOptions defines the options for the SelectFile dialog
type SelectFileOptions struct {
	Title   string   `json:"title"`
	Filters []string `json:"patterns"`
}

// SelectFile opens a file dialog with the given options
func (a *App) SelectFile(options string) (string, error) {
	var opts SelectFileOptions
	err := json.Unmarshal([]byte(options), &opts)
	if err != nil {
		return "", fmt.Errorf("invalid options for SelectFile: %w", err)
	}

	var filters []runtime.FileFilter
	if len(opts.Filters) > 0 {
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Files",
			Pattern:     strings.Join(opts.Filters, ";"),
		})
	}

	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   opts.Title,
		Filters: filters,
	})
}

// SelectDirectory opens a directory dialog
func (a *App) SelectDirectory(title string) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
}


// SaveConfig saves the configuration to build.conf
func (a *App) SaveConfig(configData string) error {
	var config map[string]string
	err := json.Unmarshal([]byte(configData), &config)
	if err != nil {
		return fmt.Errorf("invalid config format: %w", err)
	}

	file, err := os.Create("build.conf")
	if err != nil {
		return fmt.Errorf("failed to create build.conf: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range config {
		// Basic escaping for quotes in value
		escapedValue := strings.ReplaceAll(value, "\"", "\"")
		line := fmt.Sprintf("%s=\"%s\"\n", key, escapedValue)		
		_, err := writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write to build.conf: %w", err)
		}
	}

	return writer.Flush()
}

// LoadConfig loads the configuration from build.conf
func (a *App) LoadConfig() (string, error) {
	file, err := os.Open("build.conf")
	if err != nil {
		if os.IsNotExist(err) {
			return "{}", nil // Return empty JSON if file doesn't exist
		}
		return "", fmt.Errorf("failed to open build.conf: %w", err)
	}
	defer file.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(file)
	// Regex to parse KEY="value" format, handling escaped quotes
	re := regexp.MustCompile(`^([^=]+)=\"(.*)\"`)

	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) == 3 {
			key := matches[1]
			value := matches[2]
			// Basic un-escaping for quotes in value
			unescapedValue := strings.ReplaceAll(value, "\\\"", "\"")
			config[key] = unescapedValue
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read build.conf: %w", err)
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to serialize config to JSON: %w", err)
	}

	return string(configJSON), nil
}


// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// BuildImage simulates the BuildImage function with progress updates
func (a *App) BuildImage(buildStream string, logStream string, adminPassword string, azureToken string) bool {
	// 清除舊的 cancelBuild
	if a.cancelBuild != nil {
		a.cancelBuild() // 取消之前的操作
		a.cancelBuild = nil
	}

	// Get the directory where the script is located
	scriptDir, err := os.Getwd()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to get current directory: %v", err))
		return false
	}

	// Path to the sample build script
	scriptPath := filepath.Join(scriptDir, "/scripts/build_script.sh")

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		runtime.LogError(a.ctx, fmt.Sprintf("Sample build script not found: %s", scriptPath))
		return false
	}

	// Create context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Store the cancel function in a field to call later
	a.cancelBuild = cancel

	// Create log file
	logFilename := "sample_build.log"
	logPath := filepath.Join(scriptDir, logFilename)
	logFile, err := os.Create(logPath)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to create log file: %v", err))
		return false
	}
	defer logFile.Close()

	// Start the build script
	cmd := exec.CommandContext(ctx, "/bin/bash", scriptPath)
	cmd.Dir = scriptDir
	// 在結束前用 mutex 保護 cmd 避免競爭
	var cmdMu sync.Mutex

	// Get stdout and stderr pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to create stdout pipe: %v", err))
		return false
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to create stderr pipe: %v", err))
		return false
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to start command: %v", err))
		return false
	}

	// Variables to track progress
	cumulativePercent := 0
	preStatus := 0 // 0: idle, 1: failed, 2: success
	endFlag := false

	// Write initial log
	writeLog(a.ctx, logStream, logPath, "[INFO] Sample BSP image build started.")

	// Use a WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2) // We have two goroutines to wait for

	// Goroutine to read stdout
	go func() {
		defer wg.Done() // Decrement counter when goroutine finishes
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				if !endFlag {
					// 當 context 被取消時，退出 goroutine
					writeLog(a.ctx, logStream, logPath, "[ERROR] Build BSP image cancelled by user")

					cmdMu.Lock()
					if cmd != nil && cmd.Process != nil {
						// 使用 syscall.Kill 來向進程發送 SIGINT 信號 - Linux
						err := syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
						if err != nil {
							writeLog(a.ctx, logStream, logPath, "[ERROR] Send SIGINT command failed")
						} else {
							writeLog(a.ctx, logStream, logPath, "[INFO] Send SIGINT command successfully")
						}
					}
					cmdMu.Unlock()

					result := BuildResult{
						Message:   "Build image cancelled by user",
						Status:    "Build Image Failed",
						Percent:   cumulativePercent,
						PreStatus: 1,
					}
					resultJSON, _ := json.Marshal(result)
					runtime.EventsEmit(a.ctx, buildStream, string(resultJSON))
				}
				return
			default:
				line := scanner.Text()
				writeLog(a.ctx, logStream, logPath, line)

				// Process the build log and update progress
				result := processBuildLog(line, cumulativePercent, preStatus)
				cumulativePercent = result.Percent
				preStatus = result.PreStatus

				// Emit progress to frontend
				resultJSON, _ := json.Marshal(result)
				runtime.EventsEmit(a.ctx, buildStream, string(resultJSON))
			}
		}
	}()

	// Goroutine to read stderr
	go func() {
		defer wg.Done() // Decrement counter when goroutine finishes
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				line := scanner.Text()
				writeLog(a.ctx, logStream, logPath, line)

				// Process error output
				result := processBuildLog(line, cumulativePercent, preStatus)
				cumulativePercent = result.Percent
				preStatus = result.PreStatus

				// Emit progress to frontend
				resultJSON, _ := json.Marshal(result)
				runtime.EventsEmit(a.ctx, buildStream, string(resultJSON))
			}
		}
	}()

	// 先等 goroutine 處理 stdout/stderr
	wg.Wait()

	// Wait for command to complete
	err = cmd.Wait()
	if err != nil {
		writeLog(a.ctx, logStream, logPath, fmt.Sprintf("[ERROR] Build failed: %v", err))
		result := BuildResult{
			Message:   "Build failed",
			Status:    "Build Image Failed",
			Percent:   cumulativePercent,
			PreStatus: 1,
		}
		resultJSON, _ := json.Marshal(result)
		runtime.EventsEmit(a.ctx, buildStream, string(resultJSON))
		return false
	}

	endFlag = true
	
	// Check final status
	if preStatus != 2 && preStatus != 1 {
		writeLog(a.ctx, logStream, logPath, "[ERROR] Build ended unexpectedly!")
		result := BuildResult{
			Message:   "Build failed",
			Status:    "Build Image Failed",
			Percent:   cumulativePercent,
			PreStatus: 1,
		}
		resultJSON, _ := json.Marshal(result)
		runtime.EventsEmit(a.ctx, buildStream, string(resultJSON))
		return false
	}

	return true
}

// CancelBuild cancels the sample build process
func (a *App) CancelBuild() bool {
  	if a.cancelBuild != nil {
        a.cancelBuild() // Call the cancel function to stop the download
		a.cancelBuild = nil
    }
	runtime.LogInfo(a.ctx, "Sample build cancelled by user")
	return true
}


// processBuildLog processes build log output and extracts progress information
func processBuildLog(input string, cumulativePercent int, preStatus int) BuildResult {
	var result BuildResult
	result.Percent = cumulativePercent
	result.Message = input
	result.Status = "Building"
	result.PreStatus = preStatus

	// Check for progress percentage
	re := regexp.MustCompile(`BSP Build Progress:\s*(\d+)%`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		progress, _ := strconv.ParseInt(matches[1], 10, 32)
		cumulativePercent = int(progress)
		result.Percent = int(progress)
	}

	// Check for success
	if strings.Contains(input, "[SUCCESS]") {
		result.Status = "Build Image Completed"
		result.Percent = 100
		result.PreStatus = 2
	}

	// Check for error
	if strings.Contains(input, "[ERROR]") {
		result.Status = "Build Image Failed"
		result.Percent = 100
		result.PreStatus = 1
	}

	if result.PreStatus == 2 {
		result.Status = "Build Image Completed"
	} else if result.PreStatus == 1 {
		result.Status = "Build Image Failed"
	}

	return result
}

// writeLog writes a log message to file and emits it to frontend
func writeLog(ctx context.Context, logStream string, logPath string, msg string) {
	// Write to file
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Write log file error: %v", err))
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, msg)
	_, _ = f.WriteString(logLine)

	// Emit to frontend
	runtime.EventsEmit(ctx, logStream, logLine)
}

// GetSampleBuildLog reads the sample build log file
func (a *App) GetSampleBuildLog() string {
	scriptDir, err := os.Getwd()
	if err != nil {
		return "Error: Failed to get current directory"
	}

	logPath := filepath.Join(scriptDir, "sample_build.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		return "Error: Failed to read log file"
	}

	return string(content)
}