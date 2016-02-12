var ErrorAlert = React.createClass({
  render: function() {
    if (this.props.hide) {
      return null
    }
    return (
      <div className="row">
        <div className="col-lg-3"></div>
        <div className="col-lg-6">
          <div className="alert alert-danger text-center" role="alert">
            <i className="fa fa-exclamation-triangle"></i> Backend offline, please check console for details...
          </div>
        </div>
      </div>
    )
  }
});

var Dashboard = React.createClass({
  getInitialState: function() {
    return {
      index: false, market: true,
      apiOnline: true, apiFailure: ""};
  },
  loadIndex: function() {
    this.setState({index: false, market: true});
  },
  loadMarket: function() {
    this.setState({index: true, market: false});
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
  render: function() {
    return (
      <div>
        <nav className="navbar navbar-fixed-top">
          <div className="container">
            <div className="navbar-header">
              <a className="navbar-brand" href="#" onClick={this.loadIndex}>Ether APIs Dashboard</a>
            </div>
            <div id="navbar" className="navbar-collapse collapse">
              <ul className="nav navbar-nav">
                <li><a href="#" onClick={this.loadMarket}><i className="fa fa-shopping-basket"></i> API Market</a></li>
              </ul>
              <div className="navbar-right">
                <EthereumStats ajax={this.apiCall} apiurl="/api/v0/ethereum" refresh={1000}/>
              </div>
            </div>
          </div>
        </nav>
        <div className="container">
          <ErrorAlert hide={this.state.apiOnline}/>
          <Proxies hide={this.state.index}/>
          <Vault hide={this.state.index}/>
          <Market hide={this.state.market}/>
        </div>
      </div>
    );
  }
});
window.Dashboard = Dashboard
