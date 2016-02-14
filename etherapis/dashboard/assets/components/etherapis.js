// EtherAPI are a bunch of helper methods collected to provide access to
// EtherAPI restful subscriptions API.
var Services = React.createClass({
	getInitialState: function() {
		return { services: [] };
	},

	componentDidMount: function() {
		this.refreshServices();
		if(this.props.refresh > 0) {
			setInterval(this.refreshServices, this.props.refresh);
		}
	},

	refreshServices: function() {
		this.props.ajax(this.props.apiurl + "/services/" + this.props.address, function(services) {
			this.setState({services: services});
		}.bind(this));
	},

	render: function() {
		return (
			<div className="row">
				{this.state.services.map(function(service) {
					return (
						<Proxy name={service.name} service={service} online={true}/>
					);
				})}
			</div>
		);
	},
});
window.Services = Services;
