package state

// Boolean is a bool typed wrapper for Value.
//
// false <-> []byte{0x00}
// true <-> []byte{0x01}
type Boolean struct {
	Value
}

func (v Value) Boolean() Boolean {
	return Boolean{v}
}

// Get decodes and returns the stored boolean value if it exists. It will panic
// if the value exists but is not boolean type.
func (v Boolean) Get(ctx Context) (res bool) {
	v.Value.Get(ctx, &res)
	return
}

// GetSafe decodes and returns the stored boolean value. It will return an error
// if the value does not exist or not boolean.
func (v Boolean) GetSafe(ctx Context) (res bool, err error) {
	err = v.Value.GetSafe(ctx, &res)
	return
}

// Set encodes and sets the boolean argument to the state.
func (v Boolean) Set(ctx Context, value bool) {
	v.Value.Set(ctx, value)
}

// Query() retrives state value and proof from a queryable reference
func (v Boolean) Query(q ABCIQuerier) (res bool, proof *Proof, err error) {
	proof, err = v.Value.Query(q, &res)
	return
}
