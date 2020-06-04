package api

import (
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"webix/usermanager/data"

	"github.com/disintegration/imaging"
)

func updateAvatar(id int, file io.Reader, path, server string, dao *data.DAO) (*data.User, error) {
	idStr := strconv.Itoa(id)
	target, err := ioutil.TempFile(filepath.Join(path, "avatars"), "*.jpg")
	if err != nil {
		return nil, err
	}

	err = getImagePreview(file, 300, 300, target)
	if err != nil {
		return nil, err
	}

	// get existing user
	u, err := dao.Users.GetOne(id)
	if err != nil {
		return nil, err
	}

	u.Avatar = getAvatarURL(idStr, filepath.Base(target.Name()), server)
	dao.Users.Save(u)

	return u, nil
}

func getAvatarURL(id, name, server string) string {
	return server+path.Join("/users", id, "avatar", name)
}

func getImagePreview(source io.Reader, width, height int, target io.Writer) error {
	src, err := imaging.Decode(source)
	if err != nil {
		return err
	}

	dst := imaging.Thumbnail(src, width, height, imaging.Lanczos)
	err = imaging.Encode(target, dst, imaging.JPEG)

	if err != nil {
		return err
	}
	return nil
}
