package voice

const (
	// listenerPath is the object AT-SPI calls when a key changes.
	listenerPath = "/org/a11y/atspi/listeners/1"

	// accessibilityBusName and accessibilityBusPath locate the AT-SPI address service.
	accessibilityBusName = "org.a11y.Bus"
	accessibilityBusPath = "/org/a11y/bus"
	getAddressMethod     = "org.a11y.Bus.GetAddress"

	// deviceEventListenerInterface is the interface exported for key notifications.
	deviceEventListenerInterface = "org.a11y.atspi.DeviceEventListener"
	rightControlKeyName          = "Control_R"

	// atspiRegistryName owns the device controller that registers keys.
	atspiRegistryName    = "org.a11y.atspi.Registry"
	deviceControllerPath = "/org/a11y/atspi/registry/deviceeventcontroller"

	registerKeystrokeMethod   = "org.a11y.atspi.DeviceEventController.RegisterKeystrokeListener"
	deregisterKeystrokeMethod = "org.a11y.atspi.DeviceEventController.DeregisterKeystrokeListener"

	atspiAddressError      = "atspi address"
	listenerRefusedMessage = "atspi refused the Right Ctrl listener"
)
