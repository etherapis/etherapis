// Dashboard is the entire webpage content of the EtherAPIs dahsboard. It should
// more or less be fully self contained, with perhaps a few external properties.
//
// Properties:
//   apiurl: Base path for the API access - string
var Dashboard = React.createClass({
  getInitialState: function() {
    return {
      section: "index",
      apiOnline: true, apiFailure: "",
      needSync: false,
      account: "",
    };
  },

  loadIndex:      function() { this.setState({section: "index"}); },
  loadProvider:   function() { this.setState({section: "provider"}); },
  loadSubscriber: function() { this.setState({section: "subscriber"}); },
  loadMarket:     function() { this.setState({section: "market"}); },
  loadAccount:    function() { this.setState({section: "account"}); },

  // refreshAccount retrieves the current account configured
  refreshAccount: function() {
    this.props.ajax(this.props.apiurl + "/peers", function(data) { this.setState({peers: data}); }.bind(this));
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
      <div>
        <nav className="navbar navbar-fixed-top">
          <div className="container">
            <div className="navbar-header">
              <a className="navbar-brand" href="#" onClick={this.loadIndex}>Ether APIs</a>
            </div>
            <div className="navbar-collapse collapse">
              <ul className="nav navbar-nav">
                <li><a href="#" onClick={this.loadAccount}><i className="fa fa-user-secret"></i> Account</a></li>
                <li><a href="#" onClick={this.loadProvider}><i className="fa fa-cloud-upload"></i> Provided</a></li>
                <li><a href="#" onClick={this.loadSubscriber}><i className="fa fa-cloud-download"></i> Subscribed</a></li>
                <li><a href="#" onClick={this.loadMarket}><i className="fa fa-shopping-basket"></i> Market</a></li>
              </ul>
              <EthereumStats ajax={this.apiCall} apiurl={this.props.apiurl + "/ethereum"} refresh={1000} nosync={this.needSync}/>
            </div>
          </div>
        </nav>
        <div className="container">
          <BackendError hide={this.state.apiOnline}/>
          <NotSyncedWarning hide={!this.state.needSync || !this.state.apiOnline}/>

          <Tutorial hide={this.state.section != "index"}/>
          <Account hide={this.state.section != "account"}/>
          <Provider hide={this.state.section != "provider"} ajax={this.apiCall} apiurl={this.props.apiurl} address={this.state.account} />
          <Subscriber hide={this.state.section != "subscriber"}/>
          <Market hide={this.state.section != "market"}/>
        </div>
        <TestnetFooter hide={false}/>
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
      <div style={{position: "absolute", bottom: 0, width: "100%", height: "52px"}}>
          <div className="alert alert-info   text-center" role="alert" style={{marginBottom: 0}}>
            <i className="fa fa-info-circle"></i> This instance is currently configured to use the <strong>test network</strong>. Consider everything here a consequence free playground.
          </div>
      </div>
    )
  }
});
