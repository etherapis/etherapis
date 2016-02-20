// EtherAPIs is a constructor to create a state linked to an Ether APIs backend,
// receiving state updates via websocket push notifications.
function EtherAPIs(url, onupdate) {
	this.url      = url;
	this.onupdate = onupdate;

	this.reconnect();
}

// update is called whenever the EtherAPIs backend sends a state diff. The first
// message is always the entire state to use as a starting point, and subsequent
// messages are lists or state siffts to apply.
EtherAPIs.prototype.update = function(message) {
	// If the origin state is not set yet, the first message is it
	if (this.state == null) {
		this.state = JSON.parse(message);
		this.alive = Date.now();
		return;
	}
	// Otherwise the message is a list of diffs, apply individually
	var diffs = JSON.parse(message).diffs;
	for (var i=0; i<diffs.length; i++) {
		var parent = this.state;
		while (diffs[i].path.length > 1) {
			parent = parent[diffs[i].path.shift()];
		}
		if (diffs[i].node != null) {
			parent[diffs[i].path[0]] = diffs[i].node;
		} else {
			delete parent[diffs[i].path[0]];
		}
	}
	this.alive = Date.now();

	// Run the registered callback for the updated state
	if (diffs.length > 0) {
		this.onupdate(this.state);
	}
};

// reconnect drops all state data recorded, disconnects any live sockets from the
// backend and rebuilds everything from scratch.
EtherAPIs.prototype.reconnect = function() {
	// Clean up any previous state
	this.state = null;
	if (this.socket !== undefined) {
		this.socket.close();
	}
	// Connect to the backend server and watch for state updates
	this.socket = new WebSocket(this.url);
	this.socket.onmessage = function (event) {
	  this.update(event.data);
	}.bind(this);
}
