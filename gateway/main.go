package gateway

func main() {
	handler := NewHandler()
	handler.registerRoutes()
}
