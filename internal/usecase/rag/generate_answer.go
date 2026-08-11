package usecase_rag

import (
	"context"
)

func (u *ragUsecase) GenerateAnswer(ctx context.Context, name string) (string, error) {
	test, _ := u.grpc.GetAnswer(ctx, name)

	return test, nil
}
