package tcp_package

type TcpPackage struct {
	SubpackageSize int64 // 分包长度(位
	PackageSize    int64 // 包体大小
	SubpackageId   int64 // 分包id
}

type PackageContent struct {
	PackageId    int64  // 包体id
	PackageSize  int64  // 包体总大小
	PackageIndex int64  // 分片的块
	Content      []byte // 内容
}
