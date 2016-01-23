contract Channels {
	struct PaymentChannel {
		address from;
		address to;
		uint256 nonce;

		uint256 value;
		uint validUntil;

		bool valid;
	}
	mapping(bytes32 => PaymentChannel) public channels;

	event NewChannel(address indexed from, address indexed to, bytes32 channel, uint nonce);
	event Deposit(address indexed from, bytes32 indexed channel);
	event Claim(address indexed who, bytes32 indexed channel);
	event Reclaim(bytes32 indexed channel);

	function makeChannelName(address from, address to) {
		return sha3(from, to);
	}

	function createChannel(address to) {
		bytes32 channel = sha3(from, to);
		PaymentChannel ch = channels[channel];
		if(!ch.valid) {
			ch = PaymentChannel(msg.sender, to, 0, msg.value, block.timestamp + 1 days, true);
			channels[channel] = ch
		}

		NewChannel(msg.sender, to, channel, ch.nonce);
	}

	// creates a hash using the recipient and value.
	function getHash(bytes32 channel, address recipient, uint nonce, uint value) constant returns(bytes32) {
		return sha3(channel, recipient, nonce, value);
	}

	// verify a message (receipient || value) with the provided signature
	function verify(bytes32 channel, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) constant returns(bool) {
		PaymentChannel ch = channels[channel];
		return  ch.valid &&
		        ch.validUntil > block.timestamp &&
		        ch.from == ecrecover(getHash(channel, ch.to, nonce, value), v, r, s);
	}

	// claim funds
	function claim(bytes32 channel, uint nonce, uint value, uint8 v, bytes32 r, bytes32 s) {
		if( !verify(channel, nonce, value, v, r, s) ) return;
		
		PaymentChannel ch = channels[channel];
		if( ch.nonce != nonce ) return;
		
		if( value > ch.value ) {
			ch.to.send(ch.value);
			ch.value = 0;
		} else {
			ch.to.send(value);
			ch.value -= value;
		}

		// channel is no longer valid
		channels[channel].valid = false;

		Claim(ch.to, channel);
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
