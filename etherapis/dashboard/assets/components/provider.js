// Provider is the content page that displays stats about the users current
// provided APIs, proxies, configurations, etc.
var Provider = React.createClass({
	render: function() {
		// Short circuit rendering if we're not on the tutorial page
		if (this.props.hide) {
			return null
		}
		// Gather the provided services, sorted by provider account and paired up
		var addresses = Object.keys(this.props.services).sort();
		var services  = [];
		for (var i = 0; i < addresses.length; i++) {
			Array.prototype.push.apply(services, this.props.services[addresses[i]]);
		}
		var pairs = [];
		for (var i=0; i<services.length; i+=2) {
			pairs.push({first: services[i], second: services[i+1]});
		}

		return (
			<div>
				{ pairs.length == 0 ? null :
					<div className="row">
						<div className="col-lg-12">
							<h3>My services</h3>
						</div>
					</div>
				}
				{
					pairs.map(function(pair) {
						return (
							<div key={pair.first.owner + pair.first.name} className="row">
								<div className="col-lg-6">
									<Service apiurl={this.props.apiurl} service={pair.first}/>
								</div>
								{ pair.second === undefined ? null :
									<div className="col-lg-6">
										<Service apiurl={this.props.apiurl} service={pair.second}/>
									</div>
								}
							</div>
						)
					}.bind(this))
				}
				{ addresses.length == 0 ? null :
					<div className="row">
						<div className="col-lg-12">
							<h3>Register new service</h3>
							<form className="form-horizontal well" style={{margintBottom: 0}}>
									<div className="col-lg-6">
								    <div className="form-group">
								      <label className="col-lg-2 control-label">Provider</label>
								      <div className="col-lg-10">
												<button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
											    <Address address={this.props.active}/> <span className="caret"></span>
											  </button>
											  <ul className="dropdown-menu"> {
													addresses.map(function(address) {
														return (<li key={address}><a href="#"><Address address={address}/></a></li>);
													})
											  } </ul>
								      </div>
								    </div>
										<div className="form-group">
											<label className="col-lg-2 control-label">Public</label>
											<div className="col-lg-10">
												<div className="checkbox">
													<label><input type="checkbox" defaultChecked/> Advertise market availability</label>
												</div>
											</div>
										</div>
								    <div className="form-group">
								      <label className="col-lg-2 control-label">Name</label>
								      <div className="col-lg-10">
								        <input type="text" className="form-control" placeholder="Public name"/>
								      </div>
										</div>
								    <div className="form-group">
											<label className="col-lg-2 control-label">Endpoint</label>
											<div className="col-lg-10">
												<input type="text" className="form-control" placeholder="Public endpoint"/>
											</div>
								    </div>
									</div>
									<div className="col-lg-6">
								    <div className="form-group">
								      <label className="col-lg-2 control-label">Payment</label>
								      <div className="col-lg-10">
								        <div className="radio">
								          <label>
								            <input type="radio" name="serviceType" defaultChecked/>
								            Per API invocation (calls)
								          </label>
								        </div>
								        <div className="radio">
								          <label>
								            <input type="radio" name="serviceType"/>
								            Per consumed data traffic (bytes)
								          </label>
								        </div>
												<div className="radio">
													<label>
														<input type="radio" name="serviceType"/>
														Per maintained connection time (seconds)
													</label>
												</div>
								      </div>
								    </div>
										<div className="form-group">
											<label className="col-lg-2 control-label">Price</label>
											<div className="col-lg-10">
												<div className="input-group pull-right">
										      <input type="text" className="form-control" value={1} onChange={this.updateAmount}/>
										      <div className="input-group-btn">
										        <button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">{EthereumUnits[7]} <span className="caret"></span></button>
										        <ul className="dropdown-menu dropdown-menu-right">
										          <li><a href="#" onClick={this.updateUnit}>{EthereumUnits[7]}</a></li>
										          <li><a href="#" onClick={this.updateUnit}>{EthereumUnits[6]}</a></li>
										          <li><a href="#" onClick={this.updateUnit}>{EthereumUnits[5]}</a></li>
															<li><a href="#" onClick={this.updateUnit}>{EthereumUnits[4]}</a></li>
															<li><a href="#" onClick={this.updateUnit}>{EthereumUnits[1]}</a></li>
										        </ul>
										      </div>
										    </div>
											</div>
										</div>
										<div className="form-group">
											<label className="col-lg-2 control-label">Lockin</label>
											<div className="col-lg-10">
												<div className="input-group pull-right">
													<input type="text" className="form-control" value={1} onChange={this.updateAmount}/>
													<div className="input-group-btn">
														<button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Seconds <span className="caret"></span></button>
														<ul className="dropdown-menu dropdown-menu-right">
															<li><a href="#" onClick={this.updateUnit}>Seconds</a></li>
															<li><a href="#" onClick={this.updateUnit}>Minutes</a></li>
															<li><a href="#" onClick={this.updateUnit}>Hours</a></li>
															<li><a href="#" onClick={this.updateUnit}>Days</a></li>
														</ul>
													</div>
												</div>
											</div>
										</div>
									</div>
									<div className="form-group" style={{marginBottom: 0}}>
										<div className="col-lg-2 col-lg-offset-5">
											<button type="submit" className="btn btn-default" style={{width: "100%"}}>Register service</button>
										</div>
									</div>
							</form>
						</div>
					</div>
				}
			</div>
		);
	}
});
window.Provider = Provider // Expose the component

var Service = React.createClass({
	render: function() {
		return (
			<div className={this.props.service.enabled ? "panel panel-success" : "panel panel-default"}>
				<div className="panel-heading">
					<h3 className="panel-title">{this.props.service.name}: {this.props.service.enabled ? "Enabled" : "Disabled"}</h3>
				</div>
				<div className="panel-body" id="services">
					<table className="table table-condensed">
						<tbody>
							<tr><td><i className="fa fa-user"></i> Owner</td><td><Address address={this.props.service.owner}/></td></tr>
							<tr><td><i className="fa fa-link"></i> Endpoint</td><td>{this.props.service.endpoint}</td></tr>
							<tr><td>&Xi; Price</td><td>{formatBalance(this.props.service.price)}</td></tr>
							<tr><td><i className="fa fa-ban"></i> Cancellation</td><td>{moment.duration(this.props.service.cancellationTime, "seconds").humanize()} ({this.props.service.cancellationTime} secs)</td></tr>
						</tbody>
					</table>
					<table className="table table-striped table-condensed">
						<thead>
							<tr><th>Subscriber</th><th>Funds</th><th>Owed</th><th></th></tr>
						</thead>
						<tbody>
							<tr>
								<td><Address address={this.props.service.owner} small/></td>
								<td><small>{formatBalance("10000000000000000000000")}</small></td>
								<td><small>{formatBalance("10000000000000000")}</small></td>
								<td><button type="button" className="btn btn-default btn-xs">Charge</button></td>
							</tr>
							<tr>
								<td><Address address={this.props.service.owner} small/></td>
								<td><small>{formatBalance("10000000000000000000000")}</small></td>
								<td><small>{formatBalance("0")}</small></td>
								<td></td>
							</tr>
							<tr>
								<td><Address address={this.props.service.owner} small/></td>
								<td><small>{formatBalance("10000000000000000000000")}</small></td>
								<td><small>{formatBalance("13")}</small></td>
								<td><button type="button" className="btn btn-default btn-xs">Charge</button></td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		);
	}
});
