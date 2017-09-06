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

## [PaytureAPI Services](#apiServises)
* [API Service](#api)
* [Ewallet Service](#ewallet)
* [InPay Service](#inpay)
* [DigitalWallet Service](#digital)


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



Visit our [site](http://payture.com/) for more information.
You can find our contact [here](http://payture.com/kontakty/).