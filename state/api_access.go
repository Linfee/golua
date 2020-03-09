package state

import "golua/api"

// 返回指定lua类型的字符串表示
func (s *luaState) TypeName(tp api.LuaType) string {
	switch tp {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRING:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// 返回指定索引位置的值得类型
func (s *luaState) Type(idx int) api.LuaType {
	if s.stack.isValid(idx) {
		return typeOf(s.stack.get(idx))
	}
	return api.LUA_TNONE
}

func (s *luaState) IsNone(idx int) bool {
	return s.Type(idx) == api.LUA_TNONE
}

func (s *luaState) IsNil(idx int) bool {
	return s.Type(idx) == api.LUA_TNIL
}

func (s *luaState) IsNoneOrNil(idx int) bool {
	return s.Type(idx) <= api.LUA_TNIL
}

func (s *luaState) IsBoolean(idx int) bool {
	return s.Type(idx) == api.LUA_TBOOLEAN
}

func (s *luaState) IsString(idx int) bool {
	t := s.Type(idx)
	return t == api.LUA_TSTRING || t == api.LUA_TNUMBER
}

func (s *luaState) IsNumber(idx int) bool {
	_, ok := s.ToNumberX(idx)
	return ok
}

func (s *luaState) IsInteger(idx int) bool {
	val := s.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (s *luaState) ToBoolean(idx int) bool {
	val := s.stack.get(idx)
	return convertToBoolean(val)
}

// lua中，只有false和nil表示假，其它都表示真
func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func (s *luaState) ToNumber(idx int) float64 {
	n, _ := s.ToNumberX(idx)
	return n
}

func (s *luaState) ToNumberX(idx int) (float64, bool) {
	val := s.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func (s *luaState) ToInteger(idx int) int64 {
	i, _ := s.ToIntegerX(idx)
	return i
}

func (s *luaState) ToIntegerX(idx int) (int64, bool) {
	val := s.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

func (s *luaState) ToString(idx int) string {
	str, _ := s.ToStringX(idx)
	return str
}

func (s *luaState) ToStringX(idx int) (string, bool) {
	val := s.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		str := fmt.Sprintf("%v", x)
		s.stack.set(idx, s)
		return str, true
	default:
		return "", false
	}
}
