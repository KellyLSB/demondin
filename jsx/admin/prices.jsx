import React from 'react'
import { Button } from 'semantic-ui-react';
import Price from './price'

export default class Prices extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      prices: [],
      item: null,
    };

    if ('prices' in props) {
      this.state.prices = props.prices;
    }

    if ('item' in props) {
      this.state.item = props.item;
    }

    this.addPrice = this.addPrice.bind(this)
  }

  addPrice() {
    console.log("POKEY")
    this.setState((state) => {
      state.prices.push({
        ItemID: this.props.item.id,
      })

      return state
    })
  }

  render() {
    return (
      <React.Fragment>
        {this.state.prices.map((price, altID) =>
          <Price key={price.id ? price.id : altID}
            item={this.state.item}
            price={price} />
        )}

        <Button onClick={this.addPrice}>
          Add Another Price
        </Button>
      </React.Fragment>
    )
  }
}
