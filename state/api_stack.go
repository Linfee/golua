package state

func (s *luaState) GetTop() int {
	return s.stack.top
}

func (s *luaState) AbsIndex(idx int) int {
	return s.stack.absIndex(idx)
}

func (s *luaState) CheckStack(n int) bool {
	s.stack.check(n)
	return true
}

func (s *luaState) Pop(n int) {
	s.SetTop(-n - 1)
}

// 把一个值从一个位置复制到另一个位置
func (s *luaState) Copy(fromIdx, toIdx int) {
	val := s.stack.get(fromIdx)
	s.stack.set(toIdx, val)
}

// 将指定索引处的值推入栈顶
func (s *luaState) PushValue(idx int) {
	val := s.stack.get(idx)
	s.stack.push(val)
}

// 将栈顶值弹出，并替换掉指定位置
func (s *luaState) Replace(idx int) {
	val := s.stack.pop()
	s.stack.set(idx, val)
}

// 将栈顶值弹出，并插入指定位置
func (s *luaState) Insert(idx int) {
	s.Rotate(idx, 1)
}

// 删除指定索引处的值
func (s *luaState) Remove(idx int) {
	s.Rotate(idx, -1)
	s.Pop(1)
}

// 将[idx, top]索引区间的值朝栈顶方向旋转n个位置
func (s *luaState) Rotate(idx, n int) {
	t := s.stack.top - 1
	p := s.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	s.stack.reverse(p, m)
	s.stack.reverse(m+1, t)
	s.stack.reverse(p, t)
}

// 将指定索引设为栈顶，高于改索引的值都弹出
func (s *luaState) SetTop(idx int) {
	newTop := s.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}
	if n := s.stack.top - newTop; n > 0 {
		for i := 0; i < n; i++ {
			s.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			s.stack.push(nil)
		}
	}
}

func (s *luaState) PushNil()              { s.stack.push(nil) }
func (s *luaState) PushBoolean(b bool)    { s.stack.push(b) }
func (s *luaState) PushInteger(n int64)   { s.stack.push(n) }
func (s *luaState) PushNumber(n float64)  { s.stack.push(n) }
func (s *luaState) PushString(str string) { s.stack.push(str) }
