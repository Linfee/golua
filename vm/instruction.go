package vm

const MAXARG_Bx = 1<<18 - 1       // 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 131071

type Instruction uint32

// 操作码名称
func (i *Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

// 编码模式
func (i *Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

// 操作数B
func (i *Instruction) BMode() byte {
	return opcodes[i.Opcode()].argBMode
}

// 操作数C
func (i *Instruction) CMode() byte {
	return opcodes[i.Opcode()].argCMode
}

// 操作码
func (i *Instruction) Opcode() int {
	return int(*i & 0x3F)
}

func (i *Instruction) ABC() (a, b, c int) {
	a = int(*i >> 6 & 0xFF)
	b = int(*i >> 14 & 0x1FF)
	c = int(*i >> 23 & 0x1FF)
	return
}

func (i *Instruction) ABx() (a, bx int) {
	a = int(*i >> 6 & 0xFF)
	bx = int(*i >> 14)
	return a, bx
}

func (i *Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

func (i *Instruction) Ax() int {
	return int(*i >> 6)
}
