package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func expandRange(startIP, endIP string) ([]string, error) {
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	if start == nil || end == nil {
		return nil, fmt.Errorf("Invalid IP: %s - %s", startIP, endIP)
	}

	startParts := strings.Split(startIP, ".")
	endParts := strings.Split(endIP, ".")

	if len(startParts) != 4 || len(endParts) != 4 {
		return nil, fmt.Errorf("Invalid Format: %s - %s", startIP, endIP)
	}

	base := strings.Join(startParts[:3], ".")
	startOctet, _ := strconv.Atoi(startParts[3])
	endOctet, _ := strconv.Atoi(endParts[3])

	if startOctet > endOctet {
		return nil, fmt.Errorf("inicial octed bigger than the last: %s - %s", startIP, endIP)
	}

	var ips []string
	for i := startOctet; i <= endOctet; i++ {
		ips = append(ips, fmt.Sprintf("%s.%d", base, i))
	}

	return ips, nil
}

func main() {
	inputFile := flag.String("i", "", "Input file containing IP ranges in format x.x.x.1-x.x.x.2 or x.x.x.1 x.x.x.2, one per line")
	outputFile := flag.String("o", "", "Output file to store all the IP addresses expanded")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Use: go run main.go -i input.txt -o output.txt")
		return
	}

	in, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de entrada:", err)
		return
	}
	defer in.Close()

	out, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Erro ao criar arquivo de saída:", err)
		return
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var startIP, endIP string

		// Suporta "start end" ou "start-end"
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				fmt.Printf("Linha inválida (esperado start-end): %s\n", line)
				continue
			}
			startIP = strings.TrimSpace(parts[0])
			endIP = strings.TrimSpace(parts[1])
		} else {
			parts := strings.Fields(line)
			if len(parts) != 2 {
				fmt.Printf("Linha inválida (esperado start end): %s\n", line)
				continue
			}
			startIP = parts[0]
			endIP = parts[1]
		}

		ips, err := expandRange(startIP, endIP)
		if err != nil {
			fmt.Println("Erro:", err)
			continue
		}

		for _, ip := range ips {
			writer.WriteString(ip + "\n")
		}
	}

	writer.Flush()
}
