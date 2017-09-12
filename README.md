# Go-Payture-official

This is Offical Payture API for Go. We're try to make this as simple as possible for you! Explore tutorial and get started. Please note, you will need a Merchant account,  contact our support to get one. 
Here you can explore how to use our API functions!

## Install



And include to your project:
```go
import 	"github.com/Payture/Go-Payture-official/payture"
```

## Payture API tutorial
Just look at the picture below to get some conception how this package is constructed. 

## [Sequence of usage PaytureAPI](#newMerchant)

 * [Creating merchant account](#newMerchant)
 * [Get required API Service](#accessToAPI)
 * [Send transaction request](#expandTransaction)


## First Step - Creating Merchant Account <a id="newMerchant" ></a>
For getting access to API usage just create the instance of Merchant struct, set the name of the host, name of your account and your account password.  Suppose that you have Merchant account with  name: Key = "MyMerchantAccount" and password Password = "MyPassword".

Pass the 'https://sandbox.payture.com' for test as the name of Host (first parameter).
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
OrderManager is struct that provide access to API methods in Payture system allowing to manage state of payment:

| Method's name  | Definition                                                            |
| -------------- | --------------------------------------------------------------------- |
| Unblock        | Unblock fund on card, that were block earlier.                        |
| Refund         | Return fund to card, that were already charge-off earlier.            |
| Charge         | The second step in two-stage charge-off operation.                    |
| PayStatus      | Return current status of payment by OrderId (For EWallet and InPay).  |
| GetState       | Return current status of payment by OrderId (For API).                |

Other kind of api managers contained this struct internally and have accesses to all these methods too, but provide some extra functionallity like pay/block methods.

## APIManager <a id="apiManager">
APIManager is struct that provide simple access to API methods in Payture system allowing  full payment management - internally it contains OrderManager struct and suppied additional methods for pay/block operations

| Method's name  | Definition                                              |
| -------------- | ------------------------------------------------------- |
| Pay            | Method for one-stage charged-off fund.                  |
| Block          | Fist step in two-stage  charged-off operation.          |

## InpayManager <a id="inpManager">
InpayManager is struct that provide simple access to API methods in Payture system allowing pament management - internally it also contains OrderManager struct and supplied one additional method - Init. With this method you can create session for pay or block operation and by created Id of session receive access to payment template where customer must enter card's infomation to complete payment.

| Method's name  | Definition                                                              |
| -------------- | ----------------------------------------------------------------------- |
| Init           | Method for create payment session for specified template and language.  |


## EwalletManager <a id="ewManager">
EwalletManager is struct with great functionality. With this manager you can use all methods from OrderManager struct, you can create session to get template for entering card's infomation (like Init in InpayManager). But main distinction is needs to provide customer's information for payments. For this we're created methods for register customers in our system and manage information about these customers, and you can attach(delete) cards to(from) already registered customers. Futher more you can use this registered card's and customer's account for pay/block session.

| Method's name  | Definition                                                                                |
| -------------- | ----------------------------------------------------------------------------------------- |
| Init           | Method for create payment session or session for addition card to specified customer.     |
| Register       | Created customer's account in Payture system.                                             |
| Update         | Changed some infomation for specified customer.                                           |
| Delete         | Delete early added customer from Payture system.                                          |
| Check          | Check where or not customer registered in Payture system at the moment.                   |
| Add            | Attached card to specified customer.                                                      |
| Activate       | Activate card that was attached to customer, in case then auto activate not allowed.      |
| Remove         | Removed early added card to specified customer.                                           |
| GetList        | Returns all cards attached to current customer at the moment.                             |
| Pay            | Method for one-stage charged-off fund.                                                    |
| SendCode       | Send code for additional secure reasons.                                                  |


## DigitalWalletManager <a id="digManager">
DigitalWalletManager is struct that consist of these 'managers':

* AppleManager 
* AndroidManager
* MasterPassManager

All of these  provided accesses to methods:

| Method's name  | Definition                                                |
| -------------- | --------------------------------------------------------- |
| Pay            | Method for one-stage charged-off fund.                    |
| Block          | Fist step in two-stage  charged-off operation.            |


Visit our [site](http://payture.com/) for more information.
You can find our contact [here](http://payture.com/kontakty/).