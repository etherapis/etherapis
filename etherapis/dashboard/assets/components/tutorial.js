// Tutorial is the content page that explains the various components of the system.
// It should be the default page that id displayed when the dahsboard loads up.
var Tutorial = React.createClass({
  render: function() {
    // Short circuit rendering if we're not on the tutorial page
    if (this.props.hide) {
      return null
    }
    return (
      <div>
        <div className="row">
          <div className="col-lg-12">
            <h3>Welcome to EtherAPIs!</h3>
            This section it empty fow now (wow, no shit, right? :D). But it is the placeholder
            for the dashboard tutorial which should present the UI components and their use.
            It should only be written after the remainder of the dashboard reaches at least a
            tiny bit of stability to allow adding screenshots and the likes.
          </div>
        </div>
      </div>
    );
  }
});
window.Tutorial = Tutorial // Expose the component
