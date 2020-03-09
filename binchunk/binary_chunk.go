package binchunk

const (
	LUA_SIGNATURE    = "\x1bLua"            // 签名魔数
	LUAC_VERSION     = 0x53                 // lua版本
	LUAC_FORMAT      = 0                    // 格式号
	LUAC_DATA        = "\x19\x93\r\n\x1a\n" // 0x1993 是lua发布年份，后面依次是回车符换行符替换福和另一个换行符
	CINT_SIZE        = 4                    // 整数长度
	CSIZET_SIZE      = 8                    //
	INSTRUCTION_SIZE = 4                    // lua虚拟机指令字节数
	LUA_INTEGER_SIZE = 8                    // lua整数字节数
	LUA_NUMBER_SIZE  = 8                    // lua浮点数字节数
	LUAC_INT         = 0x5678               //
	LUAC_NUM         = 370.5                //
)

// 常量tag前缀, 常量表中使用
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

// 二进制chunk格式可以分为
// - 头部
// - 主函数与upvalue数量
// - 主函数原型
type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 主函数upvalue数量
	mainFunc     *Prototype // 主函数原型
}

// 头部的结构如下
type header struct {
	signature       [4]byte // 签名, 魔数, 4个字节，分别是esc，L, u, a的ASCII码, 十六进制0x1B4C7561, go语言字符串字面量是"\x1bLua"
	version         byte    // 版本号, lua5.3.4包含三个版本号, Major Version, Minor Version, Release Version, 这里的version = major * 16 + minor
	format          byte    // 格式号, lua官方实现的格式号是0
	loacData        [6]byte // 校验文件用的，等于0x1993(lua1.0发布的年份)0D(回车符)0A(换行符)1A(替换符)0A, go语言字符串字面量等于 "\x19\x93\r\n\x1a\n"
	cintSize        byte    // cint在chunk中占用的字节数
	sizetSize       byte    // size_t在chunk中占用的字节数
	instructionSize byte    // lua虚拟机指令在chunk中占用的字节数
	luaIntegerSize  byte    // lua整数在chunk中占用的字节数
	luaNumberSize   byte    // lua浮点数在chunk中占用的字节数
	luacInt         int64   // 下面n个字节存放lua整数值0x5678, n等于上面的luaIntegerSize, 主要是为了判断大小端方式
	luacNum         float64 // 下面n个字节存放lua浮点数370.5, n等于上面的luaNumberSize, 主要是为了检测chunk使用的二进制浮点格式, 一般都是 IEEE754
}

type Prototype struct {
	Source          string        // 源文件名, @开头是真正的文件，=stdin表示从标准输入读取的，其他还有从程序提供的字符串编译来的
	LineDefined     uint32        // 开始行号, 结束行号, 对于普通函数，都大于0，如果是main函数，都等于0
	LastLineDefined uint32        // 结束行号
	NumParams       byte          // 固定参数的个数，相对于可变参数而言的
	IsVararg        byte          // 是否是Vararg函数，是否是可边长参数函数
	MaxStackSize    byte          // 需要的寄存器数量
	Code            []uint32      // 指令表，每个指令4个字节
	Constants       []interface{} // 常量表，程序中出现的字面量, 包括nil, bool值，整数，浮点数，字符串
	Upvalues        []Upvalue     // Upvalue表，每个元素占2个字节
	Protos          []*Prototype  // 子函数原型表
	LineInfo        []uint32      // 行号表, 这里的行号和指令表中的指令一一对应, 记录指令对应的行号
	LocVars         []LocVar      // 局部变量表，记录局部变量名和起止索引
	UpvalueNames    []string      // Upvalue名列表, 和前面Upvalue表中的元素一一对应, 记录Upvalue在代码中的名称
}

// LineInfo, LocVars, UpvalueNames存放的都是调试信息，对程序的执行并没有影响，如果在编译时指定`-s`选项，
// 编译后的chunk中这三个表将被清空

type Upvalue struct {
	Instack byte
	Idx     byte
}

// 局部变量表结构
type LocVar struct {
	VarName string
	StartPc uint32
	EndPC   uint32
}

func UnDump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        // 校验头部
	reader.readByte()           // 跳过Upvalue数量
	return reader.readProto("") // 读取函数原型
}
