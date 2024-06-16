package dtos_test

import (
	"testing"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
)

type registerDtoTestCase struct {
	name string
	dtos.Register
	expectedErrs int
}

var testcases = []registerDtoTestCase{
	{
		name:         "empty register dto",
		Register:     dtos.Register{},
		expectedErrs: 8,
	},
	{
		name: "valid register dto",
		Register: dtos.Register{
			Email:     "test@domain.co.uk",
			Password:  "Pa22W0rd123",
			FirstName: "John",
			LastName:  "Doe",
			Username:  "JohnDoe2024",
		},
		expectedErrs: 0,
	},
	{
		name: "password too short",
		Register: dtos.Register{
			Email:     "test@domain.co.uk",
			Password:  "P1sfs",
			FirstName: "John",
			LastName:  "Doe",
			Username:  "JohnDoe2024",
		},
		expectedErrs: 1,
	},
	{
		name: "password missing numbers",
		Register: dtos.Register{
			Email:     "test@domain.co.uk",
			Password:  "aaaaaa",
			FirstName: "John",
			LastName:  "Doe",
			Username:  "JohnDoe2024",
		},
		expectedErrs: 1,
	},
	{
		name: "password missing letters",
		Register: dtos.Register{
			Email:     "test@domain.co.uk",
			Password:  "123456",
			FirstName: "John",
			LastName:  "Doe",
			Username:  "JohnDoe2024",
		},
		expectedErrs: 1,
	},
	{
		name: "spaces in username",
		Register: dtos.Register{
			Email:     "test@domain.co.uk",
			Password:  "Pa22W0rd123!",
			FirstName: "John",
			LastName:  "Doe",
			Username:  "John Doe 2024",
		},
		expectedErrs: 1,
	},
	{
		name: "no alphanumberic firstname : containing space",
		Register: dtos.Register{
			Email:     "test@domain.co.uk",
			Password:  "Pa22W0rd123!",
			FirstName: "Joh n",
			LastName:  "Doe",
			Username:  "JohnDoe2024",
		},
		expectedErrs: 1,
	},
	{
		name: "no alphanumberic lastname : containing symbol",
		Register: dtos.Register{
			Email:     "test@domain.co.uk",
			Password:  "Pa22W0rd123!",
			FirstName: "John",
			LastName:  "Do!e",
			Username:  "JohnDoe2024",
		},
		expectedErrs: 1,
	},
}

func TestRegister_Validation(t *testing.T) {
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if errs := testcase.Validate(); len(errs) != testcase.expectedErrs {
				t.Errorf("expected %v errors but got %v", testcase.expectedErrs, len(errs))
				t.Log(errs)
			}
		})
	}
}
