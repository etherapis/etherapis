contract ServiceProviders {
	// PaymentModel is the possible payment models that proxies should be able to handle.
	enum PaymentModel {PerCall, PerData, PerTime}

	struct Terms {
		PaymentModel model;
		uint         price;
		uint         cancellation;
	}

	struct Service {
		uint id;
		address owner;
		string name;
		string endpoint;

		Terms terms;

		bool enabled;
		bool deleted;
		bool exist;
	}
	Service[] services;

	// serviceOwner modifier throws if the invoker does not match the owner field in the
	// service struct.
	modifier serviceOwner(uint serviceID) { if(services[serviceID].owner == msg.sender) _ }

	// NewService Event is fired when a new service has been created.
	event NewService(string indexed name, address indexed owner, uint serviceID);
	// UpdateService is fired when a service has been updated or flagged for deletion.
	event UpdateService(uint indexed serviceID);

	// services per user.
	mapping(address => uint[]) public userServices;
	// get the length per user.
	function userServicesLength(address addr) constant returns(uint) { return userServices[addr].length; }
	// get the total length of the services.
	function servicesLength() constant returns(uint)                 { return services.length; }
	// get a service and return a tuple.
	function getService(uint serviceID) constant returns(
		string name,
		address owner,
		string endpoint,
		uint model,
		uint price,
		uint cancellation,
		bool enabled,
		bool deleted
	) {
		Service service = services[serviceID];
		return (
			service.name,
			service.owner,
			service.endpoint,
			uint(service.terms.model),
			service.terms.price,
			service.terms.cancellation,
			service.enabled,
			service.deleted
		);
	}
	// delete a service.
	function deleteService(uint serviceID) serviceOwner(serviceID) {
		services[serviceID].deleted = true;
		UpdateService(serviceID);
	}
	// enable a service
	function enableService(uint serviceID) serviceOwner(serviceID) {
		services[serviceID].enabled = true;
		UpdateService(serviceID);
	}
	// disable a service
	function disableService(uint serviceID) serviceOwner(serviceID) {
		services[serviceID].enabled = false;
		UpdateService(serviceID);
	}

	// Add a new service.
	function addService(string name, string endpoint, uint model, uint price, uint cancellation) {
		Service service = services[services.length++];
		service.exist = true;
		service.enabled = false;
		service.id = services.length-1;
		service.owner = msg.sender;
		service.name = name;
		service.endpoint = endpoint;
		service.terms.model = PaymentModel(model);
		service.terms.price = price;
		service.terms.cancellation = cancellation;

		userServices[msg.sender].push(service.id);

		NewService(name, msg.sender, service.id);
	}
}

