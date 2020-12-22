package tcp_package

func NewTcpPackage(maxSize int64, subpackageId int64) *TcpPackage {
	return &TcpPackage{
		MaxSubpackageSize: maxSize,
		SubpackageId:      subpackageId,
	}
}
