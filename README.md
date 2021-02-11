# preface
This repo heavily copies code from [shirou/gopsutil](https://github.com/shirou/gopsutil). It is an awesome project. But it is not for mobile. So I rerange the code to fit for mobile usage. BTW, I only adopt for linux server, other platforms have been omitted now. (Maybe some time later, I will add them. But it is a high probability event.)

# gopsutils-for-mobile
gopsutils is written for mobile, use gomobile to build framework for iOS and Android.

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
            ...connection success...
        }
    } catch {
    }
}

```

Attension: I set the timeout of connection to 5 seconds, do not connect server in main queue but in backgroup queue. Or It will stuck.
```
DispatchQueue.global(qos: .background).async {
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
latest 1min,5min,15min load average
=> client.ArgLoad()
{
    load1:   Float
    load5:   Float
    load15:  Float
}
```

net:
```
```
