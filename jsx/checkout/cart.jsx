import React from 'react'

import Item from './item'
import { CartContext } from './cart_context'

import { Header, Grid } from 'semantic-ui-react'

export default class Items extends React.Component {
  static contextType = CartContext
  
  constructor(props) {
    super(props);

    this.state = {
      data: []
    };

    if ('data' in props) {
      this.state.data = props.data;
    } else {
      this.fetchItems();
    }
  }

  fetchItems() {
    fetch("/shop/invoicing.json", {
      method: "GET"
    }).then((response) => response.json()).then((data) => {
      this.setState((state) => {
        state.data = data
        return state
      })
    })
  }

  render() {
    return (
      <Grid.Row>
        <Header>Badges</Header>
        {this.state.data.Badges.map((badge) =>
          <Badge key={badge.ID} data={badge} />
        )}
      </Grid.Row>
    )
  }
}
