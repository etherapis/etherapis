var Proxy = React.createClass({
  render: function() {
    return (
      <div className="col-lg-4">
        <div className={this.props.online ? "panel panel-success" : "panel panel-default"}>
          <div className="panel-heading">
            <h3 className="panel-title">{this.props.name}: {this.props.online ? "Online" : "Offline"}</h3>
          </div>
          <div className="panel-body">
            Proxy stats go here...
          </div>
        </div>
      </div>
    )
  }
});

var Proxies = React.createClass({
  render: function() {
    if (this.props.hide) {
      return null
    }
    return (
      <div>
        <div className="row">
          <div className="col-lg-12">
            <h3>Payment proxies</h3>
          </div>
        </div>
        <div className="row">
          <Proxy name={"Geolocator demo"} online={true}/>
          <Proxy name={"Streamer demo"} online={false}/>
        </div>
      </div>
    )
  }
});
window.Proxies = Proxies

var Vault = React.createClass({
  render: function() {
    if (this.props.hide) {
      return null
    }
    return (
      <div>
        <div className="row">
          <div className="col-lg-12">
            <h3>Payment vault</h3>
            <div className="panel panel-default">
              <div className="panel-body">
                <table className="table table-striped table-hover ">
                  <thead>
                    <tr>
                      <th>Service</th>
                      <th>Subscriber</th>
                      <th>Owed</th>
                      <th>Limit</th>
                      <th>Total paid</th>
                      <th>Last request</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>Geolocator demo</td>
                      <td>0x3903a55e0802f011077bb33382822937d94efb7f</td>
                      <td>314 Finney</td>
                      <td>500 Finney</td>
                      <td>1.1 Ether</td>
                      <td>5 minutes ago</td>
                    </tr>
                    <tr>
                      <td>Geolocator demo</td>
                      <td>0x9fc6fefd7f33ca29ee17f2bfec944695e5f29caf</td>
                      <td>1.1 Ether</td>
                      <td>10 Ether</td>
                      <td>101 Ether</td>
                      <td>a few seconds ago</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
});
window.Vault = Vault
