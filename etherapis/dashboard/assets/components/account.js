// Account is the content page that displays stats about the users active Ethereum
// account, with a possibility to import/export/create and other meta functionality.
var Account = React.createClass({
  render: function() {
    // Short circuit rendering if we're not on the tutorial page
    if (this.props.hide) {
      return null
    }
    return (
      <div>
        <div className="row">
          <div className="col-lg-12">
            <h3>My accounts!</h3>
          </div>
        </div>
      </div>
    );
  }
});
window.Account = Account // Expose the component
