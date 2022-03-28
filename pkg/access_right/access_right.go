package access_right

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/storage"
	"context"
	"errors"
	"strconv"
	"strings"
)

type AccessRightHandler struct {
	env config.Environment
}

func NewAccessRightHandlerHandler(env config.Environment) *AccessRightHandler {
	return &AccessRightHandler{env: env}
}

func (h *AccessRightHandler) Exec(ctx context.Context, args *storage.AccessRightQuery) (err error) {

	strArr := strings.Fields(h.env.Admins)
	for _, s := range strArr {
		if s == strconv.FormatInt(args.Id, 10) {
			return nil
		}
	}
	return errors.New("no access")
}

func (h *AccessRightHandler) Context() interface{} {
	return (*storage.AccessRightQuery)(nil)
}
