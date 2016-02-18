// Accounts is the content page that displays stats about the users active Ethereum
// account, with a possibility to import/export/create and other meta functionality.
var Accounts = React.createClass({
  render: function() {
    // Short circuit rendering if we're not on the accounts page
    if (this.props.hide) {
      return null
    }
    return (
      <div>
        {
          this.props.accounts.length == 0 ? null :
          <div className="row">
            <div className="col-lg-12">
              <h3>My accounts</h3>
            </div>
            {
              this.props.accounts.map(function(account) {
                return (
                  <div key={account} className="col-lg-6">
                    <Account apiurl={this.props.apiurl} explorer={this.props.explorer} account={account} active={this.props.active} switch={this.props.accounts.length > 1 ? this.props.switch : null} refresh={this.props.refresh}/>
                  </div>
                )
              }.bind(this))
            }
          </div>
        }
        <div className="row">
          <div className="col-lg-12">
            <h3>Add account</h3>
          </div>
          <div className="col-lg-6"><AccountCreator apiurl={this.props.apiurl} refresh={this.props.refresh}/></div>
          <div className="col-lg-6"><AccountImporter apiurl={this.props.apiurl} refresh={this.props.refresh}/></div>
        </div>
      </div>
    );
  }
});
window.Accounts = Accounts // Expose the component

// Account is a UI component displaying the summary infos of a single Ethereum
// account.
var Account = React.createClass({
  // getInitialState sets the zero values of the component.
  getInitialState: function() {
    return {
      action: "",
    };
  },
  // activate switches the dashboard to use this particular account.
  activate: function() { this.props.switch(this.props.account); },

  // abortAction restores the account UI into it's default no-action state.
  abortAction: function() { this.setState({action: ""}); },

  // confirmExport displays the account export warning message, the password input
  // field to encrypt the key with and the manual confirmation buttons.
  confirmExport: function() { this.setState({action: "export"}); },

  // confirmDelete displays the account deletion warning messages and the manual
  // confirmation buttons.
  confirmDelete: function() { this.setState({action: "delete"}); },

  // render flattens the account stats into a UI panel.
  render: function() {
    return (
      <div className={this.props.account == this.props.active && this.props.switch != null ? "panel panel-success" : "panel panel-default"}>
        <div className="panel-heading">
          <img style={{borderRadius: "50%", marginRight: "8px"}} src={blockies.create({seed: this.props.account, size: 8, scale: 2}).toDataURL()}/>
          <span style={{fontFamily: "monospace"}}>{this.props.account}</span>{this.props.account == this.props.active && this.props.switch != null ? " – Active" : null}
          <a href={this.props.explorer + this.props.account} target="_blank" className="pull-right"><i className="fa fa-external-link"></i></a>
        </div>
        <div className="panel-body">
          <div>
            Provided services:
          </div>
          <div>
            Subscribed services:
          </div>
          <div className="clearfix">
            <hr style={{margin: "10px 0"}}/>
            { this.props.account == this.props.active ? null : <a href="#" className="btn btn-sm btn-success" onClick={this.activate}><i className="fa fa-check-circle-o"></i> Activate</a>}
            <div className="pull-right">
              <a href="#" className="btn btn-sm btn-warning" onClick={this.confirmExport}><i className="fa fa-arrow-circle-o-down"></i> Export</a>
              &nbsp;
              <a href="#" className="btn btn-sm btn-danger" onClick={this.confirmDelete}><i className="fa fa-user-times"></i> Delete</a>
            </div>
          </div>
          <ExportConfirm apiurl={this.props.apiurl} account={this.props.account} hide={this.state.action != "export"} abort={this.abortAction}/>
          <DeleteConfirm apiurl={this.props.apiurl} account={this.props.account} hide={this.state.action != "delete"} abort={this.abortAction} active={this.props.active} refresh={this.props.refresh}/>
        </div>
      </div>
    );
  }
});

