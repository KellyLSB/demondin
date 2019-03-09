import React from 'react'
import { Grid, Header, Segment, Button, Icon, Label } from 'semantic-ui-react'

import { CartContext } from './cart_context'

export default class Item extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: {}
    };

    if ('data' in props) {
      this.state.data = props.data;
    }
  }
  
  toDollars(cents) {
    return (cents / 100).toLocaleString("en-US", {
      style:"currency", currency:"USD"
    })
  }
    
  toCents(dollars) {
    return Math.round(100 * parseFloat(
      `${dollars}`.replace(/[$,]/g, '')
    ))
  }
  
  currentPrice() {
    return this.state.data.prices[0];
  }
  
  printPrice() {
    var price = this.currentPrice()
    
    if(price.taxable) {
      return this.toDollars(price.price) + " + Tax"
    }
    
    return this.toDollars(price.price)
  }

  render() {
    return (
      <Segment>
        <Header size="huge">{this.state.data.name}</Header>
        <span>{this.state.data.description}</span>
        <Grid columns={2}>
					<Grid.Row>
		        <Grid.Column>
		          <Header sub>Price</Header>
		          <Label tag>{this.printPrice()}</Label>
		        </Grid.Column>
		        <Grid.Column textAlign='right'>
		          <CartContext.Consumer>
		            {(cartContext) => (
		              <Button
		                primary icon labelPosition='left'
		                onClick={cartContext.addToCart(this.state.data.ID)}>
		                <Icon name='shop' />
		                Add To Cart
		              </Button>
		            )}
		          </CartContext.Consumer>
		        </Grid.Column>
					</Grid.Row>
					<Grid.Row>
						<Grid.Column columns={2}>
							{JSON.stringify(this.state.data.options)}
							{this.state.data.options.map((option) => {
								if (option.valueType == "select") return (
									<React.Fragment>
										<label for={option.key}>{option.key}</label>
										<select name={option.key}>
											{option.values.map((value) => <option>{value}</option>)}
										</select>
									</React.Fragment>
								);

								return (
									<React.Fragment>
										<label for={option.key}>{option.key}</label>
										<input name={option.key}></input>
									</React.Fragment>
								);
							})}
						</Grid.Column>
					</Grid.Row>
        </Grid>
      </Segment>
    );
  }
}
