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

host:
```


```
