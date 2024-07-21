package main

// This generates a simple MOTD using Lipgloss by Charm. The colors are based on the hostname
// Porco 2024

import (
	"fmt"
	"os"
	"os/exec"
	"net"
	"strings"
	"strconv"
	"runtime"
	"crypto/md5"
	"encoding/binary"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/term"
)


// Style definitions.
var (

	// General.
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	)


/////////////////////////// Helper Functions 

func hashToFloat(input string) float64 {
	hash := md5.Sum([]byte(input))
	hashInt := binary.BigEndian.Uint64(hash[:])
	hashFloat := float64(hashInt % 360)
	return hashFloat + 1
}

/////////////////////////// MAIN FUNCTION 

func main() {
	
	// Get OS information
	osInfo := runtime.GOOS
	
	
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	
	// Get IP address
	ipAddr := ""
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ipAddr = ipnet.IP.String()
					break
				}
			}
		}
	}
	
	
	// Get OS version number and codename
	var osVersion, osCodename string
	if osInfo == "darwin" { // macOS
	cmd := exec.Command("sw_vers", "-productVersion")
	output, _ := cmd.Output()
	osVersion = strings.TrimSpace(string(output))
	
	cmd = exec.Command("sw_vers", "-productName")
	output, _ = cmd.Output()
	// productName := strings.TrimSpace(string(output))
	
	switch osVersion {
	case "10.10":
		osCodename = "Yosemite"
	case "10.11":
		osCodename = "El Capitan"
	case "10.12":
		osCodename = "Sierra"
	case "10.13":
		osCodename = "High Sierra"
	case "10.14":
		osCodename = "Mojave"
	case "10.15":
		osCodename = "Catalina"
	case "11.0":
		osCodename = "Big Sur"
	case "12.0":
		osCodename = "Monterey"
	case "13.0":
		osCodename = "Ventura"
	case "14.0":
		osCodename = "Sonoma"	
	default:
		osCodename = "Unknown"
	}
	
	osVersion = osVersion + " (" + osCodename + ")"
	} else { // Linux
		cmd := exec.Command("lsb_release", "-d")
		output, _ := cmd.Output()
		osVersion = strings.TrimSpace(strings.SplitN(string(output), ":", 2)[1])
	}

	// Get disk space information
	var totalSpace, freeSpace, usedPercent string
	if osInfo == "darwin" { // macOS
		cmd := exec.Command("df", "-H", "/")
		output, _ := cmd.Output()
		diskInfo := strings.Fields(strings.Split(string(output), "\n")[1])
		totalSpace = diskInfo[1]
		freeSpace = diskInfo[3]
		usedPercent = diskInfo[4]
	} else { // Linux
		cmd := exec.Command("df", "-H", "--output=size,avail,pcent", "/")
		output, _ := cmd.Output()
		diskInfo := strings.Fields(strings.Split(string(output), "\n")[1])
		totalSpace = diskInfo[0]
		freeSpace = diskInfo[1]
		usedPercent = diskInfo[2]
	}
	
	var uptime string
	if osInfo == "darwin" { // macOS
		cmd := exec.Command("uptime")
		output, _ := cmd.Output()
		uptime = strings.TrimSpace(strings.SplitN(string(output), ",", 2)[0])
	} else { // Linux
		cmd := exec.Command("uptime", "-p")
		output, _ := cmd.Output()
		uptime = strings.TrimSpace(string(output))
	}
	
	// Get CPU usage
	cmd := exec.Command("ps", "-A", "-o", "%cpu")
	output, _ := cmd.Output()
	cpuInfo := string(output)
	cpuUsage := strings.Split(cpuInfo, "\n")[1] // Assuming the first line contains the CPU usage

	// Get RAM usage
	var ramUsage string
	if osInfo == "darwin" { // macOS
		cmd := exec.Command("top", "-l", "1", "-s", "0", "-n", "0")
		output, _ := cmd.Output()
		lines := strings.Split(string(output), "\n")
		memLine := lines[4] // Assuming the fifth line contains the memory usage
		memFields := strings.Fields(memLine)
		usedMem := strings.Trim(memFields[1], "M")
		totalMem := strings.Trim(memFields[3], "M")
		used, _ := strconv.Atoi(usedMem)
		total, _ := strconv.Atoi(totalMem)
		ramUsage = fmt.Sprintf("Used: %d MB / Total: %d MB", used, total)
	} else { // Linux
		cmd := exec.Command("free", "-m")
		output, _ := cmd.Output()
		lines := strings.Split(string(output), "\n")
		memLine := lines[1] // Assuming the second line contains the memory usage
		memFields := strings.Fields(memLine)
		usedMem := memFields[2]
		totalMem := memFields[1]
		used, _ := strconv.Atoi(usedMem)
		total, _ := strconv.Atoi(totalMem)
		ramUsage = fmt.Sprintf("Used: %d MB / Total: %d MB", used, total)
	}
	
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	
	result := hashToFloat(hostname)
	
	c := colorful.Hcl(result, 0.5, 0.5)
	hostcolor := c.Hex()
	
	hostStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color(hostcolor)).
		Bold(true).
		Width((physicalWidth/3)-2).
		Padding(1).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
	
	bannerStyle := lipgloss.NewStyle().
		MarginLeft(1).
		Padding(0, 1).
		Width(((physicalWidth/3)*2)-2).
		Align(lipgloss.Left)
		
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(hostcolor)).Render
	
	
	host := hostStyle.Render(hostname)
	banner := bannerStyle.Render(" OS:", infoStyle(osInfo), "Version:", infoStyle(osVersion),
	"\n","IP Address:", infoStyle(ipAddr),
	"\n","Total:", infoStyle(totalSpace),"Free:", infoStyle(freeSpace),"Used Percent:", infoStyle(usedPercent),
	"\n","Uptime:", infoStyle(uptime),
	"\n","RAM Usage:", infoStyle(ramUsage),
	"\n","CPU Usage:", infoStyle(cpuUsage))

	row := lipgloss.JoinHorizontal(lipgloss.Center, host, banner)	

	fmt.Println(row)

}
