
# Layer 2 Frame Sender-Receiver App

This application shows the way how the communication between two layer 2 devices is handled. This application is intended to be used as a practice exercise for the Linux Networking study.

To begin with you should create a new virtual network interface `eth0:1` (or any name you prefer) using the following command:

```bash
sudo ip link add link eth0 eth0:1 type macvlan mode bridge
```

The `macvlan` mode is used to create a virtual interface. The `bridge` mode is most common as it allows the virtual interface to communicate with the physical interface.

Check the newly created virtual interface:

```bash
ip addr show eth0:1
```

Once you have created the virtual interface, you can configure its IP address and network parameters as needed, and bring it up with:

```bash
sudo ip addr add 192.168.1.2/24 dev eth0:1
sudo ip link set eth0:1 up
```

Check the status of the layer 2 interface:

```bash
ip link show
```

To place the two interfaces (`eth0:1` and `eth0:2`) under a Layer 2 LAN, you can create a bridge interface and attach the virtual interfaces to it. Here's how you can achieve this:

```bash
sudo ip link add br0 type bridge
```

Attach the virtual interfaces (`eth0:1` and `eth0:2`) to the bridge:

```bash
sudo ip link set eth0:1 master br0
sudo ip link set eth0:2 master br0
```

Bring up the bridge interface and the virtual interfaces:

```bash
sudo ip link set br0 up
sudo ip link set eth0:1 up
sudo ip link set eth0:2 up
```

By following these steps, you will have both virtual interfaces (`eth0:1` and `eth0:2`) placed under the bridge interface (`br0`), effectively creating a Layer 2 LAN.

---

**Grant Capabilities (preferred method on Linux):**

Use the `setcap` command to grant the `CAP_NET_RAW` and `CAP_NET_ADMIN` capabilities to the executable binary.

```bash
sudo setcap 'CAP_NET_RAW+eip CAP_NET_ADMIN+eip' /path/to/executable

To send an Ethernet frame from a sender to a receiver on the same machine using virtual network interfaces, you need to build two separate Go programs: one for sending frames (sender) and one for receiving frames (receiver).

**Building the Sender Program:**

1. Build the sender program using the `go build` command, specifying the output binary name if desired:
   ```bash
   go build -o sender main.go
   ```
2Ensure the sender program has the necessary permissions to send frames by granting capabilities or running as root.

**Building the Receiver Program:**

1. Build the receiver program using the `go build` command, specifying the output binary name if desired:
   ```bash
   go build -o receiver main.go
   ```
2. Ensure the receiver program has the necessary permissions to capture packets by granting capabilities or running as root.

**Running the Programs:**

1. Set up the virtual network interfaces and bridge as described in the previous steps.
2. Run the receiver program, specifying the virtual interface that will listen for frames:
   ```bash
   sudo ./receiver eth0:1
   ```
3. In a separate terminal, run the sender program, specifying the source virtual interface and the destination interface:
   ```bash
   sudo ./sender eth0:2 eth0:1
   ```
4. The sender program will send an Ethernet frame to the specified destination MAC address, which should be received by the receiver program listening on the corresponding virtual interface.
