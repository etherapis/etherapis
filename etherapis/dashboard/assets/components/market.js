var Market = React.createClass({
  render: function() {
    if (this.props.hide) {
      return null
    }
    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <h3 className="panel-title">API Market</h3>
        </div>
        <div className="panel-body">
          Yeah, you wish... :P
        </div>
      </div>
    );
  }
});

window.Market = Market