// ExportConfirm displays a warning message and requires an addtional confirmation
// from the user to ensure that no accidental account deletion happens.
var ExportConfirm = React.createClass({
  // getInitialState sets the zero values of the component.
  getInitialState: function() {
    return {
      input:    "",
      confirm:  "",
      progress: false,
    };
  },
  // updateInput and updateConfirm pulls in the users modifications from the input
  // boxes and updates the UIs internal state with it.
  updateInput:   function(event) { this.setState({input: event.target.value}); },
  updateConfirm: function(event) { this.setState({confirm: event.target.value}); },

  // exportAccount executes the actual account export, sending back the account
  // identifier along with the password to encrypt it with.
  exportAccount: function() {
    this.setState({progress: true});

    // We have no idea how much time it takes, display for 2 secs, ten hide :P
    setTimeout(function() {
      this.setState(this.getInitialState());
      this.props.abort();
    }.bind(this), 2000);
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
          Please note, an exported account is an <strong>extremely sensitive</strong> piece of information.
          Although it will be encrypted with the passphrase set below, it contains access to all funds
          stored within the account. Do not share it, do not lose it. Be cautions!
        </p>
        <div className="form-group">
          <input type="password" className="form-control pull-right" style={{width: "49%"}} placeholder="Confirm passphrase" onChange={this.updateConfirm}/>
          <input type="password" className="form-control" style={{width: "49%"}} placeholder="Passphrase" onChange={this.updateInput}/>
        </div>
        <div style={{textAlign: "center"}}>
          <p><strong>Do not forget this password, there is no way to recover it!</strong></p>
          <a href={this.props.apiurl + "/" + this.props.account + "/" + this.state.input} className={"btn btn-warning " + (this.state.input == "" || this.state.input != this.state.confirm || this.state.progress ? "disabled" : "")} onClick={this.exportAccount}>
            { this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} Export this account
          </a>
          &nbsp;&nbsp;&nbsp;
          <a href="#" className="btn btn-info" onClick={this.props.abort}>Do not export account</a>
        </div>
      </div>
    )
  }
})

// DeleteConfirm displays a warning message and requires an addtional confirmation
// from the user to ensure that no accidental account deletion happens.
var DeleteConfirm = React.createClass({
  // getInitialState sets the zero values of the component.
  getInitialState: function() {
    return {
      progress: false,
      failure:  null,
    };
  },
  // deleteAccount executes the actual account deletion, sending back the account
  // identifier to the server for irreversible removal.
  deleteAccount: function() {
    // Show the spinner until something happens
    this.setState({progress: true});

    // Execute the account deletion request
    $.ajax({type: "DELETE", url: this.props.apiurl + "/" + this.props.account, cache: false,
      success: function() {
        this.setState({progress: false, failure: null});
        this.props.refresh(this.props.active);
      }.bind(this),
      error: function(xhr, status, err) {
        this.setState({progress: false, failure: xhr.responseText});
      }.bind(this),
    });
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
          <strong>Warning!</strong> Deleting an account is a permanent and irreversible action.
          There are no automatic backups, there are no failsafes, there are no restore facilities.
          Removing a non exported account will forever forfeit access to it and any funds contained within.
        </p>
        <div style={{textAlign: "center"}}>
          <a href="#" className={"btn btn-danger " + (this.state.progress ? "disabled" : "")} onClick={this.deleteAccount}>
            { this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} <strong>Irreversibly</strong> delete account!
          </a>
          &nbsp;&nbsp;&nbsp;
          <a href="#" className="btn btn-success" onClick={this.props.abort}>I have changed my mind!</a>
        </div>
        { this.state.failure ? <div style={{textAlign: "center"}}><hr/><p className="text-danger">Failed to delete account: {this.state.failure}</p></div> : null }
      </div>
    )
  }
})

