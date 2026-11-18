package value

import "github.com/textwire/textwire/v6/pkg/fail"

// emptyGlobal is shared by every scope without global data. It's frozen, so
// sharing it is safe.
var emptyGlobal = func() *Obj {
	obj := NewObj(nil)
	obj.frozen = true
	return obj
}()

// Freeze marks v and everything reachable from it as read-only.
func Freeze(v Literal) {
	freezeWalk(v, map[Literal]bool{})
}

// NewFrozenGlobal converts native global data into a frozen object.
func NewFrozenGlobal(data map[string]any) (*Obj, *fail.Error) {
	if data == nil {
		return emptyGlobal, nil
	}

	lit := NativeToValue(data)
	if errObj, ok := lit.(*Error); ok {
		return nil, errObj.Err
	}

	// NativeToValue returns an *Obj for any map, so this is only a guard.
	obj, ok := lit.(*Obj)
	if !ok {
		return nil, fail.New(nil, "", fail.OriginTpl, fail.ErrUnknownType, lit.Type())
	}

	Freeze(obj)

	return obj, nil
}

// freezeWalk freezes arrays and objects depth first. The seen map keeps cycles
// from recursing forever.
func freezeWalk(v Literal, seen map[Literal]bool) {
	kids, ok := childValues(v)
	if !ok || seen[v] {
		return
	}

	seen[v] = true

	switch val := v.(type) {
	case *Obj:
		val.frozen = true
	case *Arr:
		val.frozen = true
	}

	for _, kid := range kids {
		freezeWalk(kid, seen)
	}
}
