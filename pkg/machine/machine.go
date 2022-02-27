package machine

import (
	"fmt"
	"os"
	"strconv"

	"github.com/godbus/dbus"
)

type Machine struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

const (
	dbusDest      = "org.freedesktop.login1"
	dbusInterface = "org.freedesktop.login1.Manager"
	dbusPath      = "/org/freedesktop/login1"
)

// New establishes a connection to the system bus and authenticates.
func New() (*Machine, error) {
	m := &Machine{}
	if err := m.initConnection(); err != nil {
		return nil, err
	}
	if err := m.CanShutdown(); err != nil {
		return nil, err
	}
	return m, nil
}

// Close closes the dbus connection
func (m *Machine) Close() {
	if m == nil {
		return
	}
	if m.conn != nil {
		m.conn.Close()
	}
}

func (m *Machine) initConnection() error {
	var err error
	m.conn, err = dbus.SystemBusPrivate()
	if err != nil {
		return err
	}
	// Only use EXTERNAL method, and hardcode the uid (not username)
	// to avoid a username lookup (which requires a dynamically linked
	// libc)
	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}
	err = m.conn.Auth(methods)
	if err != nil {
		m.conn.Close()
		return err
	}
	err = m.conn.Hello()
	if err != nil {
		m.conn.Close()
		return err
	}
	m.object = m.conn.Object("org.freedesktop.login1", dbus.ObjectPath(dbusPath))
	return nil
}

func (m *Machine) CanShutdown() error {
	var out interface{}
	if err := m.object.Call(dbusInterface+".CanPowerOff", 0).Store(&out); err != nil {
		return err
	}
	switch out.(string) {
	case "na":
		return fmt.Errorf("PowerOff not supported on this machine")
	case "no", "challenge":
		return fmt.Errorf("user does not have proper permissions for PowerOff")
	}
	return nil
}

func (m *Machine) Shutdown() error {
	m.object.Call(dbusInterface+".PowerOff", 0, false)
	return nil
}
