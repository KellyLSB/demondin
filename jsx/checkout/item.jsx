import React from 'react'
import { Grid, Header, Segment, Button, Icon, Label } from 'semantic-ui-react'

import { CartContext } from './cart_context'
import ItemOption from './item_option'

export default class Item extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
			hideForm: true,
      data: {}
    };

    if ('data' in props) {
      this.state.data = props.data;
    }

		this.onToggleForm = this.onToggleForm.bind(this)
  }

	onToggleForm() {
		this.setState((state) => {
			state.hideForm = !state.hideForm;
		});

		this.render();
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
    var price = this.currentPrice();

		if (price) {		
			console.log(price);
    	if('taxable' in price && price.taxable) {
    	  return this.toDollars(price.price) + " + Tax"
    	}
    
    	return this.toDollars(price.price)
		}

		return "Undefined"
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
		                onClick={this.onToggleForm}>
		                <Icon name='shop' />
		                Add To Cart
		              </Button>
		            )}
		          </CartContext.Consumer>
		        </Grid.Column>
					</Grid.Row>
					{ !this.state.hideForm ? (
						<Grid.Row>
							<Grid.Column columns={2}>
								{this.state.data.options.map((option) => 
									<ItemOption key={option.key} option={option} />
								)}
							</Grid.Column>
						</Grid.Row>
					) : null }
        </Grid>
      </Segment>
    );
  }
}
