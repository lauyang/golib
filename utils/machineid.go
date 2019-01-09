// 获取机器识别码
package utils

import (
	"encoding/hex"
	"net"
)

func GetMachineID() []string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	ret := make([]string, 0, len(interfaces))

	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		// return hex.EncodeToString(mac[:6])
		id := hex.EncodeToString(mac)
		if len(id) < 6 {
			continue
		}

		ret = append(ret, id)
	}

	return ret
}
