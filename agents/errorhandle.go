package agents

import "proposalsubmitters/entities"

func aggErr(err error, rpcErr *entities.RPCError) error {
	if err != nil {
		return err
	}
	if rpcErr != nil {
		return rpcErr
	}
	return nil
}
