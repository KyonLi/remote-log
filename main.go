package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"remote-log/command"
)

var mossSep = "-----------------------------------------------------------------\n"

var filePath = "/var/log/tomcat8/catalina.out"
var hostStr = "root@tc.shencai.net.cn"
var password = ""
var privateKey = ""

func printWelcomeMessage(config command.Config) {
	fmt.Println(mossSep)

	for _, server := range config.Servers {
		// If there is no tail_file for a service configuration, the global configuration is used
		if server.TailFile == "" {
			server.TailFile = config.TailFile
		}

		serverInfo := fmt.Sprintf("%s@%s:%s", server.User, server.Hostname, server.TailFile)
		fmt.Println(serverInfo)
	}
	fmt.Printf("\n%s\n", mossSep)
}

func parseConfig(filePath string, hostStr string) (config command.Config) {
	hosts := strings.Split(hostStr, ",")

	config = command.Config{}
	config.TailFile = filePath
	config.Servers = make(map[string]command.Server, len(hosts))
	for index, hostname := range hosts {
		hostInfo := strings.Split(strings.Replace(hostname, ":", "@", -1), "@")
		var port int
		if len(hostInfo) > 2 {
			port, _ = strconv.Atoi(hostInfo[2])
		}
		config.Servers["server_"+string(index)] = command.Server{
			ServerName: "server_" + string(index),
			Hostname:   hostInfo[1],
			User:       hostInfo[0],
			Port:       port,
			Password:   password,
			PrivateKey: privateKey,
		}
	}

	return
}

func main() {

	config := parseConfig(filePath, hostStr)
	printWelcomeMessage(config)

	outputs := make(chan command.Message, 255)
	var wg sync.WaitGroup

	for _, server := range config.Servers {
		wg.Add(1)
		go func(server command.Server) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Error: %s\n", err)
				}
			}()
			defer wg.Done()

			// If there is no tail_file for a service configuration, the global configuration is used
			if server.TailFile == "" {
				server.TailFile = config.TailFile
			}

			// If the service configuration does not have a port, the default value of 22 is used
			if server.Port == 0 {
				server.Port = 22
			}

			cmd := command.NewCommand(server)
			cmd.Execute(outputs)
		}(server)
	}

	if len(config.Servers) > 0 {
		go func() {
			for output := range outputs {
				content := strings.Trim(output.Content, "\r\n")
				// 去掉文件名称输出
				if content == "" || (strings.HasPrefix(content, "==>") && strings.HasSuffix(content, "<==")) {
					continue
				}

				fmt.Printf("%s\n", content)
			}
		}()
	} else {
		fmt.Println("No target host is available")
	}

	wg.Wait()
}
