# Go-Payture-official

This is Offical Payture API for Go. We're try to make this as simple as possible for you! Explore tutorial and get started. Please note, you will need a Merchant account,  contact our support to get one. 
Here you can explore how to use our API functions!

## Install



And include to your project:
```go
import 	"github.com/Payture/Go-Payture-official/payture"
```

## Payture API tutorial
## [Sequence of usage PaytureAPI](#newMerchant)

 * [Creating merchant account](#newMerchant)
 * [Get required API Service](#accessToAPI)
 * [Send transaction request](#expandTransaction)


## First Step - Creating Merchant Account <a id="newMerchant" ></a>
For getting access to API usage just create the instance of Merchant struct, set the name of the host, name of your account and your account password.  Suppose that you have Merchant account with  name: Key = "MyMerchantAccount" and password Password = "MyPassword".

Pass the 'https://sandbox.payture.com' for test as the name of Host.
```go
merchant := payture.Merchant{Key: "MyMerchantAccount", Password: "MyPassword", Host: "https://sandbox.payture.com"}
```
We're completed the first step! Go next!
***
Please note, that  Key = "'MyMerchantAccount" and Password = "MyMerchantAccount"  - fake, [our support](http://payture.com/kontakty/) help you to get one!
***

## Second Step - Get required API Service <a id="accessToAPI" ></a>
At this step you just call one of following methods (which provide proper API type for you) and pass in the Merchant struct: 
* API (this is PaytureAPI)
```go
apiManager := payture.APIService(merchant)
```
* InPay (this is PaytureInPay)
```go
inpayManager := payture.InPayService(merchant)
```
* EWallet (this is PaytureEWallet)
```go
ewManager := payture.EwalletService(merchant)
```
* Apple (this is PaytureApplePay)
```go
appleManager := payture.AppleService(merchant)
```
* Android (this is PaytureAndroidPay)
```go
androidManager := payture.AndroidService(merchant)
```
* MasterPass (this is Payture MasterPass)
```go
mpManager := payture.MasterPassService(merchant)
```

Result of this methods is one of  the supplied 'Manager' struct. See [here]() for details about this 'Managers'. After this step you get access to api methods - just call one of it!

## Third Step - Send transaction request <a id="extpandTransaction" ></a>
This is the most difficult step, but you can do it!
In the previous step we get the appopriate API Service. And now we'll starting explore how to use it.

## [PaytureAPI Services](#apiServises)
* [API Service](#api)
* [Ewallet Service](#ewallet)
* [InPay Service](#inpay)
* [DigitalWallet Service](#digital)

## [PaytureAPI Manager's structs](#apiManagers)
* [OrderManager](#ordManager)
* [APIManager](#apiManager)
* [EwalletManager](#ewManager)
* [InpayManager](#inpManager)
* [DigitalWalletManager](#digManager)

## OrderManager <a id="ordManager">
OrderManager is struct that provide access to API methods in Payture system allowing to manage state of payment ([see description of this methods](#ordManagerMethods)):

| Method's name  | Definition                                                            |
| -------------- | --------------------------------------------------------------------- |
| [Unblock](#unblock)        | Unblock fund on card, that were block earlier.                        |
| [Refund](#refund)         | Return fund to card, that were already charge-off earlier.            |
| [Charge](#charge)        | The second step in two-stage charge-off operation.                    |
| [PayStatus](#paystat)       | Return current status of payment by OrderId (For EWallet and InPay).  |
| [GetState](#getstate)        | Return current status of payment by OrderId (For API).                |

You don't need create directly this struct, because other kind of api managers contain this struct internally and have accesses to all these methods too, but provide some extra functionallity like pay/block methods.

### OrderManager's methods description <a id="ordManagerMethods">
These methods usage is very simple.
The Result of calling any of these methods is [OrderResponse](#ordResponse) struct, which contains parsed information that was returned by Payture server as response on you request.
All of this methods accept the [Payment](#payment) type struct.
So let's go deeper.
Suppose we're have the 'manager' variable that can be one of existing "Manager's" struct (APIManager, InpayManager, EWalletManager, DigitalManager), and we're have 'order' variable - that just Payment type:

```go
order := payture.Payment{Amount: "123000", OrderId: "TestOrder-0000001-000001"}
```

#### Unblock <a id="unblock">
You can call this method only in case if you done Block before. And you must provide exactly the same OrderId as you used in Block operation.
```go
orderResponse, err := manager.Unblock(order)
```

#### Refund <a id="refund">
You can call this method only in case if you done [Charge](#charge) before. Provide exactly the same OrderId as you used in Charge operation.
```go
orderResponse, err := manager.Refund(order)
```

#### Charge <a id="charge">
Like Unblock, you call this method in case then you done Block before. Provide exactly the same OrderId as you used in Block operation. Common usage of this method is the second step in two-stage charged-off funds from card.
```go
orderResponse,err := manager.Charge(order)
```

#### PayStatus <a id="paystat">
You call this method only for Ewallet and Inpay services (on EWalletManager and InpayManager structs respectively). Call this in any time for get current state of transaction in which you intrested.
```go
orderResponse,err := manager.PayStatus(order)
```

#### GetState <a id="getstate">
You call this method only for API service (on APIManager struct). Call this in any time for get current state of transaction in which you intrested.
```go
orderResponse, err := manager.GetState(order)
```


## APIManager <a id="apiManager">
APIManager is struct that provide simple access to API methods in Payture system allowing  full payment management - internally it contains [OrderManager](#ordManager) struct and supplied additional methods for pay/block operations

| Method's name  | Definition                                              |
| -------------- | ------------------------------------------------------- |
| [Pay](#apiPay)            | Method for one-stage charged-off fund.                  |
| [Block](#apiBlock)          | Fist step in two-stage  charged-off operation.          |

### APIManager's methods description <a id="apiManagerMethods">
Suppose, we're already create the 'merchant' variable. And now create the instanse of APIManager type:
```go
apiManager := payture.APIService(merchant)
card := payture.Card{EMonth: "12", EYear: "21", CardHolder: "pety petryshkin", SecureCode: "123"}
order := payture.Payment{Amount: "123000", OrderId: "TestOrder-0000001-000001"}
info := payture.PayInfo{Card: card, Order: paym, PAN: "4111111111111112"}
additional := payture.CustParams{CustomerFields: make(map[string]string)}
additional.CustomerFields["Product"] = "Lime"

```
After that we're got access to Pay/Block methods and all [methods of OrgerManager struct](#ordManager)
The result of Pay/Block operations are [APIResponses](#apiResponse) struct
#### Pay <a id="apiPay">
Pay method performs one-stage charged-off money from card. After that you can call [Refund](#refund) operation on purpose to return funds to card from which they were charged earlier.
```go
apiResponse, err := apiManager.Pay(paym, info, additional, "GoCistomer", "")
```

#### Block <a id="apiBlock">
Common usage of Block method as the fist step in two-stage charged-off fund. After this call [Charge](#charge) or you can call [Unblock](#unblock)
```go
apiResponse, err := apiManager.Block(paym, info, additional, "GoCistomer", "")
``` 

## InpayManager <a id="inpManager">
InpayManager is struct that provide simple access to API methods in Payture system allowing payment management - internally it also contains OrderManager struct and supplied one additional method - Init. With this method you can create session for pay or block operation and by created Id of session receive access to payment template where customer must enter card's infomation to complete payment.

| Method's name | Definition                                                              |
| ------------- | ----------------------------------------------------------------------- |
| Init          | Method for create payment session for specified SessionType, template and language.  |
| Pay           | Provide access to payment template for entering required information for compete pay or block operation.  |


Inpay Init method in action:
```go
inpayManager := payture.InPayService(merchant)
order := payture.Payment{Amount: "123000", OrderId: "TestOrder-0000001-000001"}
sessionType := SESSION_PAY
tag := "TemplateTag"
lang := "RU"
ip := "127.0.0.12"
additionals := payture.CustParams{}
urlReturn := "https://returnUrl.ret"
initResponse, err := inpayManager.Init(order, sessionType, tag, lang, ip, urlReturn, additionals)
```

## EwalletManager <a id="ewManager">
EwalletManager is struct with great functionality. With this manager you can use all methods from OrderManager struct, you can create session to get template for entering card's infomation (like Init in InpayManager). But main distinction is needs to provide customer's information for payments. For this we're created methods for register customers in our system and manage information about these customers, and you can attach(delete) cards to(from) already registered customers. Futher more you can use this registered card's and customer's account for pay/block session.

| Method's name  | Definition                                                                                |
| -------------- | ----------------------------------------------------------------------------------------- |
| [Init](#ewInit)       | Method for create payment session or session for addition card to specified customer.     |
| [CustomerRegister](#register) | Created customer's account in Payture system.                                     |
| [CustomerUpdate](#update)     | Changed some infomation for specified customer.                                   |
| [CustomerDelete](#delete)     | Delete early added customer from Payture system.                                  |
| [CustomerCheckBylogin](#checkbylog)       | Check where or not customer registered in Payture system at the moment. |
| [CustomerCheck](#checkcust)       | Check where or not customer registered in Payture system at the moment.       |
| [CardAdd](#add)           | Attached card to specified customer.                                                  |
| [CardActivate](#activate) | Activate card that was attached to customer, in case then auto activate not allowed.  |
| [CardRemove](#remove)     | Removed early added card to specified customer.                                       |
| [GetCardList](#getlist)   | Returns all cards attached to current customer at the moment.                         |
| [Pay](#ewpay)         | Method for one-stage charged-off fund.                                                    |
| [SendCode](#sendcode) | Send code for additional secure reasons.                                                  |

### EwalletManager's methods description <a id="ewalletManagerMethods">
```go
ewManager := payture.EwalletService(merchant)
```
#### Init <a id="ewInit">
Init method is intended to create session for futher process pay or block operation. After session was created Payture server returned response with SessionId, which you need to used in Pay method to complete transaction.
```go
order := payture.Payment{Amount: "123000", OrderId: "TestOrder-0000001-000001"}
sessionType := SESSION_PAY //  you can use instead SESSION_BLOCK or SESSION_ADD constants
custLogin := "customerLog"
tag := "TamplateTag"
lang := "RU"
ip := "127.0.0.12"
cardId := "40252318-de07-4853-b43d-4b67f2cd2077"
initResponse, err := ewManager.Init(custLogin, sessionType, ip, order Payment, tag, lang, cardId)
```

#### CustomerRegister <a id="register">
CustomerRegister method is intended to create customer's account in Payture system. After account was created you get [CustomerResponse](#customerResponse) struct with parsed response from Payture server.
```go
customer := payture.Customer{VWUserLgn : "firstTestCust", VWUserPsw : "df578fger8", PhoneNumber : "78888845", Email : "firstTestCust@testtest.rt"}
customerResponse, err := ewManager.CustomerRegister(customer)
```
#### CustomerUpdate <a id="update">
CustomerUpdate method is intended to update customer's account in Payture system. After account was updated you get [CustomerResponse](#customerResponse) struct with parsed response from Payture server.
```go
customer := payture.Customer{VWUserLgn : "firstTestCust", VWUserPsw : "df578fger8", PhoneNumber : "55544444"}
customerResponse, err := ewManager.CustomerUpdate(customer)
```
#### CustomerDelete <a id="delete">
CustomerDelete method is intended to delete customer's account in Payture system. After account was deleted you get [CustomerResponse](#customerResponse) struct with parsed response from Payture server.
```go
custLogin := "firstTestCust"
customerResponse, err := ewManager.CustomerDelete(custLogin)
```
#### CustomerCheckBylogin <a id="checkbylog">
CustomerCheckBylogin method is intended to check whether or not customer's account already registered in Payture system. This method recieved only customer's login. After checking you get [CustomerResponse](#customerResponse) struct with parsed response from Payture server.
```go
custLogin := "firstTestCust"
customerResponse, err := ewManager.CustomerCheckBylogin(custLogin)
```

#### CustomerCheck <a id="checkcust">
CustomerCheck method is intended to check whether or not customer's account already registered in Payture system. For use this method you need pass in Customer struct(login and password are required other are optional). After checking you get [CustomerResponse](#customerResponse) struct with parsed response from Payture server.
```go
customer := payture.Customer{VWUserLgn : "firstTestCust", VWUserPsw : "df578fger8"}
customerResponse, err := ewManager.CustomerCheck(customer)
```

#### CardAdd <a id="Add">
CardAdd method is intended to add specified card to already registered customer in Payture system. For use this method you need pass in Customer's account login and instance of [NotRegisteredCard](#notregcard) struct.
```go
custlogin :=  "firstTestCust"
cardForAdd := payture.NotRegisteredCard{CardHolder : "Test CardHolder", EMonth : "05", EYear : "22", SecureCode : "111", CardNumber "4111111111111112"}
cardResponse, err := ewManager.CardAdd(custlogin, cardForAdd)
```
#### CardActivate <a id="activate">
CardAdd method is intended to activate card that was added. For use this method you need pass in Customer's account login, card's Id and amount for blocking operation - then activation process was completed this amount will be unblocked.
```go
custlogin :=  "firstTestCust"
cardId := "40252318-de07-4853-b43d-4b67f2cd2077"
amount := "1001"
cardResponse, err := ewManager.CardActivate(custlogin, cardId, amount)
```
#### CardRemove <a id="remove">
CardRemove method is intended to remove card that was added for specified customer. For use this method you need pass in Customer's account login, card's Id.
```go
custlogin :=  "firstTestCust"
cardId := "40252318-de07-4853-b43d-4b67f2cd2077"
cardResponse, err := ewManager.CardRemove(custlogin, cardId)
```

#### GetCardList <a id="getlist">
GetCardList method is intended to get all cards that were added for specified customer. For use this method you need pass in Customer's account login.
```go
custlogin :=  "firstTestCust"
httpResponse, err := ewManager.GetCardList(custlogin)
```
#### Pay <a id="ewpay">
##### PayRegCard

```go
custlogin :=  "firstTestCust"
cardId := "40252318-de07-4853-b43d-4b67f2cd2077"
securecode := "111"
order := payture.Payment{Amount: "123000", OrderId: "TestOrder-0000001-000001"}
ip := "127.0.0.12"
confirmCode := "" // required only in case additional check
custFields := make(map[string]string)
payResponse := ewManager.PayRegCard(custLogin, cardId, secureCode, order, ip, confirmCode, custFields)
```

##### PayNoRegCard
```go
custlogin :=  "firstTestCust"
card := payture.NotRegisteredCard{CardHolder : "Test CardHolder", EMonth : "05", EYear : "22", SecureCode : "111", CardNumber "4111111111111112"}
order := payture.Payment{Amount: "123000", OrderId: "TestOrder-0000001-000001"}
ip := "127.0.0.12"
confirmCode := "" // required only in case additional check
custFields := make(map[string]string)
payResponse := ewManager.PayRegCard(custLogin, card, order, ip, confirmCode, custFields, true)
```

#### SendCode <a id="sendcode">
SendCode method is intended to send secure code during the payment processing to specified customer.
```go
custlogin :=  "firstTestCust"
httpResponse, err := ewManager.SendCode(custlogin)
```


## DigitalWalletManager <a id="digManager">
DigitalWalletManager is struct that consist of following 'managers' (we're suppose - as before that - 'merchant' variable already created at moment):

* AppleManager - is intended to get correct accesses for pay/block methods for Apple system. For creation instance of this struct invoke AppleService method:
```go
appleManager := payture.AppleService(merchant)
```
* AndroidManager - is intended to get correct accesses for pay/block methods for Android system.
 For creation instance of this struct invoke AppleService method:
```go
androidManager := payture.AndroidService(merchant)
```
* MasterPassManager - is intended to get correct accesses for pay/block methods for Master Pass system.
 For creation instance of this struct invoke AppleService method:
```go
mpManager := payture.MasterPassService(merchant)
```

All of these  provided accesses to methods:

| Method's name  | Definition                                                |
| -------------- | --------------------------------------------------------- |
| [Pay](#digitpay)            | Method for one-stage charged-off fund.                    |
| [Block](#digblock)          | Fist step in two-stage  charged-off operation.            |

### DigitalWalletManager's methods description <a id="appleManagerMethods">
For all of above services the way of calling specified methods are the same. Result of calling these methods are instance of DigitalResponse struct which contains parsed response from Payture server.

#### Pay <a id="digitpay">

```go
order := payture.Payment{Amount : "1000", OrderId : "43371575325653242457767057674540671"}
token := "abcdefg"
secureCode := "125"
//for apple and android - no needs for securecode
digitalResponse, err := digitalManager.Pay(order, token, nil) 

//for master pass secure code is required
digitalResponse, err := digitalManager.Pay(order, token, secureCode) 
```

#### Block <a id="digitblock">

```go
order := payture.Payment{Amount : "1000", OrderId: "43371575325653242457767057674540671"}
token := "abcdefg"
secureCode := "125"
//for apple and android - no needs for securecode
digitalResponse, err := digitalManager.Block(order, token, nil) 

//for master pass secure code is required
digitalResponse, err := digitalManager.Block(order, token, secureCode) 
```
Visit our [site](http://payture.com/) for more information.
You can find our contact [here](http://payture.com/kontakty/).