contract Channels {
	struct PaymentChannel {
		address from;
		address to;
		uint256 nonce;

		uint256 price;
		uint256 value;
		uint validUntil;

		bool valid;
	}
	mapping(bytes32 => PaymentChannel) public channels;

	event NewChannel(address indexed from, address indexed to, bytes32 channel, uint nonce, uint256 price);
	event Deposit(address indexed from, bytes32 indexed channel);
	event Redeem(bytes32 indexed channel, uint nonce);
	event Reclaim(bytes32 indexed channel);

	function makeChannelId(address from, address to) constant returns(bytes32) {
		return sha3(from, to);
	}

	function createChannel(address to, uint256 price) {
		bytes32 channel = makeChannelId(msg.sender, to);
		PaymentChannel ch = channels[channel];
		if(!ch.valid) {
			channels[channel] = PaymentChannel(msg.sender, to, 0, price, msg.value, block.timestamp + 7 days, true);
		}

		NewChannel(msg.sender, to, channel, ch.nonce, price);
	}

	// creates a hash using the recipient and value.
	function getHash(address from, address to, uint nonce, uint value) constant returns(bytes32) {
		return sha3(from, to, nonce,value);
	}

	// verify a message (receipient || value) with the provided signature
	function verifySignature(bytes32 channel, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) constant returns(bool) {
		PaymentChannel ch = channels[channel];
		return  ch.valid &&
		        ch.validUntil > block.timestamp &&
		        ch.from == ecrecover(getHash(ch.from, ch.to, nonce, value), v, r, s);
	}

	function verifyPayment(bytes32 channel, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) constant returns(bool) {
		if( !verifySignature(channel, nonce, value, v, r, s) ) return false;

		PaymentChannel ch = channels[channel];
		if( ch.nonce != nonce ) return false;

		return true;
	}

	// claim funds
	function claim(bytes32 channel, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) {
		if( !verifySignature(channel, nonce, value, v, r, s) ) return;
		
		PaymentChannel ch = channels[channel];
		
		if( ch.nonce != nonce ) return;
		
		if( value > ch.value ) {
			ch.to.send(ch.value);
			ch.value = 0;
		} else {
			ch.to.send(value);
			ch.value -= value;
		}

		Redeem(channel, ch.nonce);

		channels[channel].nonce++;
	}

	function deposit(bytes32 channel) {
		if( !isValidChannel(channel) ) throw;

		PaymentChannel ch = channels[channel]; 
		ch.value += msg.value;

		Deposit(msg.sender, channel);
	}

	// reclaim a channel
	function reclaim(bytes32 channel) {
		PaymentChannel ch = channels[channel];
		if( ch.value > 0 && ch.validUntil < block.timestamp ) {
			ch.from.send(ch.value);
			delete channels[channel];
		}
	}

	function getChannelValue(bytes32 channel) constant returns(uint256) {
		return channels[channel].value;
	}

	function getChannelNonce(bytes32 channel) constant returns(uint256) {
		return channels[channel].nonce;
	}

	function getChannelPrice(bytes32 channel) constant returns(uint256) {
		return channels[channel].price;
	}

	function getChannelOwner(bytes32 channel) constant returns(address) {
		return channels[channel].from;
	}

	function  getChannelValidUntil(bytes32 channel) constant returns(uint) {
		return channels[channel].validUntil;
	}
	function isValidChannel(bytes32 channel) constant returns(bool) {
		PaymentChannel ch = channels[channel];
		return ch.valid && ch.validUntil >= block.timestamp;
	}
}


