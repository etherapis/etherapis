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
				<ServiceCreator apiurl={this.props.apiurl} addresses={addresses} active={this.props.active} loadaccs={this.props.loadaccs} switch={this.props.switch}/>
			</div>
		);
	}
});
window.Provider = Provider // Expose the component

var Service = React.createClass({
	// getInitialState sets the zero values of the component.
	getInitialState: function() {
		return {
			action: "",
		};
	},
	// abortAction restores the service UI into it's default no-action state.
	abortAction: function(event) {
		if (event != null) {
			event.preventDefault();
		}
		this.setState({action: ""});
	},
	// confirmUnlock displays service unlocking explanation message and the manual
	// confirmation buttons.
	confirmUnlock: function(event) {
		event.preventDefault();
		this.setState({action: "unlock"});
	},
	// confirmLock displays service locking explanation and warning messages and the
	// manual confirmation buttons.
	confirmLock: function(event) {
		event.preventDefault();
		this.setState({action: "lock"});
	},
	// confirmDelete displays the services deletion warning message and the manual
	// confirmation buttons.
	confirmDelete: function(event) {
		event.preventDefault();
		this.setState({action: "delete"});
	},
	// render flattens the service stats into a UI panel.
	render: function() {
		return (
			<div className={this.props.service.enabled ? "panel panel-default" : "panel panel-warning"}>
				<div className="panel-heading">
					<div className="pull-right"><i className={this.props.service.enabled ? "fa fa-unlock" : "fa fa-lock"}></i></div>
					<h3 className="panel-title">{this.props.service.name}&nbsp;</h3>
				</div>
				<div className="panel-body" id="services">
					<table className="table table-condensed">
						<tbody>
							<tr><td className="text-center"><i className="fa fa-user"></i></td><td>Owner</td><td style={{width: "100%"}}><Address address={this.props.service.owner}/></td></tr>
							<tr><td className="text-center"><i className="fa fa-link"></i></td><td>Endpoint</td><td>{this.props.service.endpoint}</td></tr>
							<tr><td className="text-center">&Xi;</td><td>Price</td><td>{formatBalance(this.props.service.price)}</td></tr>
							<tr><td className="text-center"><i className="fa fa-ban"></i></td><td>Cancellation</td><td>{moment.duration(this.props.service.cancellationTime, "seconds").humanize()} ({this.props.service.cancellationTime} secs)</td></tr>
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
					<div className="clearfix">
						<hr style={{margin: "10px 0"}}/>
						<div className="pull-right">
							{this.props.service.enabled ?
								<a href="#" className="btn btn-sm btn-warning" onClick={this.confirmLock}><i className="fa fa-lock"></i> Disable</a> :
								<a href="#" className="btn btn-sm btn-success" onClick={this.confirmUnlock}><i className="fa fa-unlock"></i> Enable</a>
							}
							&nbsp;
							<a href="#" className="btn btn-sm btn-danger" onClick={this.confirmDelete}><i className="fa fa-times"></i> Delete</a>
						</div>
					</div>
					<UnlockConfirm apiurl={this.props.apiurl} service={this.props.service} hide={this.state.action != "unlock"} abort={this.abortAction}/>
					<LockConfirm apiurl={this.props.apiurl} service={this.props.service} hide={this.state.action != "lock"} abort={this.abortAction}/>
					<DeleteConfirm apiurl={this.props.apiurl} service={this.props.service} hide={this.state.action != "delete"} abort={this.abortAction}/>
				</div>
			</div>
		);
	}
});

