package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type PlayStatus struct {
	Status        string `json:"status"`
	NetworkName   string `json:"network_name"`
	ChannelNumber int    `json:"channel_number"`
	Timestamp     string `json:"timestamp"`
	Title         string `json:"title"`
}

type ChannelCommand struct {
	Command string `json:"command"`
	Channel int    `json:"channel"`
}

var (
	playStatusSocket string
	channelSocket    string
	listenAddr       string
)

var htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Cable Box Control</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f0f0f0;
        }
        .container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .status {
            margin: 20px 0;
            padding: 15px;
            background-color: #f8f9fa;
            border-radius: 4px;
        }
        .button-group {
            display: flex;
            gap: 10px;
            margin: 20px 0;
        }
        .direct-channel {
            display: flex;
            gap: 10px;
            margin: 20px 0;
            align-items: center;
        }
        input[type="number"] {
            padding: 10px;
            font-size: 16px;
            width: 80px;
            border: 1px solid #ccc;
            border-radius: 4px;
        }
        button {
            padding: 10px 20px;
            font-size: 16px;
            cursor: pointer;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            transition: background-color 0.2s;
        }
        button:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Cable Box Control</h1>
        <div class="status">
            <h2>Current Status</h2>
            <p><strong>Status:</strong> <span id="status">{{.Status}}</span></p>
            <p><strong>Network:</strong> <span id="network">{{.NetworkName}}</span></p>
            <p><strong>Channel:</strong> <span id="channel">{{.ChannelNumber}}</span></p>
            <p><strong>Title:</strong> <span id="title">{{.Title}}</span></p>
        </div>
        <div class="button-group">
            <button onclick="changeChannel('up')">Channel Up</button>
            <button onclick="changeChannel('down')">Channel Down</button>
        </div>
        <div class="direct-channel">
            <input type="number" id="channelNumber" min="1" placeholder="Channel #">
            <button onclick="setDirectChannel()">Go to Channel</button>
        </div>
    </div>

    <script>
        function changeChannel(direction) {
            fetch('/channel/' + direction, {
                method: 'POST'
            })
            .then(response => response.json())
            .then(data => {
                updateStatus(data);
            })
            .catch(error => console.error('Error:', error));
        }

        function setDirectChannel() {
            const channelNumber = document.getElementById('channelNumber').value;
            if (!channelNumber) {
                alert('Please enter a channel number');
                return;
            }
            
            fetch('/channel/direct', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ channel: parseInt(channelNumber) })
            })
            .then(response => response.json())
            .then(data => {
                updateStatus(data);
                document.getElementById('channelNumber').value = '';
            })
            .catch(error => console.error('Error:', error));
        }

        function updateStatus(data) {
            document.getElementById('status').textContent = data.status;
            document.getElementById('network').textContent = data.network_name;
            document.getElementById('channel').textContent = data.channel_number;
            document.getElementById('title').textContent = data.title;
        }

        // Poll for status updates every 2 seconds
        setInterval(() => {
            fetch('/status')
                .then(response => response.json())
                .then(data => {
                    updateStatus(data);
                })
                .catch(error => console.error('Error:', error));
        }, 2000);
    </script>
</body>
</html>
`

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Println("Starting Cable Box Control application")

	// Define flags
	flag.StringVar(&playStatusSocket, "status-socket", "FieldStation42/runtime/play_status.socket", "Path to the play status socket")
	flag.StringVar(&channelSocket, "channel-socket", "FieldStation42/runtime/channel.socket", "Path to the channel socket")
	flag.StringVar(&listenAddr, "listen", ":8080", "Address to listen on (e.g. ':8080' or '127.0.0.1:8080')")
	flag.Parse()

	log.Printf("Configuration loaded - Status socket: %s, Channel socket: %s, Listen address: %s",
		playStatusSocket, channelSocket, listenAddr)

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/channel/up", handleChannelUp)
	http.HandleFunc("/channel/down", handleChannelDown)
	http.HandleFunc("/channel/direct", handleChannelDirect)

	log.Printf("Server starting on %s...", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for home page from %s", r.RemoteAddr)
	tmpl := template.Must(template.New("home").Parse(htmlTemplate))
	status, err := getPlayStatus()
	if err != nil {
		log.Printf("Error getting status for home page: %v", err)
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	log.Printf("Rendering home page with status: %+v", status)
	tmpl.Execute(w, status)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received status request from %s", r.RemoteAddr)
	status, err := getPlayStatus()
	if err != nil {
		log.Printf("Error getting status: %v", err)
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	log.Printf("Sending status response: %+v", status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleChannelUp(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received channel up request from %s", r.RemoteAddr)
	if err := sendChannelCommand("up"); err != nil {
		log.Printf("Error sending channel up command: %v", err)
		http.Error(w, "Error changing channel", http.StatusInternalServerError)
		return
	}
	status, err := getPlayStatus()
	if err != nil {
		log.Printf("Error getting status after channel up: %v", err)
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	log.Printf("Channel up successful, new status: %+v", status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleChannelDown(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received channel down request from %s", r.RemoteAddr)
	if err := sendChannelCommand("down"); err != nil {
		log.Printf("Error sending channel down command: %v", err)
		http.Error(w, "Error changing channel", http.StatusInternalServerError)
		return
	}
	status, err := getPlayStatus()
	if err != nil {
		log.Printf("Error getting status after channel down: %v", err)
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	log.Printf("Channel down successful, new status: %+v", status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleChannelDirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("Invalid method %s for direct channel request from %s", r.Method, r.RemoteAddr)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Channel int `json:"channel"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding direct channel request from %s: %v", r.RemoteAddr, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received direct channel request for channel %d from %s", request.Channel, r.RemoteAddr)

	if request.Channel < 1 {
		log.Printf("Invalid channel number %d requested from %s", request.Channel, r.RemoteAddr)
		http.Error(w, "Channel number must be positive", http.StatusBadRequest)
		return
	}

	if err := sendChannelCommand(fmt.Sprintf("direct %d", request.Channel)); err != nil {
		log.Printf("Error sending direct channel command: %v", err)
		http.Error(w, "Error changing channel", http.StatusInternalServerError)
		return
	}

	status, err := getPlayStatus()
	if err != nil {
		log.Printf("Error getting status after direct channel change: %v", err)
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	log.Printf("Direct channel change successful, new status: %+v", status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func getPlayStatus() (*PlayStatus, error) {
	log.Printf("Reading play status from file: %s", playStatusSocket)
	data, err := os.ReadFile(playStatusSocket)
	if err != nil {
		log.Printf("Failed to read play status file: %v", err)
		return nil, fmt.Errorf("failed to read play status file: %v", err)
	}

	var status PlayStatus
	if err := json.Unmarshal(data, &status); err != nil {
		log.Printf("Failed to parse play status: %v", err)
		return nil, fmt.Errorf("failed to parse play status: %v", err)
	}

	log.Printf("Successfully retrieved play status: %+v", status)
	return &status, nil
}

func sendChannelCommand(command string) error {
	log.Printf("Attempting to send channel command: %s", command)

	cmd := ChannelCommand{
		Command: command,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Failed to marshal channel command: %v", err)
		return fmt.Errorf("failed to marshal channel command: %v", err)
	}

	if err := os.WriteFile(channelSocket, data, 0666); err != nil {
		log.Printf("Failed to write to channel file: %v", err)
		return fmt.Errorf("failed to write to channel file: %v", err)
	}

	log.Printf("Successfully sent channel command: %s", command)
	return nil
}
