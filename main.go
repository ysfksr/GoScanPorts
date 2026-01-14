package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

type PortScanner struct {
	host       string
	timeout    time.Duration
	retryCount int
}

func NewPortScanner(host string, timeout time.Duration, retryCount int) *PortScanner {
	return &PortScanner{
		host:       host,
		timeout:    timeout,
		retryCount: retryCount,
	}
}

func (ps *PortScanner) ScanPort(port int, wg *sync.WaitGroup, openPorts chan<- int) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", ps.host, port)

	// Try multiple times to ensure reliability
	for attempt := 0; attempt <= ps.retryCount; attempt++ {
		conn, err := net.DialTimeout("tcp", address, ps.timeout)

		if err == nil {
			// Successfully connected, verify it's truly open
			if conn != nil {
				// Set a deadline to ensure the connection is real
				conn.SetDeadline(time.Now().Add(ps.timeout))

				// Try to read/write to verify the connection
				buf := make([]byte, 1)
				conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
				_, _ = conn.Read(buf)

				conn.Close()
				openPorts <- port
				return
			}
		}

		// If not the last attempt, wait a bit before retrying
		if attempt < ps.retryCount {
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (ps *PortScanner) Scan(startPort, endPort int) []int {
	var wg sync.WaitGroup
	openPorts := make(chan int, endPort-startPort+1)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go ps.ScanPort(port, &wg, openPorts)
	}

	wg.Wait()
	close(openPorts)

	var results []int
	for port := range openPorts {
		results = append(results, port)
	}

	sort.Ints(results)
	return results
}

func getPopularPorts() []int {
	return []int{
		// Web Services
		80,   // HTTP
		443,  // HTTPS
		8080, // HTTP Alternative
		8443, // HTTPS Alternative
		8000, // HTTP Development
		3000, // Node.js/React Development
		4200, // Angular Development
		5000, // Flask/Various

		// Databases
		3306,  // MySQL
		5432,  // PostgreSQL
		27017, // MongoDB
		6379,  // Redis
		1433,  // MS SQL Server
		5984,  // CouchDB
		9200,  // Elasticsearch

		// Remote Access
		21,   // FTP
		22,   // SSH
		23,   // Telnet
		3389, // RDP
		5900, // VNC

		// Email
		25,  // SMTP
		110, // POP3
		143, // IMAP
		465, // SMTPS
		587, // SMTP (submission)
		993, // IMAPS
		995, // POP3S

		// Other Services
		53,    // DNS
		67,    // DHCP
		445,   // SMB
		5672,  // RabbitMQ
		9090,  // Prometheus
		11211, // Memcached
	}
}

func (ps *PortScanner) ScanPopularPorts() []int {
	ports := getPopularPorts()
	var wg sync.WaitGroup
	openPorts := make(chan int, len(ports))

	for _, port := range ports {
		wg.Add(1)
		go ps.ScanPort(port, &wg, openPorts)
	}

	wg.Wait()
	close(openPorts)

	var results []int
	for port := range openPorts {
		results = append(results, port)
	}

	sort.Ints(results)
	return results
}

func main() {
	host := flag.String("host", "localhost", "Host to scan (e.g., localhost, 192.168.1.1)")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	timeout := flag.Int("timeout", 2000, "Connection timeout in milliseconds")
	retries := flag.Int("retries", 2, "Number of retries for each port")
	popular := flag.Bool("popular", false, "Scan only popular ports (web, databases, etc.)")
	thx := flag.Bool("thx", false, "Show credits")

	flag.Parse()

	if *thx {
		fmt.Println("You are welcome (https://github.com/ysfksr)")
		return
	}

	if *startPort < 1 || *startPort > 65535 || *endPort < 1 || *endPort > 65535 {
		fmt.Println("Error: Port range must be between 1 and 65535")
		return
	}

	if *startPort > *endPort {
		fmt.Println("Error: Start port must be less than or equal to end port")
		return
	}

	scanner := NewPortScanner(*host, time.Duration(*timeout)*time.Millisecond, *retries)
	startTime := time.Now()

	var openPorts []int

	if *popular {
		fmt.Printf("Scanning %s for popular ports (web, databases, remote access, etc.)...\n", *host)
		openPorts = scanner.ScanPopularPorts()
	} else {
		fmt.Printf("Scanning %s from port %d to %d...\n", *host, *startPort, *endPort)
		fmt.Println("This may take a while depending on the range...")
		openPorts = scanner.Scan(*startPort, *endPort)
	}

	elapsed := time.Since(startTime)

	fmt.Printf("\nScan completed in %s\n", elapsed)
	fmt.Printf("Found %d open port(s):\n", len(openPorts))

	if len(openPorts) > 0 {
		for _, port := range openPorts {
			serviceName := getServiceName(port)
			fmt.Printf("  Port %d is open %s\n", port, serviceName)
		}
	} else {
		fmt.Println("  No open ports found in the specified range")
	}
}

func getServiceName(port int) string {
	services := map[int]string{
		20:    "(FTP Data)",
		21:    "(FTP Control)",
		22:    "(SSH)",
		23:    "(Telnet)",
		25:    "(SMTP)",
		53:    "(DNS)",
		80:    "(HTTP)",
		110:   "(POP3)",
		143:   "(IMAP)",
		443:   "(HTTPS)",
		445:   "(SMB)",
		465:   "(SMTPS)",
		587:   "(SMTP Submission)",
		993:   "(IMAPS)",
		995:   "(POP3S)",
		1433:  "(MS SQL Server)",
		3000:  "(Dev Server)",
		3306:  "(MySQL)",
		3389:  "(RDP)",
		4200:  "(Angular Dev)",
		5000:  "(Flask/Various)",
		5432:  "(PostgreSQL)",
		5672:  "(RabbitMQ)",
		5900:  "(VNC)",
		5984:  "(CouchDB)",
		6379:  "(Redis)",
		8000:  "(HTTP Dev)",
		8080:  "(HTTP Alt)",
		8443:  "(HTTPS Alt)",
		9090:  "(Prometheus)",
		9200:  "(Elasticsearch)",
		11211: "(Memcached)",
		27017: "(MongoDB)",
	}

	if name, exists := services[port]; exists {
		return name
	}
	return ""
}
