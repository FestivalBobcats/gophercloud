package users

import (
	"testing"

	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	os "github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const instanceID = "{instanceID}"

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleCreate(t)

	opts := os.BatchCreateOpts{
		os.CreateOpts{
			Databases: db.BatchCreateOpts{
				db.CreateOpts{Name: "databaseA"},
			},
			Name:     "dbuser3",
			Password: "secretsecret",
		},
		os.CreateOpts{
			Databases: db.BatchCreateOpts{
				db.CreateOpts{Name: "databaseB"},
				db.CreateOpts{Name: "databaseC"},
			},
			Name:     "dbuser4",
			Password: "secretsecret",
		},
	}

	res := Create(fake.ServiceClient(), instanceID, opts)
	th.AssertNoErr(t, res.Err)
}

func TestUserList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleList(t)

	expectedUsers := []os.User{
		os.User{
			Databases: []db.Database{
				db.Database{Name: "databaseA"},
			},
			Name: "dbuser3",
		},
		os.User{
			Databases: []db.Database{
				db.Database{Name: "databaseB"},
				db.Database{Name: "databaseC"},
			},
			Name: "dbuser4",
		},
	}

	pages := 0
	err := List(fake.ServiceClient(), instanceID).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := os.ExtractUsers(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, expectedUsers, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	os.HandleDelete(t)

	res := Delete(fake.ServiceClient(), instanceID, "{userName}")
	th.AssertNoErr(t, res.Err)
}
