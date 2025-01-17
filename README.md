# tic_tac_toe

## Tic-Tac-Toe Multiplayer Game

Description

This is a multiplayer Tic-Tac-Toe game implemented in Go. The server allows two players to connect via a TCP connection and play against each other in real-time.

## **Features**

● Connect players via TCP.

● Implements Tic-Tac-Toe game logic.

● Manages player turns.

● Checks for a winner and announces game results.

## **Installation**

1. Ensure Go version 1.20 or higher is installed.

2. Clone the repository:

`git clone https://github.com/Vova-luk/tic_tac_toe.git`

3. Navigate to the project folder:

`cd tic_tac_toe`

4. Run the server:

`go run main.go`

## Connecting Players

### On the same network (local network):

1. Ensure the server is running and use your local IP address (e.g., 192.168.x.x).

2. Find your local IP address:

Windows: `Run ipconfig` and look for IPv4 Address.

Linux/MacOS: `Run ifconfig` or `ip a`.

3. Replace line 179 with `ip := "192.168.x.x"`, where is "192.168.x.x" your local ip

4. Have the player connect via Telnet:

`telnet 192.168.x.x 8080`

### From another network (over the Internet):

1. Find your public IP address, e.g., via https://whatismyipaddress.com/ or by running:

`curl ifconfig.me`

2. Set up port forwarding on your router:

Log in to your router settings and forward port `8080` to your local machine's IP.

3. Have the other player connect using your public IP:

telnet [your public IP] 8080

License

This project is licensed under the MIT License. See the LICENSE file for details.

