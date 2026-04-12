# 🛡️ PRYVO - Encrypted P2P Desktop Chat

Welcome to **PRYVO**, a modern, decentralized, and fully private desktop chat application built with **Go** and **Vanilla Javascript**. 

Unlike traditional messaging apps, PRYVO does not store your messages on a central server. It establishes a direct **WebRTC Peer-to-Peer (P2P)** encrypted tunnel between you and your friend.

## ✨ Key Features
* **Zero Central Storage:** Your profile and chat history are stored locally in a safe SQLite database on your machine.
* **P2P Encrypted Tunnels:** Messages travel directly from IP to IP using WebRTC Data Channels.
* **Image Compression Engine:** Built-in HTML5 Canvas engine to compress profile pictures locally before P2P transmission.
* **Modern UI:** Responsive and modern chat interface built without heavy frontend frameworks.
* **Cross-Platform:** Compiles to single executable files for Windows, Linux, and macOS using Wails.

## 🛠️ Technology Stack
* **Backend Engine:** Go (Golang) + SQLite
* **Frontend:** HTML5, CSS3, Vanilla Javascript
* **Bridge & Compilation:** Wails v2
* **Networking:** WebSockets (Signaling) + WebRTC (P2P Data Channels)

## 🚀 How to Run Locally
### Prerequisites
1. Install [Go](https://golang.org/dl/)
2. Install [Node.js](https://nodejs.org/en/download/)
3. Install [Wails](https://wails.io/docs/gettingstarted/installation)

### Live Development
Clone this repository and run the Wails dev server:
\`\`\`bash
git clone https://github.com/your-username/pryvo.git
cd pryvo/pryvo-desktop
wails dev
\`\`\`

### Build Executable
To compile PRYVO into a standalone desktop application (`.exe` for Windows, bin for Linux, or `.app` for macOS):
\`\`\`bash
# Build for your current system
wails build

# Cross-compile for Windows (from Linux/Mac)
wails build -platform windows/amd64
\`\`\`
The compiled application will be available in the `build/bin/` directory.

---
*Built with ❤️ and privacy in mind by Guilherme Marques*

## License

GNU GENERAL PUBLIC LICENSE License © 2026 Guilherme Marques.
