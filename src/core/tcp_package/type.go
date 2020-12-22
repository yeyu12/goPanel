package tcp_package

type TcpPackage struct {
	MaxSubpackageSize int64 // 分包最大长度
	PackageSize       int64 // 包体大小
	SubpackageId      int64 // 分包id
}
