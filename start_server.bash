go build ./server.go
chmod +x ./server
./server | python3 analyse.py
