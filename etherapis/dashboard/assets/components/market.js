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
				<table className="table">
					<thead>
						<tr>
							<th className="text-nowrap"><i className="fa fa-bookmark"></i> Name</th>
							<th className="text-nowrap"><i className="fa fa-link"></i> Endpoint</th>
							<th className="text-nowrap"><i className="fa fa-user"></i> Owner</th>
							<th className="text-nowrap">&Xi; Price</th>
							<th className="text-nowrap"><i className="fa fa-ban"></i> Cancellation</th>
							<th className="text-nowrap"></th>
						</tr>
					</thead>

					<tbody>
						{this.props.market.map(function(service, i) {
							return (
								<tr key={"market-item-"+i} >
									<td><small>{service.name}</small></td>
									<td><small>{service.endpoint}</small></td>
									<td><Address address={service.owner} small/></td>
									<td className="text-nowrap text-center"><small>{formatBalance(service.price)}</small></td>
									<td className="text-nowrap text-center"><small>{moment.duration(service.cancellationTime, "seconds").humanize()}</small></td>
									<td><Subscribe/></td>
								</tr>
							);
						}.bind(this))}
					</tbody>
				</table>
			</div>
		);
	}
});
window.Market = Market;

var Subscribe = React.createClass({
	render: function() {
		return (
			<button type="button" className="btn btn-default btn-xs right">
				<span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Subscribe
			</button>
		);
	}
});
