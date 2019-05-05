package binchunk

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0.53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZE_SIZE       = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

// 二进制chunk格式可以分为
// - 头部
// - 主函数与upvalue数量
// - 主函数原型
type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 主函数与upvalue数量
	mainFunc     *Prototype // 主函数原型
}

// 头部的结构如下
type header struct {
	// 签名, 魔数, 4个字节，分别是esc，L, u, a的ASCII码, 十六进制0x1B4C7561, go语言字符串字面量是"\x1bLua"
	signature [4]byte
	// 版本号, lua5.3.4包含三个版本号, Major Version, Minor Version, Release Version, 这里的version = major * 16 + minor
	version byte
	// 格式号, lua官方实现的格式号是0
	format byte
	// 校验文件用的，等于0x1993(lua1.0发布的年份)0D(回车符)0A(换行符)1A(替换符)0A, go语言字符串字面量等于 "\x19\x93\r\n\x1a\n"
	loacData [6]byte
	// cint在chunk中占用的字节数
	cintSize byte
	// size_t在chunk中占用的字节数
	sizetSize byte
	// lua虚拟机指令在chunk中占用的字节数
	instructionSize byte
	// lua整数在chunk中占用的字节数
	luaIntegerSize byte
	// lua浮点数在chunk中占用的字节数
	luaNumberSize byte
	// 下面n个字节存放lua整数值0x5678, n等于上面的luaIntegerSize, 主要是为了判断大小端方式
	luacInt int64
	// 下面n个字节存放lua浮点数370.5, n等于上面的luaNumberSize, 主要是为了检测chunk使用的二进制浮点格式, 一般都是 IEEE754
	luacNum float64
}
