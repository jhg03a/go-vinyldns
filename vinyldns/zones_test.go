/*
Copyright 2018 Comcast Cable Communications Management, LLC
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vinyldns

import (
	"encoding/json"
	"testing"

	"github.com/gobs/pretty"
)

func TestZones(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones",
			code:     202,
			body:     zonesJSON,
		},
	})
	defer server.Close()

	zones, err := client.Zones()
	if err != nil {
		t.Log(pretty.PrettyFormat(zones))
		t.Error(err)
	}
	if len(zones) != 2 {
		t.Error("Expected 2 Domains")
	}
	for _, z := range zones {
		if z.Name == "" {
			t.Error("Expected zone.Name to have a value")
		}
		if z.Email == "" {
			t.Error("Expected zone.Email to have a value")
		}
		if z.Status == "" {
			t.Error("Expected zone.Status to have a value")
		}
		if z.Created == "" {
			t.Error("Expected zone.Created to have a value")
		}
		if z.ID == "" {
			t.Error("Expected zone.ID to have a value")
		}
		if z.AdminGroupID == "" {
			t.Error("Expected zone.AdminGroupID to have a value")
		}
	}
}

func TestZonesListAll(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones?maxItems=1",
			code:     200,
			body:     zonesListJSON1,
		},
		testToolsConfig{
			endpoint: "http://host.com/zones?startFrom=2&maxItems=1",
			code:     200,
			body:     zonesListJSON2,
		},
	})

	defer server.Close()

	if _, err := client.ZonesListAll(ListFilter{
		MaxItems: 200,
	}); err == nil {
		t.Error("Expected error -- MaxItems must be between 1 and 100")
	}

	zones, err := client.ZonesListAll(ListFilter{
		MaxItems: 1,
	})
	if err != nil {
		t.Error(err)
	}

	if len(zones) != 2 {
		t.Error("Expected 2 Zones; got ", len(zones))
	}

	if zones[0].ID != "1" {
		t.Error("Expected Zone.ID to be 1")
	}

	if zones[1].ID != "2" {
		t.Error("Expected Zone.ID to be 2")
	}
}

func TestZonesListAllWhenNone(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones",
			code:     200,
			body:     zonesListNoneJSON,
		},
	})

	defer server.Close()

	zones, err := client.ZonesListAll(ListFilter{})
	if err != nil {
		t.Error(err)
	}

	if len(zones) != 0 {
		t.Error("Expected 0 Zones; got ", len(zones))
	}

	j, err := json.Marshal(zones)
	if err != nil {
		t.Error(err)
	}

	if string(j) != "[]" {
		t.Error("Expected string-converted marshaled JSON to be '[]'; got ", string(j))
	}
}

func TestZone(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123",
			code:     200,
			body:     zoneJSON,
		},
	})

	defer server.Close()

	z, err := client.Zone("123")
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}

	if z.Name != "vinyldns." {
		t.Error("Expected zone.Name to have a value")
	}
	if z.Email != "some_user@foo.com" {
		t.Error("Expected zone.Email to have a value")
	}
	if z.Status != "Active" {
		t.Error("Expected zone.Status to have a value")
	}
	if z.Created != "2015-10-30T01:25:46Z" {
		t.Error("Expected zone.Created to have a value")
	}
	if z.ID != "123" {
		t.Error("Expected zone.ID to have a value")
	}
	if z.LatestSync == "" {
		t.Error("Expected zone.LatestSync to have a value")
	}
	if z.Updated == "" {
		t.Error("Expected zone.Updated to have a value")
	}
	if z.AdminGroupID == "" {
		t.Error("Expected zone.AdminGroupID to have a value")
	}
	if z.Connection.Name != "vinyldns." {
		t.Error("Expected zone.Connection.Name to have a value")
	}
	if z.Connection.KeyName != "vinyldns." {
		t.Error("Expected zone.Connection.KeyName to have a value")
	}
	if z.Connection.Key != "OBF:1:ABC" {
		t.Error("Expected zone.Connection.Key to have a value")
	}
	if z.Connection.PrimaryServer != "127.0.0.1" {
		t.Error("Expected zone.Connection.PrimaryServer to have a value")
	}
	if z.TransferConnection.Name != "vinyldns." {
		t.Error("Expected zone.TransferConnection.Name to have a value")
	}
	if z.TransferConnection.KeyName != "vinyldns." {
		t.Error("Expected zone.TransferConnection.KeyName to have a value")
	}
	if z.TransferConnection.Key != "OBF:1:ABC+5" {
		t.Error("Expected zone.TransferConnection.Key to have a value")
	}
	if z.TransferConnection.PrimaryServer != "127.0.0.1" {
		t.Error("Expected zone.TransferConnection.PrimaryServer to have a value")
	}

	rule := z.ACL.Rules[0]
	if rule.AccessLevel != "Read" {
		t.Error("Expected rule.AccessLevel to be Read")
	}
	if rule.Description != "test-acl-group-id" {
		t.Error("Expected rule.Description to be test-acl-group-id")
	}
	if rule.GroupID != "123" {
		t.Error("Expected rule.GroupId to be 123")
	}
	if rule.RecordMask != "www-*" {
		t.Error("Expected rule.RecordMask to be www-*")
	}
	for _, rt := range rule.RecordTypes {
		if rt != "A" && rt != "AAAA" && rt != "CNAME" {
			t.Error("Expected rule.RecordTypes to be A, AAAA, CNAME")
		}
	}
}

func TestZoneCreate(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones",
			code:     200,
			body:     zoneUpdateResponseJSON,
		},
	})

	defer server.Close()

	zone := &Zone{
		Name:         "test.",
		Email:        "email@email.com",
		AdminGroupID: "1234",
		Connection: &ZoneConnection{
			Name:          "connectionName",
			KeyName:       "keyName",
			Key:           "key",
			PrimaryServer: "1.2.3.4",
		},
	}

	z, err := client.ZoneCreate(zone)
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}
	if z.Zone.Name != "test." {
		t.Error("Expected zoneResponse.Zone.Name to have a value")
	}
	if z.UserID != "pclear" {
		t.Error("Expected zoneResponse.Zone.UserId to have a value")
	}
}

func TestZoneUpdate(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123",
			code:     200,
			body:     zoneUpdateResponseJSON,
		},
	})

	defer server.Close()

	zone := &Zone{
		Name:         "test.",
		Email:        "email@email.com",
		AdminGroupID: "123",
		Connection: &ZoneConnection{
			Name:          "connectionName",
			KeyName:       "keyName",
			Key:           "key",
			PrimaryServer: "1.2.3.4",
		},
	}

	z, err := client.ZoneUpdate("123", zone)
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}
	if z.Zone.Name != "test." {
		t.Error("Expected zoneResponse.Zone.Name to have a value")
	}
	if z.UserID != "pclear" {
		t.Error("Expected zoneResponse.Zone.UserId to have a value")
	}
}

func TestZoneDelete(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123",
			code:     200,
			body:     zoneUpdateResponseJSON,
		},
	})

	defer server.Close()

	z, err := client.ZoneDelete("123")
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}
	if z.Zone.Name != "test." {
		t.Error("Expected zoneResponse.Zone.Name to have a value")
	}
	if z.UserID != "pclear" {
		t.Error("Expected zoneResponse.Zone.UserId to have a value")
	}
}

func TestZoneExists_yes(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123",
			code:     200,
			body:     zoneUpdateResponseJSON,
		},
	})

	defer server.Close()

	z, err := client.ZoneExists("123")
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}
	if z != true {
		t.Error("Expected ZoneExists to be true")
	}
}

func TestZoneExists_no(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123",
			code:     404,
			body:     zoneUpdateResponseJSON,
		},
	})

	defer server.Close()

	z, err := client.ZoneExists("123")
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}
	if z != false {
		t.Error("Expected ZoneExists to be false")
	}
}

func TestZoneHistory(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123/history",
			code:     200,
			body:     zoneHistoryJSON,
		},
	})

	defer server.Close()

	z, err := client.ZoneHistory("123")
	zc := z.ZoneChanges[0]
	rs := z.RecordSetChanges[0]
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}
	if z.ZoneID != "123" {
		t.Error("Expected ZoneHistory.ZoneId to have a value")
	}
	if zc.UserID != "userId1" {
		t.Error("Expected ZoneHistory.ZoneChanges[0].UserId to have a value")
	}
	if zc.ChangeType != "Create" {
		t.Error("Expected ZoneHistory.ZoneChanges[0].ChangeType to have a value")
	}
	if zc.Status != "Complete" {
		t.Error("Expected ZoneHistory.ZoneChanges[0].Status to have a value")
	}
	if zc.ID != "change123" {
		t.Error("Expected ZoneHistory.ZoneChanges[0].Id to have a value")
	}
	if zc.Zone.Name != "vinyldnstest.sys.vinyldns.net." {
		t.Error("Expected ZoneHistory.ZoneChange.Zone.Name to have a value")
	}
	if rs.UserID != "account" {
		t.Error("Expected ZoneHistory.RecordSetChange.UserId to have a value")
	}
	if rs.ChangeType != "Create" {
		t.Error("Expected ZoneHistory.RecordSetChange.ChangeType to have a value")
	}
	if rs.Status != "Complete" {
		t.Error("Expected ZoneHistory.RecordSetChange.Status to have a value")
	}
	if rs.Created != "2015-11-02T13:59:48Z" {
		t.Error("Expected ZoneHistory.RecordSetChange.Status to have a value")
	}
	if rs.ID != "13c0f664-58c2-4b1a-9c46-086c3658f861" {
		t.Error("Expected ZoneHistory.RecordSetChange.Status to have a value")
	}
	if rs.Zone.Name != "vinyldnstest.sys.vinyldns.net." {
		t.Error("Expected ZoneHistory.RecordSetChange.Zone.Name to have a value")
	}
	if rs.RecordSet.ID != "rs123" {
		t.Error("Expected ZoneHistory.RecordSetChange.RecordSet.Id to have a value")
	}
	if rs.RecordSet.Records[0].Address != "127.0.0.1" {
		t.Error("Expected ZoneHistory.RecordSetChange.RecordSet.Records[0].Address to have a value")
	}
}

func TestZoneChange(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123/history",
			code:     200,
			body:     zoneHistoryJSON,
		},
	})

	defer server.Close()

	z, err := client.ZoneChange("123", "change123")
	if err != nil {
		t.Log(pretty.PrettyFormat(z))
		t.Error(err)
	}
	if z.UserID != "userId1" {
		t.Error("Expected ZoneChange.UserId to have a value")
	}
}

func TestRecordSets(t *testing.T) {
	server, client := testTools([]testToolsConfig{
		testToolsConfig{
			endpoint: "http://host.com/zones/123/recordsets",
			code:     200,
			body:     recordSetsJSON,
		},
	})

	defer server.Close()

	rs, err := client.RecordSets("123")
	if err != nil {
		t.Log(pretty.PrettyFormat(rs))
		t.Error(err)
	}
	if len(rs) != 2 {
		t.Error("Expected 2 Record Sets")
	}
	for _, r := range rs {
		if r.ID == "" {
			t.Error("Expected RecordSet.Id to have a value")
		}
	}
}