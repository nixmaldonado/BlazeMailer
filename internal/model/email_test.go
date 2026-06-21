package model

import "testing"

func TestEmailRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		req     EmailRequest
		wantErr bool
	}{
		{
			"valid",
			EmailRequest{From: "a@b.com", To: "c@d.com", Subject: "hi", Body: "yo"},
			false,
		},
		{
			"missing from",
			EmailRequest{To: "c@d.com", Subject: "hi", Body: "yo"},
			true,
		},
		{
			"missing to",
			EmailRequest{From: "a@b.com", Subject: "hi", Body: "yo"},
			true,
		},
		{
			"bad to",
			EmailRequest{From: "a@b.com", To: "not-an-email", Subject: "hi", Body: "yo"},
			true,
		},
		{
			"empty subject ok",
			EmailRequest{From: "a@b.com", To: "c@d.com", Body: "yo"},
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() err = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
