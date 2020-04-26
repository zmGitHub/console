# auto message 递归代替循环

1. start conv
2. create a check timeout delay job
3. create check visitor/agent no action delay job

4. convTimeoutCheckHandler:
```java
if conv end {
        delete job
        return
    }

    if conv timeout {
        end conv
        delete job
        return
    }

    delete old job
    create new check timeout delay job
```
    

5. no resp handler:
```java
if conv is end {
       delete oldJob
       return
   }
   if visitor/agent no action in duration {
       send msg
   }

   delete oldJob
   create new check job
```
