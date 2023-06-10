package user_test

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gitamped/bud/services/user"
	"github.com/gitamped/bud/services/user/stores/nosql"
	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
	"github.com/gitamped/seed/values"
	"github.com/gitamped/stem/data/nosql/dbtest"
	"github.com/gitamped/stem/docker"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func Test_User(t *testing.T) {
	b, _ := os.ReadFile("../../testdata/collections.txt")
	cols := strings.Split(string(b), "\n")
	b, _ = os.ReadFile("../../testdata/seed.txt")
	seed := string(b)

	d := dbtest.Data{
		CollectionData: cols,
		SeedAql:        seed,
	}

	log, db, teardown := dbtest.NewUnit(t, c, "testcreateuser", d)
	t.Cleanup(teardown)
	storer := nosql.NewStore(log, db)

	core := user.NewUserServicer(log, storer)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := context.Background()
			now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
			email, err := mail.ParseAddress("user@example.com")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to parse email: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to parse email.", dbtest.Success, testID)

			nu := user.CreateUserRequest{}
			nu.NewUser.Name = "John Doe"
			nu.NewUser.Email = *email
			nu.NewUser.Roles = []user.Role{user.RoleAdmin}
			nu.NewUser.Password = "gophers"
			nu.NewUser.PasswordConfirm = "gophers"

			cuUsr := core.CreateUser(nu, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})
			if cuUsr.User.Name != "John Doe" {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user %+v : got %+v.", dbtest.Failed, testID, nu, cuUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", dbtest.Success, testID)

			uu := user.UpdateUserRequest{cuUsr.User}
			uu.User.Name = "JD"
			uuUsr := core.UpdateUser(uu, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{Roles: []string{auth.RoleAdmin}},
				Values: &values.Values{Now: now},
			})

			if uuUsr.User.ID != uu.User.ID {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user %+v : got %+v.", dbtest.Failed, testID, uu, uuUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", dbtest.Success, testID)

			// User update
			uuu := user.UpdateUserRequest{cuUsr.User}
			uuu.User.Name = "JD"
			uuuUsr := core.UpdateUser(uuu, server.GenericRequest{
				Ctx: ctx,
				Claims: auth.Claims{
					RegisteredClaims: jwt.RegisteredClaims{ID: uuu.User.ID.String()},
					Roles:            []string{auth.RoleUser},
				},
				Values: &values.Values{Now: now},
			})

			if uuuUsr.User.ID != uuu.User.ID {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user %+v : got %+v.", dbtest.Failed, testID, uuu, uuuUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", dbtest.Success, testID)

			// invalid user update
			iuuu := user.UpdateUserRequest{cuUsr.User}
			iuuu.User.Name = "JD"
			iuuuUsr := core.UpdateUser(iuuu, server.GenericRequest{
				Ctx: ctx,
				Claims: auth.Claims{
					RegisteredClaims: jwt.RegisteredClaims{ID: uuid.NewString()},
					Roles:            []string{auth.RoleUser},
				},
				Values: &values.Values{Now: now},
			})

			if iuuuUsr.Error != "Unauthorized action" {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to update another user profile %+v : got %+v.", dbtest.Failed, testID, iuuu, iuuuUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould not be able to update another user profile.", dbtest.Success, testID)

			// delete user
			du := user.DeleteUserRequest{cuUsr.User}
			duUsr := core.DeleteUser(du, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if duUsr.User.ID != cuUsr.User.ID {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user %+v : got %+v.", dbtest.Failed, testID, du, duUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", dbtest.Success, testID)

		}
	}
}
