package get_create_update_user

import (
	"backend-vpn/pkg/storage"
	"context"
)

type CreateUpdateUser struct {
}

//func (pe UserDto) String() string {
//	if st, ok := pe.Msg.(fmt.Stringer); ok {
//		return "Msg: " + st.String()
//	}
//
//	return "Msg: <unknown>"
//}

type CreateUpdateHandler struct {
	//	publisher service.EventPublisher
}

func NewCreateUpdateHandler(publisher service.EventPublisher) *CreateUpdateHandler {
	return &CreateUpdateHandler{publisher}
}

func (h CreateUpdateHandler) Exec(ctx context.Context, args *storage.UserQuery) error {
	//native, err := model.MessageFromProto(args.Msg, model.GetReplicationKey(ctx), model.GetRequestContext(ctx).ID)
	//if err != nil {
	//	return err
	//}
	//
	//
	//return h.publisher.Publish(ctx, native)
	return nil
}

func (CreateUpdateHandler) Context() interface{} {
	return (*storage.UserQuery)(nil)
}
