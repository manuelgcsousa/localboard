# Localboard

Simple textbox to share text.

I needed a very simple way to share text between several devices. All I needed was a textbox. Kinda like a clipboard on your local network?

<img src="./public/example.png"/>

## Usage

- Host this on a device on your local network (you can use this simple [docker-compose](./docker-compose.yml)).
- The textbox data is saved to a text file on the server.

### WebSockets

Localboard uses WebSockets to keep the text synced across all connected devices in real-time. When one device updates the textbox, the change is instantly broadcast to all other devices with an open connection.

- The server handles incoming WebSocket connections and pushes updates to all connected clients.
- Clients automatically receive the updated text and reflect it in the textbox.

Build:
```bash
make build
```

Run:
```bash
./bin/localboard
```
