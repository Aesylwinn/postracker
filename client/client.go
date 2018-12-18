package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    log.Print("Client starting")

    // Set up packet listener
    pc, err := net.ListenPacket("udp", ":0")
    if err != nil {
        log.Fatalf("Failed to set port\n%v", err)
    }
    defer pc.Close()

    // Get server address
    dest, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
    if err != nil {
        log.Fatalf("Failed to resolve server address\n%v", err)
    }

    // Send data
    for i := 0; i < 500; i += 1 {
        buffer := []byte(fmt.Sprint(i))
        _, err := pc.WriteTo(buffer, dest)
        if err != nil {
            log.Fatalf("Failed to send payload %v\n%v", i, err)
        }
    }
}
