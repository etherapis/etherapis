// EthereumStats is a UI component that displays the ethereum network status infos
// on the main navigation bar. It lists the numer of connected peers, the current
// block number and the arrival time of the last block.
//
// Properties:
//   ajax:    API requester          - func(url string, success func(data json.RawMessage))
//   apiurl:  API ethereum endpoint  - string
//   refresh: Stats refresh interval - time in milliseconds
var EthereumStats = React.createClass({
  // getInitialState sets the zero values of the component.
  getInitialState: function() {
    return {
      syncing: false,
      peers:   [],
      head:    null,
    };
  },
  // componentDidMount is invoked when the status component finishes loading. It
  // starts the periodical refresh of its internal state based on the backend API.
  componentDidMount: function() {
    this.refreshSyncing(); setInterval(this.refreshSyncing, this.props.refresh);
    this.refreshPeers(); setInterval(this.refreshPeers, this.props.refresh);
    this.refreshHead(); setInterval(this.refreshHead, this.props.refresh);
  },
  // refreshSyncing retrieves the sync status/progress of the Ethereum client.
  refreshSyncing: function() {
    this.props.ajax(this.props.apiurl + "/syncing", function(syncing) { this.setState({syncing: syncing}); }.bind(this));
  },
  // refreshPeers retrieves the list of connected peers and injects them into
  // the local state for sub-component interpretation.
  refreshPeers: function() {
    this.props.ajax(this.props.apiurl + "/peers", function(data) { this.setState({peers: data}); }.bind(this));
  },
  // refreshHead retrieves the current head block of the Ethereum blockchain.
  refreshHead: function() {
    this.props.ajax(this.props.apiurl + "/head", function(head) { this.setState({head: head}); }.bind(this));
  },
  // render flattens the Ethereum stats into UI report objects.
  render: function() {
    return (
      <div className="navbar-inner" style={{height: "50px"}}>
        <div className="navbar-right">
          <PeerCounter peers={this.state.peers}/>
          <BlockNumber block={this.state.head}/>
          <BlockTimer block={this.state.head}/>
        </div>
        <SyncProgress syncing={this.state.syncing} head={this.state.head}/>
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
//   syncing: Ethereum syncing progress
//   head:    Head block (to decide if initial sync or not)
var SyncProgress = React.createClass({
  render: function() {
    if (this.props.syncing == false) {
      return null;
    }
    var start   = parseInt(this.props.syncing.startingBlock, 16);
    var height  = parseInt(this.props.syncing.highestBlock, 16);
    var current = parseInt(this.props.syncing.currentBlock, 16);
    var pulled  = parseInt(this.props.syncing.pulledStates, 16)
    var known   = parseInt(this.props.syncing.knownStates, 16);

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
//   peers: List of peer data - []Peer
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
//   block: Ethereum block to display the number of - Block
var BlockNumber = React.createClass({
  render: function() {
    return (
      <span className="navbar-text">
        <i className="fa fa-database"></i> {this.props.block ? parseInt(this.props.block.number, 16) : 0} blocks
      </span>
    );
  }
});

// BlockTimer is a UI component that displays the time elapsed since s particualar
// Ethereum block was mined.
//
// Properties:
//   block: Ethereum block to display the elapsed time since - Block
var BlockTimer = React.createClass({
  render: function() {
    return (
      <span className="navbar-text">
        <i className="fa fa-clock-o"></i> {this.props.block ? moment.unix(this.props.block.timestamp).fromNow() : "never synced"}
      </span>
    );
  }
});
