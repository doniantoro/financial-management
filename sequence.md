

title This is a title

participant Customer
participant Admin
participant Sales
participant FE
participant BE(MiddleWare)
participant BE
participant Redis

==Authentication Start==
Customer->FE:Insert Email & Password (1)
FE->BE:(2)
BE->BE:validate email & password(3)
alt #Pink if true
BE->BE: Generate token and refresh-token (4)
BE->Redis:Save token and refresh token to redis(5)
BE->FE:Response Success with bring token,refresh token and user data(5)
FE->Customer: response success(5)
else #LightGreen else 
BE->FE: Response failed(6)
FE->Customer:Response failed,and tell to insert other credential(7)
end 
note over FE,BE:With some interval time,fe will hit refresh token (8)
Customer->FE: Click logout button(9)
FE->BE:validate token (asumtion true) (10)
note over BE:validate token in middleware
BE->Redis: Delete tokn from redis(11)
BE->FE:Success response logout (12)
FE->Customer:Success resposponse logout(13)
==Authentication End==
==Middleware Proccess Validation==

BE->BE(MiddleWare) : Validate Token and check permission(14)
BE(MiddleWare)->Redis: Check Token (15)
alt #Pink if false
Redis->BE:(16)
BE->FE:Response failed Authentication(17)
FE->Customer: response failed(18)
else #Pink 
BE->BE: continue journey(19)
end 

==Middleware Proccess Validation==

==Apply Aplication Start==
Customer->FE:Insert data(20)
FE->BE: (21)
note over BE:validate token in middleware
BE->BE:Check is there application before(22) 

alt #Pink if There is apply before that less than 1 month, 
BE->FE:Response failed failed(23)
FE->Customer: response failed(24)
end 

BE->BE:validate input file,is that image or not
alt #Pink if false
BE->FE:response failed validation(25)
FE->Customer:Response validation (26)
end 
BE->Database:Insert data into database(27)
BE->Redis:Insert data into Redis(28)
Database->BE:(29)
BE->FE: Response Success(30)
FE->Customer: Response Success(31)

Admin->FE:Click button to see application(32)
FE->BE:(33)
note over BE:validate token in middleware
alt  #Pink if data  exist on redis,
BE->FE:Response data (34)
FE->Admin:Response data(35)
else
BE->Database:get data in database(36)
Database->BE:(37)
BE->FE:Response data(38)
FE->Admin:Response data(39)
end 

FE->Admin:show data (40)
Admin->FE:See detail data(41) 
FE->BE:(42)
note over BE:validate token in middleware
alt  #Pink if data  exist on redis,
BE->FE:Response data(43)
FE->Admin:Response data(44)
else
BE->Database:get data in database(45)
Database->BE:(46)
BE->FE:Response data(47)
FE->Admin:Response data(48)
end 

Admin->Admin:Analyze to make decision(49)
Admin->BE:Make decision(50)
note over BE:validate token in middleware
BE->BE:Do some validation(51)
BE->Database:Update data in Database(52)
alt #Pink if data exist on redis,
BE->Redis:Update data in redis(53)
end
BE->FE:Give response(54)
FE->Admin:Give response(55)
FE->Customer:give notif to customer(56)

==Apply Aplication End==

==Transaction Start==
===Transaction Offline Start ===
Sales->FE:Insert data (57)
FE->BE:(58)
note over BE:validate token in middleware
BE->BE:remain limit(59)
BE->BE:create otp (60)
BE->Redis:save data token n otp(61) 
BE->Sms Gateway:Trigger send otp (62)
Sms Gateway->Customer:send otp to user(63)
Customer->Sales:Tell otp to sales(64)
Sales->FE:Input otp (65)
FE->BE:(66)
BE->Redis:Validate otp(67)
Redis->BE:(68)
BE->Database:Insert data to database(69)
BE->Database:Insert data to redis(70)
BE->Sms Gateway:Trigger to tell status transaction (71)
BE->FE:give response(72)
FE->Sales:give response(73)
Sales->Customer:tell status transaction(74)

===Transaction Offline End ===

===Transaction Marketplace Start ===
Customer->FE:See Web(75)
FE->BE:(76)
note over BE:validate token in middleware(77)
alt  #Pink if data  exist on redis,
BE->Redis:Check data(78)
Redis->BE:(79)
else 
BE->MarketPlace:request product(80)
MarketPlace->BE:response product(81)

end 
BE->FE:response list product(82)
FE->Customer:show list product(82)
Customer->FE:Insert data to make transaction(83)
FE->BE:(84)
note over BE:validate token in middleware
BE->BE:remain limit(85)
BE->BE:create otp (86)
BE->Redis:save data token n otp (87)
BE->Sms Gateway:Trigger send otp (88)
Sms Gateway->Customer:send otp to user(89)
Customer->FE:Input otp (90)
FE->BE:(91)
BE->Redis:Validate otp(92)
Redis->BE:(93)

BE->MarketPlace:request product detail(94)
MarketPlace->BE:Response product detail(95)
BE->MarketPlace:make transaction(96)
MarketPlace->MarketPlace:Decrese bucket balance(97)
MarketPlace->BE:Response(98)
BE->Database:Insert data to database(99)
BE->Redis:Insert data to redis(100)
BE->Sms Gateway:Trigger to tell status(101) 
BE->FE:give response(102)
FE->Customer:tell status transaction(103)

===Transaction Marketplace End ===

===See Transaction Start ===

Customer->FE:See Transaction(104)
note over Customer,Redis: all of role can see transactions with different permission
FE->BE:Request Transaction(105)
note over BE:validate token in middleware
BE->Redis:check(105)
Redis->BE:Response(106)
alt  #Pink if data  exist on redis,
BE->FE:Response data(107)
FE->Customer:Response data(108)
else
BE->Database:get data in database(109)
Database->BE:(110)
BE->FE:Response data(111)
FE->Customer:Response data(112)
end 
===See Transaction Eend ===


==Transaction End==

