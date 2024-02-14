# Data-Link Layer Frame Sender-Receiver App

## Overview

This application demonstrates the communication process between two Layer 2 (Data Link Layer) devices. It is designed as a practical exercise for those studying Linux Networking, providing hands-on experience with sending and receiving Ethernet frames.

## Prerequisites

- Docker and Docker Compose installed
- Go programming language environment set up

## Setup Instructions

### Starting the Containers

To establish a test environment with a sender and receiver device, use the provided `docker-compose` file:

```bash
docker-compose up -d
```

### Accessing the Containers

To interact with the sender container:

```bash
docker exec -it sender bash
```

To interact with the receiver container:

```bash
docker exec -it receiver bash
```

### Checking Interface Information

To view details about the Data-Link layer interfaces:

```bash
ip link
```

If an interface is down, bring it up with:

```bash
sudo ip link set eth0 up
```

To check the status of the Layer 2 interfaces:

```bash
ip link show
```

## Granting Necessary Capabilities

For the executables to operate correctly on Linux, grant them the `CAP_NET_RAW` and `CAP_NET_ADMIN` capabilities:

```bash
sudo setcap 'CAP_NET_RAW+eip CAP_NET_ADMIN+eip' /path/to/executable
```

## Building the Applications

### Sender Program

1. Compile the sender program:

```bash
go build -o sender main.go
```

2. Assign the required permissions to send frames, either by granting capabilities or running as root.

### Receiver Program

1. Compile the receiver program:

```bash
go build -o receiver main.go
```

2. Assign the required permissions to capture packets, either by granting capabilities or running as root.

## Running the Applications

After successfully building both the sender and receiver programs, execute the `docker-compose` command to start the containers. The compose file is configured to mount the local directory containing the compiled binaries into the `/app` directory within each container.

To interact with the running containers and execute the sender and receiver applications:

1. Access the receiver container and initiate the receiver application, specifying the MAC address of the interface that will listen for incoming Ethernet frames:

```bash
sudo /app/receiver 02:42:ac:11:00:03
```

2. Open a separate terminal session, access the sender container, and run the sender application, providing the source and destination MAC addresses:

```bash
sudo /app/sender 02:42:ac:11:00:02 02:42:ac:11:00:03
```

3. The sender application will send an Ethernet frame to the receiver's MAC address. The receiver application, running in its own container, is set up to capture and display the details of the received frame.

4. Verify the output on the receiver side to confirm that the Ethernet frame was received correctly:

```
Listening for frames on interface eth1 ...
Reading packet...
Source MAC: fa:8d:90:cb:b1:e7
Destination MAC: ce:2f:a6:93:4b:8b
Payload: Hi
```