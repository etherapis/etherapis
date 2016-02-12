// Dashboard is the entire webpage content of the EtherAPIs dahsboard. It should
// more or less be fully self contained, with perhaps a few external properties.
//
// Properties:
//   apiurl: Base path for the API access - string
var Dashboard = React.createClass({
  getInitialState: function() {
    return {
      section: "index",
      apiOnline: true, apiFailure: ""};
  },

  loadIndex:      function() { this.setState({section: "index"}); },
  loadProvider:   function() { this.setState({section: "provider"}); },
  loadSubscriber: function() { this.setState({section: "subscriber"}); },
  loadMarket:     function() { this.setState({section: "market"}); },

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
                <li><a href="#" onClick={this.loadProvider}><i className="fa fa-cloud-upload"></i> Provided</a></li>
                <li><a href="#" onClick={this.loadSubscriber}><i className="fa fa-cloud-download"></i> Subscribed</a></li>
                <li><a href="#" onClick={this.loadMarket}><i className="fa fa-shopping-basket"></i> Market</a></li>
              </ul>
              <div className="navbar-right">
                <EthereumStats ajax={this.apiCall} apiurl={this.props.apiurl + "/ethereum"} refresh={1000}/>
              </div>
            </div>
          </div>
        </nav>
        <div className="container">
          <BackendError hide={this.state.apiOnline}/>

          <Tutorial hide={this.state.section != "index"}/>
          <Provider hide={this.state.section != "provider"}/>
          <Subscriber hide={this.state.section != "subscriber"}/>
          <Market hide={this.state.section != "market"}/>
        </div>
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
