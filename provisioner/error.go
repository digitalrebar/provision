package provisioner

type pkgError struct {
	Message string
}

func NewError(s string) *pkgError {
	return &pkgError{Message: s}
}
