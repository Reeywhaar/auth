package token

import (
	"crypto/sha1"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_HashID(t *testing.T) {
	tbl := []struct {
		id   string
		hash string
	}{
		{"myid", "6e34471f84557e1713012d64a7477c71bfdac631"},
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"blah blah", "135a1e01bae742c4a576b20fd41a683f6483ca43"},
		{"da39a3ee5e6b4b0d3255bfef95601890afd80709", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	}

	for i, tt := range tbl {
		hh := sha1.New()
		assert.Equal(t, tt.hash, HashID(hh, tt.id), "case #%d", i)
	}
}

func TestUser_Attrs(t *testing.T) {
	u := User{Name: "test", IP: "127.0.0.1"}

	u.SetBoolAttr("k1", true)
	v := u.BoolAttr("k1")
	assert.True(t, v)

	u.SetBoolAttr("k1", false)
	v = u.BoolAttr("k1")
	assert.False(t, v)
	err := u.StrAttr("k1")
	assert.NotNil(t, err)

	u.SetStrAttr("k2", "v2")
	vs := u.StrAttr("k2")
	assert.Equal(t, "v2", vs)

	u.SetStrAttr("k2", "v22")
	vs = u.StrAttr("k2")
	assert.Equal(t, "v22", vs)

	vb := u.BoolAttr("k2")
	assert.False(t, vb)
}

func TestUser_Admin(t *testing.T) {
	u := User{Name: "test", IP: "127.0.0.1"}
	assert.False(t, u.IsAdmin())
	u.SetAdmin(true)
	assert.True(t, u.IsAdmin())
	u.SetAdmin(false)
	assert.False(t, u.IsAdmin())
}

func TestUser_GetUserInfo(t *testing.T) {
	r, err := http.NewRequest("GET", "http://blah.com", nil)
	assert.Nil(t, err)
	_, err = GetUserInfo(r)
	assert.NotNil(t, err, "no user info")

	r = SetUserInfo(r, User{Name: "test", ID: "id"})
	u, err := GetUserInfo(r)
	assert.Nil(t, err)
	assert.Equal(t, User{Name: "test", ID: "id"}, u)
}

func TestUser_MustGetUserInfo(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recovered from panic")
		}
	}()

	r, err := http.NewRequest("GET", "http://blah.com", nil)
	assert.Nil(t, err)
	_ = MustGetUserInfo(r)
	assert.Fail(t, "should panic")

	r = SetUserInfo(r, User{Name: "test", ID: "id"})
	u := MustGetUserInfo(r)
	assert.Nil(t, err)
	assert.Equal(t, User{Name: "test", ID: "id"}, u)
}