// UnlockConfirm displays a warning message and requires an additional confirmation
// from the user to ensure that no accidental service unlocking - and inherently
// subscription floodgate opening - happen.
var UnlockConfirm = React.createClass({
	// getInitialState sets the zero values of the component.
	getInitialState: function() {
		return {
			progress: false,
			failure:	null,
		};
	},
	// unlockService executes the actual service unlocking, sending back the service
	// identifier to the server for subscription allowance.
	unlockService: function(event) {
		event.preventDefault();

		// Show the spinner until something happens
		this.setState({progress: true});

		// Execute the service locking
		/*$.ajax({type: "DELETE", url: this.props.apiurl + "/" + this.props.service.name, cache: false,
			error: function(xhr, status, err) {
				this.setState({progress: false, failure: xhr.responseText});
			}.bind(this),
		});*/
	},
	// render flattens the account stats into a UI panel.
	render: function() {
		// Short circuit rendering if we're not confirming deletion
		if (this.props.hide) {
			return null
		}
		return (
			<div>
				<hr/>
				<p>
					This service is currently <strong>disabled</strong>: it is hidden from the marketplace and new users
					cannot subscribe to it. By enabling it you permit global subscriptions. Please ensure your API is
					online to prevent advertising dysfunctional endpoints.
				</p>
				<div style={{textAlign: "center"}}>
					<p><strong>Availability changes incur the consensus fees of the Ethereum network.</strong></p>
					<a href="#" className={"btn btn-success " + (this.state.progress ? "disabled" : "")} onClick={this.unlockService}>
						{ this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} Enable new subscriptions
					</a>
					&nbsp;&nbsp;&nbsp;
					<a href="#" className="btn btn-info" onClick={this.props.abort}>Keep subscriptions closed</a>
				</div>
				{ this.state.failure ? <div style={{textAlign: "center"}}><hr/><p className="text-danger">Failed to delete account: {this.state.failure}</p></div> : null }
			</div>
		)
	}
})

// LockConfirm displays a warning message and requires an additional confirmation
// from the user to ensure that no accidental service locking - and inherently
// subscription prevention - happen.
var LockConfirm = React.createClass({
	// getInitialState sets the zero values of the component.
	getInitialState: function() {
		return {
			progress: false,
			failure:	null,
		};
	},
	// lockService executes the actual service locking, sending back the service
	// identifier to the server for subscription prevention.
	lockService: function(event) {
		event.preventDefault();

		// Show the spinner until something happens
		this.setState({progress: true});

		// Execute the service locking
		/*$.ajax({type: "DELETE", url: this.props.apiurl + "/" + this.props.service.name, cache: false,
			error: function(xhr, status, err) {
				this.setState({progress: false, failure: xhr.responseText});
			}.bind(this),
		});*/
	},
	// render flattens the account stats into a UI panel.
	render: function() {
		// Short circuit rendering if we're not confirming deletion
		if (this.props.hide) {
			return null
		}
		return (
			<div>
				<hr/>
				<p>
					This service is currently <strong>enabled</strong> and can accept new API subscriptions from users
					around the globe. By disabling it you can hide the service from the marketplace and prevent new
					users from subscribing. Existing ones remain valid and operational!
				</p>
				<div style={{textAlign: "center"}}>
					<p><strong>Availability changes incur the consensus fees of the Ethereum network.</strong></p>
					<a href="#" className={"btn btn-warning " + (this.state.progress ? "disabled" : "")} onClick={this.lockService}>
						{ this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} Disable new subscriptions
					</a>
					&nbsp;&nbsp;&nbsp;
					<a href="#" className="btn btn-info" onClick={this.props.abort}>Keep subscriptions open</a>
				</div>
				{ this.state.failure ? <div style={{textAlign: "center"}}><hr/><p className="text-danger">Failed to delete account: {this.state.failure}</p></div> : null }
			</div>
		)
	}
})

// DeleteConfirm displays a warning message and requires an additional confirmation
// from the user to ensure that no accidental service deletion happens.
var DeleteConfirm = React.createClass({
	// getInitialState sets the zero values of the component.
	getInitialState: function() {
		return {
			progress: false,
			failure:	null,
		};
	},
	// deleteService executes the actual service deletion, sending back the account
	// identifier to the server for irreversible removal.
	deleteService: function(event) {
		event.preventDefault();

		// Show the spinner until something happens
		this.setState({progress: true});

		// Execute the account deletion request
		/*$.ajax({type: "DELETE", url: this.props.apiurl + "/" + this.props.service.name, cache: false,
			error: function(xhr, status, err) {
				this.setState({progress: false, failure: xhr.responseText});
			}.bind(this),
		});*/
	},
	// render flattens the account stats into a UI panel.
	render: function() {
		// Short circuit rendering if we're not confirming deletion
		if (this.props.hide) {
			return null
		}
		return (
			<div>
				<hr/>
				<p>
					<strong>Warning!</strong> Deleting a service is permanent and irreversible. It will stop all associated
					API proxies (local and remote) and take this service off the marketplace. An attempt is also made to
					charge all pending subscription payments.
				</p>
				<div style={{textAlign: "center"}}>
					<p><strong>Deletion incurs the consensus fees of the Ethereum network.</strong></p>
					<a href="#" className={"btn btn-danger " + (this.state.progress ? "disabled" : "")} onClick={this.deleteService}>
						{ this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} <strong>Irreversibly</strong> delete service
					</a>
					&nbsp;&nbsp;&nbsp;
					<a href="#" className="btn btn-success" onClick={this.props.abort}>Keep service available</a>
				</div>
				{ this.state.failure ? <div style={{textAlign: "center"}}><hr/><p className="text-danger">Failed to delete service: {this.state.failure}</p></div> : null }
			</div>
		)
	}
})

