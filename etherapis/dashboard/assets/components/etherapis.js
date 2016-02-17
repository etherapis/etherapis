// EtherAPI are a bunch of helper methods collected to provide access to
// EtherAPI restful subscriptions API.
var Services = React.createClass({
	getInitialState: function() {
		return { services: [] };
	},

	componentDidMount: function() {
		this.refreshServices();

		if(this.props.refresh !== undefined) {
			setInterval(this.refreshServices, this.props.refresh);
		}
	},

	refreshServices: function() {
		var baseUrl = this.props.apiurl + "/services";
		if(this.props.address !== undefined) {
			baseUrl += "/" + this.props.address;
		}

		this.props.ajax(baseUrl, function(services) {
			this.setState({services: services});
		}.bind(this));
	},

	render: function() {
		return (
			<div>
				<div className="row">
					<div className="col-lg-12">
						<h3>Services</h3>
						{this.state.services.map(function(service, i) {
							return (
								<Service key={"service-"+i} name={service.name} service={service} online={true}/>
							);
						}.bind(this))}
					</div>
				</div>
			</div>
		);
	},
});
window.Services = Services;


var Service = React.createClass({
	render: function() {
		return (
			<div className="col-lg-4">
				<div className={this.props.service.enabled ? "panel panel-success" : "panel panel-default"}>
					<div className="panel-heading">
						<h3 className="panel-title">{this.props.name}: {this.props.service.enabled ? "Enabled" : "Disabled"}</h3>
					</div>
					<div className="panel-body" id="services">
						<table><tbody>
						<tr>
							<td>Endpoint:</td>
							<td><a href={this.props.service.endpoint}>{this.props.service.endpoint}</a></td>
						</tr>

						<tr>
							<td>Price:</td>
							<td>{this.props.service.price}</td>
						</tr>

						<tr>
							<td>Cancellation time:</td>
							<td>{this.props.service.cancellationTime}</td>
						</tr>
						</tbody></table>
					</div>
				</div>
			</div>
		);
	}
});
