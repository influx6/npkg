package natomic

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRespondGroupAtom(t *testing.T) {
	var group ResponderGroup

	var started = make(chan struct{}, 1)
	group.Start(started)
	<-started

	require.True(t, group.running())

	var signals signalDelivery
	require.NoError(t, group.AddGuaranteed(&signals, started))
	<-started

	var atom = NewAtom(&group)
	atom.Set(1)
	require.Equal(t, 1, atom.Read())

	group.Close()

	require.Equal(t, int(signals.event), 1)
}

func TestRespondGroup(t *testing.T) {
	var group ResponderGroup
	require.False(t, group.running())

	var started = make(chan struct{}, 1)
	group.Start(started)
	<-started
	require.True(t, group.running())

	var signals signalDelivery
	require.NoError(t, group.AddGuaranteed(&signals, started))
	<-started

	require.Len(t, group.responders, 1)

	var _, added = group.responders[&signals]
	require.True(t, added)

	group.Respond(event{})
	group.Respond(event{})
	group.Respond(event{})
	group.Respond(event{})

	require.Equal(t, int(signals.event), 4)

	group.Close()

	_, added = group.responders[&signals]
	require.False(t, added)
}

type event struct{}

func (event) Target() string {
	return "target"
}

func (event) Source() string {
	return "cli"
}

func (event) Type() string {
	return "event"
}

type signalDelivery struct {
	event int64
}

func (s *signalDelivery) Respond(_ Signal) {
	atomic.AddInt64(&s.event, 1)
}
