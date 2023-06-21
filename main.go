package main

import (
	"context"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/swinton/go-probot/probot"
)

func main() {
	probot.HandleEvent("pull-request", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.PullRequestEvent)
		pr := event.GetPullRequest()
		repo := event.GetRepo()

		log.Printf("Reviewing PR with title %v", pr.Number)

		// Start PR review
		prReview := &github.PullRequestReviewRequest{CommitID: pr.Head.SHA}

		// Check PR title for value, and approve if present
		if strings.Contains(*pr.Title, "Tenant Onboarding Request") {
			prReview.Event = github.String("APPROVE")
			prReview.Body = github.String("LGTM")
		}

		// Send approval to GitHub API
		review, _, err := ctx.GitHub.PullRequests.CreateReview(context.Background(), *repo.Owner.Login, *repo.Name, pr.GetNumber(), prReview)

		if err != nil {
			return err
		}

		log.Printf("Approved PR Review Request: %v", review)

		// Attempt PR merge
		if !*pr.Merged && *pr.Mergeable {

			mergeCommitMessage := "Merge onboarding request"
			mergeOptions := &github.PullRequestOptions{CommitTitle: mergeCommitMessage, SHA: *pr.Head.SHA, MergeMethod: "squash"}
			mergeResult, _, err := ctx.GitHub.PullRequests.Merge(context.Background(), *repo.Owner.Login, *repo.Name, pr.GetNumber(), mergeCommitMessage, mergeOptions)

			if err != nil {
				log.Printf("Failed to merge PR: %v", err)
				return err
			}

			log.Printf("Pull Request merged: %v ", mergeResult.Merged)

		}
		return nil
	})

	probot.Start()
}
