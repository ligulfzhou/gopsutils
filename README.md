# gopsutils-for-mobile
gopsutils is written for mobile, use gomobile to build framework for iOS or Android.

# build
```
cd PSUtils

build for iOS:
=> gomobile bind -target ios -o ../PSUtils.framework
```

# use this framework.
connect linux:
```
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

Attention:
I set the timeout of connection to 5 seconds, do not connect server in main queue but in backgroup queue.
=> DispatchQueue.global(qos: .background).async {
     ...do connection...
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