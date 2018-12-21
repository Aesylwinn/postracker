package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "time"

    "github.com/Aesylwinn/postracker/common"
)


func deliverPos(pc net.PacketConn, dest net.Addr, pos common.PositionPayload) {
    buffer := pos.Encode()
    _, err := pc.WriteTo(buffer, dest)
    if err != nil {
        log.Printf("Failed to send payload\n%v", err)
    }
}

func main() {
    log.Print("Client starting")

    // Set up packet listener
    pc, err := net.ListenPacket("udp", ":0")
    if err != nil {
        log.Fatalf("Failed to set port\n%v", err)
    }
    defer pc.Close()

    // Get server address
    serverAddress, ok := os.LookupEnv("SERVER_ADDR")
    if !ok {
        serverAddress = "127.0.0.1:5000"
    }
    log.Printf("Server: %v", serverAddress)

    dest, err := net.ResolveUDPAddr("udp", serverAddress)
    if err != nil {
        log.Fatalf("Failed to resolve server address\n%v", err)
    }

    // Set up network worker
    posChan := make(chan common.PositionPayload, 1)

    go func() {
        updateInterval := 500 * time.Millisecond
        timer := time.NewTimer(updateInterval)
        lastPos := common.NewPositionPayload(0, 0)

        // Update loop
        for {
            select {
            case newPos := <- posChan:
                deliverPos(pc, dest, newPos)
                lastPos = newPos
                if !timer.Stop() {
                    <-timer.C
                }
            case <- timer.C:
                lastPos.Refresh()
                deliverPos(pc, dest, lastPos)
            }
            timer.Reset(updateInterval)
        }
    }()

    // Handle input
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        var x, y float64
        fmt.Sscanf(scanner.Text(), "%v %v", &x, &y)
        // Non-blocking write
        select {
        case posChan <- common.NewPositionPayload(x, y):
        default:
        }
    }
}
