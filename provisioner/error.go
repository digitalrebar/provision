package provisioner

import "github.com/rackn/rocket-skates/models"

func NewError(c int64, s string) *models.Error {
	return &models.Error{Code: c, Message: s}
}
