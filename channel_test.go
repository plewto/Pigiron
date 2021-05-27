package main

import (
	// "fmt"
	"testing"

	"github.com/plewto/pigiron/op"
)

func TestChannelValidation(t *testing.T) {
	valid := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	invalid := []int{0, 17}
	for _, c := range valid {
		err := op.ValidateMIDIChannel(c)
		if err != nil {
			t.Fatalf("ValidateMIDIChannel return incorrect result for %d", c)
		}
	}
	for _, c := range invalid {
		err := op.ValidateMIDIChannel(c)
		if err == nil {
			t.Fatalf("ValidateMIDIChannel return incorrect negative for %d", c)
		}
	}
}


func TestNullSelector(t *testing.T) {
	ncs := op.MakeNullChannelSelector()
	mode := ncs.Mode()
	if mode != op.NoChannel {
		t.Fatalf(".Mode() return incorrect: %s", mode)
	}
	channels := ncs.SelectedChannels()
	if len(channels) != 0 {
		t.Fatalf(".Channels() expected empty slice, got %v", channels)
	}
	for _, c := range []int{1, 2, 3, 4, -1, 19} {
		flag := ncs.ChannelSelected(c)
		if flag {
			t.Fatalf(".CahnnelEnabeld(%d), expected false, got %v", c, flag)
		}
	}
}


func TestSingleChannelSelector(t *testing.T) {
	scs := op.MakeSingleChannelSelector()
	mode := scs.Mode()
	if mode != op.SingleChannel {
		t.Fatalf(".Mode(), expected SingleChannel, got %v", mode)
	}
	odd := []int{1, 3, 5, 7, 9, 11, 13, 15}
	invalid := []int{0, 17}
	for _, c := range odd {
		err := scs.EnableChannel(c, true)
		if err != nil {
			t.Fatalf(".EnableChannel(%d) returned unexpected error", c)
		}
		if !scs.ChannelSelected(c) {
			t.Fatalf(".ChannelSelected(%d) returned false negative", c)
		}
	}
	for _, c := range invalid {
		err := scs.EnableChannel(c, true)
		if err == nil {
			t.Fatalf(".EnableChannel(%d), did not return expected error.", c)
		}
	}
	c := 5
	scs.SelectChannel(c)
	clst := scs.SelectedChannels()
	if len(clst) != 1 || clst[0] != c {
		t.Fatalf(".SelectedChanels(), expected [1], got %v", clst)
	}
}


func TestMultiChannelSelector(t *testing.T) {
	mcs := op.MakeMultiChannelSelector()
	mode := mcs.Mode()
	if mode != op.MultiChannel {
		t.Fatalf(".Mode(), expected MultiChannel, got %v", mode)
	}
	primes := []int{2, 3, 5, 7, 11, 13}
	composite := []int{1, 4, 6, 8, 9, 10, 12, 14, 15, 16}
	for _, c := range primes {
		mcs.EnableChannel(c, true)
		if !mcs.ChannelSelected(c) {
			t.Fatalf(".ChannelSelected(%d), unexpected return.", c)
		}
	}
	for _, c := range composite {
		if mcs.ChannelSelected(c) {
			t.Fatalf(".ChannelSelected(%d), false positive.", c)
		}
	}

	clst := mcs.SelectedChannels()
	if len(clst) != len(primes) {
		t.Fatalf(".SelectedChanels() return unexpected, %v", clst)
	}

	c := 2
	mcs.EnableChannel(c, false)
	if mcs.ChannelSelected(c) {
		t.Fatalf(".SelectedChannel(%d) returns false positive", c)
	}

	c = 99
	if mcs.ChannelSelected(c) {
		t.Fatalf(".SelectedChannel(%d) returns false positive", c)
	}

	mcs.DeselectAllChannels()
	for i := 1; i < 17; i++ {
		if mcs.ChannelSelected(i) {
			t.Fatalf(".DeselectAllChannels() --> ChannelSelected(%d) fails", i)
		}
	}
}
