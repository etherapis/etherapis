var Market = React.createClass({
	getInitialState: function() {
		return { services: [] };
	},

	componentDidMount: function() {
		this.refreshMarket();
		if(this.props.interval !== undefined) {
			setInterval(this.refreshMarket, this.props.interval);
		}
	},
	
	refreshMarket: function() {
		this.props.ajax(this.props.apiurl + "/services", function(services) {
			this.setState({services: services});
		}.bind(this));
	},

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
					The API store in its full glory and awesomeness...
				</div>

				<table className="table">
					<thead>
						<tr>
							<th>Name</th>
							<th>Endpoint</th>
							<th>Owner</th>
							<th>Price</th>
							<th>Cancellation time</th>
							<th></th>
						</tr>
					</thead>

					<tbody>
						{this.state.services.map(function(service, i) {
							return (
								<MarketItem key={"market-item-"+i} name={service.name} ajax={this.props.ajax} service={service}/>
							);
						}.bind(this))}
					</tbody>
				</table>
			</div>
		);
	}
});
window.Market = Market;

var MarketItem = React.createClass({
	render: function() {
		return (
			<tr>
				<td>{this.props.service.name}</td>
				<td><a href={this.props.service.endpoint}>{this.props.service.endpoint}</a></td>
				<td>{this.props.service.owner}</td>
				<td>{this.props.service.price}</td>
				<td>{this.props.service.cancellationTime}</td>
				<td><Subscribe/></td>
			</tr>
		);
	},
});

var Subscribe = React.createClass({
	render: function() {
		return (
			<button type="button" className="btn btn-default btn-xs right">
				<span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Subscribe
			</button>
		);
	}
});
