package voice

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"
)

// Hold closes PressedChannel when Right Ctrl goes down and ReleasedChannel when it comes up.
type Hold struct {
	PressedChannel  chan struct{}
	ReleasedChannel chan struct{}
}

// deviceEvent is one AT-SPI key notification.
type deviceEvent struct {
	Type        uint32
	ID          int32
	HWCode      uint32
	Modifiers   uint32
	Timestamp   int32
	EventString string
	IsText      bool
}

// keyDef selects a key by code, symbol, name, and an unused field.
type keyDef struct {
	Keycode   int32
	Keysym    int32
	Keystring string
	Unused    int32
}

// listenMode tells AT-SPI to deliver the key from every window.
type listenMode struct {
	Sync       bool
	Preemptive bool
	AllWindows bool
}

type keyListener struct {
	pressOnce   sync.Once
	releaseOnce sync.Once
	hold        *Hold
}

// NotifyEvent marks press as type 0 and release as type 1.
func (listener *keyListener) NotifyEvent(event deviceEvent) (bool, *dbus.Error) {
	switch event.Type {
	case 0:
		listener.pressOnce.Do(func() { close(listener.hold.PressedChannel) })
	case 1:
		listener.releaseOnce.Do(func() { close(listener.hold.ReleasedChannel) })
	}
	return false, nil
}

// WaitForHold listens for Right Ctrl until the caller cancels.
func WaitForHold(ctx context.Context) (*Hold, func(), error) {
	// Open the session bus so we can ask AT-SPI where its own bus lives.
	sessionBus, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, nil, err
	}

	// Read the accessibility bus address from the session bus.
	var busAddress string
	err = sessionBus.Object(accessibilityBusName, accessibilityBusPath).
		Call(getAddressMethod, 0).
		Store(&busAddress)
	if err != nil {
		sessionBus.Close()
		return nil, nil, fmt.Errorf("%s: %w", atspiAddressError, err)
	}

	// Connect to the accessibility bus that delivers key events.
	accessibilityBus, err := dbus.Connect(busAddress)
	if err != nil {
		sessionBus.Close()
		return nil, nil, err
	}

	// PressedChannel and ReleasedChannel close once, when Right Ctrl goes down and up.
	hold := &Hold{
		PressedChannel:  make(chan struct{}),
		ReleasedChannel: make(chan struct{}),
	}
	listener := &keyListener{hold: hold}
	// Export the listener so AT-SPI can call NotifyEvent on this path.
	if err := accessibilityBus.Export(listener, listenerPath, deviceEventListenerInterface); err != nil {
		accessibilityBus.Close()
		sessionBus.Close()
		return nil, nil, err
	}

	// Match Right Ctrl by name, with an empty key code and key symbol.
	rightControlKeys := []keyDef{{Keystring: rightControlKeyName}}
	deviceController := accessibilityBus.Object(atspiRegistryName, deviceControllerPath)

	// Register press and release (event mask 3) from every window.
	var registered bool
	registerCall := deviceController.Call(
		registerKeystrokeMethod,
		0,
		dbus.ObjectPath(listenerPath),
		rightControlKeys,
		uint32(0),
		uint32(3),
		listenMode{AllWindows: true},
	)
	if err := registerCall.Store(&registered); err != nil || !registered {
		accessibilityBus.Close()
		sessionBus.Close()
		if err == nil {
			err = errors.New(listenerRefusedMessage)
		}
		return nil, nil, err
	}

	// Drop the registration and both buses the first time cleanup runs.
	var stopOnce sync.Once
	cleanup := func() {
		stopOnce.Do(func() {
			_ = deviceController.Call(
				deregisterKeystrokeMethod,
				0,
				dbus.ObjectPath(listenerPath),
				rightControlKeys,
				uint32(0),
				uint32(3),
			).Err
			accessibilityBus.Close()
			sessionBus.Close()
		})
	}
	// Cancel the wait if the caller is done before the key is used.
	go func() {
		<-ctx.Done()
		cleanup()
	}()

	return hold, cleanup, nil
}
