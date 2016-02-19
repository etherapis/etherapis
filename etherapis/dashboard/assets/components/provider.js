// Provider is the content page that displays stats about the users current
// provided APIs, proxies, configurations, etc.
var Provider = React.createClass({
	render: function() {
		// Short circuit rendering if we're not on the tutorial page
		if (this.props.hide) {
			return null
		}
		return (
			<div>
				<div className="row">
					<div className="col-lg-12">
						<Proxies/>
			      <Services ajax={this.props.ajax} apiurl={this.props.apiurl} address={this.props.address} refresh={1000}/>
						<Vault/>
					</div>
				</div>
			</div>
		);
	}
});
window.Provider = Provider // Expose the component

var Proxy = React.createClass({
	render: function() {
		return (
			<div className="col-lg-4">
				<div className={this.props.online ? "panel panel-success" : "panel panel-default"}>
					<div className="panel-heading">
						<h3 className="panel-title">{this.props.name}: {this.props.online ? "Online" : "Offline"}</h3>
					</div>
					<div className="panel-body" id="proxies">
			      Body here...
					</div>
				</div>
			</div>
		)
	}
});

var Proxies = React.createClass({
	render: function() {
		return (
    <div className="row">
      <div className="col-lg-12">
        <h3>Payment proxies</h3>
        <Proxy name="Proxy 1" online={true}/>
        <Proxy name="Proxy 2" online={false}/>
      </div>
    </div>
		)
	}
});

var Vault = React.createClass({
	render: function() {
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
