package mesos

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	MaxPort    = 65535 - 1024
	RandomPort = 0
)

func PortMaping(reqPorts []int, effectivePorts map[int]bool) []int {
	var portsMap = make(map[string]bool)
	for _, v := range reqPorts {
		if v == RandomPort {
			for {
				v := randPort()
				if !portsMap[v] && effectivePorts[v] {
					portsMap[v] = true
					break
				}

				v = randPort()
			}
		} else {
			portsMap[v] = true
		}
	}

	return portsMap
}

func PortsEnv(reqPorts []int) map[string]string {
	var envMap = make(map[string]string)
	for i, v := range reqPorts {
		envMap[fmt.Sprintf("PORT%s", strconv.Itoa(i))] = strconv.Itoa(v)
	}

	return envMap
}

func randPort() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(MaxPort) + 1025
}
