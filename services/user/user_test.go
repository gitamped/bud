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

			// query user by id
			qu := user.QueryUserByIDRequest{cuUsr.User.ID.String()}
			quUsr := core.QueryUserByID(qu, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if quUsr.User.ID != cuUsr.User.ID && quUsr.User.Email != cuUsr.User.Email {
				t.Fatalf("\t%s\tTest %d:\tShould be able to query user by id %+v : got %+v.", dbtest.Failed, testID, qu, quUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to query user by id.", dbtest.Success, testID)

			// query user by email
			que := user.QueryUserByEmailRequest{Email: cuUsr.User.Email.Address}
			queUsr := core.QueryUserByEmail(que, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if queUsr.User.ID != cuUsr.User.ID && queUsr.User.Email != cuUsr.User.Email {
				t.Fatalf("\t%s\tTest %d:\tShould be able to query user by email %+v : got %+v.", dbtest.Failed, testID, que, queUsr)
			}

			t.Logf("\t%s\tTest %d:\tShould be able to query user by email.", dbtest.Success, testID)

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
