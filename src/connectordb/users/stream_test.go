/**
Copyright (c) 2016 The ConnectorDB Contributors
Licensed under the MIT license.
**/
package users

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Note: The schema here must be same as the minified version
	streamtestType = "{\"type\":\"string\"}"
)

func TestCreateStream(t *testing.T) {
	for _, testdb := range testdatabases {
		_, dev, stream, err := CreateUDS(testdb)
		require.Nil(t, err)

		err = testdb.CreateStream(&StreamMaker{Stream: Stream{Name: stream.Name, Schema: streamtestType, DeviceID: dev.DeviceID}})
		assert.NotNil(t, err, "Created stream with duplicate name")

		// Test with invalid schema
		err = testdb.CreateStream(&StreamMaker{Stream: Stream{Name: "tcs_001", Schema: "{", DeviceID: dev.DeviceID}})
		assert.NotNil(t, err, "Created stream with invalid schema")

		// Test with embedded objects
		/*
			err = testdb.CreateStream("tcs_002", `{
			"type":"object",
			"properties":{
					"foo":{
						"type":"object"
					}
				}
			}`, dev.DeviceID)
			assert.NotNil(t, err, "Created stream with object schema")
		*/
	}
}

func TestUpdateStream(t *testing.T) {

	for _, testdb := range testdatabases {
		_, _, stream, err := CreateUDS(testdb)
		require.Nil(t, err)

		stream.Nickname = "true"
		stream.Schema = streamtestType
		stream.Datatype = "mytype"

		err = testdb.UpdateStream(stream)
		assert.Nil(t, err, "Could not update stream %v", err)

		stream2, err := testdb.ReadStreamByID(stream.StreamID)
		require.Nil(t, err, "got an error when trying to get a stream that should exist %v", err)

		if !reflect.DeepEqual(stream, stream2) {
			t.Errorf("The original and updated objects don't match orig: %v updated %v", stream, stream2)
		}

		err = testdb.UpdateStream(nil)
		assert.Equal(t, err, InvalidPointerError, "Function safeguards failed")
	}
}

func TestDeleteStream(t *testing.T) {

	for _, testdb := range testdatabases {
		_, _, stream, err := CreateUDS(testdb)
		require.Nil(t, err)

		err = testdb.DeleteStream(stream.StreamID)
		require.Nil(t, err, "Error when attempted delete %v", err)

		_, err = testdb.ReadStreamByID(stream.StreamID)
		require.NotNil(t, err, "The stream with the selected ID should have errored out, but it was not")
	}
}

func TestReadStreamByDevice(t *testing.T) {

	for _, testdb := range testdatabases {
		_, dev, _, err := CreateUDS(testdb)
		require.Nil(t, err)

		testdb.CreateStream(&StreamMaker{Stream: Stream{Name: "TestReadStreamByDevice2", Schema: streamtestType, DeviceID: dev.DeviceID}})

		streams, err := testdb.ReadStreamsByDevice(dev.DeviceID)
		require.Nil(t, err)
		require.Len(t, streams, 2, "didn't get enough streams")
	}
}

func TestReadStreamsByUser(t *testing.T) {
	for _, testdb := range testdatabases {

		user, _, stream, err := CreateUDS(testdb)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.NotNil(t, stream)

		fmt.Printf("User Id: %v\n", user.UserID)

		// create a bunch of devices
		for i := 0; i < 10; i++ {
			device, err := CreateTestDevice(testdb, user)

			require.Nil(t, err)

			fmt.Printf("Device Id: %v\n", device.DeviceID)

			// create a bunch of streams
			for j := 0; j < 10; j++ {
				name := GetNextName()
				err = testdb.CreateStream(&StreamMaker{Stream: Stream{Name: name, Schema: "{\"type\":\"number\"}", DeviceID: device.DeviceID}})
				require.Nil(t, err)
			}

			// And create 2 special downlink streams
			for j := 0; j < 2; j++ {
				name := GetNextName()
				err = testdb.CreateStream(&StreamMaker{Stream: Stream{Name: name, Schema: "{\"type\":\"number\"}", DeviceID: device.DeviceID, Downlink: true}})
				require.Nil(t, err)
			}
		}

		// Test selecting them
		streams, err := testdb.ReadStreamsByUser(user.UserID, false, false, false)
		require.Nil(t, err, "Retrieved streams was nil")

		// We need to add in the other missing stream
		require.Equal(t, 120+1, len(streams), "Wrong number of streams returned")

		// Now get the downlink streams
		streams, err = testdb.ReadStreamsByUser(user.UserID, false, true, false)
		require.Nil(t, err, "Retrieved streams was nil")

		require.Equal(t, 20, len(streams), "Wrong number of streams returned")
	}
}
