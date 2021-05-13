package main

func main() {
	server := NewServer("192.168.1.77", 8888)

	server.Start()
}
