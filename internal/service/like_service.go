package service

type LikeService interface {
	LikePost(userID, postID uint) error
	UnlikePost(likeID uint) error
}
