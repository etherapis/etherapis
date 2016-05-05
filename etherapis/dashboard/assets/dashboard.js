// Dashboard is the entire webpage content of the EtherAPIs dahsboard. It should
// more or less be fully self contained, with perhaps a few external properties.
//
// Properties:
//	 apiurl: Base path for the API access - string
var Dashboard = React.createClass({
	// getInitialState sets the zero values of the component.
	getInitialState: function() {
		return {
			section: "index",
			account: "",
			console: false,

			server: null,       // State data injected by the Ether APIs backend
			alive:  Date.now(), // Timestamp (ms) of the last state update

			footer: "52px",
		};
	},
	// componentDidMount is invoked when the status component finishes loading. It
	// loads up the Ethereum account used by the Ether APIs server.
	componentDidMount: function() {
		// Connect to the backend to receive state updates
		var server = new EtherAPIs("ws://" + window.location.href.split("/")[2] + this.props.apiurl + "/", function(state) {
			// Sanity check the the active address wans't deleted
			if (state.accounts[this.state.account] === undefined) {
				this.state.account = Object.keys(state.accounts)[0]
			}
			this.setState({server: state, alive: Date.now()});
		}.bind(this));

		// Start a heartbeat mechanism to notice backend failures
		setInterval(function() {
			if (Date.now() - this.state.alive > 3000) {
				// Drop any previous state to clear the UI and try to reconnect
				if (this.state.server != null) {
					this.setState({server: null});
				}
				server.reconnect();

				// Bumt the timer so we're not reconnecting like crazy
				this.state.alive = Date.now();
			}
		}.bind(this), 1000);
	},

	loadIndex:      function() { this.setState({section: "index"}); },
	loadAccount:    function() { this.setState({section: "account"}); },
	loadProvider:   function() { this.setState({section: "provider"}); },
	loadSubscriber: function() { this.setState({section: "subscriber"}); },
	loadMarket:     function() { this.setState({section: "market"}); },

	// toggleConsole shows or hides the Geth console.
	toggleConsole: function() { this.setState({console: !this.state.console}); },

	// switchAccount switches out the currently active account to the one specified.
	switchAccount: function(address) {
		this.setState({account: address});
	},
	render: function() {
		// Do not render anything until we get our state from the backend
		if (this.state.server == null) {
			return (
				<div className="row">
					<div className="col-lg-3"></div>
					<div className="col-lg-6">
						<div className="well" style={{textAlign: "center"}}>
						  <span><i className="fa fa-refresh fa-spin"></i>&nbsp;&nbsp;Connecting to Ether APIs backend...</span>
						</div>
					</div>
				</div>
			)
		}
		// State available, render the full dashboard
		return (
			<div style={{height: "100%"}}>
				<nav className="navbar navbar-fixed-top">
					<div className="container">
						<div className="navbar-header">
							<a className="navbar-brand" href="#" onClick={this.loadIndex}>Ether APIs</a>
						</div>
						<div className="navbar-collapse collapse">
							<ul className="nav navbar-nav">
								<li><a href="#" onClick={this.loadAccount}><i className="fa fa-user"></i> Account</a></li>
								<li><a href="#" onClick={this.loadProvider}><i className="fa fa-cloud-upload"></i> Provided</a></li>
								<li><a href="#" onClick={this.loadSubscriber}><i className="fa fa-cloud-download"></i> Subscribed</a></li>
								<li><a href="#" onClick={this.loadMarket}><i className="fa fa-shopping-basket"></i> Market</a></li>
							</ul>
							<div className="navbar-inner" style={{height: "50px", float: "right"}}>
								<div className="navbar-right">
									<span className="navbar-text">
										<a href="#" onClick={this.toggleConsole}><i className="fa fa-bug"></i></a>
									</span>
								</div>
							</div>
							<EthereumStats data={this.state.server.ethereum}/>
						</div>
					</div>
				</nav>
				<div className="container" style={{minHeight: "100%", margin: "0 auto -" + this.state.footer}}>
					<StaleChainWarning head={this.state.server.ethereum.head} threshold={180}/>

					<Tutorial hide={this.state.section != "index"}/>
					<Accounts hide={this.state.section != "account"} apiurl={this.props.apiurl + "/accounts"} explorer={"http://testnet.etherscan.io/"} accounts={this.state.server.accounts} active={this.state.account} switch={this.switchAccount}/>
					<Provider hide={this.state.section != "provider"} apiurl={this.props.apiurl + "/services"} accounts={this.state.server} active={this.state.account} services={this.state.server.services} loadaccs={this.loadAccount} switch={this.switchAccount}/>
					<Subscriber hide={this.state.section != "subscriber"}/>
					<Market hide={this.state.section != "market"} apiurl={this.props.apiurl + "/market"} market={this.state.server.market}/>
					<div style={{height: this.state.footer}}></div>
					<Console hide={!this.state.console} toggle={this.toggleConsole} apiurl={this.props.apiurl + "/console"} />
				</div>
				<TestnetFooter hide={false} height={this.state.footer}/>
			</div>
		);
	}
});
window.Dashboard = Dashboard // Expose the component

