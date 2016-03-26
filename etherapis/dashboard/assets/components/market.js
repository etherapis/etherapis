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
				<table className="table table-condensed">
					<thead>
						<tr>
							<th className="text-nowrap"><i className="fa fa-bookmark"></i> Name</th>
							<th className="text-nowrap"><i className="fa fa-star"></i> Rating</th>
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
									<td className="text-nowrap"><small>{service.name}</small></td>
									<td><RatingBar rating={6 * service.name.length}/></td>
									<td><small>{service.endpoint}</small></td>
									<td><Address address={service.owner} small/></td>
									<td className="text-nowrap text-center"><small>{formatBalance(service.price)}</small></td>
									<td className="text-nowrap text-center"><small>{moment.duration(service.cancellation, "seconds").humanize()}</small></td>
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

var RatingBar =  React.createClass({
	render: function() {
		var color = "progress-bar-danger";
		if (this.props.rating > 40) {
			color = "progress-bar-warning";
		}
		if (this.props.rating > 60) {
			color = "progress-bar-info";
		}
		if (this.props.rating > 85) {
			color = "progress-bar-success";
		}
		return (
			<div className="progress" style={{marginTop: "6px", marginBottom: 0, position: "relative", height: "10px"}}>
				<div className={"progress-bar " + color} style={{width: this.props.rating + "%"}}>
					<span style={{position: "absolute", display: "block", width: "100%", marginTop: "-5px"}}><small>{this.props.rating + "%"}</small></span>
				</div>
			</div>
		)
	}
})

var Subscribe = React.createClass({
	render: function() {
		return (
			<button type="button" className="btn btn-default btn-xs right">
				<i className="fa fa-plus"></i> Subscribe
			</button>
		);
	}
});
