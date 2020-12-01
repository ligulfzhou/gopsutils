# gopsutils-ssh
gopsutils is written for mobile usage.


# build this framework for iOS
```
cd PSUtils

gomobile bind -target ios -o ../PSUtils.framework
```


# use this framework.

```
connect:

let client = PSUtilsNewPSUtils(user, password, host, keyPath, keyString, port)
if let client = client {
    print("client initialized...")
    var conn: ObjCBool = false
    do {
        try client.connect(&conn)
        print("connection status: \(conn.boolValue)")
        if conn.boolValue == true {
        } else {
        }
    } catch {
    }
}
```

cpu:
```
cpu physical count: 
=> client.CPUCount(false)

cpu logical count: 
=> client.CPUCount(true)

cpu information:
=> client.CPUInfo()

cpu time: 
=> 
```

host:
```
host information: 
(Hostname,Uptime,BootTime,Procs,OS,Platform,PlatformFamily,PlatformVersion,KernelVersion,KernelArch,VirtualizationSystem,VirtualizationRole,HostID)  
=> client.HostInfoStat()
```

mem:
```
```

disk:
```
```

load:
```
```

net:
```
```