// StaleChainWarning is a UI component displaying a warning message when the local
// blockchain is older than a certain threshold.
//
// Properties:
//	 head: Head block based on which to show or hide the warning
var StaleChainWarning = React.createClass({
	render: function() {
		if (moment().unix() < this.props.head.time + this.props.threshold) {
			return null
		}
		return (
			<div className="row">
				<div className="col-lg-3"></div>
				<div className="col-lg-6">
					<div className="alert alert-warning text-center" role="alert">
						<i className="fa fa-spinner fa-spin"></i> <strong>Local state too old.</strong> Please wait for network sync...
					</div>
				</div>
			</div>
		)
	}
});

// TestnetFooter is a UI component that displays a small information message when
// the node is connected to the Morden test network.
var TestnetFooter = React.createClass({
	render: function() {
		if (this.props.hide) {
			return null
		}
		return (
			<div style={{width: "100%", height: this.props.height}}>
					<div className="alert alert-info text-center" role="alert" style={{marginBottom: 0}}>
						<i className="fa fa-info-circle"></i> This instance is currently configured to use the <strong>test network</strong>. Consider everything here a consequence free playground.
					</div>
			</div>
		)
	}
});

// Console is a UI component that displays a Geth console attached to the
// underlying node.
var Console = React.createClass({
	// getInitialState sets the zero values of the component.
	getInitialState: function() {
		return {
			hidden:      true,
			initialized: false,
		};
	},
	// componentWillReceiveProps runs when the console is being toggled on or off.
	// The method resets the init flag to false to force a initializing any newly
	// created elements.
	componentWillReceiveProps: function(props) {
		if (props.hide != this.state.hidden) {
			this.setState({hidden: props.hide, initialized: false});
		}
	},
	// componentDidUpdate is invoked when the status component changes after a load. It
	// initializes the console to a draggable and resizable component, also starting
	// up the actual internal console plugin.
	componentDidUpdate: function() {
		// Only initialize the console if we've just displayed it
		if (this.props.hide || this.state.initialized) {
			return
		}
		this.setState({initialized: true});

		// Make the console UI a floating live dialog window
	  $("#console").draggable();
		$("#console").resizable();

		// Initialize the content with the console plugin
		var jqconsole = $('#console-content').jqconsole('Welcome to the Geth console!\n\n', '⡢ ');
		var startPrompt = function () {
			jqconsole.Prompt(true, function (input) {
				var form = new FormData();
				form.append("command", input);
				form.append("action", "exec");

				$.ajax({type: "POST", url: this.props.apiurl, cache: false, data: form, processData: false, contentType: false,
					success: function(data) {
						jqconsole.Write(data + '\n', 'jqconsole-output');
						startPrompt();
					}.bind(this),
					error: function(xhr, status, err) {
						jqconsole.Write(xhr.responseText, 'text-danger');
						startPrompt();
					}.bind(this),
				});
			}.bind(this));
		}.bind(this);

		// Replace the default indenting tab with command completion
		jqconsole.SetControlKeyHandler(function(e) {
			if (e.which === 9) {
				var form = new FormData();
				form.append("command", jqconsole.GetPromptText());
				form.append("action", "hint");

				$.ajax({type: "POST", url: this.props.apiurl, cache: false, data: form, processData: false, contentType: false,
					success: function(data) {
						if (/\r|\n/.exec(data)) {
							jqconsole.Write('\n' + data + '\n', 'text-muted');
						} else {
							jqconsole.SetPromptText(data);
						}
						startPrompt();
					}.bind(this),
					error: function(xhr, status, err) {
						jqconsole.Write(xhr.responseText, 'text-danger');
						startPrompt();
					}.bind(this),
				});
				return false;
			}
		}.bind(this));

		// Close the console on Ctrl+D
		jqconsole.RegisterShortcut('D', function() {
			this.props.toggle();
		}.bind(this));

		startPrompt();
	},
	// render generates the floating console component.
	render: function() {
		if (this.props.hide) {
			return null
		}
		return (
			<div id="console" className="well">
				<div id="console-content"></div>
			</div>
		)
	}
});
