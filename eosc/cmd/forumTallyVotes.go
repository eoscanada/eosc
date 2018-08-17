// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/forum"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forumTallyVotesCmd = &cobra.Command{
	Use:   "tally-votes [proposer] [proposal_name]",
	Short: "Tally votes according to the `type` of the proposal.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetAccount := toAccount(viper.GetString("forum-cmd-target-contract"), "--target-contract")

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal_name")

		api := getAPI()

		// proposal, err := getForumProposal(api, targetAccount, proposer, proposalName)
		// errorCheck("getting proposal", err)

		votes, err := getForumVotesRows(api, targetAccount, proposer, proposalName)
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

		fmt.Printf("Vote tally for proposer %q's proposal %q:\n", proposer, proposalName)
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

func getForumProposal(api *eos.API, contract, proposer eos.AccountName, proposalName eos.Name) (out *forum.Propose, err error) {
	for {
		resp, err := api.GetTableRows(eos.GetTableRowsRequest{})
		if err != nil {
			return nil, err
		}

		for _, row := range resp.Rows {
			_ = row
		}

		if !resp.More {
			break
		}
	}

	return
}

func getForumVotesRows(api *eos.API, contract, proposer eos.AccountName, proposalName eos.Name) (out []*forumVoteEntry, err error) {

	for {
		// TODO: Optimize by querying the secondary index..
		resp, err := api.GetTableRows(eos.GetTableRowsRequest{
			JSON:  true,
			Code:  string(contract),
			Scope: string(proposer),
			Table: string("vote"),
		})
		if err != nil {
			return nil, err
		}

		var votes []*forum.Vote
		if err := json.Unmarshal(resp.Rows, &votes); err != nil {
			return nil, err
		}

		for _, vote := range votes {
			entry := &forumVoteEntry{
				vote: vote,
			}

			if vote.ProposalName != proposalName {
				continue
			}

			// TODO: optimize with getting only the bare minimum, like:
			// cleosk get table eosio eosio voters -L cancancan234 --limit 1
			acctResp, err := api.GetAccount(entry.vote.Voter)
			if err != nil {
				return nil, err
			}
			entry.account = acctResp

			out = append(out, entry)
		}

		if !resp.More {
			break
		}
	}

	return
}

type forumVoteEntry struct {
	vote    *forum.Vote
	account *eos.AccountResp
}

func init() {
	forumCmd.AddCommand(forumTallyVotesCmd)

	forumTallyVotesCmd.Flags().StringP("reserved", "", "", "reserved option")

	for _, flag := range []string{"reserved"} {
		if err := viper.BindPFlag("forum-taly-votes-cmd-"+flag, forumTallyVotesCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}
