package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte // 存放要被解析的chunk数据
}

func (r *reader) readByte() byte {
	res := r.data[0]
	r.data = r.data[1:]
	return res
}

func (r *reader) readUint32() uint32 {
	res := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return res
}

func (r *reader) readUint64() uint64 {
	res := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return res
}

// 读取lua整数
func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

// 读取lua浮点数
func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

// 读取字符串
func (r *reader) readString() string {
	// 假设为短字符串
	size := uint(r.readByte())
	if size == 0 {
		// null字符串
		return ""
	}
	if size == 0xFF { // 长字符串
		size = uint(r.readUint64())
	}
	bytes := r.readBytes(size - 1)
	return string(bytes)
}

// 从字节流读取n个字节
func (r *reader) readBytes(n uint) []byte {
	bytes := r.data[:n]
	r.data = r.data[n:]
	return bytes
}

// 检查头部
func (r *reader) checkHeader() {
	if string(r.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	} else if r.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	} else if r.readByte() != LUAC_FORMAT {
		panic("format mismatch")
	} else if string(r.readBytes(6)) != LUAC_DATA {
		panic("corrupted")
	} else if r.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	} else if r.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	} else if r.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	} else if r.readByte() != LUA_INTEGER_SIZE {
		panic("lua_integer size mismatch!")
	} else if r.readByte() != LUA_NUMBER_SIZE {
		panic("lua_number size mismatch!")
	} else if r.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	} else if r.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

// 读取函数原型
func (r *reader) readProto(parentSource string) *Prototype {
	source := r.readString()
	if source == "" {
		// 子函数原型要从父函数原型哪里获取文件名
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

// 读取指令表
func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

// 读取常量表
func (r *reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

// 读取一个常量
func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.readByte() != 0
	case TAG_INTEGER:
		return r.readLuaInteger()
	case TAG_NUMBER:
		return r.readLuaNumber()
	case TAG_SHORT_STR:
		return r.readString()
	case TAG_LONG_STR:
		return r.readString()
	default:
		panic("corrupted!")
	}
}

// 读取Upvalue
func (r *reader) readUpvalues() []Upvalue {
	res := make([]Upvalue, r.readUint32())
	for i := range res {
		res[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return res
}

// 读取函数原型表
func (r *reader) readProtos(parentSourdce string) []*Prototype {
	res := make([]*Prototype, r.readUint32())
	for i := range res {
		res[i] = r.readProto(parentSourdce)
	}
	return res
}

// 读取行号表
func (r *reader) readLineInfo() []uint32 {
	res := make([]uint32, r.readUint32())
	for i := range res {
		res[i] = r.readUint32()
	}
	return res
}

// 读取局部变量表
func (r *reader) readLocVars() []LocVar {
	res := make([]LocVar, r.readUint32())
	for i := range res {
		res[i] = LocVar{
			VarName: r.readString(),
			StartPc: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return res
}

func (r *reader) readUpvalueNames() []string {
	res := make([]string, r.readUint32())
	for i := range res {
		res[i] = r.readString()
	}
	return res
}
