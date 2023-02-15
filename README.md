This project is aimed to secretly write and read api keys. <br />
basic path:<br />
write method: User passes {write command, key, value} --> service layer starts encryption process (in future decides in which storage to pass data) --
-> storage passes to dedicated to exact storage <br />
Supports: <br /> 
multiple local storages and any other form of saving data <br />
          can be added http server.



