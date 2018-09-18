// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/forum"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumTallyVotesCmd = &cobra.Command{
	Use:   "tally-votes [proposal_name]",
	Short: "Tally votes according to the `type` of the proposal.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")
		proposalName := toName(args[0], "proposal_name")
		api := getAPI()

		votes, err := getForumVotesRows(api, targetAccount, proposalName)
		errorCheck("getting proposal", err)

		tallyStaked := make(map[uint8]int64)
		tallyAccounts := make(map[uint8]int64)
		var totalStaked int64
		for _, vote := range votes {
			tallyStaked[vote.vote.Vote] = tallyStaked[vote.vote.Vote] + int64(vote.account.VoterInfo.Staked)
			totalStaked += int64(vote.account.VoterInfo.Staked)
			tallyAccounts[vote.vote.Vote] = tallyAccounts[vote.vote.Vote] + 1

		}
		totalStakedEOS := eos.NewEOSAsset(totalStaked)

		fmt.Printf("Vote tally for proposal %q:\n", proposalName)
		fmt.Printf("* %d accounts voted\n", len(votes))
		fmt.Printf("* %s staked total\n", totalStakedEOS.String())

		output := []string{
			"Vote value | Num accounts | EOS staked",
			"---------- | ------------ | ----------",
		}
		for k, stakedForVote := range tallyStaked {
			accountsForVote := tallyAccounts[k]
			output = append(output, fmt.Sprintf("%d | %d | %s", k, accountsForVote, eos.NewEOSAsset(stakedForVote).String()))
		}
		fmt.Println(columnize.SimpleFormat(output))
	},
}

func getForumVotesRows(api *eos.API, contract eos.AccountName, proposalName eos.Name) (out []*forumVoteEntry, err error) {
	// lowerBound := "first"
	// for {
	// 	// TODO: Optimize by querying the secondary index..
	// 	resp, err := api.GetTableRows(eos.GetTableRowsRequest{
	// 		Code:       string(contract),
	// 		Scope:      string(contract),
	// 		Table:      string("vote"),
	// 		Index:      "sec",  // Secondary Index `by_proposal` - https://github.com/eoscanada/eosio.forum/blob/master/include/forum.hpp#L99-L115
	// 		KeyType:    "i128", // `by_proposal` is uint128 - Compute as https://github.com/eoscanada/eosio.forum/blob/master/include/forum.hpp#L72-L74
	// 		LowerBound: "first",
	// 		Limit:      1000,
	// 		JSON:       true,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	var votes []*forumVoteEntry
	// 	if err := json.Unmarshal(resp.Rows, &votes); err != nil {
	// 		return nil, err
	// 	}

	// 	for _, vote := range votes {
	// 		// TODO: optimize with getting only the bare minimum, like:
	// 		// cleosk get table eosio eosio voters -L cancancan234 --limit 1
	// 		acctResp, err := api.GetAccount(entry.vote.Voter)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		entry.account = acctResp

	// 		out = append(out, entry)
	// 	}

	// 	if !resp.More {
	// 		break
	// 	}
	// }

	return
}

type forumVoteEntry struct {
	vote    *forum.Vote
	account *eos.AccountResp
}

func init() {
	forumCmd.AddCommand(forumTallyVotesCmd)
}
