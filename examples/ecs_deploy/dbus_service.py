import gobject
import dbus
import dbus.service
import dbus.mainloop.glib
class DemoException(dbus.DBusException):
    _dbus_error_name = 'com.example.DemoException'
class SomeObject(dbus.service.Object):
    @dbus.service.method("com.emc.vipr.fabric.Value",
                         in_signature='s', out_signature='a(ssissss)')
    def getDiskInfoListByService(self, hello_message):
        print (str(hello_message))
        return [["1","1", 100 ,"Good","1","1","1"]]
    @dbus.service.method("com.emc.vipr.fabric.Value",
                         in_signature='', out_signature='(ss)')
    def getNodeMode(self):
        return ["OPER_MODE","asdd"]
    def HelloWorld(self, hello_message):
        print (str(hello_message))
        return ["Hello", " from example-service.py", "with unique name",
                session_bus.get_unique_name()]
    @dbus.service.method("com.example.SampleInterface",
                         in_signature='', out_signature='')
    def Exit(self):
        mainloop.quit()
if __name__ == '__main__':
    dbus.mainloop.glib.DBusGMainLoop(set_as_default=True)
    session_bus = dbus.SystemBus()
    name = dbus.service.BusName("com.emc.vipr.fabric.hal", session_bus)
    object = SomeObject(session_bus, '/GlobalValue')
    mainloop = gobject.MainLoop()
    print "Running example service."
#    print usage
    mainloop.run()
