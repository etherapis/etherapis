contract ServiceProviders {
	struct Terms {
		uint price;
		uint cancellationTime;
	}

	struct Service {
		uint id;
		address owner;
		string name;
		string endpoint;

		Terms terms;

		bool enabled;
		bool exist;
	}
	Service[] services;

	event NewService(string indexed name, address indexed owner, uint serviceId);
	event UpdateService(uint indexed serviceId);

	// services per user.
	mapping(address => uint[]) public userServices;
	// get the length per user.
	function userServicesLength(address addr) constant returns(uint) { return userServices[addr].length; }
	// get the total length of the services.
	function servicesLength() constant returns(uint)                 { return services.length; }
	// get a service and return a tuple.
	function getService(uint serviceId) constant returns(string name, address owner, string endpoint, uint price, uint cancellationTime, bool enabled) {
		Service service = services[serviceId];
		return (service.name, service.owner, service.endpoint, service.terms.price, service.terms.cancellationTime, service.enabled);
	}
	// Add a new service.
	function addService(string name, string endpoint, uint price, uint cancellationTime) {
		Service service = services[services.length++];
		service.exist = true;
		service.enabled = true;
		service.id = services.length-1;
		service.owner = msg.sender;
		service.name = name;
		service.endpoint = endpoint;
		service.terms.price = price;
		service.terms.cancellationTime = cancellationTime;

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

	event NewSubscription(address indexed from, uint indexed serviceId, bytes32 subscriptionId, uint nonce);
	event Deposit(address indexed from, bytes32 indexed subscriptionId);
	event Redeem(bytes32 indexed subscriptionId, uint nonce);
	event Cancel(bytes32 indexed subscriptionId, uint closedAt);
	event Reclaim(bytes32 indexed subscriptionId);

	modifier isOwner(bytes32 subscriptionId) { if( subscriptions[subscriptionId].from != msg.sender) throw; }

	function makeSubscriptionId(address from, uint serviceId) constant returns(bytes32) {
		return sha3(from, serviceId);
	}

	// creates a hash using the recipient and value.
	function getHash(address from, uint serviceId, uint nonce, uint value) constant returns(bytes32) {
		return sha3(from, serviceId, nonce,value);
	}

	// verify a message (receipient || value) with the provided signature
	function verifySignature(bytes32 subscriptionId, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) constant returns(bool) {
		Subscription ch = subscriptions[subscriptionId];
		return ch.exist && ch.from == ecrecover(getHash(ch.from, ch.service.id, nonce, value), v, r, s);
	}

	function verifyPayment(bytes32 subscriptionId, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) constant returns(bool) {
		if( !verifySignature(subscriptionId, nonce, value, v, r, s) ) return false;

		Subscription ch = subscriptions[subscriptionId];
		if( ch.closedAt >= now ) return false;
		if( ch.nonce != nonce ) return false;

		return true;
	}

	// claim funds
	function claim(bytes32 subscriptionId, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) {
		if( !verifySignature(subscriptionId, nonce, value, v, r, s) ) return;

		Subscription ch = subscriptions[subscriptionId];

		if( ch.nonce != nonce ) return;
		if( ch.service.owner != msg.sender ) throw;

		if( value > ch.value ) {
			ch.service.owner.send(ch.value);
			ch.value = 0;
		} else {
			ch.service.owner.send(value);
			ch.value -= value;
		}

		Redeem(subscriptionId, ch.nonce);

		subscriptions[subscriptionId].nonce++;
	}

	function deposit(bytes32 subscriptionId) {
		if( !isValidSubscription(subscriptionId) ) throw;

		Subscription ch = subscriptions[subscriptionId]; 
		ch.value += msg.value;

		Deposit(msg.sender, subscriptionId);
	}

	function cancel(bytes32 subscriptionId) isOwner(subscriptionId) {
		Subscription ch = subscriptions[subscriptionId];

		uint closedAt = now + ch.service.terms.cancellationTime;

		subscriptions[subscriptionId].cancelled = true;
		subscriptions[subscriptionId].closedAt = closedAt;

		Cancel(subscriptionId, closedAt);
	}

	// reclaim a subscriptionId
	function reclaim(bytes32 subscriptionId) {
		Subscription ch = subscriptions[subscriptionId];
		if( ch.closedAt <= block.timestamp ) {
			ch.from.send(ch.value);
			delete subscriptions[subscriptionId];
		}
	}

	function getSubscription(bytes32 subscriptionId) constant returns(address from, uint serviceId, uint nonce, uint value, bool cancelled, uint closedAt) {
		Subscription ch = subscriptions[subscriptionId];

		return (ch.from, ch.service.id, ch.nonce, ch.value, ch.cancelled, ch.closedAt);
	}

	function getSubscriptionValue(bytes32 subscriptionId) constant returns(uint256) {
		return subscriptions[subscriptionId].value;
	}

	function getSubscriptionNonce(bytes32 subscriptionId) constant returns(uint256) {
		return subscriptions[subscriptionId].nonce;
	}

	function getSubscriptionOwner(bytes32 subscriptionId) constant returns(address) {
		return subscriptions[subscriptionId].from;
	}

	function getSubscriptionServiceId(bytes32 subscriptionId) constant returns(uint) {
		return subscriptions[subscriptionId].service.id;
	}

	function getSubscriptionClosedAt(bytes32 subscriptionId) constant returns(uint) {
		return subscriptions[subscriptionId].closedAt;
	}

	function isValidSubscription(bytes32 subscriptionId) constant returns(bool) {
		Subscription ch = subscriptions[subscriptionId];
		return ch.exist && ch.closedAt < block.timestamp;
	}
}

contract EtherApis is Subscriptions, ServiceProviders {
	mapping(address => bytes32[]) public userSubscriptions;
	function userSubscriptionsLength(address addr) constant returns(uint) {
		return userSubscriptions[addr].length;
	}

	function subscribe(uint serviceId) {
		bytes32 subscriptionId = makeSubscriptionId(msg.sender, serviceId);
		Subscription ch = subscriptions[subscriptionId];

		if( !ch.exist )  {
			Service service = services[serviceId];

			subscriptions[subscriptionId] = Subscription(msg.sender, service, 0, msg.value, false, 0, true);
			userSubscriptions[msg.sender].push(subscriptionId);

			NewSubscription(msg.sender, serviceId, subscriptionId, ch.nonce);
		}

	}
}