// AccountCreator is a UI component for generating a brand new pristing account.
var AccountCreator = React.createClass({
  // getInitialState sets the zero values of the component.
  getInitialState: function() {
    return {
      progress: false,
      failure:  null,
    };
  },
  // createAccount executes the actual account creation, sending an account
  // generation request to the backend server.
  createAccount: function() {
    // Show the spinner until something happens
    this.setState({progress: true});

    // Do a simple account creation post request
    var form = new FormData();
    form.append("action", "create");

    $.ajax({type: "POST", url: this.props.apiurl, cache: false, data: form, processData: false, contentType: false,
      success: function(data) {
        this.setState({progress: false, failure: null});
        this.props.refresh(data);
      }.bind(this),
      error: function(xhr, status, err) {
        this.setState({progress: false, failure: xhr.responseText});
      }.bind(this),
    });
  },
  // render flattens the account stats into a UI panel.
  render: function() {
    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <i className="fa fa-user-secret"></i> Create new account
        </div>
        <div className="panel-body">
          <p>
            Creating an account generates a brand new, empty Ethereum account that can be used both for providing API services
            to others, as well as for subscribing to the APIs provided by others.
          </p>
          <p>
            Please note, that in order to interact with the Ethereum blockchain, the new account needs to hold at least a minimal
            amount of Ether. You can either transfer Ether from another account or obtain it via an exchange.
          </p>
          <div style={{textAlign: "center"}}>
            <a href="#" className={"btn btn-success " + (this.state.progress ? "disabled" : "")} onClick={this.createAccount}>
              { this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} Create account
            </a>
          </div>
          { this.state.failure ? <div style={{textAlign: "center"}}><hr/><p className="text-danger">Failed to create account: {this.state.failure}</p></div> : null }
        </div>
      </div>
    )
  }
})

// AccountImporter is a UI component for importing an already existing account.
var AccountImporter = React.createClass({
  // getInitialState sets the zero values of the component.
  getInitialState: function() {
    return {
      filename: "",
      fileblob: null,
      password: "",
      progress: false,
      failure:  null,
    };
  },
  // updateFile sets the file to be uploaded for importing.
  updateFile: function(event) {
    this.setState({
      filename: event.target.value,
      fileblob: event.target.files[0],
    });
  },
  // updatePassword pulls in the users modifications from the password box.
  updatePassword: function(event) { this.setState({password: event.target.value}); },

  // importAccount executes the actual account import, sending back the account
  // file and the password to decrypt it with.
  importAccount: function(event) {
    // Don't refresh the page, we don't want that
    event.preventDefault();

    // Show the spinner until something happens
    this.setState({progress: true});

    // Upload the form manually via AJAX queries
    var form = new FormData();
    form.append("action", "import");
    form.append("account", this.state.fileblob);
    form.append("password", this.state.password);

    $.ajax({type: "POST", url: this.props.apiurl, cache: false, data: form, processData: false, contentType: false,
      success: function(data) {
        this.setState(this.getInitialState());
        this.props.refresh(data);
      }.bind(this),
      error: function(xhr, status, err) {
        this.setState({progress: false, failure: xhr.responseText});
      }.bind(this),
    });
  },
  // render flattens the account stats into a UI panel.
  render: function() {
    return (
      <div className="panel panel-default">
        <div className="panel-heading">
          <i className="fa fa-user-plus"></i> Import existing account
        </div>
        <div className="panel-body">
          <p>Please select a previously exported account to import:</p>
          <form>
            <div className="form-group">
              <input type="password" className="form-control pull-right" style={{width: "49%"}} placeholder="Passphrase" value={this.state.password} onChange={this.updatePassword}/>
              <div className="input-group" style={{width: "49%"}}>
                <span className="input-group-btn">
                  <span className="btn btn-default btn-file">
                    Browse&hellip; <input type="file" onChange={this.updateFile}/>
                  </span>
                </span>
                <input type="text" className="form-control" value={this.state.filename} disabled/>
              </div>
            </div>
            <p>
              Importing will decrypt the uploaded account key with the provided password, and will reencrypt it using
              its own master password before saving it into its keystore.
            </p>
            <div style={{textAlign: "center"}}>
              <button type="submit" className="btn btn-success" disabled={this.state.filename == "" || this.state.progress} onClick={this.importAccount}>
                { this.state.progress ? <i className="fa fa-spinner fa-spin"></i> : null} Import account
              </button>
            </div>
            { this.state.failure ? <div style={{textAlign: "center"}}><hr/><p className="text-danger">Failed to import account: {this.state.failure}</p></div> : null }
          </form>
        </div>
      </div>
    )
  }
})
