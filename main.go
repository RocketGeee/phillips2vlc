// main.go

package main

import (
    "fmt"
    "os"
    "os/exec"
    "time"

    "github.com/karalabe/hid"
)

// Constants for USB device
const (
    VendorID  = 0x0911
    ProductID = 0x1844
)

func main() {
    // Initialize USB device
    device, err := initializeUSBDevice(VendorID, ProductID)
    if err != nil {
        fmt.Println("Error initializing USB device:", err)
        return
    }
    defer device.Close()

    // Main event loop
    for {
        // Read foot pedal events
        data, err := readFootPedalEvent(device)
        if err != nil {
            fmt.Println("Error reading from USB device:", err)
            continue
        }

        // Process the foot pedal event
        processEvent(data)
    }
}

// initializeUSBDevice initializes the USB device
func initializeUSBDevice(vendorID, productID uint16) (*hid.Device, error) {
    devices := hid.Enumerate(vendorID, productID)
    if len(devices) == 0) {
        return nil, fmt.Errorf("no device found")
    }
    device, err := devices[0].Open()
    if err != nil {
        return nil, err
    }
    return device, nil
}

// readFootPedalEvent reads an event from the foot pedal
func readFootPedalEvent(device *hid.Device) ([]byte, error) {
    data := make([]byte, 8)
    _, err := device.Read(data)
    if err != nil {
        return nil, err
    }
    return data, nil
}

// processEvent processes the foot pedal event and controls VLC
func processEvent(data []byte) {
    switch data[0] {
    case 289: // Code for play button
        executeVLCCommand("play")
    // Add other button codes here as needed
    default:
        fmt.Println("Unknown event:", data)
    }
}

// executeVLCCommand sends a command to VLC
func executeVLCCommand(command string) {
    cmd := exec.Command("vlc", "--remote", command)
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error executing VLC command:", err)
    }
}
