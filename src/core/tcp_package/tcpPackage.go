package tcp_package

import (
	"errors"
	"goPanel/src/common"
	"math"
	"sort"
)

func NewTcpPackage(subpackageSize int64, subpackageId int64) *TcpPackage {
	return &TcpPackage{
		SubpackageSize: subpackageSize,
		SubpackageId:   subpackageId,
	}
}

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
	var packageIndex int64 = 0
	for i = 0; i < t.PackageSize; i += t.SubpackageSize {
		packageIndexSize, err := common.IntToBytes(packageIndex)
		if err != nil {
			return []string{}, err
		}

		tmpData := byteId
		tmpData = append(tmpData, bytePackageSize...)
		tmpData = append(tmpData, packageIndexSize...)
		// 处理边界问题
		if (i + t.SubpackageSize) <= t.PackageSize {
			tmpData = append(tmpData, byteData[i:(i+t.SubpackageSize)]...)
		} else {
			tmpData = append(tmpData, byteData[i:]...)
		}

		packageIndex++
		retData = append(retData, string(tmpData))
	}

	return retData, nil
}

// 拆包
func (t *TcpPackage) TcpUnPacking(data []byte) (*PackageContent, error) {
	packageId, err := common.BytesToInt(data[:8])
	if err != nil {
		return nil, err
	}
	packageSize, err := common.BytesToInt(data[8:16])
	if err != nil {
		return nil, err
	}
	packageIndex, err := common.BytesToInt(data[16:24])
	if err != nil {
		return nil, err
	}
	packageContent := data[24:]

	content := &PackageContent{
		PackageId:    packageId,
		PackageSize:  packageSize,
		PackageIndex: packageIndex,
		Content:      packageContent,
	}

	return content, nil
}

// 合包
// 如果包体接收完成则合并，包体接收未完成则返回失败
func (t *TcpPackage) TcpJoinPackage(data map[int64]*PackageContent) ([]byte, error) {
	var content []byte

	if len(data) == 0 {
		return content, errors.New("The packet is empty！")
	}

	// 排序要合包的数据，转为数组
	// 拿到第一个获取到包体头信息，判断接收的数据是否完成
	var dataIndex []int
	for index, _ := range data {
		dataIndex = append(dataIndex, int(index))
	}
	sort.Ints(dataIndex)

	// 计算包体数据是否完整
	subpackageIndex := math.Ceil(float64(data[int64(dataIndex[0])].PackageSize) / float64(t.SubpackageSize))
	if int64(len(dataIndex)) != int64(subpackageIndex) {
		return content, errors.New("Incomplete packet reception!")
	}

	for _, item := range dataIndex {
		content = append(content, data[int64(item)].Content...)
	}

	return content, nil
}
