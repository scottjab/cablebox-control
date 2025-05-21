package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

const (
	playStatusSocket = "/home/scottjab/FieldStation42/runtime/play_status.socket"
	channelSocket    = "/home/scottjab/FieldStation42/runtime/channel.socket"
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
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/channel/up", handleChannelUp)
	http.HandleFunc("/channel/down", handleChannelDown)
	http.HandleFunc("/channel/direct", handleChannelDirect)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("home").Parse(htmlTemplate))
	status, err := getPlayStatus()
	if err != nil {
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, status)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	status, err := getPlayStatus()
	if err != nil {
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleChannelUp(w http.ResponseWriter, r *http.Request) {
	if err := sendChannelCommand("up"); err != nil {
		http.Error(w, "Error changing channel", http.StatusInternalServerError)
		return
	}
	status, err := getPlayStatus()
	if err != nil {
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleChannelDown(w http.ResponseWriter, r *http.Request) {
	if err := sendChannelCommand("down"); err != nil {
		http.Error(w, "Error changing channel", http.StatusInternalServerError)
		return
	}
	status, err := getPlayStatus()
	if err != nil {
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleChannelDirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Channel int `json:"channel"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Channel < 1 {
		http.Error(w, "Channel number must be positive", http.StatusBadRequest)
		return
	}

	cmd := ChannelCommand{
		Command: "direct",
		Channel: request.Channel,
	}

	conn, err := os.OpenFile(channelSocket, os.O_WRONLY, 0)
	if err != nil {
		http.Error(w, "Error opening channel socket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	data, err := json.Marshal(cmd)
	if err != nil {
		http.Error(w, "Error marshaling command", http.StatusInternalServerError)
		return
	}

	if _, err := conn.Write(data); err != nil {
		http.Error(w, "Error writing to channel socket", http.StatusInternalServerError)
		return
	}

	// Give the system a moment to process the command
	time.Sleep(100 * time.Millisecond)

	status, err := getPlayStatus()
	if err != nil {
		http.Error(w, "Error getting status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func getPlayStatus() (*PlayStatus, error) {
	conn, err := os.OpenFile(playStatusSocket, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening play status socket: %v", err)
	}
	defer conn.Close()

	data, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("error reading play status: %v", err)
	}

	var status PlayStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("error parsing play status: %v", err)
	}

	return &status, nil
}

func sendChannelCommand(command string) error {
	conn, err := os.OpenFile(channelSocket, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("error opening channel socket: %v", err)
	}
	defer conn.Close()

	cmd := ChannelCommand{
		Command: command,
		Channel: -1,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("error marshaling command: %v", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("error writing to channel socket: %v", err)
	}

	// Give the system a moment to process the command
	time.Sleep(100 * time.Millisecond)
	return nil
}
