# Ethersim

Ethersim is a simulator for the original ethernet protocol proposed by
MetCalfe and Boggs in 1976. The protocol was originally developed for 
a *broadcast* network and delivered messages with *high probability*. 
The simulation runs in discrete ticks and breaks the simulation down 
into three main components: *devices, transceivers, and network edges.*


https://github.com/user-attachments/assets/5b1373b8-6342-446c-9890-c3e65397e918


## Functionality

- *3.1* Networks have the topology of an unrooted tree
- *3.2* Messages may collide with other messages and transceivers cooperate
- *3.3* Messages are sent from a source device to a destination device
- *3.4* Messages are delivered with high probability

- *3.5.1* Transceivers detect when messages are being sent and defer sending their own
- *3.5.2* Transceivers detect when their transmitted messages are interfering with others and back off
- *3.5.3* Transceivers only forward complete messages to their devices
- ~~*3.5.4* Messages are sent with checksums for error correction~~
- *3.5.5* Transceivers jam the network to gaurantee consensus

## Running the Simulator

Locally:

```sh
~/ethersim> $ go run .
```

As a web application:

```sh
~/ethersim> $ go run github.com/hajimehoshi/wasmserve@latest .

```


