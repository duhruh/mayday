package repository

import "context"

// Init -
func Init(ctx context.Context, o Observation, t Type) error {
	if err := t.Init(ctx); err != nil {
		return err
	}
	if err := o.Init(ctx); err != nil {
		return err
	}
	return nil
}
