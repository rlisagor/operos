/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React from 'react';

import Icon from 'material-ui/Icon';
import Card, {CardContent, CardActions} from 'material-ui/Card';
import Typography from 'material-ui/Typography';
import {withStyles} from 'material-ui/styles';
import classNames from 'classnames';
import {green} from 'material-ui/colors';
import Button from 'material-ui/Button';
import Dialog, {
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from 'material-ui/Dialog';


const styles = theme => ({
  upgradeStatus: {
    display: 'flex',
    margin: '0.5em 0'
  },
  upgradeIcon: {
    marginRight: 6
  },
  uptoDate: {
    color: green[500]
  },
  notUptoDate: {
    color: theme.palette.primary[500]
  }
});


class UpgradeCard extends React.Component {
  constructor() {
    super();
    this.state = {
      confirmation: false
    };
  }

  handleConfirm() {
    this.setState({
      confirmation: false
    });
  }

  handleCancel() {
    this.setState({
      confirmation: false
    });
  }

  handleUpgradeButton() {
    this.setState({
      confirmation: true
    });
  }

  render() {
    const {classes, className} = this.props;

    const msgUptoDate = (
      <div className={classes.upgradeStatus}>
        <Icon className={classNames(classes.upgradeIcon, classes.uptoDate)}>
          check_circle
        </Icon>
        Operos is up to date
      </div>
    );

    const msgNotUptodate = (
      <div className={classes.upgradeStatus}>
        <Icon className={classNames(classes.upgradeIcon, classes.notUptoDate)}>
          warning
        </Icon>
        <div>
          <div>
            Operos update available: v0.3.x
            (<a href="https://www.paxautoma.com/operos/docs/0.3.0" target="_blank">release notes</a>)
          </div>
          <div>Current version: v0.2.x</div>
        </div>
      </div>
    );

    return (
      <Card className={className}>
        <CardContent>
          <Typography type="headline" component="h2">
            Software upgrade
          </Typography>
          {msgNotUptodate}
        </CardContent>
        <CardActions>
          <Button dense color="primary" onClick={() => this.handleUpgradeButton()}>
            Upgrade Operos
          </Button>
        </CardActions>
        {this.renderConfirmationDialog()}
      </Card>
    );
  }

  renderConfirmationDialog() {
    return (
      <Dialog
          open={this.state.confirmation}
          onClose={() => this.handleCancel()}
      >
        <DialogTitle>Begin upgrade?</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Please click "Begin" to start the upgrade process. Each machine
            in the cluster will be rebooted, one by one. The controller will
            also be rebooted. During this time, the controller will be
            unavailable for a brief period.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => this.handleConfirm()} color="primary">
            Begin
          </Button>
          <Button onClick={() => this.handleCancel()} color="secondary">
            Cancel
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
}

export default withStyles(styles)(UpgradeCard);
