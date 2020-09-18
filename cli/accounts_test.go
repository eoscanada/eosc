package cli

import (
	"testing"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	pubKey1, _ := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	pubKey2, _ := ecc.NewPublicKey("EOS4xenzB8vAWjwHxnk8eGLkPumXDAEA1Sgq11U2muX3kJ8n7v2KA")

	tests := []struct {
		name   string
		input  string
		expect *eos.Authority
	}{
		{
			name:  "full",
			input: "3=EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV+2,abourget@perm+3,bob",
			expect: &eos.Authority{
				Threshold: uint32(3),
				Waits:     []eos.WaitWeight{},
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      eos.AccountName("abourget"),
							Permission: eos.PermissionName("perm"),
						},
						Weight: 3,
					},
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      eos.AccountName("bob"),
							Permission: eos.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: pubKey1,
						Weight:    2,
					},
				},
			},
		},
		{
			name:  "full of spaces",
			input: "3  =  EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV  +  2  ,  abourget@perm + 3 , bob",
			expect: &eos.Authority{
				Threshold: uint32(3),
				Waits:     []eos.WaitWeight{},
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      eos.AccountName("abourget"),
							Permission: eos.PermissionName("perm"),
						},
						Weight: 3,
					},
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      eos.AccountName("bob"),
							Permission: eos.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: pubKey1,
						Weight:    2,
					},
				},
			},
		},
		{
			name:  "single account",
			input: "abourget",
			expect: &eos.Authority{
				Threshold: uint32(1),
				Waits:     []eos.WaitWeight{},
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      eos.AccountName("abourget"),
							Permission: eos.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []eos.KeyWeight{},
			},
		},
		{
			name:  "single key",
			input: "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expect: &eos.Authority{
				Threshold: uint32(1),
				Waits:     []eos.WaitWeight{},
				Accounts:  []eos.PermissionLevelWeight{},
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: pubKey1,
						Weight:    1,
					},
				},
			},
		},
		{
			name:  "sorted keys",
			input: "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV, EOS4xenzB8vAWjwHxnk8eGLkPumXDAEA1Sgq11U2muX3kJ8n7v2KA",
			expect: &eos.Authority{
				Threshold: uint32(1),
				Waits:     []eos.WaitWeight{},
				Accounts:  []eos.PermissionLevelWeight{},
				Keys: []eos.KeyWeight{
					eos.KeyWeight{
						PublicKey: pubKey2,
						Weight:    1,
					},
					eos.KeyWeight{
						PublicKey: pubKey1,
						Weight:    1,
					},
				},
			},
		},
		{
			name:  "sorted accounts",
			input: "alex, bob",
			expect: &eos.Authority{
				Threshold: uint32(1),
				Waits:     []eos.WaitWeight{},
				Accounts: []eos.PermissionLevelWeight{
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      eos.AccountName("alex"),
							Permission: eos.PermissionName("active"),
						},
						Weight: 1,
					},
					eos.PermissionLevelWeight{
						Permission: eos.PermissionLevel{
							Actor:      eos.AccountName("bob"),
							Permission: eos.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []eos.KeyWeight{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ParseShortFormAuth(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expect, res)
		})
	}
}
