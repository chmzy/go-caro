package converter

import (
	modelrepo "go-caro/internal/repository/pending_album/model"
	modelserv "go-caro/internal/service/pending_album/model"
)

func ToPendingAlbumFromRepo(album []modelrepo.AlbumPost) []modelserv.AlbumPost {
	out := make([]modelserv.AlbumPost, 0, len(album))
	for _, post := range album {
		p := modelserv.AlbumPost{
			ID:      post.ID,
			AlbumID: post.AlbumID,
			Author:  post.Author,
			MsgLink: modelserv.Link(post.MsgLink),
		}
		out = append(out, p)
	}

	return out
}
