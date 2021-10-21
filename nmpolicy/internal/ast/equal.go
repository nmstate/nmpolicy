package ast

func deepEqualStringPtr(lhs, rhs *string) bool {
	if lhs == rhs {
		return true
	}
	if rhs != nil && lhs != nil {
		return *lhs == *rhs
	}
	return false
}

func (lhs Terminal) DeepEqual(rhs Terminal) bool {
	return deepEqualStringPtr(rhs.Identity, lhs.Identity) &&
		deepEqualStringPtr(rhs.String, lhs.String)
}
