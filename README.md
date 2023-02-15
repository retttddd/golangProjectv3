This project is aimed to secretly write and read api keys.
basic path:
write: User passes {write command, key, value} -
-> service layer starts encryption process (in future decides in which storage to pass data) --
-> storage passes to dedicated to exact storage
Supports: multiple local storages and any other form of saving data
          can be added http server.

FLOW DIAGRAM
https://docs.google.com/drawings/d/1tWcCpbe8mgTYYUC4GNBO3ofqPE5wAEsUSBzxqz_X_jQ/edit?usp=sharing

