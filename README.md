# Gopher Gala - Ether APIs - [etherapis.io](http://etherapis.io)

With the advent of cloud computing, everyone is trying to create some sort of online service that you can access and interact with via an online API. As the popular iOS slogan goes *"there's an app for that"*, in the API world it's a similar mantra: *"there's an API for that"*, irrelevant what you're trying to build: machine translation, mapping service, email sending, video rendering and streaming, etc., you'll find a large or small company that does exactly that, and the only thing you need to do is to hook into their system.

However, the online service and API ecosystem has a less attractive side, a problem yet unsolved, which leads to a high barrier of entry for small projects: *trust*. New and geographically remote projects are inherently untrusted, which makes payment for said services very problematic.

## Pain point

There are two main types of API providers: large corporations and small to medium scale projects, the two of which have very different offers when it comes to paid APIs.

 * Large corporations (e.g Google, Amazon) require a consumer to only have a credit card on file, and whenever the user consumes some API, it is measured and counted against the user's quota. At the end of the payment cycle, the client is charged with the outstanding amount. This however has the drawbacks that I need to provide all my personal and payment information to the API provider, as well as there's a risk that I consume more than I would like (e.g. forgot to turn off the dang test VMs), causing a huge charge that I have to pay nonetheless.
 * Small and medium companies on the other hand do not have the necessary resources to "lend" API capacity for an entire month and only afterwards get their own costs covered. Additionally, they do not have the means to force clients with bad payment information to actually come through. These two monetary shortcomings force the smaller API providers into a pre-paid subscription model, where an API consumer needs to pay a potentially hefty sum up front, which also expires at the end of the payment cycle. Opposed to large corporations (which can be somewhat trusted), a user is not comfortable with giving away personal infos to unknown API providers from far away countries, neither is he comfortable with paying a lot of money up front, just to realize the API is unreliable, or worst the provider a scammer.

A prominent solution currently available is offered by [mashape](https://www.mashape.com/): an API marketplace, broker and payment gateway acting as an escrow service between providers and consumers. Providers register their API endpoints into mashape's centralized brokerage service (configuring an associated pricing model), which on the other hand exposes those hidden APIs to the general public via their own servers. As all API invocations pass through them, those can be individually authorized, accounted for and charged at the end of the payment cycle. The value mashape brings to the API ecosystem is a payment escrow service, where API providers and consumers don't have to trust each other, but rather they all trust mashape itself to do the right thing. In exchange, mashape charges a 20% flat fee on all transactions that go through its brokarage service.

The problem with the current API ecosystem is that it requires blind trust in one or both of the API participants and true pay-per-use can only be achieved by the big players in the industry. All the existing solutions to these problems take the shape of centralization, which again requires a mutually trusted party, but also introduce a potential point of failure, a potential bottleneck and last but not least, significant privacy concerns.

## Solution

The proposed solution that we're pushing for is the re-decentralization of online services and APIs (like they originally were meant to be), but in a fully trustless way, where neither of the API participants has to know the other party or even trust it with funds (pre-lent or pre-paid).

By binding the financial contracts to the Ethereum blockchain, the API consumer and provider can remain in charge or their communication and the payment execution is enforced and secured by the censorship resistant Ethereum blockchain. Payment is made on a per-call basis, resurrecting the dream of a *"pay for what you use"*: an API consumer only pays for what he actually uses, at the exact moment when he uses it; whereas an API provider can at any point redeem the payments made to his service, without the consumer being able to refuse the already authorized payments.

Payments are done over a secure, on the Ethereum blockchain, payment channel. The payment channel work in such a way that you won't burden the network with micro payments or yourself with unecassary fees. Signing of micropayments happen client side, off chain, by the consumer and are included in the authorisation header of each API call. The signed transaction can then be verified using the Ethereum state by the service provider (off chain). Once a micropayment has been signed by the consumer it can no longer be undone, this is great since both parties can't hold one another for ransom or extort the other party by not holding up their end. Payments can be verified using the following algorithm. Given that `H` is the hashed output of `(consumer || provider || nonce || amount)` and `S` being the signature of `H` using standard ECDSA with the secp256k1 curve, one can derive the sender (and therefor the proof) and verify that a payment channel is 1. owned by said signer; and 2. contains enough funds. Using `S` one can redeem a cheque and collect their reward for their provided service. This works great since we don't have to burden the network with useless transaction as we can make use a sliding window technique where we keep increasing `amount` provided for `H`. Once a cheque has been collected we increase the nonce invalidating **any** other transaction with the given nonce provided in `H`, this is to say that a particular signed cheque can only be claimed once, ever; it is therefore wise to obviously redeem the cheque with the highest possible `amount`.

## Use cases (recipies ala the current Ethereum website)

So, what could this decentralized and trustless API ecosystem support:

 * Stateless APIs: each call authorizes payment for its own data processing requirement
   * E.g. machine translation, geolocation, image rendering
 * Streaming APIs: payment is authorized for some amount of future processing, renewed until both parties desire
   * E.g. music streaming services, video rendering services
 * Private APIs: the system can handle the payment for privately negotaited and executed calls
   * E.g. business to business private API execution

## Value proposition

So why is this project better that anything out there?

 * Decentralized: there isn't a single point of failure, single point of bottleneck, or single point of trust
 * Trustless: nethier the API provider or consumer needs trust, everything is enforced by the blockchain
 * Private: interacting peers are anonymous, communication is direct between them, no middleman
 * Secure: Ethereum ensures all transactions are final, inalterable and non censurable

How much does such a setup cost?

 * Registering a new API, subscribing to an API, charging all accumulated payments is ~$0.0015 (blockchain fees)
 * Searching for APIs and making payments for API calls is completely free (done off chain)
 * Payment contract charges a service fee of 1% (opposed to the 20% of existing competitors)
