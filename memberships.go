// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	youtube "google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"os"
)

type MembershipSubscription struct {
	// DisplayName: The channel's display name.
	DisplayName string `csv:"name"`
	
	// ProfileImageUrl: The channel's avatar URL.
	ProfileImageUrl string `csv:"profile_image_url"`

	// ChannelUrl: The channel's URL.
	ChannelUrl string `csv:"channelUrl"`

	// HighestAccessibleLevelDisplayName: Display name for the highest level
	// that the user has access to at the moment.
	HighestAccessibleLevelDisplayName string `cs:"highestAccessibleLevelDisplayName"`

	// MemberSince: The date and time when the user became a continuous
	// member across all levels.
	MemberSince string `cs:"memberSince,omitempty"`

	// MemberTotalDurationMonths: The cumulative time the user has been a
	// member across all levels in complete months (the time is rounded down
	// to the nearest integer).
	MemberTotalDurationMonths int64 `csv:"memberTotalDurationMonths"`
}

func init() {
	registerCommand("memberships-list", youtube.YoutubeChannelMembershipsCreatorScope, MembershipsLists)
}

func MembershipsLists(client *http.Client, argv []string) {
	if len(argv) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: memberships-list channelId")
		return
	}

	channelId := argv[0]

	service, err := youtube.New(client)

	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	var membershipListParts []string
	var memberships *youtube.MemberListResponse

	
	memberships, err = service.Members.List(membershipListParts).FilterByMemberChannelId(channelId).Mode("all_current").MaxResults(1000).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve membership list: %v", err)
	}

	var membersSubscriptionsImported []*MembershipSubscription

	for _, membership := range memberships.Items {
		membershipSnippet := membership.Snippet
		
		memberDetails := &MembershipSubscription{
			DisplayName: membershipSnippet.MemberDetails.DisplayName,
			ProfileImageUrl: membershipSnippet.MemberDetails.ProfileImageUrl,
			ChannelUrl: membershipSnippet.MemberDetails.ChannelUrl,
			HighestAccessibleLevelDisplayName: membershipSnippet.MembershipsDetails.HighestAccessibleLevelDisplayName,
			MemberSince: membershipSnippet.MembershipsDetails.MembershipsDuration.MemberSince,
			MemberTotalDurationMonths: membershipSnippet.MembershipsDetails.MembershipsDuration.MemberTotalDurationMonths,
		}

		membersSubscriptionsImported = append(membersSubscriptionsImported, memberDetails)
	}

	membershipsListFile, err := os.OpenFile("memberships-list.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		panic(err)
	}
	
	defer membershipsListFile.Close()

	err = gocsv.MarshalFile(&membersSubscriptionsImported, membershipsListFile)

	if err != nil {
		panic(err)
	}
}