// ServiceCreator is a UI component for generating a brand new pristine service.
var ServiceCreator = React.createClass({
	// getInitialState sets the zero values of the component.
	getInitialState: function() {
		return {
			public:   true,
			name:     "",
			endpoint: "",
			payment:  "call",
			price:    "",
			denom:    EthereumUnits[4],
			cancel:   "",
			scale:    "Seconds",
			progress: false,
			failure:  null,
		};
	},
	// loadAccounts navigates to the accounts page.
	loadAccounts: function(event) {
		event.preventDefault();
		this.props.loadaccs();
	},
	// the method set below pulls in the users modifications from the input boxes
	// and updates the UIs internal state with it.
	updatePublic:   function(event) {
		if (event.target.checked) {
			this.setState({public: true});
		} else {
			this.setState({public: false, name: "", endpoint: ""});
		}
	},
	updateName:     function(event) { this.setState({name: event.target.value}); },
	updateEndpoint: function(event) { this.setState({endpoint: event.target.value}); },
	updatePrice:    function(event) { this.setState({price: event.target.value}); },
	updateCancel:   function(event) { this.setState({cancel: event.target.value}); },

	updateDenom: function(event) {
		event.preventDefault();
		this.setState({denom: event.target.textContent});
	},
	updateScale: function(event) {
		event.preventDefault();
		this.setState({scale: event.target.textContent});
	},
	// registerService executes the actual service registration.
	registerService: function(event) {
		event.preventDefault();

		// Show the spinner until something happens
		this.setState({progress: true});

		// Assemble and send the value transfer request
		var form = new FormData();
		form.append("name", this.state.name);
		form.append("url", this.state.endpoint);
		form.append("price", weiAmount(this.state.price, this.state.denom));
		form.append("cancel", secondsDuration(this.state.cancel, this.state.scale));

		$.ajax({type: "POST", url: this.props.apiurl + "/" + this.props.active, cache: false, data: form, processData: false, contentType: false,
			success: function(data) {
				this.setState(this.getInitialState());
			}.bind(this),
			error: function(xhr, status, err) {
				this.setState({progress: false, failure: xhr.responseText});
			}.bind(this),
		});
	},
	// render flattens the account stats into a UI panel.
	render: function() {
		// Short circuit if no accounts are known
		if (this.props.addresses.length == 0) {
			return (
				<div className="row">
					<div className="col-lg-12">
						<h3>No available accounts</h3>
						<p>
							Ether APIs was unable to locate any Ethereum accounts with which to provide new APIs, or for which to manage
							exsiting ones. In order to create new services or manage your	already registered ones, please switch to the
							&nbsp;<a href="#"><button className="btn btn-xs btn-default" onClick={this.loadAccounts}>
								<i className="fa fa-user"></i> Account
							</button></a>&nbsp;
							tab and either generate a new Ethereum account or import an already existing one.
						</p>
					</div>
				</div>
			);
		}
		// Create an account switcher
		var switcher = function(address) {
			return function(event) {
				event.preventDefault();
				this.props.switch(address);
			}.bind(this)
		}.bind(this)
		// Otherwise display the service creation form
		return (
			<div className="row">
				<div className="col-lg-12">
					<h3>Register service</h3>
					<p>
						Registering a service will create a new offering in the Ether APIs decentralized marketplace. Once created, the
						details of the service and the terms of use are set in stone (or rather the blockchain) and cannot ever be changed;
						only a new one created in its place. Users may also create private, non-advertised APIs where only the terms of
						payment are specified (to be enforcable via Ether APIs), leaving it to the owner to provide the accessability
						information to select customers.
					</p>
					<form className="form-horizontal well">
						<div className="col-lg-6">
							<div className="form-group">
								<label className="col-lg-2 control-label">Provider</label>
								<div className="col-lg-10">
									<button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" style={{width: "100%"}}>
										<Address address={this.props.active}/>&nbsp;&nbsp;&nbsp;<span className="caret"></span>
									</button>
									<ul className="dropdown-menu"> {
										this.props.addresses.map(function(address) {
											return (<li key={address}><a href="#" onClick={switcher(address)}><Address address={address}/></a></li>);
										}.bind(this))
									} </ul>
								</div>
							</div>
							<div className="form-group">
								<label className="col-lg-2 control-label">Public</label>
								<div className="col-lg-10">
									<div className="checkbox">
										<label><input type="checkbox" defaultChecked={this.state.public} onChange={this.updatePublic}/> Advertise marketplace availability</label>
									</div>
								</div>
							</div>
							<div className="form-group">
								<label className="col-lg-2 control-label">Name</label>
								<div className="col-lg-10">
									<input type="text" className="form-control" disabled={!this.state.public} placeholder={this.state.public ? "Public name" : "Private services cannot advertise name"} value={this.state.name} onChange={this.updateName}/>
								</div>
							</div>
							<div className="form-group">
								<label className="col-lg-2 control-label">Endpoint</label>
								<div className="col-lg-10">
									<input type="text" className="form-control" disabled={!this.state.public} placeholder={this.state.public ? "Public endpoint" : "Private services cannot advertise endpoint"} value={this.state.endpoint} onChange={this.updateEndpoint}/>
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
										<input type="text" className="form-control" placeholder="Unit price" value={this.state.price} onChange={this.updatePrice}/>
										<div className="input-group-btn">
											<button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">{this.state.denom} <span className="caret"></span></button>
											<ul className="dropdown-menu dropdown-menu-right">
												<li><a href="#" onClick={this.updateDenom}>{EthereumUnits[7]}</a></li>
												<li><a href="#" onClick={this.updateDenom}>{EthereumUnits[6]}</a></li>
												<li><a href="#" onClick={this.updateDenom}>{EthereumUnits[5]}</a></li>
												<li><a href="#" onClick={this.updateDenom}>{EthereumUnits[4]}</a></li>
												<li><a href="#" onClick={this.updateDenom}>{EthereumUnits[1]}</a></li>
											</ul>
										</div>
									</div>
								</div>
							</div>
							<div className="form-group">
								<label className="col-lg-2 control-label">Lockin</label>
								<div className="col-lg-10">
									<div className="input-group pull-right">
										<input type="text" className="form-control" placeholder="Cancellation time" value={this.state.time} onChange={this.updateCancel}/>
										<div className="input-group-btn">
											<button type="button" className="btn btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">{this.state.scale} <span className="caret"></span></button>
											<ul className="dropdown-menu dropdown-menu-right">
												<li><a href="#" onClick={this.updateScale}>Seconds</a></li>
												<li><a href="#" onClick={this.updateScale}>Minutes</a></li>
												<li><a href="#" onClick={this.updateScale}>Hours</a></li>
												<li><a href="#" onClick={this.updateScale}>Days</a></li>
											</ul>
										</div>
									</div>
								</div>
							</div>
						</div>
						<div className="form-group" style={{marginBottom: 0, textAlign: "center"}}>
							<p><strong>Registering a service incurs the consensus fees of the Ethereum network.</strong></p>
							<div className="col-lg-2 col-lg-offset-5">
								<a href="#" className={"btn btn-default " + ((this.state.public && (this.state.name == "" || this.state.endpoint == "")) || this.state.price == "" || this.state.cancel == "" ? "disabled" : "")} style={{width: "100%"}}  onClick={this.registerService}>
									{ this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} Register service
								</a>
							</div>
						</div>
						{ this.state.failure ? <div style={{textAlign: "center"}}><p className="text-danger">Failed to register service: {this.state.failure}</p></div> : null }
					</form>
				</div>
			</div>
		);
	}
})
