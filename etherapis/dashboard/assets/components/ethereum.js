// EthereumStats is a UI component that displays the ethereum network status infos
// on the main navigation bar. It lists the numer of connected peers, the current
// block number and the arrival time of the last block.
//
// Properties:
//	 ajax:		API requester					- func(url string, success func(data json.RawMessage))
//	 apiurl:	API ethereum endpoint	- string
//	 refresh: Stats refresh interval - time in milliseconds
var EthereumStats = React.createClass({
	render: function() {
		return (
			<div className="navbar-inner" style={{height: "50px"}}>
				<div className="navbar-right">
					<PeerCounter peers={this.props.data.peers}/>
					<BlockNumber block={this.props.data.head}/>
					<BlockTimer block={this.props.data.head}/>
				</div>
				<SyncProgress syncing={this.props.data.syncing} head={this.props.data.head}/>
			</div>
		);
	}
});
window.EthereumStats = EthereumStats // Expose the component

// SyncProgress is a UI component that displays the current sync status of the
// Ethereum node as a progress bar in the navbar. The bar is only displayed if
// there's actually a sync in progress.
//
// Properties:
//	 syncing: Ethereum syncing progress
//	 head:		Head block (to decide if initial sync or not)
var SyncProgress = React.createClass({
	render: function() {
		// Stop rendering if we're not syncing
		if (this.props.syncing.currentBlock >= this.props.syncing.highestBlock) {
			return null;
		}
		var start   = this.props.syncing.startingBlock;
		var height  = this.props.syncing.highestBlock;
		var current = this.props.syncing.currentBlock;
		var pulled  = this.props.syncing.pulledStates
		var known   = this.props.syncing.knownStates;

		var progress = 100 * (current - start) / (height - start);

		var label = "Sync: "
		if (this.props.head.number == 0) {
			label = "Initial sync: "
		}
		if (known - pulled > 0) {
			label += (height - current) + " blocks, " + (known - pulled) + "/" + known + " states left";
		} else {
			label += (height - current) + " blocks left";
		}
		return (
			<div style={{paddingTop: "15px"}}>
				<div className="progress" style={{position: "relative", height: "20px"}}>
					<div className="progress-bar progress-bar-striped active" role="progressbar" style={{width: progress + "%"}}>
						<span style={{position: "absolute", display: "block", width: "100%"}}><small>{label}</small></span>
					</div>
				</div>
			</div>
		);
	}
});

// PeerCounter is a UI component that displays the number of Ethereum peers we
// are currently connected to.
//
// Properties:
//	 peers: List of peer data - []Peer
var PeerCounter = React.createClass({
	render: function() {
		return (
			<span className="navbar-text">
				<i className="fa fa-rss"></i> {this.props.peers.length} peers
			</span>
		);
	}
});

// BlockNumber is a UI component that displays the current Ethereum block number.
//
// Properties:
//	 block: Ethereum block to display the number of - Block
var BlockNumber = React.createClass({
	render: function() {
		return (
			<span className="navbar-text">
				<i className="fa fa-database"></i> {this.props.block.number} blocks
			</span>
		);
	}
});

// BlockTimer is a UI component that displays the time elapsed since s particualar
// Ethereum block was mined.
//
// Properties:
//	 block: Ethereum block to display the elapsed time since - Block
var BlockTimer = React.createClass({
	render: function() {
		return (
			<span className="navbar-text">
				<i className="fa fa-clock-o"></i> {this.props.block ? moment.unix(this.props.block.time).fromNow() : "never synced"}
			</span>
		);
	}
});
