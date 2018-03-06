package main

import "testing"

func Test_getStatusMessage(t *testing.T) {
	statusMessage := "Hi Wicker Park, our ants are ready to swarm!"
	if got := getStatusMessage(); got != statusMessage {
		t.Errorf("getStatusMessage() = %v, want %v", got, statusMessage)
	}

}
