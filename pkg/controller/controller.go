package controller

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type ErrHandlerNotFound struct {
	Type reflect.Type
}

func (e ErrHandlerNotFound) Error() string {
	var typeName string
	switch e.Type.Kind() {
	case reflect.Ptr:
		typeName = "*" + e.Type.Elem().Name()
	default:
		typeName = e.Type.Name()
	}
	return fmt.Sprintf("handler not found for type %s:%s", e.Type.PkgPath(), typeName)
}

func (ErrHandlerNotFound) Is(target error) bool {
	_, okPointer := target.(*ErrHandlerNotFound)
	_, okStruct := target.(ErrHandlerNotFound)
	return okPointer || okStruct
}

type Controller interface {
	// Exec accepts args of arbitary type and look up for handler accepting exactly same args type.
	// It must return ErrHandlerNotFound in case no handler able to handle provided args.
	Exec(ctx context.Context, args interface{}) error
}

type HandlerSpec interface {
	Context() interface{}
}

type handler struct {
	spec        HandlerSpec
	handlerFunc reflect.Value
}

type ControllerImpl struct {
	handlers map[reflect.Type]handler
}

func New(handlers ...HandlerSpec) *ControllerImpl {
	impl := &ControllerImpl{
		handlers: make(map[reflect.Type]handler),
	}
	for i := range handlers {
		if err := impl.RegisterHandler(handlers[i]); err != nil {
			panic(fmt.Errorf("invalid handler at %d: %w", i, err))
		}
	}

	return impl
}

func (c *ControllerImpl) RegisterHandler(hnd HandlerSpec) error {
	method := reflect.ValueOf(hnd).MethodByName("Exec")
	if !method.IsValid() {
		return fmt.Errorf("handler missed Exec()")
	}

	methodType := method.Type()
	if methodType.NumIn() != 2 {
		return fmt.Errorf("handler shold receive only 2 args instead of %d", methodType)
	}

	in0 := methodType.In(0)
	in0Value := reflect.TypeOf(context.Background())
	if !in0Value.Implements(in0) {
		return fmt.Errorf("1st arg should be of type context.Context instead of %v", in0)
	}

	argType := reflect.TypeOf(hnd.Context())
	if methodType.In(1) != argType {
		return fmt.Errorf("2nd arg (%v) should match Context() result %v", methodType.In(1), argType)
	}

	if methodType.NumOut() != 1 {
		return fmt.Errorf("handler should return error only, but it returns %d values", methodType.NumOut())
	}

	out0 := methodType.Out(0)
	out0Value := reflect.TypeOf(errors.New(""))
	if !out0Value.Implements(out0) {
		return fmt.Errorf("handler should return error instead of %v", out0)
	}

	c.handlers[argType] = handler{
		spec:        hnd,
		handlerFunc: method,
	}
	return nil
}

func (c ControllerImpl) Handlers() []HandlerSpec {
	result := make([]HandlerSpec, len(c.handlers))
	i := 0
	for _, v := range c.handlers {
		result[i] = v.spec
		i++
	}

	return result
}

// Exec implements Controller.Exec by routing args for specific handler or
// returning ErrHandlerNotFound
func (c *ControllerImpl) Exec(ctx context.Context, args interface{}) error {
	argsType := reflect.TypeOf(args)
	handler, ok := c.handlers[argsType]
	if !ok {
		return &ErrHandlerNotFound{Type: argsType}
	}

	rawErr := handler.handlerFunc.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(args)})
	err, isNotNil := rawErr[0].Interface().(error)
	if !isNotNil {
		return nil
	}

	return err
}
