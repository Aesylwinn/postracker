package main

import (
    "log"
    "net"
    "github.com/Aesylwinn/postracker/common"
)

func handlePacket(addr net.Addr, buffer []byte) {
    pos := common.DecodePosition(buffer)
    log.Print("Message:", pos)
}

func main() {
    log.Print("Server starting")

    // Set up packet listener
    pc, err := net.ListenPacket("udp", ":5000")
    if err != nil {
        log.Fatalf("Failed to set port\n%v", err)
    }
    defer pc.Close()

    // Listening loop
    for {
        // Read a packet
        buffer := make([]byte, 512)
        num, addr, err := pc.ReadFrom(buffer)
        if err != nil {
            log.Printf("Error reading packet\n%v", err)
        } else {
            handlePacket(addr, buffer[:num])
        }
    }
}
