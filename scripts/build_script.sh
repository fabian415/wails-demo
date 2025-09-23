#!/bin/bash

# Sample Build Script for BSP Image Building
# This script simulates the build process with progress updates

echo "Starting BSP Image Build Process..."

# Function to simulate progress
simulate_progress() {
    local total_steps=10
    local current_step=0
    
    while [ $current_step -le $total_steps ]; do
        local percentage=$((current_step * 10))
        
        # Simulate different build stages
        case $current_step in
            0)
                echo "[INFO] Initializing build environment..."
                ;;
            1)
                echo "[INFO] Downloading dependencies..."
                ;;
            2)
                echo "[INFO] Configuring build parameters..."
                ;;
            3)
                echo "[INFO] Compiling kernel modules..."
                ;;
            4)
                echo "[INFO] Building device tree..."
                ;;
            5)
                echo "[INFO] Compiling user space applications..."
                ;;
            6)
                echo "[INFO] Creating root filesystem..."
                ;;
            7)
                echo "[INFO] Packaging BSP image..."
                ;;
            8)
                echo "[INFO] Generating checksums..."
                ;;
            9)
                echo "[INFO] Finalizing build..."
                ;;
            10)
                echo "[SUCCESS] BSP Build Progress: 100%"
                echo "[SUCCESS] Script build completed successfully!"
                echo "Final image generated at: /tmp/sample_bsp_image.img"
                break
                ;;
        esac
        
        # Output progress in the format expected by the Go code
        if [ $current_step -lt $total_steps ]; then
            echo "BSP Build Progress: ${percentage}%"
        fi
        
        # Simulate work time
        sleep 2
        
        current_step=$((current_step + 1))
    done
}

# Check if we should simulate an error
if [ "$1" = "error" ]; then
    echo "[INFO] Starting build with error simulation..."
    echo "BSP Build Progress: 30%"
    sleep 2
    echo "BSP Build Progress: 60%"
    sleep 2
    echo "[ERROR] Build failed at step 6: Compilation error in kernel module"
    exit 1
fi

# Run the simulation
simulate_progress

echo "Build script completed."
