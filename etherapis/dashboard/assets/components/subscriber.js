// Subscriber is the content page that displays stats about the users current
// subscriptions and allows various operations on them.
var Subscriber = React.createClass({
	render: function() {
		// Short circuit rendering if we're not on the tutorial page
		if (this.props.hide) {
			return null
		}
		return (
			<div>
				<div className="row">
					<div className="col-lg-12">
						<h3>My subscriptions!</h3>
					</div>
				</div>
			</div>
		);
	}
});
window.Subscriber = Subscriber // Expose the component
