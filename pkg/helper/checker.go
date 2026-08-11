package helper

func IsLocalIP(addr string) bool {
	return len(addr) >= 9 && (addr[:9] == "localhost" || addr[:9] == "127.0.0.1")
}
