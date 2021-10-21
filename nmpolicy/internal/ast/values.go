package ast

func CurrentStateIdentity() Terminal {
	literal := "currentState"
	return Terminal{
		Identity: &literal,
	}
}
