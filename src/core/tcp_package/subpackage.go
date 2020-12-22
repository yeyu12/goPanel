package tcp_package

import (
	"fmt"
	"goPanel/src/common"
)

// 分包
func (t *TcpPackage) TcpSubpackage(data string) ([]string, error) {
	var retData []string
	byteData := []byte(data)
	t.PackageSize = int64(len(byteData))
	byteId, err := common.IntToBytes(t.SubpackageId)
	if err != nil {
		return retData, err
	}
	bytePackageSize, err := common.IntToBytes(int64(len(byteData)))
	if err != nil {
		return retData, err
	}

	var i int64 = 0
	for i = 0; i < t.PackageSize; i += t.MaxSubpackageSize {
		tmpData := byteId
		tmpData = append(tmpData, bytePackageSize...)

		// 处理边界问题
		if (i + t.MaxSubpackageSize) <= t.PackageSize {
			tmpData = append(tmpData, byteData[i:(i+t.MaxSubpackageSize)]...)
		} else {
			tmpData = append(tmpData, byteData[i:]...)
		}

		fmt.Println(common.BytesToInt(tmpData[:8]))
		fmt.Println(common.BytesToInt(tmpData[8:16]))
		fmt.Println(string(tmpData[16:]))

		retData = append(retData, string(tmpData))
	}

	return retData, nil
}
