package main

import (
	"context"
	comment "toktik/kitex_gen/douyin/comment"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// ActionComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ActionComment(ctx context.Context, req *comment.ActionCommentRequest) (resp *comment.ActionCommentResponse, err error) {
	// TODO: Your code here...
	return
}

// ListComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ListComment(ctx context.Context, req *comment.ListCommentRequest) (resp *comment.ListCommentResponse, err error) {
	// TODO: Your code here...
	return
}
