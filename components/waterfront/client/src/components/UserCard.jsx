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
import {graphql, gql, compose} from 'react-apollo';

import Avatar from 'material-ui/Avatar';
import Divider from 'material-ui/Divider';
import Icon from 'material-ui/Icon';
import Card, {CardHeader} from 'material-ui/Card';
import List, {ListItem, ListItemIcon, ListItemText} from 'material-ui/List';
import {withRouter} from 'react-router';
import {withStyles} from 'material-ui/styles';

const styles = {
  card: {
    minWidth: 200
  }
};

class UserCard extends React.Component {
  constructor() {
    super();
    this.state = {
      loggingOut: false
    };
  }

  handleLogout() {
    this.setState({
      loggingOut: true
    });

    this.props.mutate().catch(err => {
      this.setState({
        loggingOut: false
      });
    });

    this.props.closePopover();
  }

  handleCredentials() {
    this.props.history.push('/access');
    this.props.closePopover();
  }

  render() {
    const {classes, data: {login_info: {user}}} = this.props;

    return (
      <Card elevation={0} className={classes.card}>
        <CardHeader
            avatar={<Avatar><Icon>person</Icon></Avatar>}
            title={user.username}
        />
        <Divider />
        <List>
          <ListItem button>
            <ListItemIcon><Icon>lock</Icon></ListItemIcon>
            <ListItemText primary="Credentials" onClick={() => this.handleCredentials()} />
          </ListItem>
          <ListItem button onClick={() => this.handleLogout()} disabled={this.state.loggingOut}>
            <ListItemIcon><Icon>exit_to_app</Icon></ListItemIcon>
            <ListItemText primary="Log out" />
          </ListItem>
        </List>
      </Card>
    );
  }
}

export default withRouter(withStyles(styles)(compose(
  graphql(gql`
    query {
      login_info {
        user {
          username
        }
      }
    }
  `),
  graphql(gql`
    mutation logout {
      logout {
        logged_in
      }
    }  
  `)
)(UserCard)));
