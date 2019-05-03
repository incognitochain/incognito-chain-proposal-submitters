package entities

func AggErr(err error, rpcErr *RPCError) error {
	if err != nil {
		return err
	}
	if rpcErr != nil {
		return rpcErr
	}
	return nil
}

type ErrorSaver struct {
	err error
}

func (s *ErrorSaver) Save(errs ...error) error {
	if s.err != nil {
		return s.err
	}
	for _, err := range errs {
		if err != nil {
			s.err = err
			return s.err
		}
	}
	return nil
}

func (s *ErrorSaver) Get() error {
	return s.err
}

func CheckError(errs ...error) error {
	errSaver := &ErrorSaver{}
	return errSaver.Save(errs...)
}