contract Subscriptions {
	struct Subscription {
		address from;
		ServiceProviders.Service service;
		uint256 nonce;
		uint256 value;

		bool cancelled;
		uint closedAt;

		bool exist;
	}
	mapping(bytes32 => Subscription) subscriptions;

	event NewSubscription(address indexed from, uint indexed serviceID, bytes32 subscriptionID, uint nonce);
	event Deposit(address indexed from, bytes32 indexed subscriptionID);
	event Redeem(bytes32 indexed subscriptionID, uint nonce);
	event Cancel(bytes32 indexed subscriptionID, uint closedAt);
	event Reclaim(bytes32 indexed subscriptionID);

	modifier isOwner(bytes32 subscriptionID) { if( subscriptions[subscriptionID].from != msg.sender) throw; }

	function makeSubscriptionID(address from, uint serviceID) constant returns(bytes32) {
		return sha3(from, serviceID);
	}

	// creates a hash using the recipient and value.
	function getHash(address from, uint serviceID, uint nonce, uint value) constant returns(bytes32) {
		return sha3(from, serviceID, nonce,value);
	}

	// verify a message (receipient || value) with the provided signature
	function verifySignature(bytes32 subscriptionID, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) constant returns(bool) {
		Subscription ch = subscriptions[subscriptionID];
		return ch.exist && ch.from == ecrecover(getHash(ch.from, ch.service.id, nonce, value), v, r, s);
	}

	function verifyPayment(bytes32 subscriptionID, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) constant returns(bool) {
		if( !verifySignature(subscriptionID, nonce, value, v, r, s) ) return false;

		Subscription ch = subscriptions[subscriptionID];
		if( ch.closedAt >= now ) return false;
		if( ch.nonce != nonce ) return false;

		return true;
	}

	// claim funds
	function claim(bytes32 subscriptionID, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) {
		if( !verifySignature(subscriptionID, nonce, value, v, r, s) ) return;

		Subscription ch = subscriptions[subscriptionID];

		if( ch.nonce != nonce ) return;
		if( ch.service.owner != msg.sender ) throw;

		if( value > ch.value ) {
			ch.service.owner.send(ch.value);
			ch.value = 0;
		} else {
			ch.service.owner.send(value);
			ch.value -= value;
		}

		Redeem(subscriptionID, ch.nonce);

		subscriptions[subscriptionID].nonce++;
	}

	function deposit(bytes32 subscriptionID) {
		if( !isValidSubscription(subscriptionID) ) throw;

		Subscription ch = subscriptions[subscriptionID];
		ch.value += msg.value;

		Deposit(msg.sender, subscriptionID);
	}

	function cancel(bytes32 subscriptionID) isOwner(subscriptionID) {
		Subscription ch = subscriptions[subscriptionID];

		uint closedAt = now + ch.service.terms.cancellation;

		subscriptions[subscriptionID].cancelled = true;
		subscriptions[subscriptionID].closedAt = closedAt;

		Cancel(subscriptionID, closedAt);
	}

	// reclaim a subscriptionID
	function reclaim(bytes32 subscriptionID) {
		Subscription ch = subscriptions[subscriptionID];
		if( ch.closedAt <= block.timestamp ) {
			ch.from.send(ch.value);
			delete subscriptions[subscriptionID];
		}
	}

	function getSubscription(bytes32 subscriptionID) constant returns(address from, uint serviceID, uint nonce, uint value, bool cancelled, uint closedAt) {
		Subscription ch = subscriptions[subscriptionID];

		return (ch.from, ch.service.id, ch.nonce, ch.value, ch.cancelled, ch.closedAt);
	}

	function getSubscriptionValue(bytes32 subscriptionID) constant returns(uint256) {
		return subscriptions[subscriptionID].value;
	}

	function getSubscriptionNonce(bytes32 subscriptionID) constant returns(uint256) {
		return subscriptions[subscriptionID].nonce;
	}

	function getSubscriptionOwner(bytes32 subscriptionID) constant returns(address) {
		return subscriptions[subscriptionID].from;
	}

	function getSubscriptionServiceID(bytes32 subscriptionID) constant returns(uint) {
		return subscriptions[subscriptionID].service.id;
	}

	function getSubscriptionClosedAt(bytes32 subscriptionID) constant returns(uint) {
		return subscriptions[subscriptionID].closedAt;
	}

	function isValidSubscription(bytes32 subscriptionID) constant returns(bool) {
		Subscription ch = subscriptions[subscriptionID];
		return ch.exist && ch.closedAt < block.timestamp;
	}
}

contract EtherAPIs is Subscriptions, ServiceProviders {
	mapping(address => bytes32[]) public userSubscriptions;
	function userSubscriptionsLength(address addr) constant returns(uint) {
		return userSubscriptions[addr].length;
	}

	function subscribe(uint serviceID) {
		bytes32 subscriptionID = makeSubscriptionID(msg.sender, serviceID);
		Subscription ch = subscriptions[subscriptionID];

		if( !ch.exist )  {
			Service service = services[serviceID];

			subscriptions[subscriptionID] = Subscription(msg.sender, service, 0, msg.value, false, 0, true);
			userSubscriptions[msg.sender].push(subscriptionID);

			NewSubscription(msg.sender, serviceID, subscriptionID, ch.nonce);
		}

	}
}
