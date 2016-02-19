// Dashboard is the entire webpage content of the EtherAPIs dahsboard. It should
// more or less be fully self contained, with perhaps a few external properties.
//
// Properties:
//   apiurl: Base path for the API access - string
var Dashboard = React.createClass({
  // getInitialState sets the zero values of the component.
  getInitialState: function() {
    return {
      section:    "index",
      apiOnline:  true,
      apiFailure: "",
      needSync:   false,
      accounts:   [],
      account:    "",

      footer: "52px",
    };
  },
  // componentDidMount is invoked when the status component finishes loading. It
  // loads up the Ethereum account used by the Ether APIs server.
  componentDidMount: function() {
    this.refreshAccounts();
  },

  loadIndex:      function() { this.setState({section: "index"}); },
  loadProvider:   function() { this.setState({section: "provider"}); },
  loadSubscriber: function() { this.setState({section: "subscriber"}); },
  loadMarket:     function() { this.setState({section: "market"}); },
  loadAccount:    function() { this.setState({section: "account"}); },

  // refreshAccounts retrieves the current accounts configured, and if successful,
  // also loads up any associated data.
  refreshAccounts: function(selected) {
    this.apiCall(this.props.apiurl + "/accounts", function(accounts) {
      var primary = null;
      for (var i = 0; i < accounts.length; i++) {
        if (accounts[i] == selected) {
          primary = selected;
        }
      }
      if (primary == null) {
        primary = accounts[0];
      }
      this.setState({accounts: accounts.sort(), account: primary});
    }.bind(this));
  },
  // switchAccount switches out the currently active account to the one specified.
  switchAccount: function(account) {
    this.setState({account: account});
  },

  apiCall: function(url, success) {
    $.ajax({url: url, dataType: 'json', cache: false,
      success: function(data) {
        if (this.state.apiFailure == url) {
          this.setState({apiOnline: true, apiFailure: ""});
        }
        success(data)
      }.bind(this),
      error: function(xhr, status, err) {
        this.setState({apiOnline: false, apiFailure: url});
        console.error(url, status, err.toString());
      }.bind(this)
    });
  },
  needSync: function(needed) {
    this.setState({needSync: needed});
  },
  render: function() {
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
              <EthereumStats ajax={this.apiCall} apiurl={this.props.apiurl + "/ethereum"} refresh={1000} nosync={this.needSync}/>
            </div>
          </div>
        </nav>
        <div className="container" style={{minHeight: "100%", margin: "0 auto -" + this.state.footer}}>
          <BackendError hide={this.state.apiOnline}/>
          <NotSyncedWarning hide={!this.state.needSync || !this.state.apiOnline}/>

          <Tutorial hide={this.state.section != "index"}/>
          <Accounts hide={this.state.section != "account"} apiurl={this.props.apiurl} ajax={this.apiCall} explorer={"http://testnet.etherscan.io/address/"} accounts={this.state.accounts} active={this.state.account} switch={this.switchAccount} refresh={this.refreshAccounts}/>
          <Provider hide={this.state.section != "provider"} ajax={this.apiCall} apiurl={this.props.apiurl} address={this.state.accounts} />
          <Subscriber hide={this.state.section != "subscriber"}/>
          <Market hide={this.state.section != "market"} ajax={this.apiCall} apiurl={this.props.apiurl} interval={1000}/>
          <div style={{height: this.state.footer}}></div>
        </div>
        <TestnetFooter hide={false} height={this.state.footer}/>
      </div>
    );
  }
});
window.Dashboard = Dashboard // Expose the component

// BackendError is a UI component displaying an error message when an API ajax
// request fails for some reason. It is just a simple "danger" alert. Figuring
// out when to display it and when the error was resolved and it cn be hidden
// is left to the user.
//
// Properties:
//   hide: Flag whether to hide or show the error - bool
var BackendError = React.createClass({
  render: function() {
    if (this.props.hide) {
      return null
    }
    return (
      <div className="row">
        <div className="col-lg-3"></div>
        <div className="col-lg-6">
          <div className="alert alert-danger text-center" role="alert">
            <i className="fa fa-exclamation-triangle"></i> <strong>Backend offline!</strong> Please check console for details...
          </div>
        </div>
      </div>
    )
  }
});

// NotSyncedWarning is a UI component displaying a warning message when the local
// blockchain is older than a certain threshold.
//
// Properties:
//   hide: Flag whether to hide or show the error - bool
var NotSyncedWarning = React.createClass({
  render: function() {
    if (this.props.hide) {
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
